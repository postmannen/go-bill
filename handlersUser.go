package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/postmannen/go-bill/data"
)

//The default handler for the / main page
func (d *webData) mainPage(w http.ResponseWriter, r *http.Request) {
	//start a web page based on template
	err := d.tpl.ExecuteTemplate(w, "mainCompletePage", d)
	if err != nil {
		log.Println("mainPage: template execution error = ", err)
	}
}

//The web handler for adding persons
func (d *webData) addUsersWeb(w http.ResponseWriter, r *http.Request) {
	err := d.tpl.ExecuteTemplate(w, "addUserCompletePage", "some data")
	if err != nil {
		log.Println("addUsersWeb: template execution error = ", err)
	}

	r.ParseForm()
	u := data.User{}
	getFormValuesUserInfo(&u, r)

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
	r.ParseForm()
	ip := r.RemoteAddr
	//query the userDB for all users and put the returning slice with result in p
	uAllUsers := data.QueryAllUserInfo(d.PDB)
	fmt.Println("---ALL USERS FROM DATABASE = ", uAllUsers)

	//Execute the web for modify users, range over allUsers to make the select user drop down menu
	err := d.tpl.ExecuteTemplate(w, "modifyUserCompletePage", uAllUsers)
	if err != nil {
		fmt.Fprint(w, "template execution error = ", err)
	}

	//Execute the modifyUserSelection drop down menu template
	err = d.tpl.ExecuteTemplate(w, "modifyUserSelection", uAllUsers)
	if err != nil {
		fmt.Fprint(w, "template execution error = ", err)
	}

	//Get the value (number) of the chosen user from form dropdown menu <select name="users">
	num, _ := strconv.Atoi(r.FormValue("users"))
	var singleUser data.User

	//Find the selected single user chosen in dropdown in the slice of all users
	for i := range uAllUsers {
		//Iterate over the complete struct of users until the chosen user is found
		if uAllUsers[i].Number == num {
			//Store the index nr in slice of the chosen user
			singleUser = uAllUsers[i]
			d.IndexUser = i
		}
	}
	err = d.tpl.ExecuteTemplate(w, "modifyUser", singleUser) //bruk bare en spesifik slice av struct og send til html template
	if err != nil {
		log.Println(ip, "modifyUsersWeb: error = ", err)
	}

	uForm := data.User{}
	//get all the values from the user info fileds of the the
	getFormValuesUserInfo(&uForm, r)

	//TODO: should work with singleUser as input below, but it doesn't. Have to investigate more about this later.
	changed := checkUserFormChanged(uForm, &uAllUsers[d.IndexUser])

	//if any of the values was changed....update information into database
	if changed {
		data.UpdateUser(d.PDB, uAllUsers[d.IndexUser])
	}
}

//The web handler for modifying the admin
func (d *webData) modifyAdminWeb(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr
	adminID := 0
	//query the userDB for all users and put the returning slice with result in p
	u := data.QuerySingleUserInfo(d.PDB, adminID)

	//Execute the web for modify users, range over p to make the select user drop down menu
	err := d.tpl.ExecuteTemplate(w, "modifyUserCompletePage", u)
	if err != nil {
		fmt.Fprint(w, "Error: modifyAdminWeb: template execution error = ", err)
	}

	//Write out all the info of the selected user to the web

	err = d.tpl.ExecuteTemplate(w, "modifyUser", u) //bruk bare en spesifik slice av struct og send til html template
	if err != nil {
		log.Println(ip, "modifyAdminWeb: error = ", err)
	}

	r.ParseForm()

	//create a variable based on user to hold the values parsed from the modify web
	uForm := data.User{}
	//get all the values like name etc. from the form, and put them in u
	getFormValuesUserInfo(&uForm, r)

	changed := checkUserFormChanged(uForm, &u)

	//Check what values that are changed

	//if any of the values was changed....update information into database
	if changed {
		data.UpdateUser(d.PDB, u)

		//Execute the redirect to modifyAdmin to refresh page
		err := d.tpl.ExecuteTemplate(w, "redirectTomodifyAdmin", u)
		if err != nil {
			fmt.Fprint(w, "Error: modifyAdminWeb: template execution error = ", err)
		}
	}
}

//takes user info taken from form, and compares it with the original values
func checkUserFormChanged(uForm data.User, originalUser *data.User) (changed bool) {
	changed = false
	if uForm.FirstName != originalUser.FirstName && uForm.FirstName != "" {
		originalUser.FirstName = uForm.FirstName
		changed = true
	}
	if uForm.LastName != originalUser.LastName && uForm.LastName != "" {
		originalUser.LastName = uForm.LastName
		changed = true
	}
	if uForm.Mail != originalUser.Mail && uForm.Mail != "" {
		originalUser.Mail = uForm.Mail
		changed = true
	}
	if uForm.Address != originalUser.Address && uForm.Address != "" {
		originalUser.Address = uForm.Address
		changed = true
	}
	if uForm.PostNrAndPlace != originalUser.PostNrAndPlace && uForm.PostNrAndPlace != "" {
		originalUser.PostNrAndPlace = uForm.PostNrAndPlace
		changed = true
	}
	if uForm.PhoneNr != originalUser.PhoneNr && uForm.PhoneNr != "" {
		originalUser.PhoneNr = uForm.PhoneNr
		changed = true
	}
	if uForm.OrgNr != originalUser.OrgNr && uForm.OrgNr != "" {
		originalUser.OrgNr = uForm.OrgNr
		changed = true
	}
	if uForm.CountryID != originalUser.CountryID && uForm.CountryID != "" {
		originalUser.CountryID = uForm.CountryID
		changed = true
	}
	if uForm.BankAccount != originalUser.BankAccount && uForm.BankAccount != "" {
		originalUser.BankAccount = uForm.BankAccount
		changed = true
	}
	return changed
}

//The web handler to show and print out all registered users in the database
func (d *webData) showUsersWeb(w http.ResponseWriter, r *http.Request) {
	p := data.QueryAllUserInfo(d.PDB)
	err := d.tpl.ExecuteTemplate(w, "showUserCompletePage", p)
	if err != nil {
		log.Println("showUsersWeb: template execution error = ", err)
	}
	fmt.Fprint(w, err)
}

//The web handler to delete a person
func (d *webData) deleteUserWeb(w http.ResponseWriter, r *http.Request) {
	p := data.QueryAllUserInfo(d.PDB)
	err := d.tpl.ExecuteTemplate(w, "deleteUserCompletePage", p)
	if err != nil {
		log.Println("showUsersWeb: template execution error = ", err)
	}

	//parse the html form and get all the data
	r.ParseForm()
	fn, _ := strconv.Atoi(r.FormValue("users"))
	data.DeleteUser(d.PDB, fn)
}

//getFormValuesUserInfo will get all the user data from form.
func getFormValuesUserInfo(u *data.User, r *http.Request) {
	u.FirstName = r.FormValue("firstName")
	u.LastName = r.FormValue("lastName")
	u.Mail = r.FormValue("mail")
	u.Address = r.FormValue("address")
	u.PostNrAndPlace = r.FormValue("poAddr")
	u.PhoneNr = r.FormValue("phone")
	u.OrgNr = r.FormValue("orgNr")
	u.CountryID = "0"
	u.BankAccount = r.FormValue("bankAccount")
}
