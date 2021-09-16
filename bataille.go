package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
    http.HandleFunc("/", TimeHandler)
		http.HandleFunc("/add", AddHandler)
		http.HandleFunc("/entries", EntriesHandler)
    fmt.Println("Server started at port 4567")
    log.Fatal(http.ListenAndServe(":4567", nil))
}

// GET /
func TimeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
    fmt.Fprintf(w, "Hello, there\nMethod is allowed\n")
	} else {
		currentTime := time.Now()
    fmt.Fprintf(w, "Il est %dh%d.\n",currentTime.Hour(), currentTime.Local().Minute())
	}
}

// POST /add
func AddHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method != "POST" {
    fmt.Fprintf(w, "Hello, there\nMethod not supported")
	} else {
    author:= r.Form.Get("author")
		entry:= r.Form.Get("entry")
			if (len(author) <=0 || len(entry)<=0) {
				fmt.Fprintf(w, "Missing parameters")
			}else{
				addEntry(entry)
				fmt.Fprintf(w,author+": "+entry)
			}
	}
}

func addEntry(entry string){
	fmt.Println(entry)
	file, err := os.OpenFile("test.txt", os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()

	_, err2 := file.WriteString(entry + "\n")
	if err2 != nil {
		log.Fatalf("failed writing to file: %s", err2)
	}
}

// GET /entries

func EntriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
    fmt.Fprintf(w, "Hello, there\nOnly GET method is allowed")
	}else{
		data, err := os.ReadFile("entries.txt")
	if err != nil {
		log.Panicf("failed reading data from file: %s", err)
	}
	fmt.Printf("%s", data)
	fmt.Fprintf(w,"%s", data)
	}
}
