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
		newBill.CreatedDate = fmt.Sprint(t.Format("2006-01-02 15:04:05"))
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
	for i := 0; i < len(BillsForUser); i++ {
		for ii := 0; ii < len(BillsForUser); ii++ {
			if BillsForUser[i].BillID > BillsForUser[ii].BillID {
				tmp := BillsForUser[ii]
				BillsForUser[ii] = BillsForUser[i]
				BillsForUser[i] = tmp
			}
		}
	}

	//draw the bill select box in the window
	err := tmpl["init.html"].ExecuteTemplate(w, "billLinesComplete", BillsForUser)
	if err != nil {
		log.Println("webBillLines: template execution error = ", err)
	}

	r.ParseForm()
	fmt.Println("r.Form = ", r.Form)

	if r.FormValue("userActionButton") == "choose bill" {
		billID, _ := strconv.Atoi(r.FormValue("billID"))
		log.Println("INFO: webBillLines: billID =", billID)
		fmt.Println("billID = ", billID)
		d.CurrentBillID = billID

	}

	//get all the billLines for current billID
	billLines := data.QueryBillLines(d.PDB, d.CurrentBillID)

	//create all the billLines on the screen
	err = tmpl["init.html"].ExecuteTemplate(w, "createBillLines", billLines)
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

	//WORKING BELOW HERE
	fmt.Println("\n---Current FORM VALUES", r.Form)

	for _, v := range billLines {
		fmt.Println("--- iterating original values : ", v.BillID, v.LineID, v.Description)
	}

	r.ParseForm()
	var numbers []int
	//find the all the unique billLine numbers in the form, and store them in numbers[]
	for k, v := range r.Form {
		fmt.Printf("--- k = %v of type %T , and v = %v of type %T\n", k, k, v, v)
		reL := regexp.MustCompile("[a-zA-Z]+")
		reN := regexp.MustCompile("[0-9]+")
		letter := reL.FindString(k)
		numberStr := reN.FindString(k)
		number, _ := strconv.Atoi(numberStr)
		fmt.Printf("-----letter = %v, and number = %v\n", letter, number)

		found := false
		//check if number is allready in the numbers slice, if NOT.....add it
		for _, vv := range numbers {
			fmt.Printf("***trying to compare vv=%v and number=%v \n", vv, number)
			if number == vv {
				found = true
				fmt.Println("The numbers are equal")
			}
		}
		if !found {
			numbers = append(numbers, number)
		}
	}
	fmt.Println("numbers = ", numbers)

	//fill a tmp slice of data.BillLines struct with the values from the http request
	var TMPlines data.BillLines
	var TMPbillLines []data.BillLines
	for _, num := range numbers {
		for k, v := range r.Form {
			fmt.Printf("--- k = %v of type %T , and v = %v of type %T\n", k, k, v, v)
			reL := regexp.MustCompile("[a-zA-Z]+")
			reN := regexp.MustCompile("[0-9]+")
			letter := reL.FindString(k)
			numberStr := reN.FindString(k)
			number, _ := strconv.Atoi(numberStr)
			//compare all the unique line numbers in the numbers[] slice with all the numbers
			//postfixed at the end of the http Request input name parameters.
			//if found add value to the tmp struct of type data.BillLines
			if num == number {
				TMPlines.BillID = d.CurrentBillID
				if letter == "billLineID" {
					myVal, err := strconv.Atoi(v[0])
					if err != nil {
						fmt.Println("ERROR: strconv billLineID : ", err)
					}
					TMPlines.LineID = myVal
				}
				if letter == "billLineDescription" {
					TMPlines.Description = v[0]
				}
				if letter == "billLineQuantity" {
					myVal, err := strconv.Atoi(v[0])
					if err != nil {
						fmt.Println("ERROR: strconv billLineQuantity : ", err)
					}
					TMPlines.Quantity = myVal
				}
				if letter == "billLineDiscountPercentage" {
					myVal, err := strconv.Atoi(v[0])
					if err != nil {
						fmt.Println("ERROR: strconv billLinePercentage : ", err)
					}
					TMPlines.DiscountPercentage = myVal
				}
				if letter == "billLineVatUsed" {
					myVal, err := strconv.Atoi(v[0])
					if err != nil {
						fmt.Println("ERROR: strconv billLineVatUsed : ", err)
					}
					TMPlines.VatUsed = myVal
				}
				if letter == "billLinePriceExVat" {
					myVal, err := strconv.ParseFloat(v[0], 64)
					if err != nil {
						fmt.Println("ERROR: strconv billLinePriceExVat : ", err)
					}
					TMPlines.PriceExVat = myVal
				}
			}
		}
		TMPbillLines = append(TMPbillLines, TMPlines)
	}
	fmt.Println("-*- TMPbillLines : ", TMPbillLines)
	fmt.Println("-*-    billLines : ", billLines)

	//going to compare this slice with the original values from DB, to know what to update
	//range over the numbers slice to do the comparison
	for _, num := range numbers {
		for _, line := range billLines {
			if line.LineID == num {
				for _, line2 := range TMPbillLines {
					if line2.LineID == num {
						if line.LineID != line2.LineID {
							fmt.Printf("LineID for Line %v have changed to %v\n", num, line2.LineID)
						}
						if line.Description != line2.Description {
							fmt.Printf("Description for Line %v have changed to %v\n", num, line2.Description)
						}
						if line.Quantity != line2.Quantity {
							fmt.Printf("Quantity for Line %v have changed to %v\n", num, line2.Quantity)
						}
						if line.DiscountPercentage != line2.DiscountPercentage {
							fmt.Printf("DiscountPercentage for Line %v have changed to %v\n", num, line2.DiscountPercentage)
						}
						if line.VatUsed != line2.VatUsed {
							fmt.Printf("VatUsed for Line %v have changed to %v\n", num, line2.VatUsed)
						}
						if line.PriceExVat != line2.PriceExVat {
							fmt.Printf("PriceExVat for Line %v have changed to %v\n", num, line2.PriceExVat)
						}
					}
				}
			}
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
