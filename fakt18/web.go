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
func mainPage(w http.ResponseWriter, r *http.Request) {
	//start a web page based on template
	err := tmpl["init.html"].ExecuteTemplate(w, "mainCompletePage", "put the data here")
	if err != nil {
		log.Println("mainPage: template execution error = ", err)
	}
}

//The web handler for adding persons
func addUsersWeb(w http.ResponseWriter, r *http.Request) {
	//read template file, and execute template defined within file, and send "some data" to the template
	err := tmpl["init.html"].ExecuteTemplate(w, "addUserCompletePage", "some data")
	if err != nil {
		log.Println("addUsersWeb: template execution error = ", err)
	}

	//r.ParseForm() lets you grab all the inputs and states from the webpage. Use FormValue to grab the specific values
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
		pid, _ := queryDBForLastCustomerUID(pDB)
		//increment the user index nr by one for the new used to add
		pid++
		println("addUsersWeb: UID = ", pid)
		addUserToDB(pDB, u)
	} else {
		//fmt.Fprintf(w, "Minimum needed is firstname")
	}
}

//The web handler for modifying a person
func modifyUsersWeb(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr
	//query the userDB for all users and put the returning slice with result in p
	p := queryDBForAllUserInfo(pDB)

	//Execute the web for modify users, give slice 'p' as input to the web page
	//the web will range over p to make the select user drop down menu
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
			log.Println(ip, "modifyUsersWeb: p[i].FirstName, p[i].LastName = ", p[i].FirstName, p[i].LastName)
			//Store the index nr in slice of the chosen user
			indexNR = i
			err := tmpl["init.html"].ExecuteTemplate(w, "modifyUserSingle", p[i]) //bruk bare en spesifik slice av struct og send til html template
			log.Println(ip, "modifyUsersWeb: error = ", err)
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
				log.Println(ip, "modifyUsersWeb: u.FirstName and FirstName are not the same ", u.FirstName, "***", p[indexNR].FirstName)
				p[indexNR].FirstName = u.FirstName
				changed = true
			}
			if u.LastName != p[indexNR].LastName && u.LastName != "" {
				log.Println(ip, "modifyUsersWeb: u.LastName and LastName are not the same ", u.LastName, "***", p[indexNR].LastName)
				p[indexNR].LastName = u.LastName
				changed = true
			}
			if u.Mail != p[indexNR].Mail && u.Mail != "" {
				log.Println(ip, "modifyUsersWeb: u.Mail and Mail are not the same ", u.Mail, "***", p[indexNR].Mail)
				p[indexNR].Mail = u.Mail
				changed = true
			}
			if u.Address != p[indexNR].Address && u.Address != "" {
				log.Println(ip, "modifyUsersWeb: u.Address and Address are not the same ", u.Address, "***", p[indexNR].Address)
				p[indexNR].Address = u.Address
				changed = true
			}
			if u.PostNrAndPlace != p[indexNR].PostNrAndPlace && u.PostNrAndPlace != "" {
				log.Println(ip, "modifyUsersWeb: u.PostNrAndPlace and PostNrAndPlace are not the same ", u.PostNrAndPlace, "***", p[indexNR].PostNrAndPlace)
				p[indexNR].PostNrAndPlace = u.PostNrAndPlace
				changed = true
			}
			if u.PhoneNr != p[indexNR].PhoneNr && u.PhoneNr != "" {
				log.Println(ip, "modifyUsersWeb: u.PhoneNr and PhoneNr are not the same ", u.PhoneNr, "***", p[indexNR].PhoneNr)
				p[indexNR].PhoneNr = u.PhoneNr
				changed = true
			}
			if u.OrgNr != p[indexNR].OrgNr && u.OrgNr != "" {
				log.Println(ip, "modifyUsersWeb: u.OrgNr and OrgNr are not the same ", u.OrgNr, "***", p[indexNR].OrgNr)
				p[indexNR].OrgNr = u.OrgNr
				changed = true
			}
			if u.CountryID != p[indexNR].CountryID && u.CountryID != "" {
				log.Println(ip, "modifyUsersWeb: coIDu.CountryID and CountryID are not the same ", u.CountryID, "***", p[indexNR].CountryID)
				p[indexNR].CountryID = u.CountryID
				changed = true
			}
		}
	} else {
		log.Println(ip, "modifyUsersWeb: The value of checkbox was not set")
	}

	log.Println(ip, "modifyUsersWeb: The person beeing modified have this original info = ", p[indexNR])

	//if any of the values was changed....update information into database
	if changed {
		//updateUserInDB(pDB, p[indexNR].Number, p[indexNR].FirstName, p[indexNR].LastName, p[indexNR].Mail, p[indexNR].Address, p[indexNR].PostNrAndPlace, p[indexNR].PhoneNr, p[indexNR].OrgNr, p[indexNR].CountryID)
		updateUserInDB(pDB, p[indexNR])
	}

}

//The web handler to show and print out all registered users in the database
func showUsersWeb(w http.ResponseWriter, r *http.Request) {
	//query the database for all information and store them in the struct 'p'
	p := queryDBForAllUserInfo(pDB)
	err := tmpl["init.html"].ExecuteTemplate(w, "showUserCompletePage", p)
	if err != nil {
		log.Println("showUsersWeb: template execution error = ", err)
	}
	fmt.Fprint(w, err)
}

//The web handler to delete a person
func deleteUserWeb(w http.ResponseWriter, r *http.Request) {
	p := queryDBForAllUserInfo(pDB)
	err := tmpl["init.html"].ExecuteTemplate(w, "deleteUserCompletePage", p)
	if err != nil {
		log.Println("showUsersWeb: template execution error = ", err)
	}

	//parse the html form and get all the data
	r.ParseForm()
	fn, _ := strconv.Atoi(r.FormValue("users"))
	//checkBox := r.Form["sure"]
	//Call the function to delete the selected user
	/*if checkBox != nil {
		if checkBox[0] == "ok" {
			deleteUserInDB(pDB, fn)
		}
	}*/
	deleteUserInDB(pDB, fn)
}

//************************* CREATE BILLS *******************************

//The web handler to create bills
func billCreateWebSelectUser(w http.ResponseWriter, r *http.Request) {
	data := webData{}

	data.Users = queryDBForAllUserInfo(pDB)

	//creates the header and the select box from templates
	err := tmpl["init.html"].ExecuteTemplate(w, "createBillCompletePage", data) //%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
	if err != nil {
		log.Println("createBillUserSelection: template execution error = ", err)
	}

	ip := r.RemoteAddr

	//Parse all the variables in the html form to get all the data
	r.ParseForm()

	//Get the value (number) of the chosen user from form dropdown menu <select name="users">
	num, _ := strconv.Atoi(r.FormValue("users"))

	//if sentence to keep the chosen user ID. Reason is that it resets to 0 when the page is redrawn after "choose" is pushed
	//put the value in chooseUserButton which is a global variable
	if r.FormValue("chooseUserButton") == "choose" {
		activeUserID = num
		//reset currentBillID so a new user dont inherit the last bill used for another user.
		currentBillID = 0
	}
	log.Println("billCreateWeb: The number active now in the user select box: num = ", num)
	log.Println("billCreateWeb: The number chosen in the user select box:activeUserID = ", activeUserID)

	//Write out all the info of the selected user to the web
	for i := range data.Users {
		log.Println(ip, "modifyUsersWeb: p[i].Number = ", data.Users[i].Number)
		//Iterate over the complete struct of users until the chosen user is found
		if data.Users[i].Number == activeUserID {
			log.Println(ip, "modifyUsersWeb: p[i].FirstName, p[i].LastName = ", data.Users[i].FirstName, data.Users[i].LastName)
			//Store the index nr in slice of the chosen user
			indexNR = i
			err := tmpl["init.html"].ExecuteTemplate(w, "billShowUserSingle", data.Users[i]) //%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
			log.Println(ip, "modifyUsersWeb: error = ", err)
		}
	}

	//Get the last used bill_id from DB
	totalBillNumbers, totalLineCount := queryDBForLastBillID(pDB)
	log.Println(ip, "billCreateWeb: totalBillNumbers = ", totalBillNumbers)
	log.Println(ip, "billCreateWeb: totalLineCount = ", totalLineCount)

	//Check which of the two input buttons where pushed. They both have name=userActionButton,
	//and the value can be read with r.FormValue("userActionButton")
	r.ParseForm()
	log.Println("r.Form shows = ", r.Form)
	log.Println("r.FormValue = ", r.FormValue("userActionButton"))
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
		totalBillNumbers, totalLineCount := queryDBForLastBillID(pDB)
		log.Println("billCreateWeb: totalBillNumbers = ", totalBillNumbers)
		log.Println("billCreateWeb: totalLineCount = ", totalLineCount)

		//create a new bill_id in bills database
		//use the next available bill number
		//use the chosen user id for user_id
		newBill := Bill{}
		newBill.BillID = totalBillNumbers + 1
		newBill.UserID = activeUserID

		//create a new bill and return the new billID to use later
		currentBillID = addBillToDB(pDB, newBill)
		log.Println("billCreateWeb: newBillID = ", currentBillID)

		billLine := BillLines{}
		billLine.BillID = currentBillID
		billLine.LineID = 1
		billLine.Description = "noe tekst"

		addBillLineToDB(pDB, billLine)
	}
}

func billCreateWebBillEdit(w http.ResponseWriter, r *http.Request) {
	fmt.Println("INFO: Active user ID when call for bills = ", activeUserID)
	BillsForUser := []Bill{}
	BillsForUser = queryDBForBillsForUser(pDB, activeUserID)
	fmt.Println("INFO: billCreateWebBillEdit: BillsForUser = ", BillsForUser)

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

	err := tmpl["init.html"].ExecuteTemplate(w, "editBillCompletePage", BillsForUser)
	if err != nil {
		log.Println("billCreateWebBillEdit: template execution error = ", err)
	}

	r.ParseForm()
	fmt.Println("r.ParseForm = ", r.Form)

	buttonAction := r.FormValue("userActionButton")
	billID, _ := strconv.Atoi(r.FormValue("billID"))

	if buttonAction == "choose bill" {
		fmt.Println(buttonAction, "pressed")
		log.Println("INFO: billCreateWebBillEdit: billID =", billID)
		fmt.Println("billID = ", billID)
		currentBillID = billID

	}

	billLines := queryDBForBillLinesInfo(pDB, billID)
	fmt.Println("***********billLines =", billLines)
	fmt.Println("billCreateWebBillEdit: queryDBForBillLinesInfo: billLines = ", billLines)
	err = tmpl["init.html"].ExecuteTemplate(w, "createBillLines", billLines)
	if err != nil {
		log.Println("createBillUserSelection: createBillLines: template execution error = ", err)
	}

	r.ParseForm()
	fmt.Println("...........form inneholder = ", r.Form)
	fmt.Println("currentBillID inneholder =", currentBillID)

	buttonAction = r.FormValue("billLineActionButton")
	if buttonAction == "add line" {
		fmt.Println("Du trykket add Line")

		billLine := BillLines{}
		billLine.BillID = currentBillID
		fmt.Println("#######billid some benyttes er =", currentBillID)
		//create a random number for the bill line....for now....
		rand.Seed(time.Now().UnixNano())
		billLine.LineID = rand.Intn(10000)
		billLine.Description = "noe tekst"

		addBillLineToDB(pDB, billLine)
	}
}
