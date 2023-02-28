package main

func main() {
	deck := loadDeckFromFile("MyDeck")
	deck.shuffle()
	deck.print()
}
