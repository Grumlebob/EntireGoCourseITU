package main

import "fmt"

type contactInfo struct {
	email   string
	zipcode int
}

type person struct {
	firstName   string
	lastName    string
	contactInfo //Har automatisk field+type af samme navn
}

func main() {
	somePerson := person{
		firstName: "John",
		lastName:  "Hansen",
		contactInfo: contactInfo{
			email:   "john@example.com",
			zipcode: 123,
		},
	}

	personPointer := &somePerson //'&' = shift+6. Give me the memory address of the value this variable is pointing at
	//Ovenstående kan også udelades, da func (p *person) tager imod enten type af *person eller person
	personPointer.updateName("Jim")            //Virker ikke
	personPointer.updateNameByReference("Jim") //Virker
	somePerson.print()
}

//GO er et pass by value language. Så (p person) får values men opdatere ikke sin reference til person.updateName
func (p person) updateName(NewfirstName string) {
	p.firstName = NewfirstName
}

//For at få pass by reference til at virke
func (pointerToThePerson *person) updateNameByReference(NewfirstName string) {
	(*pointerToThePerson).firstName = NewfirstName //'*' = Give me the value this memory address is pointing at
}

func (p person) print() {
	//fmt.Println(person) //Printer values
	fmt.Printf("%+v", p) //printer alle fields+value
}
