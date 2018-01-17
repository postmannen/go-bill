package main

import (
	"html/template"
	"net/http"

	"github.com/postmannen/fakt/fakt20/db"
	"github.com/postmannen/fakt/fakt20/web"

	_ "github.com/mattn/go-sqlite3"
)

var tmpl map[string]*template.Template //map to hold all templates
var indexNR int                        //to store the index nr. in slice where chosen person is stored
var myWeb web

func init() {
	//initate the templates
	tmpl = make(map[string]*template.Template)
	tmpl["init.html"] = template.Must(template.ParseFiles("static/templates.html"))
}

func main() {
	//create DB and store pointer in pDB
	db.PDB = db.Create()
	defer db.PDB.Close()

	//HandleFunc takes a handle (ResponseWriter) as first parameter,
	//and pointer to Request function as second parameter
	http.HandleFunc("/sp", web.ShowUsersWeb)
	http.HandleFunc("/ap", web.addUsersWeb)
	http.HandleFunc("/mp", web.modifyUsersWeb)
	http.HandleFunc("/", web.mainPage)
	http.HandleFunc("/du", web.deleteUserWeb)
	http.HandleFunc("/createBillSelectUser", web.webBillSelectUser)
	http.HandleFunc("/editBill", web.webBillLines)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":7000", nil)

}
