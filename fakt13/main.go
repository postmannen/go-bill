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
	Changed the top menu to use links insted of input box's
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
	http.HandleFunc("/styles.css", serveCSS)
	http.ListenAndServe(":7000", nil)

}

//Query the database for all users, and return a slice of struct with all users
func queryDBForAllUserInfo(pDB *sql.DB) []User {
	//get total rows in database
	_, num := queryDBForNextCustomerUID(pDB)
	p := []User{}

	for i := 1; i <= num; i++ {
		//append the row to slice
		p = append(p, queryDBForSingleUserInfo(pDB, i))
	}
	return p
}

//Query the database for the info of a single user. Takes user ID of type int as input, returns struct
func queryDBForSingleUserInfo(db *sql.DB, uid int) User {

	rows, err := db.Query("select * from user where pid=?", uid)
	checkErr(err)

	var pid int
	//variables to store the rows.Scan below
	var firstname, lastname, mail, address, postnrandplace, phonenr, orgnr string
	//Next prepares the next result row for reading with the Scan method. It returns true on success,
	//or false if there is no next result row or an error happened while preparing it.
	//Err should be consulted to distinguish between the two cases.
	for rows.Next() {
		//Scan copies the columns in the current row into the values pointed at by dest.
		//The number of values in dest must be the same as the number of columns in Rows of database.
		rows.Scan(&pid, &firstname, &lastname, &mail, &address, &postnrandplace, &phonenr, &orgnr)

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

	defer rows.Close()
	return m
}

//input *sql.DB and returns pid of type int, and total number of rows as type int
func queryDBForNextCustomerUID(db *sql.DB) (int, int) {
	rows, err := db.Query("select pid from user")
	checkErr(err)
	defer rows.Close()
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
	for i := range num {
		if highestNr < num[i] {
			highestNr = num[i]
			countLines++
		}
	}
	highestNr++
	return highestNr, countLines
}

//Update user in Database
func updateUserInDB(db *sql.DB, number int, first string, last string, mail string, adr string, ponr string, phone string, orgnr string) {
	tx, err := db.Begin()
	checkErr(err)

	log.Println("The org nr. sendt to updateUserDB function = ", orgnr)
	stmt, err := tx.Prepare("UPDATE user SET pid=?,firstname=?,lastname=?,mail=?,address=?,postnrandplace=?,phonenr=?,orgnr=? WHERE pid=?")
	checkErr(err)
	defer stmt.Close()
	log.Println("Number in updateUserInDB function = ", number)
	ape, err := stmt.Exec(number, first, last, mail, adr, ponr, phone, orgnr, number)
	//log.Println("VALUES TO DB : number = ",number,"first = ",first,"last = ",last,"mail = ",mail,"adr = ",adr,"ponr = ",ponr,"phone = ",phone,)
	log.Println("TEKST FRA APE : ", ape)
	tx.Commit()
	checkErr(err)

}

//Adds user to Database
func addUserToDB(db *sql.DB, number int, first string, last string, mail string, adr string, ponr string, phone string, orgnr string) {
	tx, err := db.Begin()
	checkErr(err)

	stmt, err := tx.Prepare("insert into user(pid,firstname,lastname,mail,address,postnrandplace,phonenr,orgnr) values(?,?,?,?,?,?,?,?)")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(number, first, last, mail, adr, ponr, phone, orgnr)
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

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS bills (
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
	checkErr(err)

	return db
}

//if error !=nil, print error message to web page
func checkErr(err error, args ...string) {
	if err != nil {
		log.Printf("ERROR : %q: %s\n", err, args)
	}
}

//serveCSS : The CSS file needs its own handler. TODO: check out http.ServeFiles
func serveCSS(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/styles.css")

}

//The default handler for the / main page
func mainPage(w http.ResponseWriter, r *http.Request) {
	err := tmpl["init.html"].ExecuteTemplate(w, "mainCompletePage", "put the data here")
	if err != nil {
		log.Println("template execution error = ", err)
	}
}

//The web handler for adding persons
func addUsersWeb(w http.ResponseWriter, r *http.Request) {
	//read template file, and execute template defined within file, and send "some data" to the template
	err := tmpl["init.html"].ExecuteTemplate(w, "addUserCompletePage", "some data")
	if err != nil {
		log.Println("template execution error = ", err)
	}

	r.ParseForm()
	fn := r.FormValue("firstName")
	ln := r.FormValue("lastName")
	ma := r.FormValue("mail")
	ad := r.FormValue("address")
	pa := r.FormValue("poAddr")
	pn := r.FormValue("phone")
	on := r.FormValue("orgNr")

	if fn != "" {
		pid, _ := queryDBForNextCustomerUID(pDB)
		println("UID = ", pid)
		addUserToDB(pDB, pid, fn, ln, ma, ad, pa, pn, on)
	} else {
		//fmt.Fprintf(w, "Minimum needed is firstname")
	}
}

//The web handler for modifying a person
func modifyUsersWeb(w http.ResponseWriter, r *http.Request) {
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
	fn, _ := strconv.Atoi(r.FormValue("users"))

	//Write out all the info of the selected user to the web
	for i := range p {
		//Iterate over the complete struct of users until the chosen user is found
		if p[i].Number == fn {
			log.Println("Du valgte ", p[i].FirstName, p[i].LastName)
			//Store the index nr in slice of the chosen user
			indexNR = i
			err := tmpl["init.html"].ExecuteTemplate(w, "showUserSingle", p[i]) //bruk bare en spesifik slice av struct og send til html template
			log.Println(err)
		}
	}

	r.ParseForm()
	fn2 := r.FormValue("firstName")
	ln2 := r.FormValue("lastName")
	ma2 := r.FormValue("mail")
	ad2 := r.FormValue("address")
	pa2 := r.FormValue("poAddr")
	pn2 := r.FormValue("phone")
	on2 := r.FormValue("orgNr")
	checkBox := r.Form["sure"]
	changed := false

	if checkBox != nil {
		if checkBox[0] == "ok" {
			fmt.Printf("Verdien av checkbox er = %v ,og typen er = %T\n\n", checkBox[0], checkBox[0])
			//Check what values that are changed
			if fn2 != p[indexNR].FirstName && fn2 != "" {
				log.Println("fn2 and FirstName are not the same ", fn2, "***", p[indexNR].FirstName)
				p[indexNR].FirstName = fn2
				changed = true
			}
			if ln2 != p[indexNR].LastName && ln2 != "" {
				log.Println("ln2 and LastName are not the same ", ln2, "***", p[indexNR].LastName)
				p[indexNR].LastName = ln2
				changed = true
			}
			if ma2 != p[indexNR].Mail && ma2 != "" {
				log.Println("ma2 and Mail are not the same ", ma2, "***", p[indexNR].Mail)
				p[indexNR].Mail = ma2
				changed = true
			}
			if ad2 != p[indexNR].Address && ad2 != "" {
				log.Println("ad2 and Address are not the same ", ad2, "***", p[indexNR].Address)
				p[indexNR].Address = ad2
				changed = true
			}
			if pa2 != p[indexNR].PostNrAndPlace && pa2 != "" {
				log.Println("pa2 and PostNrAndPlace are not the same ", pa2, "***", p[indexNR].PostNrAndPlace)
				p[indexNR].PostNrAndPlace = pa2
				changed = true
			}
			if pn2 != p[indexNR].PhoneNr && pn2 != "" {
				log.Println("pn2 and PhoneNr are not the same ", pn2, "***", p[indexNR].PhoneNr)
				p[indexNR].PhoneNr = pn2
				changed = true
			}
			if on2 != p[indexNR].OrgNr && on2 != "" {
				log.Println("on2 and OrgNr are not the same ", on2, "***", p[indexNR].OrgNr)
				p[indexNR].OrgNr = on2
				changed = true
			}
		}
	} else {
		log.Println("Verdien av checkbox var ikke satt")
	}

	log.Println("Personen du nå prøver å oppdatere har info = ", p[indexNR])

	//if any of the values was changed....update information into database
	if changed {
		updateUserInDB(pDB, p[indexNR].Number, p[indexNR].FirstName, p[indexNR].LastName, p[indexNR].Mail, p[indexNR].Address, p[indexNR].PostNrAndPlace, p[indexNR].PhoneNr, p[indexNR].OrgNr)
	}

}

//The web handler to show and print out all registered users in the database
func showUsersWeb(w http.ResponseWriter, r *http.Request) {
	//query the database for all information and store them in the struct 'p'
	p := queryDBForAllUserInfo(pDB)
	err := tmpl["init.html"].ExecuteTemplate(w, "showUserCompletePage", p)
	if err != nil {
		log.Println("template execution error = ", err)
	}
	fmt.Fprint(w, err)
}
