package main

import (
	"fmt"
	"net/http"
)

func main1() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/form", formHandler)
	fmt.Println("Starting server at port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting the server ")
	}
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "This is home! \n")
}
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "About\n")
}
func formHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html") // додали цей рядок
	if r.Method == http.MethodGet {
		fmt.Fprintf(w, `
		<form method="POST">
			<input type="text" name="name" placeholder="Enter your name">
			<button type="submit">Submit</button>
		</form>
	`)
	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(w, "Error parsing form : %v", err)
			return
		}
		name := r.FormValue("name")
		fmt.Fprintf(w, "Hello, %s!", name)

	}
}
