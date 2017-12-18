package main

import "fmt"

var person = map[int]*struct {
	firstName string
	lastName  string
	mail      string
}{
	1: {"apekatt", "hest", "apekatt@hest.no"},
	2: {"Anne", "Annesen", "anne@annesen.no"},
	3: {"Knut", "Knutsen", "knut@knutsen.no"},
	4: {"Ole", "Olesen", "ole@olesen.no"},
	5: {"Per", "Person", "per@person.no"},
}

func main() {
	fmt.Println(person)
	fmt.Println("firstName på person med index nr. 1 : ", person[1].firstName)

	for i, v := range person {
		fmt.Println(i, v.firstName, v.lastName, v.mail)
	}

}
