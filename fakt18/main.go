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

//Bill struct specifications
type Bill struct {
	BillID      int
	UserID      int
	CreatedDate string
	DueDate     string
	Comment     string
	TotalExVat  float64
	TotalIncVat float64
	Paid        int
}

//BillLines struct. Fields must be export (starting Capital letter) to be passed to template
type BillLines struct {
	BillID             int
	LineID             int
	ItemID             int
	Description        string
	Quantity           int
	DiscountPercentage int
	VatUsed            int
	PriceExVat         float64
	//just create some linenumbers for testing
}

//webData struct, used to feed data to the web templates
type webData struct {
	Users  []User
	BLines []BillLines
}

var pDB *sql.DB                        //The pointer to use with the Database
var tmpl map[string]*template.Template //map to hold all templates
var indexNR int                        //to store the index nr. in slice where chosen person is stored
var activeUserID int                   //to store the active user beeing worked on in the different web pages
var currentBillID int                  //to store the active bill id beeing worked on in different web pages

func init() {
	//initate the templates
	tmpl = make(map[string]*template.Template)
	tmpl["init.html"] = template.Must(template.ParseFiles("templates.html"))
}

func main() {
	//create DB and store pointer in pDB
	pDB = createDB()
	defer pDB.Close()

	//HandleFunc takes a handle (ResponseWriter) as first parameter,
	//and pointer to Request function as second parameter
	http.HandleFunc("/sp", showUsersWeb)
	http.HandleFunc("/ap", addUsersWeb)
	http.HandleFunc("/mp", modifyUsersWeb)
	http.HandleFunc("/", mainPage)
	http.HandleFunc("/du", deleteUserWeb)
	http.HandleFunc("/createBillSelectUser", billCreateWebSelectUser)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":7000", nil)

}
