package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/postmannen/go-bill/data"
)

//The web handler to the user selection in create bills
func (d *webData) webBillSelectUser() http.HandlerFunc {
	var init sync.Once
	var tpl *template.Template
	init.Do(func() {
		tpl = template.Must(template.ParseFiles("public/billTemplates.html"))
	})

	return func(w http.ResponseWriter, r *http.Request) {
		d.Users = data.QueryAllUserInfo(d.PDB)
		ip := r.RemoteAddr

		//creates the header and the select box from templates
		err := tpl.ExecuteTemplate(w, "createBillCompletePage", d)
		if err != nil {
			log.Println("createBillUserSelection: template execution error = ", err)
		}

		//Parse all the variables in the html form to get all the data
		r.ParseForm()

		//'if' sentence to keep the chosen user ID. Reason is that it resets to 0 when the page is redrawn after "choose" is pushed
		//put the value in chooseUserButton which is a global variable
		if r.FormValue("chooseUserButton") == "choose" {
			d.ActiveUserID, _ = strconv.Atoi(r.FormValue("users"))
			//reset data.CurrentBillID so a new user dont inherit the last bill used for another user.
			d.CurrentBillID = 0
		}

		log.Println("billCreateWeb: The numberPart chosen in the user select box:data.activeUserID = ", d.ActiveUserID)

		//Iterate the slice of struct for all the users found in db to find the data for the user selected
		for i := range d.Users {
			if d.Users[i].Number == d.ActiveUserID && d.Users[i].Number != 0 {
				log.Println(ip, "modifyUsersWeb: user p[i].FirstName, p[i].LastName, found = ", d.Users[i].FirstName, d.Users[i].LastName)
				//Store the index nr in slice of the chosen user
				d.IndexUser = i
				//store all the info of the current user in the struct for feeding variables to the templates
				d.CurrentUser = d.Users[i]
				err := tpl.ExecuteTemplate(w, "billShowUser", d)
				if err != nil {
					log.Println(ip, "modifyUsersWeb: error = ", err)
				}
			}
		}

		//Get the last used bill_id from DB
		highestBillID, totalLineCount := data.QueryLastBillID(d.PDB)
		log.Println(ip, "billCreateWeb: highestBillID = ", highestBillID, ", and totaltLineCount = ", totalLineCount)

		//Check which of the two input buttons where pushed. They both have name=userActionButton,
		//and the value can be read with r.FormValue("userActionButton")
		r.ParseForm()
		buttonAction := r.FormValue("userActionButton")

		//if the manage bills button were pushed
		if buttonAction == "manage bills" {
			err = tpl.ExecuteTemplate(w, "redirectToEditBill", nil)
			if err != nil {
				log.Println("createBillUserSelection: createBillLines: template execution error = ", err)
			}
		}

		if buttonAction == "add new bill" {
			fmt.Println(buttonAction, "pressed")

			//get the last used bill id
			highestBillID, totalLineCount := data.QueryLastBillID(d.PDB)
			log.Println("billCreateWeb: highestBillID = ", highestBillID, ", totaltLineCount = ", totalLineCount)

			newBill := data.Bill{}
			newBill.BillID = highestBillID + 1
			newBill.UserID = d.ActiveUserID
			t := time.Now()
			//										   yyyy-mm-dd
			newBill.CreatedDate = fmt.Sprint(t.Format("2006-01-02"))
			//create a new bill and return the new billID to use later
			d.CurrentBillID = data.AddBill(d.PDB, newBill)
			log.Println("billCreateWeb: newBillID = ", d.CurrentBillID)
		}
	}
}

func (d *webData) webBillLines() http.HandlerFunc {
	var init sync.Once
	var tpl *template.Template
	init.Do(func() {
		tpl = template.Must(template.ParseFiles("public/billTemplates.html"))
	})

	return func(w http.ResponseWriter, r *http.Request) {
		//exported/capitalized since its passed with the execution of the template below
		UserBills := data.QueryBillsForUser(d.PDB, d.ActiveUserID)

		//Sort the bills so the last bill_id is first in the slice, and then shown on top of the listing
		UserBills = sortBills(UserBills)

		//draw the bill select box in the window
		err := tpl.ExecuteTemplate(w, "billSelectBox", UserBills)
		if err != nil {
			log.Println("webBillLines: template execution error = ", err)
		}

		r.ParseForm()
		fmt.Println("r.Form = ", r.Form)

		if r.FormValue("userActionButton") == "choose bill" {
			d.CurrentBillID, _ = strconv.Atoi(r.FormValue("billID"))
		}

		//get all the billLines for current billID from db
		storedBillLines := data.QueryBillLines(d.PDB, d.CurrentBillID)

		//add a default nr.1 bill line if none exist
		if len(storedBillLines) == 0 && d.CurrentBillID != 0 {
			fmt.Fprintf(w, "Len was 0, value = %v\n", len(storedBillLines))
			billLine := data.BillLines{
				BillID: d.CurrentBillID,
				LineID: 1,
			}
			data.AddBillLine(d.PDB, billLine)
			//rerun gathering av bill line data for selected bill to get new data
			storedBillLines = data.QueryBillLines(d.PDB, d.CurrentBillID)
		}

		//Find all the data on the current bill id
		var CurrentBill data.Bill
		for i, v := range UserBills {
			if v.BillID == d.CurrentBillID {
				CurrentBill = UserBills[i]
			}
		}

		//update the total sums in main bill, and write it to db
		updateBillTotalExVat(&CurrentBill, d.CurrentBillID, storedBillLines)
		updateBillTotalIncVat(&CurrentBill, storedBillLines)
		data.UpdateBill(d.PDB, CurrentBill)
		//TESTING
		d.CurrentBill = CurrentBill

		err = tpl.ExecuteTemplate(w, "showBillInfo", d.CurrentBill)
		if err != nil {
			log.Println("webBillLines: template execution error = ", err)
		}

		r.ParseForm()

		//check all the data in r.Form,
		//create tmpBill of type data.Bill to hold all the bill data in r.Form
		var tmpBill data.Bill
		for k, v := range r.Form {
			reNumOnly := regexp.MustCompile("^[0-9]+$")

			//convert the string read from the r.Form into v to v1 of int which is used in struct
			var v1 int
			//check if the string only contains numbers
			if reNumOnly.Match([]byte(v[0])) {
				v1, err = strconv.Atoi(v[0])
				if err != nil {
					log.Printf("ERROR: strconv.Atoi for v[0] failed : %v", err)
				}
				log.Printf("\n---Conversion v1=%v %T and v[0]=%v %T \n\n", v1, v1, v[0], v[0])
			}
			if k == "CreatedDate" {
				tmpBill.CreatedDate = v[0]
			}
			if k == "DueDate" {
				tmpBill.DueDate = v[0]
			}
			if k == "Comment" {
				tmpBill.Comment = v[0]
			}
			if k == "Paid" {
				tmpBill.Paid = v1
			}
		}

		//compare the values of the bill struct from DB and the tmp struct from r.Form
		//to decide if to update DB with new values from the form
		changed := false
		if CurrentBill.Comment != tmpBill.Comment {
			changed = true
			CurrentBill.Comment = tmpBill.Comment
		}
		if CurrentBill.CreatedDate != tmpBill.CreatedDate {
			changed = true
			CurrentBill.CreatedDate = tmpBill.CreatedDate
		}
		if CurrentBill.DueDate != tmpBill.DueDate {
			changed = true
			CurrentBill.DueDate = tmpBill.DueDate
		}
		if CurrentBill.Paid != tmpBill.Paid {
			changed = true
			CurrentBill.Paid = tmpBill.Paid
		}

		if r.FormValue("billModifyButton") == "modify" {
			if changed {
				data.UpdateBill(d.PDB, CurrentBill)

				//doing a redirect so it redraws the page with the new line. Not sure if this is the best way....
				err = tpl.ExecuteTemplate(w, "redirectToEditBill", "some data")
				if err != nil {
					log.Println("createBillUserSelection: createBillLines: template execution error = ", err)
				}
			}
		}

		//if r.FormValue("billPrintView") == "print view" {
		//	err = tmpl["bill.html"].ExecuteTemplate(w,)
		//}

		updateLineExVatTotal(storedBillLines)
		//store all the bill lines in bill_lines db, to get ex vat total written to db
		data.UpdateBillLine(d.PDB, storedBillLines)

		d.CurrentBillLines = storedBillLines
		//create all the billLines on the screen
		err = tpl.ExecuteTemplate(w, "createBillLines", d)
		if err != nil {
			log.Println("createBillUserSelection: createBillLines: template execution error = ", err)
		}

		r.ParseForm()

		//The name of the buttons are postfixed with LineID. Separate the numbers and the letters from the map of r.Form
		//to get the ID of which LineID the button belonged to
		buttonValue, buttonNumbers := separateStrNumForButton(r)

		//using the buttonValue instead of r.FormValue since r.FormValue initiates a new parseform and
		//replaces the values from the last r.ParseForm
		//add a new billLine to db, and redraw window
		if buttonValue == "add" {
			billLine := data.BillLines{}
			billLine.BillID = d.CurrentBillID
			fmt.Println("#######billid some benyttes er =", d.CurrentBillID)
			//create a random numberPart for the bill line....for now....
			rand.Seed(time.Now().UnixNano())
			billLine.LineID = rand.Intn(10000)
			billLine.Description = "noe tekst"
			data.AddBillLine(d.PDB, billLine)
			//doing a redirect so it redraws the page with the new line. Not sure if this is the best way....
			err = tpl.ExecuteTemplate(w, "redirectToEditBill", "some data")
			if err != nil {
				log.Println("createBillUserSelection: createBillLines: template execution error = ", err)
			}

		}

		//delete a billLine and redraw window
		if buttonValue == "delete" {
			//num, err := strconv.Atoi(buttonNumbers)
			if err != nil {
				fmt.Printf("ERROR strconv.Atoi : %v\n", err)
			}
			data.DeleteBillLine(d.PDB, d.CurrentBillID, buttonNumbers)

			//doing a redirect so it redraws the page with the new line. Not sure if this is the best way....
			err = tpl.ExecuteTemplate(w, "redirectToEditBill", "some data")
			if err != nil {
				log.Println("createBillUserSelection: createBillLines: template execution error = ", err)
			}
		}

		modifyButtonPushed := false
		if buttonValue == "modify" {
			modifyButtonPushed = true
		}

		r.ParseForm()

		//find the all the unique billLine numbers in the form, and store them in []int
		lineNumbers := findBillLineNumbersInForm(r) //slice of all linenumbers in bill

		formBillLines := getBillLineFormValues(lineNumbers, r, d.CurrentBillID)
		log.Println("-*- formBillLines : ", formBillLines)
		log.Println("-*-    billLines : ", storedBillLines)

		//going to compare this slice with the original values from DB, to know what to update
		//range over the numbers slice to get all the unique line numbers
		//then range StoredBillLines, and range formBillLines to compare if any values are changed.
		changed = false
		changed = checkIfBillLineChanged(lineNumbers, storedBillLines, formBillLines)

		if changed && modifyButtonPushed {
			data.UpdateBillLine(d.PDB, formBillLines)

			err = tpl.ExecuteTemplate(w, "redirectToEditBill", "some data")
			if err != nil {
				log.Println("createBillUserSelection: createBillLines: template execution error = ", err)
			}
		}
	}
}

//printBill
func (d *webData) printBill() http.HandlerFunc {
	var init sync.Once
	var tpl *template.Template
	init.Do(func() {
		tpl = template.Must(template.ParseFiles("public/billTemplates.html"))
	})

	return func(w http.ResponseWriter, r *http.Request) {
		d.CurrentAdmin = data.QuerySingleUserInfo(d.PDB, 0)
		d.CurrentBill.TotalVat = d.CurrentBill.TotalIncVat - d.CurrentBill.TotalExVat
		err := tpl.ExecuteTemplate(w, "printBillComplete", d)
		if err != nil {
			log.Println("webBillLines: template execution error = ", err)
		}
	}
}

//------------------support functions----------------------------
//separateStrNumForButton , takes *http.Request as input, and returns string and int
func separateStrNumForButton(r *http.Request) (string, int) {
	var buttonNumbers string
	var buttonValue string
	var num int
	var err error
	for k, v := range r.Form {
		//fmt.Println("---VERDIER ---- ", k, " : ", v)
		re := regexp.MustCompile("[a-zA-Z]+")
		buttonLetters := re.FindString(k)
		re = regexp.MustCompile("[0-9]+")
		if buttonLetters == "billLineAddButton" {
			buttonValue = v[0]
			buttonNumbers = re.FindString(k)
			num, err = strconv.Atoi(buttonNumbers)
			if err != nil {
				log.Printf("ERROR: strconv.Atoi : %v", err)
			}
		}
		if buttonLetters == "billLineDeleteButton" {
			buttonValue = v[0]
			buttonNumbers = re.FindString(k)
			num, err = strconv.Atoi(buttonNumbers)
			if err != nil {
				log.Printf("ERROR: strconv.Atoi : %v", err)
			}
		}
		if buttonLetters == "billLineModifyButton" {
			buttonValue = v[0] //value is a slice of strings, get the first value
			buttonNumbers = re.FindString(k)
			num, err = strconv.Atoi(buttonNumbers)
			if err != nil {
				log.Printf("ERROR: strconv.Atoi : %v", err)
			}
		}
	}

	return buttonValue, num
}

//sortBills sorts the bills so the last bill_id is first in the slice
func sortBills(bills []data.Bill) []data.Bill {
	for i := 0; i < len(bills); i++ {
		for ii := 0; ii < len(bills); ii++ {
			if bills[i].BillID > bills[ii].BillID {
				tmp := bills[ii]
				bills[ii] = bills[i]
				bills[i] = tmp
			}
		}
	}
	return bills
}

//Finds all the numbers used in html names in form
//used to get all the unique bill lines
func findBillLineNumbersInForm(r *http.Request) (numbers []int) {
	for k, v := range r.Form {
		fmt.Printf("--- k = %v of type %T , and v = %v of type %T\n", k, k, v, v)
		reLetters := regexp.MustCompile("[a-zA-Z]+")
		reNum := regexp.MustCompile("[0-9]+")
		letterPart := reLetters.FindString(k)
		numberStr := reNum.FindString(k)
		numberPart, _ := strconv.Atoi(numberStr)
		log.Printf("-----letterPart = %v, and numberPart = %v\n", letterPart, numberPart)

		found := false
		//check if numberPart is allready in the numbers slice, if NOT.....add it
		for _, vv := range numbers {
			//fmt.Printf("***trying to compare vv=%v and numberPart=%v \n", vv, numberPart)
			if numberPart == vv {
				found = true
				fmt.Println("The numbers are equal")
			}
		}
		if !found {
			numbers = append(numbers, numberPart)
		}
	}
	return numbers
}

//range billLines to compare storedLines.X with formBillLines.X to see if any values have changed. Return changed = true if changed
func checkIfBillLineChanged(lineNRs []int, storedLines []data.BillLines, formLines []data.BillLines) (changed bool) {
	for _, num := range lineNRs {
		for _, line := range storedLines {
			if line.LineID == num {
				for _, line2 := range formLines {
					if line2.LineID == num {
						if line.LineID != line2.LineID {
							changed = true
						}
						if line.Description != line2.Description {
							changed = true
						}
						if line.Quantity != line2.Quantity {
							changed = true
						}
						if line.DiscountPercentage != line2.DiscountPercentage {
							changed = true
						}
						if line.VatUsed != line2.VatUsed {
							changed = true
						}
						if line.PriceExVat != line2.PriceExVat {
							changed = true
						}
					}
				}
			}
		}
	}
	return changed
}

//getBillLineFormValues parses all the data in the form, compares them with the current billID, and returns
//a slice with all the values entered in the form.
//all fields and buttons in the form have name values postfixed with the {{.LineID}}, so this function
//separates the first part of the name and the {{.LineID}} to know what fields to update
func getBillLineFormValues(lineNumbers []int, r *http.Request, billID int) (formBillLines []data.BillLines) {
	var tempLines data.BillLines

	for _, num := range lineNumbers {
		fmt.Println("-$- Outerloop, num of lineNumbers = ", num)
		//iterate all the data in form
		for k, v := range r.Form {

			//split out the letter and number part of button name
			reLetters := regexp.MustCompile("[a-zA-Z]+")
			reNum := regexp.MustCompile("[0-9]+")
			letterPart := reLetters.FindString(k)    //get name "billLineModifyButton"
			numberStr := reNum.FindString(k)         //get the line nr.that the button belonged to. Nr is postfixed in the name
			numberPart, _ := strconv.Atoi(numberStr) //convert the nr got`en from form to int, so it can be used later

			if num == numberPart {
				tempLines.BillID = billID
				if letterPart == "billLineID" {
					myVal, err := strconv.Atoi(v[0])
					if err != nil {
						log.Println("ERROR: strconv billLineID : ", err)
					}
					tempLines.LineID = myVal
					fmt.Printf("--- templLines.BillID er satt til %v\n", tempLines.BillID)
				}
				if letterPart == "billLineDescription" {

					tempLines.Description = v[0]
					fmt.Printf("--- templLines.Description er satt til %v\n", v[0])
				}
				if letterPart == "billLineQuantity" {
					myVal, err := strconv.Atoi(v[0])
					if err != nil {
						log.Println("ERROR: strconv billLineQuantity : ", err)
					}
					tempLines.Quantity = myVal
					fmt.Printf("--- templLines.Quantity er satt til %v\n", tempLines.Quantity)
				}
				if letterPart == "billLineDiscountPercentage" {
					myVal, err := strconv.Atoi(v[0])
					if err != nil {
						log.Println("ERROR: strconv billLineDiscountPercentage : ", err)
					}
					tempLines.DiscountPercentage = myVal
					fmt.Printf("--- templLines.DiscountPercentage er satt til %v\n", tempLines.DiscountPercentage)
				}
				if letterPart == "billLineVatUsed" {
					myVal, err := strconv.Atoi(v[0])
					if err != nil {
						log.Println("ERROR: strconv billLineVatUsed : ", err)
					}
					tempLines.VatUsed = myVal
					fmt.Printf("--- templLines.VatUsed er satt til %v\n", tempLines.VatUsed)
				}
				if letterPart == "billLinePriceExVat" {
					myVal, err := strconv.ParseFloat(v[0], 64)
					if err != nil {
						log.Println("ERROR: strconv billLinePriceExVat : ", err)
					}
					tempLines.PriceExVat = myVal
				}
			}
		}
		formBillLines = append(formBillLines, tempLines)
	}
	return
}

//updateBillTotalExVat updates the bill field total price ex vat,
//also writes the update info to correct field in db
func updateBillTotalExVat(bill *data.Bill, billID int, billLines []data.BillLines) {
	//TODO: Fix so the total values are made from line total to avoid doing the same calculations in several functions
	bill.TotalExVat = 0
	for _, v := range billLines {
		v.PriceExVat *= float64(v.Quantity)
		v.PriceExVat -= v.PriceExVat / 100 * float64(v.DiscountPercentage)
		bill.TotalExVat += v.PriceExVat
	}

	/* NOTE : Removing db db write, will put in a write for all fields when all data is processed
	//add the TotalExVat to db here
	if bill.TotalExVat != 0 {
		data.UpdateBillPriceExVat(pDB, bill.TotalExVat, bill.BillID)
	}
	*/
}

//updateBillTotalIncVat updates the bill field total price ex vat,
//also writes the update info to correct field in db
func updateBillTotalIncVat(bill *data.Bill, billLines []data.BillLines) {
	//TODO: Fix so the total values are made from line total to avoid doing the same calculations in several functions
	var lineIncVat float64
	bill.TotalIncVat = 0
	for _, v := range billLines {
		if v.VatUsed == 0 {
			lineIncVat = v.PriceExVat
			lineIncVat -= lineIncVat / 100 * float64(v.DiscountPercentage)
		} else {
			lineIncVat = v.PriceExVat + (v.PriceExVat / 100.0 * float64(v.VatUsed))
			lineIncVat -= lineIncVat / 100 * float64(v.DiscountPercentage)
		}
		lineIncVat *= float64(v.Quantity)
		bill.TotalIncVat += lineIncVat

	}
}

//updateLineExVatTotal updates the struct data for the total ex vat pr. line
func updateLineExVatTotal(b []data.BillLines) {
	for i := 0; i < len(b); i++ {
		sum := b[i].PriceExVat * float64(b[i].Quantity)
		sum = sum - (sum / 100 * float64(b[i].DiscountPercentage))
		b[i].PriceExVatTotal = sum
	}
}