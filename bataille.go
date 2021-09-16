package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
		http.HandleFunc("/hit", HitHandler)
		http.HandleFunc("/board", BoardHandler)
		http.HandleFunc("/boats", BoatsHandler)
    fmt.Println("Server started at port 4567")
    log.Fatal(http.ListenAndServe(":4567", nil))
}

// POST /hit
func HitHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method != "POST" {
    fmt.Fprintf(w, "Hello, there\nMethod not supported")
	} else {
			get_hit_x:= r.Form.Get("hit_x")
			get_hit_y:=r.Form.Get("hit_y")
			hit_x := 0
			hit_y := 0
			if x, err := strconv.Atoi(get_hit_x); err == nil {
				fmt.Printf("i=%d,", x)
				hit_x = x
			}
			if y, err := strconv.Atoi(get_hit_y); err == nil {
				fmt.Printf("i=%d,", y)
				hit_y = y
			}
			if (hit_x > 9 || hit_x < 0 || hit_y > 9 || hit_y < 0) {
				fmt.Fprintf(w, "Valeur invalide")
			}
			// check if it hits the boat on the board
	}
}

// GET /board
func BoardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
    fmt.Fprintf(w, "Hello, there\nOnly GET method is allowed")
	}else{
    fmt.Fprintf(w, "Board")
	}
}

// GET /boats
func BoatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
    fmt.Fprintf(w, "Hello, there\nOnly GET method is allowed")
	}else{
    fmt.Fprintf(w, "How many boats")
	}
}
