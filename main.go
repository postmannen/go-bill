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

func newWebData() *webData {

}

type server struct {
	addr   string
	router *mux.Router
	data   webData
}

func newServer() *server {
	return &server{
		addr:   ":8080",
		router: mux.NewRouter(),
	}
}

func (s *server) routes() {
	s.router.HandleFunc("/sp", s.data.showUsersWeb)
	s.router.HandleFunc("/ap", s.data.addUsersWeb)
	s.router.HandleFunc("/mp", s.data.modifyUsersWeb)
	s.router.HandleFunc("/modifyAdmin", s.data.modifyAdminWeb)
	s.router.HandleFunc("/", s.data.mainPage)
	s.router.HandleFunc("/du", s.data.deleteUserWeb)
	s.router.HandleFunc("/createBillSelectUser", s.data.webBillSelectUser)
	s.router.HandleFunc("/editBill", s.data.webBillLines)
	s.router.HandleFunc("/eBill", s.data.editBill)
	s.router.HandleFunc("/printBill", s.data.printBill)
	s.router.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
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

	//initialize db, and create if not exist
	s.data.PDB = data.Create()
	defer s.data.PDB.Close()

	//should be changed based on language
	s.data.Currency = "$"

	s.routes()
	http.ListenAndServe(s.addr, s.router)

}
