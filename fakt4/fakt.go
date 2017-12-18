/*
2: Tested person struct, and string to int conversion for person ID.
    Iterate the slice of struct
3: Added invoice nr. to personStruct
4: Added some menus with print and add person as options
	Added getPersonNextNr function to look up the next available person nr
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type personStruct struct {
	firstName string
	lastName  string
	number    int
	invoice   []int
}

func (p personStruct) printStruct() {
	fmt.Println(p.firstName)

}

type mittInterface interface {
	printStruct()
}

var person []personStruct

func printPerson() {
	for i := range person {
		fmt.Printf("Kunde nr: %v\t\tFornavn: %v\t\t Etternavn: %v\n", person[i].number, person[i].firstName, person[i].lastName)
	}
	fmt.Println("*******************************")
	fmt.Println(person)
	fmt.Println("*******************************")
}

//func addPerson(first string, last string, number int, invoice []int) {
func addPerson(first string, last string, number int, invoice []int) {
	person = append(person, personStruct{first, last, number, invoice})
}

func addPerson2() {
	fn := getStringFromKeyboard("Firstname : ")
	ln := getStringFromKeyboard("Lastname : ")
	num := getIntFromKeyboard("Enter number for person : ")

	person = append(person, personStruct{fn, ln, num, []int{}})
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
		if person[i].number == inputTextToInt {
			personSliceIndex = i
			//fmt.Println("The slice index number for chosen person = ", personSliceIndex)
		}
	}

	return personSliceIndex
}

func getPersonNextNr() {
	highestNr := 0
	for i := range person {
		fmt.Printf("%v\t%T", person[i].number, person[i].number)
		if highestNr < person[i].number {
			highestNr = person[i].number
		}
	}
	highestNr++
	fmt.Println("The next number is = ", highestNr)
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

func main() {
	addPerson("Donald", "Duck", 1, []int{1, 2, 3})
	addPerson("Mikke", "Mus", 2, []int{1, 2, 3})
	addPerson("Onkel", "Skrue", 3, []int{1, 2, 3})
	addPerson("Christopher", "Walker", 21, []int{})

	for true {
		topMenu()
		time.Sleep(time.Second * 2)
	}

}
