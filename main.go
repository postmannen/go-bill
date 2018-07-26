package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/postmannen/go-bill/data"
)

type server struct {
	address string
	wData   *webData
}

//webData struct, used to feed data to the web templates
type webData struct {
	Users            []data.User
	ActiveUserID     int //to store the active user beeing worked on in the different web pages
	CurrentBillID    int //to store the active bill id beeing worked on in different web pages
	CurrentAdmin     data.User
	CurrentUser      data.User
	CurrentBill      data.Bill
	CurrentBillLines []data.BillLines
	PDB              *sql.DB
	IndexUser        int    //to store the index nr. in slice where the chosen user is stored
	Currency         string //TODO: Make this linked to chosen language for admin user
	tpl              *template.Template
}

var tmpl map[string]*template.Template //map to hold all templates

func init() {
	/*//initate the templates
	tmpl = make(map[string]*template.Template)
	tmpl["user.html"] = template.Must(template.ParseFiles("public/userTemplates.html"))
	tmpl["bill.html"] = template.Must(template.ParseFiles("public/billTemplates.html"))*/
}

func newServer() *server {
	//Load the template files
	t, err := template.ParseFiles("public/userTemplates.html",
		"public/billTemplates.html")
	if err != nil {
		fmt.Println("error: Parsing templates: ", err)
	}

	return &server{
		address: "localhost:8080",
		wData: &webData{
			tpl: t,
		},
	}
}

func (s *server) handlers() {
	http.HandleFunc("/sp", s.wData.showUsersWeb)
	http.HandleFunc("/ap", s.wData.addUsersWeb)
	http.HandleFunc("/mp", s.wData.modifyUsersWeb)
	http.HandleFunc("/modifyAdmin", s.wData.modifyAdminWeb)
	http.HandleFunc("/", s.wData.mainPage)
	http.HandleFunc("/du", s.wData.deleteUserWeb)
	http.HandleFunc("/createBillSelectUser", s.wData.webBillSelectUser)
	http.HandleFunc("/editBill", s.wData.webBillLines)
	http.HandleFunc("/eBill", s.wData.editBill)
	http.HandleFunc("/printBill", s.wData.printBill)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
}

func main() {
	s := newServer()

	//create DB and store pointer in pDB
	s.wData.PDB = data.Create()
	defer s.wData.PDB.Close()
	s.wData.Currency = "$"

	//execute all the handlers
	s.handlers()
	http.ListenAndServe(s.address, nil)
}
