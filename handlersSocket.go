package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

//socketHandler is the handler who controls all the serverside part
//of the websocket. The other handlers like the rootHandle have to
//load a page containing the JS websocket code to start up the
//communication with the serside websocket.
//This handler is used with all the other handlers if they open a
//websocket on the client side.
func (d *webData) socketHandler() http.HandlerFunc {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	var init sync.Once
	var tpl *template.Template
	var err error

	init.Do(func() {
		tpl, err = template.ParseFiles("public/userTemplates.html",
			"public/billTemplates.html",
			"public/socketTemplates.gohtml")
		if err != nil {
			log.Printf("error: ParseFiles : %v\n", err)
		}
	})

	return func(w http.ResponseWriter, r *http.Request) {
		//upgrade the handler to a websocket connection
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("error: websocket Upgrade: ", err)
		}

		//divID is to keep track of the sections sendt to the
		//socket to be shown in the browser.
		//divID := 0

		for {
			//read the message from browser
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("error: websocket ReadMessage: ", err)
				return
			}

			//print message to console
			fmt.Printf("Client=%v typed : %v \n", conn.RemoteAddr(), string(msg))

			//loop through the map and check if there is a key in the map that match with
			//the msg comming in on the websocket from browser.
			//If there is no match, whats in msg will be sendt directly back ovet the socket,
			//to be printed out in the client browser.
			strMsg := string(msg)
			if strMsg != "" {
				tplName, ok := d.msgToTemplate[strMsg]
				if ok {
					//Declare a bytes.Buffer to hold the data for the executed template.
					var tplData bytes.Buffer
					//tplData is a bytes.Buffer, which is a type io.Writer. Here we choose
					//execute the template, but passing the output into tplData insted of
					//'w'. Then we can take the data in tplData and send them over the socket.
					tpl.ExecuteTemplate(&tplData, tplName, d)
					d := tplData.String()
					//New-lines between the html tags in the template source code
					//is shown in the browser. Trimming awat the new-lines in each line
					//in the template data.
					d = strings.TrimSpace(d)
					msg = []byte(d)
				}
				d.DivID++
			}

			//write message back on the socket to the browser
			err = conn.WriteMessage(msgType, msg)
			if err != nil {
				fmt.Println("error: WriteMessage failed :", err)
				return
			}

		}
	}
}
