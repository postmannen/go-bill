package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

/****************************
*	WEBHANDLERS				*
*							*
****************************/

//The default handler for the / main page
func (d *webData) mainPage(w http.ResponseWriter, r *http.Request) {
	//start a web page based on template
	err := tmpl["init.html"].ExecuteTemplate(w, "mainCompletePage", "put the data here")
	if err != nil {
		log.Println("mainPage: template execution error = ", err)
	}
}

//The web handler for adding persons
func (d *webData) addUsersWeb(w http.ResponseWriter, r *http.Request) {
	err := tmpl["init.html"].ExecuteTemplate(w, "addUserCompletePage", "some data")
	if err != nil {
		log.Println("addUsersWeb: template execution error = ", err)
	}

	r.ParseForm()
	u := User{}
	u.FirstName = r.FormValue("firstName")
	u.LastName = r.FormValue("lastName")
	u.Mail = r.FormValue("mail")
	u.Address = r.FormValue("address")
	u.PostNrAndPlace = r.FormValue("poAddr")
	u.PhoneNr = r.FormValue("phone")
	u.OrgNr = r.FormValue("orgNr")
	u.CountryID = "0"

	if u.FirstName != "" {
		pid, _ := queryDBForLastCustomerUID(d.PDB)
		//increment the user index nr by one for the new used to add
		pid++
		println("addUsersWeb: UID = ", pid)
		addUserToDB(d.PDB, u)
	}
}

//The web handler for modifying a person
func (d *webData) modifyUsersWeb(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr
	//query the userDB for all users and put the returning slice with result in p
	p := queryDBForAllUserInfo(d.PDB)

	//Execute the web for modify users, range over p to make the select user drop down menu
	err := tmpl["init.html"].ExecuteTemplate(w, "modifyUserCompletePage", p)
	if err != nil {
		fmt.Fprint(w, "template execution error = ", err)
	}

	//Parse all the variables in the html form to get all the data
	r.ParseForm()
	//Get the value (number) of the chosen user from form dropdown menu <select name="users">
	num, _ := strconv.Atoi(r.FormValue("users"))

	//Write out all the info of the selected user to the web
	for i := range p {
		log.Println(ip, "modifyUsersWeb: p[i].Number = ", p[i].Number)
		//Iterate over the complete struct of users until the chosen user is found
		if p[i].Number == num {
			log.Println(ip, "modifyUsersWeb: p[i].FirstName, p[i].LastName , found user = ", p[i].FirstName, p[i].LastName)
			//Store the index nr in slice of the chosen user
			indexNR = i
			err := tmpl["init.html"].ExecuteTemplate(w, "modifyUserSingle", p[i]) //bruk bare en spesifik slice av struct og send til html template
			if err != nil {
				log.Println(ip, "modifyUsersWeb: error = ", err)
			}
		}
	}

	//create a variable based on user to hold the values parsed from the modify web
	u := User{}
	r.ParseForm()
	u.FirstName = r.FormValue("firstName")
	u.LastName = r.FormValue("lastName")
	u.Mail = r.FormValue("mail")
	u.Address = r.FormValue("address")
	u.PostNrAndPlace = r.FormValue("poAddr")
	u.PhoneNr = r.FormValue("phone")
	u.OrgNr = r.FormValue("orgNr")
	u.CountryID = r.FormValue("countryId")
	checkBox := r.Form["sure"]
	changed := false

	if checkBox != nil {
		if checkBox[0] == "ok" {
			fmt.Printf("modifyUsersWeb: Verdien av checkbox er = %v ,og typen er = %T\n\n", checkBox[0], checkBox[0])
			//Check what values that are changed
			if u.FirstName != p[indexNR].FirstName && u.FirstName != "" {
				p[indexNR].FirstName = u.FirstName
				changed = true
			}
			if u.LastName != p[indexNR].LastName && u.LastName != "" {
				p[indexNR].LastName = u.LastName
				changed = true
			}
			if u.Mail != p[indexNR].Mail && u.Mail != "" {
				p[indexNR].Mail = u.Mail
				changed = true
			}
			if u.Address != p[indexNR].Address && u.Address != "" {
				p[indexNR].Address = u.Address
				changed = true
			}
			if u.PostNrAndPlace != p[indexNR].PostNrAndPlace && u.PostNrAndPlace != "" {
				p[indexNR].PostNrAndPlace = u.PostNrAndPlace
				changed = true
			}
			if u.PhoneNr != p[indexNR].PhoneNr && u.PhoneNr != "" {
				p[indexNR].PhoneNr = u.PhoneNr
				changed = true
			}
			if u.OrgNr != p[indexNR].OrgNr && u.OrgNr != "" {
				p[indexNR].OrgNr = u.OrgNr
				changed = true
			}
			if u.CountryID != p[indexNR].CountryID && u.CountryID != "" {
				p[indexNR].CountryID = u.CountryID
				changed = true
			}
		}
	} else {
		log.Println(ip, "modifyUsersWeb: The value of checkbox was not set")
	}

	//if any of the values was changed....update information into database
	if changed {
		updateUserInDB(d.PDB, p[indexNR])
	}

}

//The web handler to show and print out all registered users in the database
func (d *webData) showUsersWeb(w http.ResponseWriter, r *http.Request) {
	p := queryDBForAllUserInfo(d.PDB)
	err := tmpl["init.html"].ExecuteTemplate(w, "showUserCompletePage", p)
	if err != nil {
		log.Println("showUsersWeb: template execution error = ", err)
	}
	fmt.Fprint(w, err)
}

//The web handler to delete a person
func (d *webData) deleteUserWeb(w http.ResponseWriter, r *http.Request) {
	p := queryDBForAllUserInfo(d.PDB)
	err := tmpl["init.html"].ExecuteTemplate(w, "deleteUserCompletePage", p)
	if err != nil {
		log.Println("showUsersWeb: template execution error = ", err)
	}

	//parse the html form and get all the data
	r.ParseForm()
	fn, _ := strconv.Atoi(r.FormValue("users"))
	deleteUserInDB(d.PDB, fn)
}

//************************* CREATE BILLS *******************************

//The web handler to the user selection in create bills
func (d *webData) webBillSelectUser(w http.ResponseWriter, r *http.Request) {
	d.Users = queryDBForAllUserInfo(d.PDB)
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
			indexNR = i
			err := tmpl["init.html"].ExecuteTemplate(w, "billShowUserSingle", d.Users[i])
			if err != nil {
				log.Println(ip, "modifyUsersWeb: error = ", err)
			}
		}
	}

	//Get the last used bill_id from DB
	highestBillNR, totalLineCount := queryDBForLastBillID(d.PDB)
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
		highestBillNR, totalLineCount := queryDBForLastBillID(d.PDB)
		log.Println("billCreateWeb: highestBillNR = ", highestBillNR, ", totaltLineCount = ", totalLineCount)

		newBill := Bill{}
		newBill.BillID = highestBillNR + 1
		newBill.UserID = d.ActiveUserID
		t := time.Now()
		newBill.CreatedDate = fmt.Sprint(t.Format("2006-01-02 15:04:05"))
		//create a new bill and return the new billID to use later
		d.CurrentBillID = addBillToDB(d.PDB, newBill)
		log.Println("billCreateWeb: newBillID = ", d.CurrentBillID)

		billLine := BillLines{}
		billLine.BillID = d.CurrentBillID
		billLine.LineID = 1
		billLine.Description = "noe tekst"

		addBillLineToDB(d.PDB, billLine)
	}
}

func (d *webData) webBillLines(w http.ResponseWriter, r *http.Request) {
	fmt.Println("INFO: webBillLines: Active user ID when call for bills = ", d.ActiveUserID)
	BillsForUser := queryDBForBillsForUser(d.PDB, d.ActiveUserID)
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

	fmt.Println("--------skal execute template")
	err := tmpl["init.html"].ExecuteTemplate(w, "billLinesComplete", BillsForUser)
	if err != nil {
		log.Println("webBillLines: template execution error = ", err)
	}

	fmt.Println("--------Har executed template og skal parse form")
	r.ParseForm()
	fmt.Println("r.Form = ", r.Form)

	if r.FormValue("userActionButton") == "choose bill" {
		billID, _ := strconv.Atoi(r.FormValue("billID"))
		log.Println("INFO: webBillLines: billID =", billID)
		fmt.Println("billID = ", billID)
		d.CurrentBillID = billID

	}

	fmt.Println("--------HER SKAL DEN HENTE DATA FOR BILL LINES OG SKRIVE UT BILL LINES")
	fmt.Println("data.CurrentBillID inneholder = ", d.CurrentBillID)
	billLines := queryDBForBillLinesInfo(d.PDB, d.CurrentBillID)
	fmt.Println("webBillLines: queryDBForBillLinesInfo: billLines = ", billLines)

	err = tmpl["init.html"].ExecuteTemplate(w, "createBillLines", billLines)
	if err != nil {
		log.Println("createBillUserSelection: createBillLines: template execution error = ", err)
	}

	r.ParseForm()
	fmt.Println("data.CurrentBillID inneholder =", d.CurrentBillID)

	if r.FormValue("billLineActionButton") == "add line" {
		billLine := BillLines{}
		billLine.BillID = d.CurrentBillID
		fmt.Println("#######billid some benyttes er =", d.CurrentBillID)
		//create a random number for the bill line....for now....
		rand.Seed(time.Now().UnixNano())
		billLine.LineID = rand.Intn(10000)
		billLine.Description = "noe tekst"
		addBillLineToDB(d.PDB, billLine)
		//doing a redirect so it redraws the page with the new line. Not sure if this is the best way....
		err = tmpl["init.html"].ExecuteTemplate(w, "redirectToEditBill", "some data")
		if err != nil {
			log.Println("createBillUserSelection: createBillLines: template execution error = ", err)
		}

		fmt.Println("------Finnished with the add line 'if' sentence")
	}
	fmt.Println("-------Finnished with the bill lines function")
}

//TODO: Move select bill dropdown to select user window
//Sjekk ut : Lage en template som bare tegner en bill line, og så kan man heller kjøre den templaten flere ganger
//	med forskjellige data fra slice som input. Da blir det kanskje mulig å oppdatere siden fortløpende ???

//ERROR: Den viser bill linjene hvis user bare har 1 bill. Har den 2 så viser den ikke lenger.
//			Prøv å legg alle variablene over i global og se om det kan være noe den mangler når den
//			prøver å tegne opp alle bill linjene.......evt. sjekk logger først.