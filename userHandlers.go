package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"text/template"

	"github.com/gorilla/schema"

	"github.com/postmannen/go-bill/data"
)

//The handler for the / main page
func (d *webData) mainPage() http.HandlerFunc {
	var init sync.Once
	var tpl *template.Template //template
	init.Do(func() {
		tpl = template.Must(template.ParseFiles("public/userTemplates.html"))
	})

	return func(w http.ResponseWriter, r *http.Request) {
		//start a web page based on template
		err := tpl.ExecuteTemplate(w, "mainPage", "put the data here")
		if err != nil {
			log.Println("mainPage: template execution error = ", err)
		}
	}
}

//The handler for adding persons
func (d *webData) addUsers() http.HandlerFunc {
	var init sync.Once
	var tpl *template.Template //template
	init.Do(func() {
		tpl = template.Must(template.ParseFiles("public/userTemplates.html"))
	})

	return func(w http.ResponseWriter, r *http.Request) {
		//seems like the ParseForm have to before template execution with
		//POST methods, to be able to grab the values. The returned map
		//0 when it is after TemplateExecution.
		err := r.ParseForm()
		if err != nil {
			log.Printf("error: parseform : %v \n", err)
		}

		err = tpl.ExecuteTemplate(w, "addUserPage", nil)
		if err != nil {
			log.Println("addUsersWeb: template execution error = ", err)
		}

		//temp variable for holding the parsed user values from the r.Form
		var u data.User

		//use gorilla schema to parse the values of the form, and put them into
		//a temp variable 'u'
		err = formDecoder.Decode(&u, r.PostForm)
		if err != nil {
			log.Printf("error: formDecoder : %v \n", err)
		}
		fmt.Println(u)

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
}

//The web handler for modifying a person
//Using websocket
func (d *webData) modifyUsers() http.HandlerFunc {
	var init sync.Once
	var tpl *template.Template
	init.Do(func() {
		tpl = template.Must(template.ParseFiles("public/userTemplates.html", "public/websocket.html"))
	})

	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Printf("error: parseform : %v \n", err)
		}

		ip := r.RemoteAddr
		//query the userDB for all users and put the returning slice with result in p
		d.Users = data.QueryAllUserInfo(d.PDB)

		//Execute the web for modify users, range over p to make the select user drop down menu
		err = tpl.ExecuteTemplate(w, "modifyUserPage", d.Users)
		if err != nil {
			fmt.Fprint(w, "template execution error = ", err)
		}

		//Execute the modifyUserSelection drop down menu template
		err = tpl.ExecuteTemplate(w, "modifyUserSelection", d.Users)
		if err != nil {
			fmt.Fprint(w, "template execution error = ", err)
		}

		//Delete user from db when button is pushed
		fmt.Println("---The FORM = ", r.Form)
		fmt.Println("---single value from form = ", r.Form["users"])
		buttonPushed := r.FormValue("submitButton")
		//if the manage bills button were pushed
		if buttonPushed == "Delete user" {
			userID, _ := strconv.Atoi(r.FormValue("users"))
			fmt.Printf("---userID = %v, and type = %T\n", userID, userID)
			data.DeleteUser(d.PDB, userID)
		}
		//WORKING HERE !!!!!!!!!!
		if buttonPushed == "Choose selected" {
			userID, _ := strconv.Atoi(r.FormValue("users"))
			//first range all users and make sure none is marked as selected
			for i, v := range d.Users {
				fmt.Println("v = ", v)
				d.Users[i].Selected = ""
				//if the id from the dropdown is the same as the one we found, set selected
				if d.Users[i].Number == userID {
					fmt.Println("*************FOUND THE USER*******************")
					d.Users[i].Selected = "selected"
				}
				fmt.Println("--- ", d.Users[i])
			}
		}

		//Get the value (number) of the chosen user from form dropdown menu <select name="users">
		num, _ := strconv.Atoi(r.FormValue("users"))

		//Write out all the info of the selected user to the web
		for i := range d.Users {
			log.Println(ip, "modifyUsersWeb: d.Users[i].Number = ", d.Users[i].Number)
			//Iterate over the complete struct of users until the chosen user is found
			if d.Users[i].Number == num {
				log.Println(ip, "modifyUsersWeb: p[i].FirstName, p[i].LastName , found user = ", d.Users[i].FirstName, d.Users[i].LastName)
				//Store the index nr in slice of the chosen user
				d.IndexUser = i
				err := tpl.ExecuteTemplate(w, "modifyUser", d.Users[i])
				if err != nil {
					log.Println(ip, "modifyUsersWeb: error = ", err)
				}
			}

		}

		//create a variable based on user to hold the values parsed from the modify web
		//temp variable for holding the parsed user values from the r.Form
		var u data.User

		//use gorilla schema to parse the values of the form, and put them into
		//a temp variable 'u'
		err = formDecoder.Decode(&u, r.Form)
		if err != nil {
			log.Printf("error: formDecoder : %v \n", err)
		}

		changed := false

		//check if the values in the form where changed by comparing them to the original ones
		if u.FirstName != d.Users[d.IndexUser].FirstName && u.FirstName != "" {
			d.Users[d.IndexUser].FirstName = u.FirstName
			changed = true
		}
		if u.LastName != d.Users[d.IndexUser].LastName && u.LastName != "" {
			d.Users[d.IndexUser].LastName = u.LastName
			changed = true
		}
		if u.Mail != d.Users[d.IndexUser].Mail && u.Mail != "" {
			d.Users[d.IndexUser].Mail = u.Mail
			changed = true
		}
		if u.Address != d.Users[d.IndexUser].Address && u.Address != "" {
			d.Users[d.IndexUser].Address = u.Address
			changed = true
		}
		if u.PostNrAndPlace != d.Users[d.IndexUser].PostNrAndPlace && u.PostNrAndPlace != "" {
			d.Users[d.IndexUser].PostNrAndPlace = u.PostNrAndPlace
			changed = true
		}
		if u.PhoneNr != d.Users[d.IndexUser].PhoneNr && u.PhoneNr != "" {
			d.Users[d.IndexUser].PhoneNr = u.PhoneNr
			changed = true
		}
		if u.OrgNr != d.Users[d.IndexUser].OrgNr && u.OrgNr != "" {
			d.Users[d.IndexUser].OrgNr = u.OrgNr
			changed = true
		}
		if u.CountryID != d.Users[d.IndexUser].CountryID && u.CountryID != "" {
			d.Users[d.IndexUser].CountryID = u.CountryID
			changed = true
		}
		if u.BankAccount != d.Users[d.IndexUser].BankAccount && u.BankAccount != "" {
			d.Users[d.IndexUser].BankAccount = u.BankAccount
			changed = true
		}

		//if any of the values was changed....update information into database
		if changed {
			data.UpdateUser(d.PDB, d.Users[d.IndexUser])
		}
	}
}

//The web handler for modifying the admin user
func (d *webData) modifyAdmin() http.HandlerFunc {
	var init sync.Once
	var tpl *template.Template //template
	init.Do(func() {
		tpl = template.Must(template.ParseFiles("public/userTemplates.html"))
	})

	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Printf("error: parseform : %v \n", err)
		}

		ip := r.RemoteAddr
		adminID := 0
		//query the userDB for all users and put the returning slice with result in p
		p := data.QuerySingleUserInfo(d.PDB, adminID)

		//Execute the web for modify users, range over p to make the select user drop down menu
		err = tpl.ExecuteTemplate(w, "modifyUserPage", p)
		if err != nil {
			fmt.Fprint(w, "Error: modifyAdminWeb: template execution error = ", err)
		}

		//Write out all the info of the selected user to the web

		err = tpl.ExecuteTemplate(w, "modifyUser", p) //bruk bare en spesifik slice av struct og send til html template
		if err != nil {
			log.Println(ip, "modifyAdminWeb: error = ", err)
		}

		//create a variable based on user to hold the values parsed from the modify web
		//temp variable for holding the parsed user values from the r.Form
		var u data.User
		var formDecoder = schema.NewDecoder()

		//use gorilla schema to parse the values of the form, and put them into
		//a temp variable 'u'
		err = formDecoder.Decode(&u, r.Form)
		if err != nil {
			log.Printf("error: formDecoder : %v \n", err)
		}

		//check if any of the values of the form is changed compared to the original values
		changed := false
		if u.FirstName != p.FirstName && u.FirstName != "" {
			p.FirstName = u.FirstName
			changed = true
		}
		if u.LastName != p.LastName && u.LastName != "" {
			p.LastName = u.LastName
			changed = true
		}
		if u.Mail != p.Mail && u.Mail != "" {
			p.Mail = u.Mail
			changed = true
		}
		if u.Address != p.Address && u.Address != "" {
			p.Address = u.Address
			changed = true
		}
		if u.PostNrAndPlace != p.PostNrAndPlace && u.PostNrAndPlace != "" {
			p.PostNrAndPlace = u.PostNrAndPlace
			changed = true
		}
		if u.PhoneNr != p.PhoneNr && u.PhoneNr != "" {
			p.PhoneNr = u.PhoneNr
			changed = true
		}
		if u.OrgNr != p.OrgNr && u.OrgNr != "" {
			p.OrgNr = u.OrgNr
			changed = true
		}
		if u.CountryID != p.CountryID && u.CountryID != "" {
			p.CountryID = u.CountryID
			changed = true
		}
		if u.BankAccount != p.BankAccount && u.BankAccount != "" {
			p.BankAccount = u.BankAccount
			changed = true
		}

		//if any of the values was changed....update information into database
		if changed {
			data.UpdateUser(d.PDB, p)

			//Execute the redirect to modifyAdmin to refresh page
			err := tpl.ExecuteTemplate(w, "redirectTomodifyAdmin", p)
			if err != nil {
				fmt.Fprint(w, "Error: modifyAdminWeb: template execution error = ", err)
			}
		}
	}
}

//The web handler to show and print out all registered users in the database
//Using websocket template
func (d *webData) manageUsers() http.HandlerFunc {
	var init sync.Once
	var tpl *template.Template //template
	init.Do(func() {
		tpl = template.Must(template.ParseFiles("public/userTemplates.html", "public/websocket.html"))
	})

	return func(w http.ResponseWriter, r *http.Request) {
		//The idea here is to draw the manageUsersPage profile as the initial page,
		//and add templates based on which action chosen in the JS via the websocket.
		d.Users = data.QueryAllUserInfo(d.PDB)
		err := tpl.ExecuteTemplate(w, "manageUsersPage", d)
		if err != nil {
			log.Println("showUsersWeb: template execution error = ", err)
		}
		fmt.Fprint(w, err)
	}
}
