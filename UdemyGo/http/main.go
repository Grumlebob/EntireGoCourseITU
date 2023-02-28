package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type logWriter struct{}

func main() {

	resp, err := http.Get("http://google.com")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	/* "simpel dårlig måde at læse Read"
	bs := make([]byte, 999999) //initiliaser en slice med x empty pladser - så når den gror at den ikke skal dobble tit.
	resp.Body.Read(bs)
	fmt.Println(string(bs))
	*/

	//io.Copy(os.Stdout, resp.Body) //Tager en (writer, reader) - Husk writer kan ex sende til terminal.

	lw := logWriter{}

	io.Copy(lw, resp.Body)
}

func (logWriter) Write(bs []byte) (int, error) {
	fmt.Println(string(bs))
	fmt.Println("Just wrote this many bytes:", len(bs))
	return len(bs), nil
}
