/*
2: Tested person struct, and string to int conversion for person ID.
    Iterate the slice of struct
3: Added invoice nr. to personStruct
4: Added some menus with print and add person as options
	Added getPersonNextNr function to look up the next available person nr
5: Added auto next number for person
	Added more variables describing person
	Added /sp for show person info
	Added /ap for add person
	input only get added when "add" button is pushed
*/
package main

import (
	"bufio"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

//PersonStruct is used for all customers and users
type PersonStruct struct { //some
	FirstName      string
	LastName       string
	Number         int
	Invoice        []int
	Mail           string
	Address        string
	PostNrAndPlace string
	PhoneNr        string
}

func (p PersonStruct) printStruct() {
	fmt.Println(p.FirstName)

}

type mittInterface interface {
	printStruct()
}

var person []PersonStruct

func printPerson() {
	for i := range person {
		fmt.Printf("Kunde nr: %v\t\tFornavn: %v\t\t Etternavn: %v\n", person[i].Number, person[i].FirstName, person[i].LastName)
	}
	//fmt.Println(person)
}

//func addPerson(first string, last string, number int, invoice []int) {
//addPerson function are only used initially to add some dummy persons
func addPerson(first string, last string, number int, invoice []int, mail string, adr string, ponr string, phone string) {
	person = append(person, PersonStruct{first, last, number, invoice, mail, adr, ponr, phone})
}

func addPerson2() {
	fn := getStringFromKeyboard("Firstname : ")
	ln := getStringFromKeyboard("Lastname : ")
	ml := getStringFromKeyboard("Mail : ")
	ad := getStringFromKeyboard("Address : ")
	po := getStringFromKeyboard("Post nr. and place : ")
	pn := getStringFromKeyboard("Phone nr. : ")
	//num := getIntFromKeyboard("Enter number for person : ")

	num := getPersonNextNr()
	person = append(person, PersonStruct{fn, ln, num, []int{}, ml, ad, po, pn})
}

func choosePerson() int {
	//Read input from console
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Choose person : ")
	inputText, err := reader.ReadString('\n')
	//read and wait for input, wait for "\n" before proceeding
	if err != nil {
		fmt.Printf("The input of text failed,with the error : %v\n", err)
	}

	inputText = strings.TrimSuffix(inputText, "\n")
	//remove endline "\n" from string, so strconv.Atoi() works

	inputTextToInt, err := strconv.Atoi(inputText) //convert string to int
	if err != nil {
		fmt.Printf("The string to int conversion did not work, error : %v\n", err)
	}

	//find the corresponding slice index nr. for the customer nr.
	var personSliceIndex int
	for i := range person {
		if person[i].Number == inputTextToInt {
			personSliceIndex = i
			//fmt.Println("The slice index number for chosen person = ", personSliceIndex)
		}
	}

	return personSliceIndex
}

func getPersonNextNr() int {
	highestNr := 0
	for i := range person {
		//fmt.Printf("%v\t%T\n", person[i].number, person[i].number)
		if highestNr < person[i].Number {
			highestNr = person[i].Number
		}
	}
	highestNr++
	//fmt.Printf("The next number is = %T	%v\n ", highestNr, highestNr)

	return highestNr
}

func getIntFromKeyboard(text string) int {
	//Read input from console
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(text)
	inputText, err := reader.ReadString('\n')
	//read and wait for input, wait for "\n" before proceeding
	if err != nil {
		fmt.Printf("The input of text failed,with the error : %v\n", err)
	}

	inputText = strings.TrimSuffix(inputText, "\n")
	//remove endline "\n" from string, so strconv.Atoi() works

	inputTextToInt, err := strconv.Atoi(inputText) //convert string to int
	if err != nil {
		fmt.Printf("The string to int conversion did not work, error : %v\n", err)
	}
	return inputTextToInt
}

func topMenu() int {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	fmt.Println("1. Show persons")
	fmt.Println("2. Add person")
	fmt.Println("3. Add invoice nr. to person")
	switch getIntFromKeyboard("Enter number : ") {
	case 1:
		printPerson()
	case 2:
		addPerson2()
	case 3:
		getPersonNextNr()
	}
	return 0
}

func getStringFromKeyboard(text string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(text)
	inputText, err := reader.ReadString('\n')

	if err != nil {
		fmt.Printf("The input of text failed,with the error : %v\n", err)
	}
	inputText = strings.TrimSuffix(inputText, "\n")

	return inputText
}

func addPersonsWeb(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Persons</h1>")
	//fmt.Fprintf(w, "<form>")
	t := template.New("Add Person")
	t.Parse(`
		<form>
		<div>
		<input type="text" name="firstName" placeholder="Fornavn">
		<input type="text" name="lastName" placeholder="Etternavn">
		<input type="text" name="mail" placeholder="Mail">
		<input type="text" name="address" placeholder="Adresse">
		<input type="text" name="poAddr" placeholder="Post nr og sted">
		<input type="text" name="phone" placeholder="Telefon nr.">
		<input type="submit" name="button1" value="Add">			
		</div>
		</form>

		<form action="/sp" method="post">
		<input type="submit" name="button2" value="show persons">
		</form>
	`)

	//fmt.Fprintf(w, "<form>")
	fmt.Println("KNAPP = ", r.Form.Get("button1"))

	err := t.Execute(w, person)
	if err != nil {
		fmt.Fprint(w, "Error msg : ", err)
	}
	r.ParseForm()
	fn := r.FormValue("firstName")
	ln := r.FormValue("lastName")
	ma := r.FormValue("mail")
	ad := r.FormValue("address")
	pa := r.FormValue("poAddr")
	pn := r.FormValue("phone")

	if fn != "" {
		addPerson(fn, ln, getPersonNextNr(), []int{1, 2, 3}, ma, ad, pa, pn)
		//fmt.Fprintf(w, name)
	} else {
		fmt.Fprintf(w, "Fyll ut verdi på fornavn")
	}
}

func showPersonsWeb(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "<h1>Persons</h1>")
	fmt.Fprintf(w, "<form>")
	t := template.New("Show Persons")
	t.Parse(`
		<form action="/ap" method="post">
		<input type="submit" name="button3" value="add persons">
		</form>

		{{range .}}
		<div>
			<input title="The name" type="text" name="navn" placeholder={{.FirstName}}>
			<input title="The lastname" type="text" name="etternavn" placeholder={{.LastName}}>
			<input title="The customer nr" type="text" name="nummer" placeholder={{.Number}}>
			<input title="The mail address" type="text" name="mail" placeholder={{.Mail}}>
			<input title="The address" type="text" name="adresse" placeholder={{.Address}}>
			<input title="The Post nr and place" type="text" name="post nr og sted" placeholder={{.PostNrAndPlace}}>
			<input title="The phone number" type="text" name="telefon nummer" placeholder={{.PhoneNr}}>			
			</div>
		{{end}}
		`)
	fmt.Fprintf(w, "</form>")

	err := t.Execute(w, person)
	fmt.Fprint(w, err)
}

func main() {
	addPerson("Donald", "Duck", getPersonNextNr(), []int{1, 2, 3}, "info@andeby.no", "veien 1", "1 Andeby", "1111111")
	addPerson("Mikke", "Mus", getPersonNextNr(), []int{1, 2, 3}, "info@andeby.no", "veien 1", "1 Andeby", "1111111")

	/*	for true {
		topMenu()
		time.Sleep(time.Second * 2)
	}*/

	http.HandleFunc("/sp", showPersonsWeb)
	http.HandleFunc("/ap", addPersonsWeb)
	http.ListenAndServe(":8000", nil)

}
