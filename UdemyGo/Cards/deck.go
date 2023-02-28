package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type deck []string

func newDeck() deck {
	cards := deck{}

	cardSuits := []string{"Spades", "Hearts", "Diamonds", "Clubs"}
	cardNumbers := []string{"Ace", "2", "3", "4", "5", "6", "7", "8", "9", "10", "Jack", "Queen", "King"}

	for _, suit := range cardSuits {
		for _, number := range cardNumbers {
			cards = append(cards, number+" of "+suit)
		}

	}

	return cards
}

// En function som starter med () er en reciever function
// Any variable of type deck now gets access to the print method
func (d deck) print() {
	for _, card := range d {
		fmt.Println(card)
	}
}

func (d deck) shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d), func(i, j int) { d[i], d[j] = d[j], d[i] })
}

// Split input deck to two slices. One for the hand other for rest of deck. (deck, deck) are the 2 return values
func deal(d deck, handSize int) (deck, deck) {
	return d[:handSize], d[handSize:] //['0':handSize] , ["handSize:'End]
}

func (d deck) toString() string {
	return strings.Join([]string(d), ",") //[]string(d) -> Type conversion til String med deck d  - Og så Join dem med ,
}

func (d deck) saveToFile(filename string) error {
	return os.WriteFile(filename, []byte(d.toString()), 0666) //permission 0666 : Anyone can read and write
}

func loadDeckFromFile(filename string) deck {
	byteSlice, err := os.ReadFile(filename) //Når vi ser en byteSlice, kan man tænke på det som en String
	if err != nil {                         //Hvis der er en error, så print og luk program
		fmt.Println("Error:", err)
		os.Exit(1)

	}
	sliceOfStrings := strings.Split(string(byteSlice), ",") //Converter byteSlice til en String. Og så laver jeg split på String og får []string
	return deck(sliceOfStrings)                             //Converter slice til Deck
}
