package main

import (
	"os"
	"testing"
)

func Test_newDeck(t *testing.T) {
	d := newDeck()

	if len(d) != 52 {
		t.Errorf("expected 52, got %v", len(d))
	}

	if d[0] != "Ace of Spades" {
		t.Errorf("expected Ace of spades, but got %v", d[0])
	}

	if d[len(d)-1] != "King of Clubs" {
		t.Errorf("expected King of Clubs, but got %v", d[len(d)-1])
	}
}

func TestSaveToFileAndLoadDeckFromFile(t *testing.T) {
	os.Remove("_decktesting") //Fjern tidligere test fil.

	deck := newDeck()
	deck.saveToFile("_decktesting")

	loadedDeck := loadDeckFromFile("_decktesting")
	if len(loadedDeck) != 52 {
		t.Errorf("expected 52, got %v", len(loadedDeck))
	}

	os.Remove("_decktesting")
}
