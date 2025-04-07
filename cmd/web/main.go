package main

import (
	"fmt"
	"net/http"
	"github.com/zachrundle/hotel-website/pkg/handlers"
)

const portNumber = ":8080"



// main is the main application function
func main() {
  http.HandleFunc("/", handlers.Home)
  http.HandleFunc("/about", handlers.About)

  
  fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))
  http.ListenAndServe(portNumber, nil)

}
 