/*
 2: Tested person struct, and string to int conversion for person ID.
    Iterate the slice of struct
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type personStruct struct {
	firstName string
	lastName  string
	number    int
}

var person []personStruct

func printPerson() {
	for i := range person {
		fmt.Printf("Kunde nr: %v\tFornavn: %v\t Etternavn: %v\n", person[i].number, person[i].firstName, person[i].lastName)
	}
}

func addPerson(first string, last string, number int) {
	person = append(person, personStruct{first, last, number})
}

func choosePerson() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Choose person : ")
	inputText, err := reader.ReadString('\n') //read and wait for input, wait for "\n" before proceeding
	if err != nil {
		fmt.Printf("The input of text failed,with the error : %v\n", err)
	}
	inputText = strings.TrimSuffix(inputText, "\n") //remove endline "\n" from string, so strconv.Atoi() works
	fmt.Printf("The string read with reader.ReadString = %v and the type is %T\n", inputText, inputText)

	intputTextToInt, err := strconv.Atoi(inputText) //ignored the error variable with _
	if err != nil {
		fmt.Printf("The string to int conversion did not work, error : %v\n", err)
	}

	fmt.Printf("The converted string = %v \tType is %T\n", intputTextToInt, intputTextToInt)
	//	fmt.Printf("You choose : %v and the type is %T\n", inputText, inputText)
}

func main() {
	addPerson("Donald", "Duck", 1)
	addPerson("Mikke", "Mus", 2)
	addPerson("Onkel", "Skrue", 3)
	printPerson()
	choosePerson()
}
