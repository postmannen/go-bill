/*
2: Tested person struct, and string to int conversion for person ID.
    Iterate the slice of struct
3: Added invoice nr. to personStruct
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
	invoice   []int
}

var person []personStruct

func printPerson() {
	for i := range person {
		fmt.Printf("Kunde nr: %v\tFornavn: %v\t Etternavn: %v\n", person[i].number, person[i].firstName, person[i].lastName)
	}
	fmt.Println("*******************************")
	fmt.Println(person)
	fmt.Println("*******************************")
}

func addPerson(first string, last string, number int, invoice []int) {
	person = append(person, personStruct{first, last, number, invoice})
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

func main() {
	addPerson("Donald", "Duck", 1, []int{1, 2, 3})
	addPerson("Mikke", "Mus", 2, []int{1, 2, 3})
	addPerson("Onkel", "Skrue", 3, []int{1, 2, 3})
	addPerson("Christopher", "Walker", 4, []int{5, 6, 7})
	printPerson()
	fmt.Println("Slice index of person = ", choosePerson())

}
