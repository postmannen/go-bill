package main

import (
	"log"
	"net/http"
	"text/template"
)

func (d *webData) editBill(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("public/bill.html")
	if err != nil {
		log.Println("error: template.ParseFiles: ", err)
	}
	t.Execute(w, d)

}
