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
	Learn to pass values to web HandleFunc's, to make the exported user struct not exported

Ideas:
	Make the primary keys uid and bill ID random numbers, so you can sync the database
	 between different devices without getting a conflict.
	 Sorting can be done on a dummy index value that don't have to be unique



*/
package main

import (
	"database/sql"
	"html/template"
	"net/http"

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