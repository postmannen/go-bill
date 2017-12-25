/*
2: Tested person struct, and string to int conversion for person ID.
    Iterate the slice of struct
3: Added invoice nr. to personStruct
4: Added some menus with print and add person as options
	Added getPersonNextNr function to look up the next available person nr
5: Added auto next number for person
	Added more variables describing person
	Added /sp for show person info
	Added /ap for add person
	input only get added when "add" button is pushed
6:	Add sqlite with add data functions
	Removed some of the old not needed code
7:	Added query functions towards sqlite db
	Removed the not needed invoice variable from struct, and DB.
	Added dropdown list for /md (modify page)
8:	Added templates.html.
	Added templates to the handlers
	Nested the templates within the templates file,
	 so only 1 call is needed for each web page, and not seperate calls for header, menu..etc...
9:	Add modify person http, functions and database updates
	 TODO: Use a temp struct instead of all the single variables in the modify http handler function
			 Look into replacing the if to check update of fields with a switch/case
10: Cleanup
		Rename fakt.go to main.go
		Update comments
	Made the variable indexNR global to store the selected user in the modify form and function
11: Cleanup
		renamed where the name person(s) where used to user(s)
		removed unused code
12: Tested with html and CSS, but dont really understand how to align boxes, text etc.
13:	Added orgnr to user table
		ERROR : Does not update org nr in modify section
			The problem is only with the modify function, add works ok
	Changed the top menu to use links insted of input box's, and dropdown with css
	showAllUsers : replaced the input boxes with a table.
14: Changing the css styling of the pages
	Added delete user function
	Fixed an Error for counting highest nr, and number of rows which were
	 not visible until the delete function were implemented
	Rewrote the function to get next index number to get last index number,
	 and return highest user uid, and count of total uid's
15: Rewrote the DB table names to use all small letters, and underscore to seperate words
	 Changed all the code to reflect changes
16: Wrote the first db template to use in "template-database-creation.sql"
	Rewrote the addUser* functions to use type User (struct) instead of single variables of type int and string
	Rewrote the modifyUser* functions to use type User (struct) instead of single variables of type int and string

TODO:
	Keep the number 0 in the deleted user row, incase the last user is deleted
	 then a new used added will get that number
	Change the user pages to 1 page with add modify and delete

Ideas:
	Make the primary keys uid and bill ID random numbers, so you can sync the database
	 between different devices without getting a conflict.
	 Sorting can be done on a dummy index value that don't have to be unique



*/
package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

//User is used for all customers and users
type User struct { //some
	Number         int
	FirstName      string
	LastName       string
	Mail           string
	Address        string
	PostNrAndPlace string
	PhoneNr        string
	OrgNr          string
	CountryID      string
}

var pDB *sql.DB                        //The pointer to use with the Database
var tmpl map[string]*template.Template //map to hold all templates
var indexNR int                        //to store the index nr. in slice where chosen person is stored

func init() {
	//initate the templates
	tmpl = make(map[string]*template.Template)
	tmpl["init.html"] = template.Must(template.ParseFiles("templates.html"))
}

func main() {
	//create DB and store pointer in pDB
	pDB = createDB()
	defer pDB.Close()

	http.HandleFunc("/sp", showUsersWeb)
	http.HandleFunc("/ap", addUsersWeb)
	http.HandleFunc("/mp", modifyUsersWeb)
	http.HandleFunc("/", mainPage)
	http.HandleFunc("/du", deleteUserWeb)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":7000", nil)

}

/****************************
*	DATABASE FUNCTIONS		*
*							*
****************************/

//Query the database for all users, and return a slice of struct with all users
func queryDBForAllUserInfo(pDB *sql.DB) []User {
	//get total rows in database
	num, countLines := queryDBForLastCustomerUID(pDB)
	p := []User{}
	fmt.Println("queryDBForAllUserInfo : queryDBForAllUserInfo highestNR ER = ", num)
	fmt.Println("queryDBForAllUserInfo : queryDBForAllUserInfo countlines = ", countLines)

	for i := 1; i <= num; i++ {
		//append the row to slice
		pTemp := queryDBForSingleUserInfo(pDB, i)
		if pTemp.Number != 0 {
			p = append(p, queryDBForSingleUserInfo(pDB, i))
		}
	}
	return p
}

//Query the database for the info of a single user. Takes user ID of type int as input, returns struct of single user
func queryDBForSingleUserInfo(db *sql.DB, uid int) User {

	rows, err := db.Query("select * from user where user_id=?", uid)
	checkErr(err)

	var pid int
	//variables to store the rows.Scan below
	var firstname, lastname, mail, address, postnrandplace, phonenr, orgnr, countryID string
	//Next prepares the next result row for reading with the Scan method. It returns true on success,
	//or false if there is no next result row or an error happened while preparing it.
	//Err should be consulted to distinguish between the two cases.
	for rows.Next() {
		//Scan copies the columns in the current row into the values pointed at by dest.
		//The number of values in dest must be the same as the number of columns in Rows of database.
		rows.Scan(&pid, &firstname, &lastname, &mail, &address, &postnrandplace, &phonenr, &orgnr, &countryID)

	}
	m := User{}
	m.Number = pid
	m.FirstName = firstname
	m.LastName = lastname
	m.Mail = mail
	m.Address = address
	m.PostNrAndPlace = postnrandplace
	m.PhoneNr = phonenr
	m.OrgNr = orgnr
	m.CountryID = countryID

	defer rows.Close()
	return m
}

//input *sql.DB and returns the highest uid number, and line count of rows in DB
func queryDBForLastCustomerUID(db *sql.DB) (int, int) {
	rows, err := db.Query("select user_id from user")
	checkErr(err)
	defer rows.Close()
	//Prepare the slice to store numbers read from DB
	var num []int
	for rows.Next() {
		var pid int
		//The number of values below must be the same amount
		//as the number of rows in the DB
		err := rows.Scan(&pid) //puts data into the address of the variable
		checkErr(err)
		num = append(num, pid)
	}

	highestNr := 0
	countLines := 0
	//iterate the slice, and find the highest number, and number of lines.
	for i := range num {
		if highestNr < num[i] {
			highestNr = num[i]
			log.Println("queryDBForLastCustomerUID : highestNr = ", highestNr)
			countLines++
		}
	}
	log.Println("queryDBForLastCustomerUID: highestNr = ", highestNr)
	log.Println("queryDBForLastCustomerUID: countLines = ", countLines)
	return highestNr, countLines
}

//Update user in Database, takes pointer to db and type User struct as input
func updateUserInDB(db *sql.DB, u User) {
	tx, err := db.Begin()
	checkErr(err)

	log.Println("The org nr. sendt to updateUserDB function = ", u.OrgNr)
	stmt, err := tx.Prepare("UPDATE user SET user_id=?,first_name=?,last_name=?,mail=?,address=?,post_nr_place=?,phone_nr=?,org_nr=?,country_id=? WHERE user_id=?")
	checkErr(err)
	defer stmt.Close()
	log.Println("updateUserInDB : Number in updateUserInDB function = ", u.Number)
	log.Println("************", u.Number, u.FirstName, u.LastName, u.Mail, u.Address, u.PostNrAndPlace, u.PhoneNr, u.OrgNr, u.CountryID, u.Number, "*************")
	//number is passed an extra time at the end of DB statement to fill the variable for the Query, which is done by number of user
	_, err = stmt.Exec(u.Number, u.FirstName, u.LastName, u.Mail, u.Address, u.PostNrAndPlace, u.PhoneNr, u.OrgNr, u.CountryID, u.Number)

	tx.Commit()
	checkErr(err)

}

//Adds user to Database. takes pointer to DB, and type User struct as input
func addUserToDB(db *sql.DB, u User) {
	//start db session
	tx, err := db.Begin()
	checkErr(err)

	//create statement to insert values to DB
	stmt, err := tx.Prepare("insert into user(user_id,first_name,last_name,mail,address,post_nr_place,phone_nr,org_nr,country_id) values(?,?,?,?,?,?,?,?,?)")
	checkErr(err)
	//At the end of function close the DB
	defer stmt.Close()

	//execute the statement on the DB
	_, err = stmt.Exec(u.Number, u.FirstName, u.LastName, u.Mail, u.Address, u.PostNrAndPlace, u.PhoneNr, u.OrgNr, u.CountryID)
	//commit to DB
	tx.Commit()
	checkErr(err)

}

//creates the database
func createDB() *sql.DB {
	//1. Open connection

	db, err := sql.Open("sqlite3", "./fakt.db") //return types = *DB, error
	checkErr(err)
	//2. fail-fast if can't connect to DB
	checkErr(db.Ping())
	//3. create table
	_, err = db.Exec(`create table if not exists user (
						pid integer PRIMARY KEY, 
						firstname string not null,
						lastname string,
						mail string,
						address string,
						postnrandplace string,
						phonenr string,
						orgnr string)
						;`)
	checkErr(err)

	/*_, err = db.Exec(`CREATE TABLE IF NOT EXISTS bills (
						pid integer PRIMARY KEY,
						billNR string not null,
						service string not null,
						amount string not null,
						priceExVat integer not null,
						discount integer,
						priceIncVat integer,
						totalSumExVat integer,
						totalSumIncVat integer)
	`)
	checkErr(err)*/

	return db
}

//if error !=nil, print error message to web page
func checkErr(err error, args ...string) {
	if err != nil {
		log.Printf("ERROR : %q: %s\n", err, args)
	}
}

//Delete a row in DB, takes pointer to db, and index number uid which corresponds to column 1 in DB for input
func deleteUserInDB(db *sql.DB, number int) {
	tx, err := db.Begin()
	checkErr(err)
	log.Println("deleteUserInDB: The index number of the person to delete is = ", number)

	//Make the sql statement to execute
	stmt, err := tx.Prepare("DELETE FROM user WHERE user_id=?")
	checkErr(err)

	defer stmt.Close()
	//prepare the statement with a value for the "?"
	_, err = stmt.Exec(number)
	tx.Commit()
	checkErr(err)

}

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
	/*fn := r.FormValue("firstName")
	ln := r.FormValue("lastName")
	ma := r.FormValue("mail")
	ad := r.FormValue("address")
	pa := r.FormValue("poAddr")
	pn := r.FormValue("phone")
	on := r.FormValue("orgNr")
	coID := "0"*/

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
