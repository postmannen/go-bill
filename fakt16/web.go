package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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

//The web handler to create bills
func billCreateWeb(w http.ResponseWriter, r *http.Request) {
	p := queryDBForAllUserInfo(pDB)

	err := tmpl["init.html"].ExecuteTemplate(w, "createBillCompletePage", p)
	if err != nil {
		log.Println("createBillUserSelection: template execution error = ", err)
	}

	ip := r.RemoteAddr

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
			err := tmpl["init.html"].ExecuteTemplate(w, "billShowUserSingle", p[i]) //bruk bare en spesifik slice av struct og send til html template
			log.Println(ip, "modifyUsersWeb: error = ", err)
		}
	}

	//a struct for the bill lines. Fields must be export (starting Capital letter) to be passed to template
	type billLines struct {
		ItemID             int
		LineNR             int
		Description        string
		BillID             int
		Quantity           int
		DiscountPercentage int
		VatUsed            int
		PriceExVat         float64
		//just create some linenumbers for testing
	}

	/*
		//create a slice of type billLines to hold all the billLines
		line := []billLines{}
		line = append(line, billLines{LineNR: 1, Description: "Noe en", Quantity: 1, DiscountPercentage: 0, VatUsed: 25, PriceExVat: 999})
		line = append(line, billLines{LineNR: 2, Description: "Noe to", Quantity: 1, DiscountPercentage: 0, VatUsed: 25, PriceExVat: 999})
		line = append(line, billLines{LineNR: 3, Description: "Noe tre", Quantity: 1, DiscountPercentage: 0, VatUsed: 25, PriceExVat: 999})
		line = append(line, billLines{LineNR: 4, Description: "Noe fire", Quantity: 1, DiscountPercentage: 0, VatUsed: 25, PriceExVat: 999})
		fmt.Println(line)

		err = tmpl["init.html"].ExecuteTemplate(w, "createBillLines", line)
		if err != nil {
			log.Println("createBillCompletePage: template execution error = ", err)
		}
	*/

}