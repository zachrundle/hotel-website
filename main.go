package main

import (
  "net/http"
  "fmt"
  "errors"
)

const portNumber = ":8080"

// Home is the home page handler
func Home(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "This is the home page")
}

// About is the about page handler
func About(w http.ResponseWriter, r *http.Request) {
  sum,_ := addValues(2, 2)
  _, _ = fmt.Fprintf(w, "This is the about page and 2 + 2 is %d", sum)
}

// addValues add to integers and returns the same
func addValues(x, y int) (int, error) {
  sum := x + y
  return sum, nil
}

func Divide(w http.ResponseWriter, r *http.Request) {
  x := float32(100.0)
  y := float32(10.0)
  
  f, err := divideValues(x, y)
  if err != nil {
    fmt.Fprintf(w, "Cannot divide by 0")
    return
  }
  fmt.Fprintf(w, fmt.Sprintf("%f divided by %f", x, y, f))
}

func divideValues(x, y float32) (float32, error) {
  if y <= 0 {
    err := errors.New("cannot divide by zero")
    return 0, err
  }

  result := x / y
  return result, nil
}
// main is the main application function
func main() {
  http.HandleFunc("/", Home)
  http.HandleFunc("/about", About)
  http.HandleFunc("/divide", Divide)

  
  fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))
  http.ListenAndServe(portNumber, nil)

}
 