package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/postmannen/go-bill/data"
)

//The web handler to the user selection in create bills
func (d *webData) webBillSelectUser(w http.ResponseWriter, r *http.Request) {
	d.Users = data.QueryAllUserInfo(d.PDB)
	ip := r.RemoteAddr

	//creates the header and the select box from templates
	err := d.tpl.ExecuteTemplate(w, "createBillCompletePage", d)
	if err != nil {
		log.Println("createBillUserSelection: template execution error = ", err)
	}

	//Parse all the variables in the html form to get all the data
	r.ParseForm()

	//'if' sentence to keep the chosen user ID. Reason is that it resets to 0 when the page is redrawn after "choose" is pushed
	//put the value in chooseUserButton which is a global variable
	if r.FormValue("chooseUserButton") == "choose" {
		//Get the value (numberPart) of the chosen user from form dropdown menu <select name="users">
		d.ActiveUserID, _ = strconv.Atoi(r.FormValue("users"))
		//reset data.CurrentBillID so a new user dont inherit the last bill used for another user.
		d.CurrentBillID = 0
	}

	log.Println("billCreateWeb: The numberPart chosen in the user select box:data.activeUserID = ", d.ActiveUserID)

	//Iterate the slice of struct for all the users found in db to find the data for the user selected
	for i := range d.Users {
		if d.Users[i].Number == d.ActiveUserID && d.Users[i].Number != 0 {
			log.Println(ip, "modifyUsersWeb: p[i].FirstName, p[i].LastName, found = ", d.Users[i].FirstName, d.Users[i].LastName)
			//Store the index nr in slice of the chosen user
			d.IndexUser = i
			//store all the info of the current user in the struct for feeding variables to the templates
			d.CurrentUser = d.Users[i]
			err := d.tpl.ExecuteTemplate(w, "billShowUser", d)
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
		err = d.tpl.ExecuteTemplate(w, "redirectToEditBill", "some data")
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
	}
}

func (d *webData) webBillLines(w http.ResponseWriter, r *http.Request) {
	BillsForUser := data.QueryBillsForUser(d.PDB, d.ActiveUserID)

	//Sort the bills so the last bill_id is first in the slice, and then shown on top of the listing
	BillsForUser = sortBills(BillsForUser)

	//draw the bill select box in the window
	err := d.tpl.ExecuteTemplate(w, "billSelectBox", BillsForUser)
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
	for i, v := range BillsForUser {
		if v.BillID == d.CurrentBillID {
			CurrentBill = BillsForUser[i]
		}
	}

	//update the total sums in main bill, and write it to db
	updateBillTotalExVat(&CurrentBill, d.CurrentBillID, storedBillLines)
	updateBillTotalIncVat(&CurrentBill, storedBillLines)
	data.UpdateBill(d.PDB, CurrentBill)
	//TESTING
	d.CurrentBill = CurrentBill

	err = d.tpl.ExecuteTemplate(w, "showBillInfo", d.CurrentBill)
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

			//doing a redirect so it redraws the page with the new line.
			err = d.tpl.ExecuteTemplate(w, "redirectToEditBill", "some data")
			if err != nil {
				log.Println("createBillUserSelection: createBillLines: template execution error = ", err)
			}
		}
	}

	updateLineExVatTotal(storedBillLines)
	//store all the bill lines in bill_lines db, to get ex vat total written to db
	data.UpdateBillLine(d.PDB, storedBillLines)

	d.CurrentBillLines = storedBillLines
	err = d.tpl.ExecuteTemplate(w, "createBillLines", d)
	if err != nil {
		log.Println("createBillUserSelection: createBillLines: template execution error = ", err)
	}

	r.ParseForm()

	//The name of the buttons are postfixed with LineID. Separate the numbers and the letters
	//from the elements in the map of r.Form, to get the ID of which LineID the button belonged to
	buttonValue, buttonNumbers := separateStrNumForButton(r)

	//Add a new billLine to db, and redraw window.
	if buttonValue == "add" {
		billLine := data.BillLines{}
		billLine.BillID = d.CurrentBillID
		fmt.Println("#######billid some benyttes er =", d.CurrentBillID)
		//create a random numberPart for the bill line....for now....
		rand.Seed(time.Now().UnixNano())
		billLine.LineID = rand.Intn(10000)
		billLine.Description = "noe tekst"
		data.AddBillLine(d.PDB, billLine)
		//doing a redirect so it redraws the page with the new line.
		err = d.tpl.ExecuteTemplate(w, "redirectToEditBill", "some data")
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
		err = d.tpl.ExecuteTemplate(w, "redirectToEditBill", "some data")
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

		err = d.tpl.ExecuteTemplate(w, "redirectToEditBill", "some data")
		if err != nil {
			log.Println("createBillUserSelection: createBillLines: template execution error = ", err)
		}
	}
}

//printBill
func (d *webData) printBill(w http.ResponseWriter, r *http.Request) {
	d.CurrentAdmin = data.QuerySingleUserInfo(d.PDB, 0)
	d.CurrentBill.TotalVat = d.CurrentBill.TotalIncVat - d.CurrentBill.TotalExVat
	err := d.tpl.ExecuteTemplate(w, "printBillComplete", d)
	if err != nil {
		log.Println("webBillLines: template execution error = ", err)
	}
}
