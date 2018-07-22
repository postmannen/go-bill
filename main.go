package main

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"

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
	return &webData{}
}

type server struct {
	addr   string      //the adress and port to listen on
	router *mux.Router //use gorilla mux for our router
	data   webData     //put all the user data into the server struct
	//msgToTemplate is a reference to know what html template to
	//be used based on which msg comming in from the client browser.
	msgToTemplate map[string]string
}

var formDecoder = schema.NewDecoder()

func newServer() *server {
	return &server{
		addr:          ":8080",
		router:        mux.NewRouter(),
		msgToTemplate: make(map[string]string),
	}
}

func (s *server) routes() {
	s.router.HandleFunc("/", s.data.mainPage())
	s.router.HandleFunc("/showUsers", s.data.showUsers())
	s.router.HandleFunc("/addUser", s.data.addUsers())
	s.router.HandleFunc("/modifyUser", s.data.modifyUsers())
	s.router.HandleFunc("/modifyAdmin", s.data.modifyAdmin())
	s.router.HandleFunc("/deleteUser", s.data.deleteUser())
	s.router.HandleFunc("/createBillSelectUser", s.data.selectUserForBill())
	s.router.HandleFunc("/editBill", s.data.editBill())
	s.router.HandleFunc("/printBill", s.data.printBill())
	s.router.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
}

func (s *server) templates() {
	//key = name sendt to backend from JS at client browser.
	//value = name of template to use
	s.msgToTemplate = map[string]string{
		//user templates
		"topMenu":             "topMenu",
		"addUser":             "addUser",
		"modifyUserSelection": "modifyUserSelection",
		"deleteUserSelection": "deleteUserSelection",
		"showAllUsers":        "showAllUsers",
		"modifyUser":          "modifyUser",
		//bill templates
		"createBillLines":         "createBillLines",
		"createBillUserSelection": "createBillUserSelection",
		"billShowUser":            "billShowUser",
		"showBillInfo":            "showBillInfo",
		"editBillSelectBox":       "editBillSelectBox",
	}
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
