package main

import "fmt"

type bot interface {
	getGreeting() string
}

//For at implemente bot, skal den blot have en getGreeting() string metode
type englishBot struct{}
type spanishBot struct{}

func main() {
	eb := englishBot{}
	sb := spanishBot{}

	printGreeting(eb)
	printGreeting(sb)
}

func printGreeting(b bot) {
	fmt.Println(b.getGreeting())
}

func (englishBot) getGreeting() string {
	return "Hey mate"
}

func (spanishBot) getGreeting() string {
	return "O'la!"
}
