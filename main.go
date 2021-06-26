package main

import (
	"anytimes/http"
	"fmt"
	"os"
)

func main() {
	token := os.Getenv("ANYTIME_SESSION_TOKEN")
	sig := os.Getenv("ANYTIME_SESSION_SIG")
	fmt.Printf("Token %s Sig %s \n", token, sig)
	http.ReserverTime()
}
