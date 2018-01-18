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
	Invoice        []int
	Mail           string
	Address        string
	PostNrAndPlace string
	PhoneNr        string
}

//Not used
func (p PersonStruct) printStruct() {
	fmt.Println(p.FirstName)

}

//Not used
type mittInterface interface {
	printStruct()
}

var person []PersonStruct

func init() {
	addPerson(getPersonNextNr(), "Donald", "Duck", []int{1, 2, 3}, "info@andeby.no", "veien 1", "1 Andeby", "1111111")
	addPerson(getPersonNextNr(), "Mikke", "Mus", []int{1, 2, 3}, "info@andeby.no", "veien 1", "1 Andeby", "1111111")
}

func main() {
	//create DB and store pointer in pDB
	pDB := createDB()
	defer pDB.Close()
	//fmt.Printf("Typen av variabel er %T", pDB)

	addUserToDB(pDB, 1, "Donald", "Duck", []int{1, 2, 3}, "info@andeby.no", "veien 1", "1 Andeby", "1111111")
	addUserToDB(pDB, 2, "Mikke", "Mus", []int{1, 2, 3}, "info@andeby.no", "veien 1", "1 Andeby", "1111111")

	queryDB(pDB)

	http.HandleFunc("/sp", showPersonsWeb)
	http.HandleFunc("/ap", addPersonsWeb)
	http.ListenAndServe(":8000", nil)

}

func queryDB(db *sql.DB) {
	rows, err := db.Query("select * from user")
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		checkErr(err)
		fmt.Printf("SCAN : %T : %v\n", &name, &name)
	}

}

func addUserToDB(db *sql.DB, number int, first string, last string, invoice []int, mail string, adr string, ponr string, phone string) {
	tx, err := db.Begin()
	checkErr(err)

	stmt, err := tx.Prepare("insert into user(pid,firstname,lastname,mail,address,postnrandplace,phonenr) values(?,?,?,?,?,?,?)")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(number, first, last, mail, adr, ponr, phone)
	tx.Commit()
	checkErr(err)

}

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
	fmt.Printf("%T\n", db)

	return db
}

func checkErr(err error, args ...string) {
	if err != nil {
		fmt.Printf("ERROR : %q: %s\n", err, args)
	}
}

//func addPerson(first string, last string, number int, invoice []int) {
//addPerson function are only used initially to add some dummy persons
func addPerson(number int, first string, last string, invoice []int, mail string, adr string, ponr string, phone string) {
	person = append(person, PersonStruct{number, first, last, invoice, mail, adr, ponr, phone})
}

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

//The web for adding persons
func addPersonsWeb(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Persons</h1>")
	//fmt.Fprintf(w, "<form>")
	t := template.New("Add Person")
	t.Parse(`
		<form action="/sp" method="post">
		<input type="submit" name="button2" value="show persons">
		</form>

		<form>
		<div>
		<input type="text" name="firstName" placeholder="Fornavn">
		<input type="text" name="lastName" placeholder="Etternavn">
		<input type="text" name="mail" placeholder="Mail">
		<input type="text" name="address" placeholder="Adresse">
		<input type="text" name="poAddr" placeholder="Post nr og sted">
		<input type="text" name="phone" placeholder="Telefon nr.">
		<input type="submit" name="button1" value="Add">			
		</div>
		</form>
	`)

	//fmt.Fprintf(w, "<form>")

	err := t.Execute(w, person)
	if err != nil {
		fmt.Fprint(w, "Error msg : ", err)
	}
	r.ParseForm()
	fn := r.FormValue("firstName")
	ln := r.FormValue("lastName")
	ma := r.FormValue("mail")
	ad := r.FormValue("address")
	pa := r.FormValue("poAddr")
	pn := r.FormValue("phone")

	if fn != "" {
		addPerson(getPersonNextNr(), fn, ln, []int{1, 2, 3}, ma, ad, pa, pn)
		//fmt.Fprintf(w, name)
	} else {
		fmt.Fprintf(w, "Fyll ut verdi på fornavn")
	}
}

func showPersonsWeb(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "<h1>Persons</h1>")
	fmt.Fprintf(w, "<form>")
	t := template.New("Show Persons")
	t.Parse(`
		<form action="/ap" method="post">
		<input type="submit" name="button3" value="add persons">
		</form>

		{{range .}}
		<div>
			<input title="The customer nr" type="text" name="nummer" placeholder={{.Number}}>
			<input title="The name" type="text" name="navn" placeholder={{.FirstName}}>
			<input title="The lastname" type="text" name="etternavn" placeholder={{.LastName}}>
			<input title="The mail address" type="text" name="mail" placeholder={{.Mail}}>
			<input title="The address" type="text" name="adresse" placeholder={{.Address}}>
			<input title="The Post nr and place" type="text" name="post nr og sted" placeholder={{.PostNrAndPlace}}>
			<input title="The phone number" type="text" name="telefon nummer" placeholder={{.PhoneNr}}>			
			</div>
		{{end}}
		`)
	fmt.Fprintf(w, "</form>")

	err := t.Execute(w, person)
	fmt.Fprint(w, err)
}
