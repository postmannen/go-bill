package main

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/mattn/go-sqlite3"
	"github.com/postmannen/go-bill/data"
)

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
}

type server struct {
	addr   string
	router *mux.Router
}

func newServer() *server {
	return &server{
		addr:   ":8080",
		router: mux.NewRouter(),
	}
}

var tmpl map[string]*template.Template //map to hold all templates

func init() {
	//initate the templates
	tmpl = make(map[string]*template.Template)
	tmpl["user.html"] = template.Must(template.ParseFiles("public/userTemplates.html"))
	tmpl["bill.html"] = template.Must(template.ParseFiles("public/billTemplates.html"))
}

func main() {
	s := newServer()

	//create DB and store pointer in pDB
	wData := webData{}
	wData.PDB = data.Create()
	defer wData.PDB.Close()
	wData.Currency = "$"

	//openBrowser()

	//HandleFunc takes a handle (ResponseWriter) as first parameter,
	//and pointer to Request function as second parameter
	s.router.HandleFunc("/sp", wData.showUsersWeb)
	s.router.HandleFunc("/ap", wData.addUsersWeb)
	s.router.HandleFunc("/mp", wData.modifyUsersWeb)
	s.router.HandleFunc("/modifyAdmin", wData.modifyAdminWeb)
	s.router.HandleFunc("/", wData.mainPage)
	s.router.HandleFunc("/du", wData.deleteUserWeb)
	s.router.HandleFunc("/createBillSelectUser", wData.webBillSelectUser)
	s.router.HandleFunc("/editBill", wData.webBillLines)
	s.router.HandleFunc("/eBill", wData.editBill)
	s.router.HandleFunc("/printBill", wData.printBill)
	s.router.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.ListenAndServe(s.addr, s.router)

}
