package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os/exec"
	"runtime"

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
	wData.Currency = "$"

	//openBrowser()

	//HandleFunc takes a handle (ResponseWriter) as first parameter,
	//and pointer to Request function as second parameter
	http.HandleFunc("/sp", wData.showUsersWeb)
	http.HandleFunc("/ap", wData.addUsersWeb)
	http.HandleFunc("/mp", wData.modifyUsersWeb)
	http.HandleFunc("/modifyAdmin", wData.modifyAdminWeb)
	http.HandleFunc("/", wData.mainPage)
	http.HandleFunc("/du", wData.deleteUserWeb)
	http.HandleFunc("/createBillSelectUser", wData.webBillSelectUser)
	http.HandleFunc("/editBill", wData.webBillLines)
	http.HandleFunc("/eBill", wData.editBill)
	http.HandleFunc("/printBill", wData.printBill)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.ListenAndServe(":7000", nil)

}

func openBrowser() {
	fmt.Println(runtime.GOOS)

	switch runtime.GOOS {
	case "darwin":
		fmt.Println("The OS which is chosen is MacOs")
		cmd := exec.Command("open", "http://localhost:7000")
		cmd.Run()

	}
}
