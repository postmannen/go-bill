package main

import (
	"net/http"
	"sync"
	"text/template"
)

//is this one of any use ???
func (d *webData) editBill() http.HandlerFunc {
	var init sync.Once
	var tpl *template.Template
	init.Do(func() {
		tpl = template.Must(template.ParseFiles("public/billTemplates.html"))
	})

	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, d)
	}

}
