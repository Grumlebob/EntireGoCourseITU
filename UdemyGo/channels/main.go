package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	links := []string{
		"http://google.com",
		"http://facebook.com",
		"http://reddit.com",
		"http://golang.org",
	}

	channel := make(chan (string))

	for _, l := range links {
		go checkLink(l, channel) //subroutine
	}

	/*
		for i := 0; i < len(links); i++ {
			fmt.Println(<-channel) //Blocking call, wait for message
		}
	*/

	//Repeating routine.
	//loop i channel: venter på at Channel har fået indhold og assigner det til variable l.
	//Så basically "watch channel" uendeligt
	for l := range channel {
		go func(link string) { //Function literal / lambda -  Gør at vores time.sleep ikke blokkere os ved at have den i CheckLink metoden
			time.Sleep(3 * time.Second)
			checkLink(link, channel)
		}(l) //parantes til sidst betyder "call vores ukendte func". Vi skal passe L som parameter grundet pass by value.
		//We never try to share a variable from main to child routine directly. We pass it as an argument (pass by value)
	}

}

func checkLink(link string, channel chan string) {
	_, err := http.Get(link)
	if err != nil {
		fmt.Println("Link is down: ", link)
		channel <- link
		return
	}
	fmt.Println("Link is up: ", link)
	channel <- link
}
