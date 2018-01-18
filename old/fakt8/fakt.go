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
*/
package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

//PersonStruct is used for all customers and users
type PersonStruct struct { //some
	Number         int
	FirstName      string
	LastName       string
	Mail           string
	Address        string
	PostNrAndPlace string
	PhoneNr        string
}

var person []PersonStruct              //The slice of struct to hold users
var pDB *sql.DB                        //The pointer to use with the Database
var tmpl map[string]*template.Template //map to hold all templates

func init() {
	tmpl = make(map[string]*template.Template)
	tmpl["init.html"] = template.Must(template.ParseFiles("templates.html"))

}

func main() {
	//create DB and store pointer in pDB
	pDB = createDB()
	defer pDB.Close()

	//fmt.Println(queryDBForAllUserInfo(pDB))

	http.HandleFunc("/sp", showPersonsWeb)
	http.HandleFunc("/ap", addPersonsWeb)
	http.HandleFunc("/mp", modifyPersonWeb)
	http.HandleFunc("/", mainPage)
	http.ListenAndServe(":8000", nil)

}

func mainPage(w http.ResponseWriter, r *http.Request) {
	err := tmpl["init.html"].ExecuteTemplate(w, "header", "put the data here")
	if err != nil {
		fmt.Println("template execution error = ", err)
	}
	err = tmpl["init.html"].ExecuteTemplate(w, "top1", "put the data here")
	if err != nil {
		fmt.Println("template execution error = ", err)
	}

}

//Query the database for all users, and return a slice of struct with all users
func queryDBForAllUserInfo(pDB *sql.DB) []PersonStruct {
	_, num := queryDBForNextCustomerPID(pDB)
	p := []PersonStruct{}

	for i := 1; i <= num; i++ {
		p = append(p, queryDBForSingleUserInfo(pDB, i))
	}
	return p
}

//Query the database for the info of a single user. Takes user ID of type int as input, returns struct
func queryDBForSingleUserInfo(db *sql.DB, uid int) PersonStruct {
	rows, err := db.Query("select * from user where pid=?", uid)
	checkErr(err)

	var pid int
	var firstname, lastname, mail, address, postnrandplace, phonenr string
	for rows.Next() {
		rows.Scan(&pid, &firstname, &lastname, &mail, &address, &postnrandplace, &phonenr)
		//fmt.Println(pid, firstname, lastname, mail, address, postnrandplace, phonenr)
	}
	m := PersonStruct{}
	m.Number = pid
	m.FirstName = firstname
	m.LastName = lastname
	m.Mail = mail
	m.Address = address
	m.PostNrAndPlace = postnrandplace
	m.PhoneNr = phonenr

	defer rows.Close()
	return m
}

//input *sql.DB and returns pid of type int, and total number of rows as type int
func queryDBForNextCustomerPID(db *sql.DB) (int, int) {
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
		//n, _ := strconv.Atoi(pid)
		num = append(num, pid)
		//fmt.Println(pid)
	}
	//fmt.Println("The number slice contains : ", num)

	highestNr := 0
	countLines := 0
	for i := range num {
		//fmt.Printf("%v\t%T\n", person[i].number, person[i].number)
		if highestNr < num[i] {
			highestNr = num[i]
			countLines++
		}
	}
	highestNr++
	return highestNr, countLines
}

//Update user in Database
func updateUserInDB(db *sql.DB, number int, first string, last string, mail string, adr string, ponr string, phone string) {
	tx, err := db.Begin()
	checkErr(err)

	stmt, err := tx.Prepare("insert into user(pid,firstname,lastname,mail,address,postnrandplace,phonenr) values(?,?,?,?,?,?,?)")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(number, first, last, mail, adr, ponr, phone)
	tx.Commit()
	checkErr(err)

}

//Adds user to Database
func addUserToDB(db *sql.DB, number int, first string, last string, mail string, adr string, ponr string, phone string) {
	tx, err := db.Begin()
	checkErr(err)

	stmt, err := tx.Prepare("insert into user(pid,firstname,lastname,mail,address,postnrandplace,phonenr) values(?,?,?,?,?,?,?)")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(number, first, last, mail, adr, ponr, phone)
	tx.Commit()
	checkErr(err)

}

//creates the database
func createDB() *sql.DB {
	//1. Open connection

	db, err := sql.Open("sqlite3", "./fakt.db") //return types = *DB, error
	checkErr(err)
	//defer db.Close()
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
						phonenr string)
						;`)
	checkErr(err)
	//fmt.Printf("%T\n", db)

	return db
}

//if error !=nil, print error message to web page
func checkErr(err error, args ...string) {
	if err != nil {
		fmt.Printf("ERROR : %q: %s\n", err, args)
	}
}

//get the next available ID for user
func getPersonNextNr() int {
	highestNr := 0
	for i := range person {
		//fmt.Printf("%v\t%T\n", person[i].number, person[i].number)
		if highestNr < person[i].Number {
			highestNr = person[i].Number
		}
	}
	highestNr++
	//fmt.Printf("The next number is = %T	%v\n ", highestNr, highestNr)

	return highestNr
}

//The web handler for adding persons
func addPersonsWeb(w http.ResponseWriter, r *http.Request) {

	err := tmpl["init.html"].ExecuteTemplate(w, "addUserCompletePage", "some data")
	if err != nil {
		fmt.Println("template execution error = ", err)
	}

	r.ParseForm()
	fn := r.FormValue("firstName")
	ln := r.FormValue("lastName")
	ma := r.FormValue("mail")
	ad := r.FormValue("address")
	pa := r.FormValue("poAddr")
	pn := r.FormValue("phone")

	if fn != "" {
		pid, _ := queryDBForNextCustomerPID(pDB)
		println("PID = ", pid)
		addUserToDB(pDB, pid, fn, ln, ma, ad, pa, pn)
		//fmt.Fprintf(w, name)
	} else {
		fmt.Fprintf(w, "Minimum needed is firstname")
	}
}

//The web handler for modifying a person
func modifyPersonWeb(w http.ResponseWriter, r *http.Request) {
	p := queryDBForAllUserInfo(pDB)

	err := tmpl["init.html"].ExecuteTemplate(w, "modifyUserCompletePage", p)
	if err != nil {
		fmt.Println("template execution error = ", err)
	}

	//p = queryDBForAllUserInfo(pDB)
	//	err := t.Execute(w, p) //execute and put all the data into the template, w is the writer, and p is slice with user info

	fmt.Fprint(w, err)
	r.ParseForm()
	//fn := r.FormValue("users")
	fn := r.FormValue("users")
	fmt.Fprint(w, fn)
}

//The web handler to show and print out all registered users in the database
func showPersonsWeb(w http.ResponseWriter, r *http.Request) {

	p := queryDBForAllUserInfo(pDB)
	err := tmpl["init.html"].ExecuteTemplate(w, "showUserCompletePage", p)
	if err != nil {
		fmt.Println("template execution error = ", err)
	}
	//err := t.Execute(w, p) //execute and put all the data into the template
	fmt.Fprint(w, err)
}
