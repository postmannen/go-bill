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
		err := tpl.ExecuteTemplate(w, "mainCompletePage", "put the data here")
		if err != nil {
			log.Println("mainPage: template execution error = ", err)
		}
	}
}

//The handler for adding persons
func (d *webData) addUsersWeb() http.HandlerFunc {
	var init sync.Once
	var tpl *template.Template //template
	init.Do(func() {
		tpl = template.Must(template.ParseFiles("public/userTemplates.html"))
	})

	return func(w http.ResponseWriter, r *http.Request) {
		err := tpl.ExecuteTemplate(w, "addUserCompletePage", "some data")
		if err != nil {
			log.Println("addUsersWeb: template execution error = ", err)
		}

		//temp variable for holding the parsed user values from the r.Form
		var u data.User
		var formDecoder = schema.NewDecoder()

		err = r.ParseForm()
		if err != nil {
			log.Printf("error: parseform : %v \n", err)
		}

		//use gorilla schema to parse the values of the form, and put them into
		//a temp variable 'u'
		err = formDecoder.Decode(&u, r.Form)
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
func (d *webData) modifyUsersWeb() http.HandlerFunc {
	var init sync.Once
	var tpl *template.Template
	init.Do(func() {
		tpl = template.Must(template.ParseFiles("public/userTemplates.html"))
	})

	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		//query the userDB for all users and put the returning slice with result in p
		p := data.QueryAllUserInfo(d.PDB)

		//Execute the web for modify users, range over p to make the select user drop down menu
		err := tpl.ExecuteTemplate(w, "modifyUserCompletePage", p)
		if err != nil {
			fmt.Fprint(w, "template execution error = ", err)
		}

		//Execute the modifyUserSelection drop down menu template
		err = tpl.ExecuteTemplate(w, "modifyUserSelection", p)
		if err != nil {
			fmt.Fprint(w, "template execution error = ", err)
		}

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
				d.IndexUser = i
				err := tpl.ExecuteTemplate(w, "modifyUser", p[i])
				if err != nil {
					log.Println(ip, "modifyUsersWeb: error = ", err)
				}
			}

		}

		//create a variable based on user to hold the values parsed from the modify web
		//temp variable for holding the parsed user values from the r.Form
		var u data.User
		var formDecoder = schema.NewDecoder()

		err = r.ParseForm()
		if err != nil {
			log.Printf("error: parseform : %v \n", err)
		}

		//use gorilla schema to parse the values of the form, and put them into
		//a temp variable 'u'
		err = formDecoder.Decode(&u, r.Form)
		if err != nil {
			log.Printf("error: formDecoder : %v \n", err)
		}

		changed := false

		//check if the values in the form where changed by comparing them to the original ones
		if u.FirstName != p[d.IndexUser].FirstName && u.FirstName != "" {
			p[d.IndexUser].FirstName = u.FirstName
			changed = true
		}
		if u.LastName != p[d.IndexUser].LastName && u.LastName != "" {
			p[d.IndexUser].LastName = u.LastName
			changed = true
		}
		if u.Mail != p[d.IndexUser].Mail && u.Mail != "" {
			p[d.IndexUser].Mail = u.Mail
			changed = true
		}
		if u.Address != p[d.IndexUser].Address && u.Address != "" {
			p[d.IndexUser].Address = u.Address
			changed = true
		}
		if u.PostNrAndPlace != p[d.IndexUser].PostNrAndPlace && u.PostNrAndPlace != "" {
			p[d.IndexUser].PostNrAndPlace = u.PostNrAndPlace
			changed = true
		}
		if u.PhoneNr != p[d.IndexUser].PhoneNr && u.PhoneNr != "" {
			p[d.IndexUser].PhoneNr = u.PhoneNr
			changed = true
		}
		if u.OrgNr != p[d.IndexUser].OrgNr && u.OrgNr != "" {
			p[d.IndexUser].OrgNr = u.OrgNr
			changed = true
		}
		if u.CountryID != p[d.IndexUser].CountryID && u.CountryID != "" {
			p[d.IndexUser].CountryID = u.CountryID
			changed = true
		}
		if u.BankAccount != p[d.IndexUser].BankAccount && u.BankAccount != "" {
			p[d.IndexUser].BankAccount = u.BankAccount
			changed = true
		}

		//if any of the values was changed....update information into database
		if changed {
			data.UpdateUser(d.PDB, p[d.IndexUser])
		}
	}
}

//The web handler for modifying the admin user
func (d *webData) modifyAdminWeb() http.HandlerFunc {
	var init sync.Once
	var tpl *template.Template //template
	init.Do(func() {
		tpl = template.Must(template.ParseFiles("public/userTemplates.html"))
	})

	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		adminID := 0
		//query the userDB for all users and put the returning slice with result in p
		p := data.QuerySingleUserInfo(d.PDB, adminID)

		//Execute the web for modify users, range over p to make the select user drop down menu
		err := tpl.ExecuteTemplate(w, "modifyUserCompletePage", p)
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

		err = r.ParseForm()
		if err != nil {
			log.Printf("error: parseform : %v \n", err)
		}

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
func (d *webData) showUsersWeb() http.HandlerFunc {
	var init sync.Once
	var tpl *template.Template //template
	init.Do(func() {
		tpl = template.Must(template.ParseFiles("public/userTemplates.html"))
	})

	return func(w http.ResponseWriter, r *http.Request) {
		p := data.QueryAllUserInfo(d.PDB)
		err := tpl.ExecuteTemplate(w, "showUserCompletePage", p)
		if err != nil {
			log.Println("showUsersWeb: template execution error = ", err)
		}
		fmt.Fprint(w, err)
	}
}

//The web handler to delete a person
func (d *webData) deleteUserWeb() http.HandlerFunc {
	var init sync.Once
	var tpl *template.Template
	init.Do(func() {
		tpl = template.Must(template.ParseFiles("public/userTemplates.html"))
	})

	return func(w http.ResponseWriter, r *http.Request) {
		p := data.QueryAllUserInfo(d.PDB)
		err := tpl.ExecuteTemplate(w, "deleteUserCompletePage", p)
		if err != nil {
			log.Println("showUsersWeb: template execution error = ", err)
		}

		//parse the html form and get all the data
		r.ParseForm()
		fn, _ := strconv.Atoi(r.FormValue("users"))
		data.DeleteUser(d.PDB, fn)
	}
}
