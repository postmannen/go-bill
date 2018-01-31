package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/postmannen/fakt/data"
)

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
	u := data.User{}
	u.FirstName = r.FormValue("firstName")
	u.LastName = r.FormValue("lastName")
	u.Mail = r.FormValue("mail")
	u.Address = r.FormValue("address")
	u.PostNrAndPlace = r.FormValue("poAddr")
	u.PhoneNr = r.FormValue("phone")
	u.OrgNr = r.FormValue("orgNr")
	u.CountryID = "0"

	if u.FirstName != "" {
		pid, _ := data.QueryForLastUID(d.PDB)
		//increment the user index nr by one for the new used to add
		pid++
		fmt.Println("------pid ---------- = ", pid)
		println("addUsersWeb: UID = ", pid)
		u.Number = pid
		data.AddUser(d.PDB, u)
	}
}

//The web handler for modifying a person
func (d *webData) modifyUsersWeb(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr
	//query the userDB for all users and put the returning slice with result in p
	p := data.QueryAllUserInfo(d.PDB)

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
			d.indexNR = i
			err := tmpl["init.html"].ExecuteTemplate(w, "modifyUserSingle", p[i]) //bruk bare en spesifik slice av struct og send til html template
			if err != nil {
				log.Println(ip, "modifyUsersWeb: error = ", err)
			}
		}
	}

	//create a variable based on user to hold the values parsed from the modify web
	u := data.User{}
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
			if u.FirstName != p[d.indexNR].FirstName && u.FirstName != "" {
				p[d.indexNR].FirstName = u.FirstName
				changed = true
			}
			if u.LastName != p[d.indexNR].LastName && u.LastName != "" {
				p[d.indexNR].LastName = u.LastName
				changed = true
			}
			if u.Mail != p[d.indexNR].Mail && u.Mail != "" {
				p[d.indexNR].Mail = u.Mail
				changed = true
			}
			if u.Address != p[d.indexNR].Address && u.Address != "" {
				p[d.indexNR].Address = u.Address
				changed = true
			}
			if u.PostNrAndPlace != p[d.indexNR].PostNrAndPlace && u.PostNrAndPlace != "" {
				p[d.indexNR].PostNrAndPlace = u.PostNrAndPlace
				changed = true
			}
			if u.PhoneNr != p[d.indexNR].PhoneNr && u.PhoneNr != "" {
				p[d.indexNR].PhoneNr = u.PhoneNr
				changed = true
			}
			if u.OrgNr != p[d.indexNR].OrgNr && u.OrgNr != "" {
				p[d.indexNR].OrgNr = u.OrgNr
				changed = true
			}
			if u.CountryID != p[d.indexNR].CountryID && u.CountryID != "" {
				p[d.indexNR].CountryID = u.CountryID
				changed = true
			}
		}
	} else {
		log.Println(ip, "modifyUsersWeb: The value of checkbox was not set")
	}

	//if any of the values was changed....update information into database
	if changed {
		data.UpdateUser(d.PDB, p[d.indexNR])
	}
}

//The web handler to show and print out all registered users in the database
func (d *webData) showUsersWeb(w http.ResponseWriter, r *http.Request) {
	p := data.QueryAllUserInfo(d.PDB)
	err := tmpl["init.html"].ExecuteTemplate(w, "showUserCompletePage", p)
	if err != nil {
		log.Println("showUsersWeb: template execution error = ", err)
	}
	fmt.Fprint(w, err)
}

//The web handler to delete a person
func (d *webData) deleteUserWeb(w http.ResponseWriter, r *http.Request) {
	p := data.QueryAllUserInfo(d.PDB)
	err := tmpl["init.html"].ExecuteTemplate(w, "deleteUserCompletePage", p)
	if err != nil {
		log.Println("showUsersWeb: template execution error = ", err)
	}

	//parse the html form and get all the data
	r.ParseForm()
	fn, _ := strconv.Atoi(r.FormValue("users"))
	data.DeleteUser(d.PDB, fn)
}

//TODO: Move select bill dropdown to select user window
//Sjekk ut : Lage en template som bare tegner en bill line, og så kan man heller kjøre den templaten flere ganger
//	med forskjellige data fra slice som input. Da blir det kanskje mulig å oppdatere siden fortløpende ???

//ERROR: Den viser bill linjene hvis user bare har 1 bill. Har den 2 så viser den ikke lenger.
//			Prøv å legg alle variablene over i global og se om det kan være noe den mangler når den
//			prøver å tegne opp alle bill linjene.......evt. sjekk logger først.
