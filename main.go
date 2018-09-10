package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "net/http/pprof"

	"github.com/pkg/profile"

	_ "github.com/mattn/go-sqlite3"
)

type server struct {
	address string
	wData   *webData
}

//webData struct, used to feed data to the web templates
type webData struct {
	Users            []User
	ActiveUserID     int //to store the active user beeing worked on in the different web pages
	CurrentBillID    int //to store the active bill id beeing worked on in different web pages
	CurrentAdmin     User
	CurrentUser      User
	CurrentBill      Bill
	CurrentBillLines []BillLines
	PDB              *sql.DB
	IndexUser        int    //to store the index nr. in slice where the chosen user is stored
	Currency         string //TODO: Make this linked to chosen language for admin user
	tpl              *template.Template
	msgToTemplate    map[string]string //Map of all templates to be called via websocket
	DivID            int
}

func newServer() *server {
	//Load the template files
	t, err := template.ParseFiles("public/userTemplates.html",
		"public/billTemplates.html", "public/socketTemplates.gohtml")
	if err != nil {
		fmt.Println("error: Parsing templates: ", err)
	}

	return &server{
		address: "localhost:8080",
		wData: &webData{
			tpl:   t,
			DivID: 0,
		},
	}
}

func (s *server) handlers() {
	http.HandleFunc("/echo", s.wData.socketHandler())
	http.HandleFunc("/showUser", s.wData.showUsersWeb)
	http.HandleFunc("/addUser", s.wData.addUsersWeb)
	http.HandleFunc("/modifyUser", s.wData.modifyUsersWeb)
	http.HandleFunc("/modifyAdmin", s.wData.modifyAdminWeb)
	http.HandleFunc("/", s.wData.mainPage)
	http.HandleFunc("/deleteUser", s.wData.deleteUserWeb)
	http.HandleFunc("/createBillSelectUser", s.wData.webBillSelectUser)
	http.HandleFunc("/editBill", s.wData.webBillLines)
	http.HandleFunc("/printBill", s.wData.printBill)
	http.HandleFunc("/newBill", s.wData.newBill)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
}

const databaseFileName = "fakt.db"

func main() {
	defer profile.Start().Stop()
	s := newServer()
	s.wData.msgToTemplate = make(map[string]string)
	s.wData.msgToTemplate = map[string]string{
		"addButton":     "buttonTemplate1",
		"addTemplate":   "socketTemplate1",
		"addParagraph":  "paragraphTemplate1",
		"userSelection": "createBillUserSelection",
	}

	//create DB and store pointer in pDB
	s.wData.PDB = Create(databaseFileName)
	defer s.wData.PDB.Close()
	s.wData.Currency = "$"

	//execute all the handlers
	s.handlers()
	http.ListenAndServe(s.address, nil)
}
