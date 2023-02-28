package main

import "fmt"

func main() {
	//Create
	colors := map[string]string{
		"red":   "rød",
		"green": "grøn",
	}
	//Add
	colors["black"] = "sort"
	//Remove
	delete(colors, "red")

	printMap(colors)
}

//Iterating through map
func printMap(m map[string]string) {
	for keyColor, valueTranslate := range m {
		fmt.Println(keyColor, valueTranslate)
	}
}
