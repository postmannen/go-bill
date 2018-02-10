package main

import (
	"database/sql"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/postmannen/fakt/data"
)

//webData struct, used to feed data to the web templates
type webData struct {
	Users            []data.User
	BLines           []data.BillLines
	Bills            []data.Bill
	ActiveUserID     int //to store the active user beeing worked on in the different web pages
	CurrentBillID    int //to store the active bill id beeing worked on in different web pages
	CurrentBill      data.Bill
	CurrentBillLines []data.BillLines
	PDB              *sql.DB
	IndexUser        int //to store the index nr. in slice where the chosen user is stored
}

var tmpl map[string]*template.Template //map to hold all templates

func init() {
	//initate the templates
	tmpl = make(map[string]*template.Template)
	tmpl["user.html"] = template.Must(template.ParseFiles("public/userTemplates.html"))
	tmpl["bill.html"] = template.Must(template.ParseFiles("public/billTemplates.html"))
}

func main() {

	//create DB and store pointer in pDB
	wData := webData{}
	wData.PDB = data.Create()
	defer wData.PDB.Close()

	//HandleFunc takes a handle (ResponseWriter) as first parameter,
	//and pointer to Request function as second parameter
	http.HandleFunc("/sp", wData.showUsersWeb)
	http.HandleFunc("/ap", wData.addUsersWeb)
	http.HandleFunc("/mp", wData.modifyUsersWeb)
	http.HandleFunc("/", wData.mainPage)
	http.HandleFunc("/du", wData.deleteUserWeb)
	http.HandleFunc("/createBillSelectUser", wData.webBillSelectUser)
	http.HandleFunc("/editBill", wData.webBillLines)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.ListenAndServe(":7000", nil)

}
