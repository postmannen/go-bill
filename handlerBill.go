package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/postmannen/fakt/data"
)

//The web handler to the user selection in create bills
func (d *webData) webBillSelectUser(w http.ResponseWriter, r *http.Request) {
	d.Users = data.QueryAllUserInfo(d.PDB)
	ip := r.RemoteAddr

	//creates the header and the select box from templates
	err := tmpl["init.html"].ExecuteTemplate(w, "createBillCompletePage", d)
	if err != nil {
		log.Println("createBillUserSelection: template execution error = ", err)
	}

	//Parse all the variables in the html form to get all the data
	r.ParseForm()

	//'if' sentence to keep the chosen user ID. Reason is that it resets to 0 when the page is redrawn after "choose" is pushed
	//put the value in chooseUserButton which is a global variable
	if r.FormValue("chooseUserButton") == "choose" {
		//Get the value (number) of the chosen user from form dropdown menu <select name="users">
		d.ActiveUserID, _ = strconv.Atoi(r.FormValue("users"))
		//reset data.CurrentBillID so a new user dont inherit the last bill used for another user.
		d.CurrentBillID = 0
	}

	log.Println("billCreateWeb: The number chosen in the user select box:data.activeUserID = ", d.ActiveUserID)

	//Iterate the slice of struct for all the users found in db to find the data for the user selected
	for i := range d.Users {
		if d.Users[i].Number == d.ActiveUserID && d.Users[i].Number != 0 {
			log.Println(ip, "modifyUsersWeb: p[i].FirstName, p[i].LastName, found = ", d.Users[i].FirstName, d.Users[i].LastName)
			//Store the index nr in slice of the chosen user
			d.indexNR = i
			err := tmpl["init.html"].ExecuteTemplate(w, "billShowUserSingle", d.Users[i])
			if err != nil {
				log.Println(ip, "modifyUsersWeb: error = ", err)
			}
		}
	}

	//Get the last used bill_id from DB
	highestBillNR, totalLineCount := data.QueryLastBillID(d.PDB)
	log.Println(ip, "billCreateWeb: highestBillNR = ", highestBillNR, ", and totaltLineCount = ", totalLineCount)

	//Check which of the two input buttons where pushed. They both have name=userActionButton,
	//and the value can be read with r.FormValue("userActionButton")
	r.ParseForm()
	log.Println("r.Form shows = ", r.Form)
	buttonAction := r.FormValue("userActionButton")
	log.Println(ip, "billCreateWeb: userActionButton = ", buttonAction)

	//if the manage bills button were pushed
	if buttonAction == "manage bills" {
		err = tmpl["init.html"].ExecuteTemplate(w, "redirectToEditBill", "some data")
		if err != nil {
			log.Println("createBillUserSelection: createBillLines: template execution error = ", err)
		}
	}

	if buttonAction == "add new bill" {
		fmt.Println(buttonAction, "pressed")

		//get the last used bill id
		highestBillNR, totalLineCount := data.QueryLastBillID(d.PDB)
		log.Println("billCreateWeb: highestBillNR = ", highestBillNR, ", totaltLineCount = ", totalLineCount)

		newBill := data.Bill{}
		newBill.BillID = highestBillNR + 1
		newBill.UserID = d.ActiveUserID
		t := time.Now()
		//										   yyyy-mm-dd
		newBill.CreatedDate = fmt.Sprint(t.Format("2006-01-02"))
		//create a new bill and return the new billID to use later
		d.CurrentBillID = data.AddBill(d.PDB, newBill)
		log.Println("billCreateWeb: newBillID = ", d.CurrentBillID)

		billLine := data.BillLines{}
		billLine.BillID = d.CurrentBillID
		billLine.LineID = 1
		billLine.Description = "noe tekst"

		data.AddBillLine(d.PDB, billLine)
	}
}

//**********************************************************

func (d *webData) webBillLines(w http.ResponseWriter, r *http.Request) {
	fmt.Println("INFO: webBillLines: Active user ID when call for bills = ", d.ActiveUserID)
	BillsForUser := data.QueryBillsForUser(d.PDB, d.ActiveUserID)
	fmt.Println("INFO: webBillLines: BillsForUser = ", BillsForUser)

	//Sort the bills so the last bill_id is first in the slice, and then shown on top of the listing
	BillsForUser = sortBills(BillsForUser)

	//draw the bill select box in the window
	err := tmpl["init.html"].ExecuteTemplate(w, "billLinesComplete", BillsForUser)
	if err != nil {
		log.Println("webBillLines: template execution error = ", err)
	}

	r.ParseForm()
	fmt.Println("r.Form = ", r.Form)

	if r.FormValue("userActionButton") == "choose bill" {
		d.CurrentBillID, _ = strconv.Atoi(r.FormValue("billID"))
	}

	//get all the billLines for current billID from db
	currentBillsLines := data.QueryBillLines(d.PDB, d.CurrentBillID)

	//Find all the data on the current bill id
	var CurrentBill data.Bill
	for i, v := range BillsForUser {
		if v.BillID == d.CurrentBillID {
			CurrentBill = BillsForUser[i]
		}
	}

	//get the total sum ex vat of all the bill lines
	CurrentBill.TotalExVat = 0
	for _, v := range currentBillsLines {
		CurrentBill.TotalExVat += v.PriceExVat
	}

	//add the TotalExVat to db here
	if CurrentBill.TotalExVat != 0 {
		data.UpdateBillPriceExVat(d.PDB, CurrentBill.TotalExVat, CurrentBill.BillID)
	}

	//draw the bill select box in the window
	err = tmpl["init.html"].ExecuteTemplate(w, "showBills", CurrentBill)
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
			fmt.Printf("\n-----******------- Changed = %v\n\n", changed)
			fmt.Printf("-*-Orig  bill lines = %v\n", CurrentBill)
			fmt.Printf("-*- Tmp  bill lines = %v\n\n", tmpBill)
			data.UpdateBill(d.PDB, CurrentBill)

			//doing a redirect so it redraws the page with the new line. Not sure if this is the best way....
			err = tmpl["init.html"].ExecuteTemplate(w, "redirectToEditBill", "some data")
			if err != nil {
				log.Println("createBillUserSelection: createBillLines: template execution error = ", err)
			}
		}
	}
	//WORKING ABOVE HERE *********************

	//create all the billLines on the screen
	err = tmpl["init.html"].ExecuteTemplate(w, "createBillLines", currentBillsLines)
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
		//create a random number for the bill line....for now....
		rand.Seed(time.Now().UnixNano())
		billLine.LineID = rand.Intn(10000)
		billLine.Description = "noe tekst"
		data.AddBillLine(d.PDB, billLine)
		//doing a redirect so it redraws the page with the new line. Not sure if this is the best way....
		err = tmpl["init.html"].ExecuteTemplate(w, "redirectToEditBill", "some data")
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
		err = tmpl["init.html"].ExecuteTemplate(w, "redirectToEditBill", "some data")
		if err != nil {
			log.Println("createBillUserSelection: createBillLines: template execution error = ", err)
		}
	}

	modifyButtonPushed := false
	if buttonValue == "modify" {
		modifyButtonPushed = true
	}

	//for _, v := range billLines {
	//	fmt.Println("--- iterating original values : ", v.BillID, v.LineID, v.Description)
	//}

	r.ParseForm()
	var numbers []int
	//find the all the unique billLine numbers in the form, and store them in numbers[]
	for k, v := range r.Form {
		fmt.Printf("--- k = %v of type %T , and v = %v of type %T\n", k, k, v, v)
		reLetters := regexp.MustCompile("[a-zA-Z]+")
		reNum := regexp.MustCompile("[0-9]+")
		letter := reLetters.FindString(k)
		numberStr := reNum.FindString(k)
		number, _ := strconv.Atoi(numberStr)
		log.Printf("-----letter = %v, and number = %v\n", letter, number)

		found := false
		//check if number is allready in the numbers slice, if NOT.....add it
		for _, vv := range numbers {
			//fmt.Printf("***trying to compare vv=%v and number=%v \n", vv, number)
			if number == vv {
				found = true
				fmt.Println("The numbers are equal")
			}
		}
		if !found {
			numbers = append(numbers, number)
		}
	}

	//-------Edit the bill lines------------
	//fill a tmp slice of data.BillLines struct with the values from the http request
	var TMPlines data.BillLines
	var TMPbillLines []data.BillLines
	//itarate the unique bill line numbers
	for _, num := range numbers {
		//iterate all the data in form
		for k, v := range r.Form {
			reLetters := regexp.MustCompile("[a-zA-Z]+")
			reNum := regexp.MustCompile("[0-9]+")
			letter := reLetters.FindString(k)
			numberStr := reNum.FindString(k)
			number, _ := strconv.Atoi(numberStr)
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
			//compare all the unique line numbers in the numbers[] slice with all the numbers
			//postfixed at the end of the http Request input name parameters.
			//if found add value to the tmp struct of type data.BillLines
			if num == number {
				TMPlines.BillID = d.CurrentBillID
				if letter == "billLineID" {
					TMPlines.LineID = v1
				}
				if letter == "billLineDescription" {
					TMPlines.Description = v[0]
				}
				if letter == "billLineQuantity" {
					TMPlines.Quantity = v1
				}
				if letter == "billLineDiscountPercentage" {
					TMPlines.DiscountPercentage = v1
				}
				if letter == "billLineVatUsed" {
					TMPlines.VatUsed = v1
				}
				if letter == "billLinePriceExVat" {
					myVal, err := strconv.ParseFloat(v[0], 64)
					if err != nil {
						log.Println("ERROR: strconv billLinePriceExVat : ", err)
					}
					TMPlines.PriceExVat = myVal
				}
			}
		}
		TMPbillLines = append(TMPbillLines, TMPlines)
	}
	log.Println("-*- TMPbillLines : ", TMPbillLines)
	log.Println("-*-    billLines : ", currentBillsLines)

	//going to compare this slice with the original values from DB, to know what to update
	//range over the numbers slice to get all the unique line numbers
	//then range billLines, and range TMPbillLines to compare billLines.X with TMPbillLines.X
	changed = false
	for _, num := range numbers {
		for _, line := range currentBillsLines {
			if line.LineID == num {
				for _, line2 := range TMPbillLines {
					if line2.LineID == num {
						if line.LineID != line2.LineID {
							log.Printf("LineID for Line %v have changed to %v\n", num, line2.LineID)
							changed = true
						}
						if line.Description != line2.Description {
							log.Printf("Description for Line %v have changed to %v\n", num, line2.Description)
							changed = true
						}
						if line.Quantity != line2.Quantity {
							log.Printf("Quantity for Line %v have changed to %v\n", num, line2.Quantity)
							changed = true
						}
						if line.DiscountPercentage != line2.DiscountPercentage {
							log.Printf("DiscountPercentage for Line %v have changed to %v\n", num, line2.DiscountPercentage)
							changed = true
						}
						if line.VatUsed != line2.VatUsed {
							log.Printf("VatUsed for Line %v have changed to %v\n", num, line2.VatUsed)
							changed = true
						}
						if line.PriceExVat != line2.PriceExVat {
							log.Printf("PriceExVat for Line %v have changed to %v\n", num, line2.PriceExVat)
							changed = true
						}
					}
				}
			}
		}
	}
	if changed && modifyButtonPushed {
		data.UpdateBillLine(d.PDB, TMPbillLines)

		err = tmpl["init.html"].ExecuteTemplate(w, "redirectToEditBill", "some data")
		if err != nil {
			log.Println("createBillUserSelection: createBillLines: template execution error = ", err)
		}
	}

	//Create a DB function with query to update the changed field
}

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
