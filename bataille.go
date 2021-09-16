package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)
type battleShipBoard struct {
	Board [10][10]boardState
}
type boardState int

const (
	stateEmpty   boardState = iota // 0
	stateShip    boardState = iota // 1
	stateHit     boardState = iota // 2
	stateAttempt boardState = iota // 3
)

func random(min int, max int) int {
	return rand.Intn(max-min) + min
}

func main() {
	http.HandleFunc("/hit", HitHandler)
	http.HandleFunc("/board", BoardHandler)
	http.HandleFunc("/boats", BoatsHandler)
	fmt.Println("Server started at port 4567")
	log.Fatal(http.ListenAndServe(":4567", nil))
	rand.Seed(time.Now().UnixNano())

	grid := [10][10]int{{}, {}}
	fmt.Println(grid[0])
	board := makeBoard()
	for i := 0; i < 9; i++ {
		fmt.Println(board.Board[i])
	}

}

func (b *battleShipBoard) Draw() string {
	str := ""
	str += fmt.Sprint("_|A|B|C|D|E|F|G|H|I|J|_\n")
	for y, stripe := range b.Board {
		str += fmt.Sprintf("%d|", y)
		for _, x := range stripe {
			str += fmt.Sprintf("%s|", x)
		}
		str += fmt.Sprintf("%d\n", y)
	}
	str += fmt.Sprint("_|A|B|C|D|E|F|G|H|I|J|_\n")
	return str
}

func cordsToNumbers(in string) (X, Y int) {
	in = strings.ToLower(in)

	i1 := int(in[0])
	if i1 > 96 && i1 < 107 {
		X = i1 - 97
	} else {
		return -1, -1
	}

	i2, err := strconv.ParseInt(string(in[1]), 10, 64)
	if err != nil {
		return -1, -1
	}

	if i2 < 11 {
		return X, int(i2)
	}
	return -1, -1
}

func combineBoard(b1, b2 battleShipBoard) string {
	b1r, b2r := b1.Draw(), b2.Draw()

	b1rs, b2rs := strings.Split(b1r, "\n"), strings.Split(b2r, "\n")

	str := ""
	for k, _ := range b1rs {
		str += fmt.Sprintf("%s     %s\n", b1rs[k], b2rs[k])
	}

	return str
}

func makeBoard() battleShipBoard {
	a := battleShipBoard{}
	// ri, _ := cr.Int(cr.Reader, big.NewInt(math.MaxInt64))
	// rand.Seed(ri.Int64())

	a = placeShip(5, a)
	a = placeShip(4, a)
	a = placeShip(3, a)
	a = placeShip(3, a)
	a = placeShip(2, a)
	return a
}

func placeShip(size int, bo battleShipBoard) battleShipBoard {

	for {
		board := bo
		sideways := rand.Int() % 2

		if sideways == 0 { // ship goes up
			X := rand.Int() % 10
			Y := rand.Int() % 10
			if Y+size > 10 {
				continue
			}

			for y := Y; y < Y+size; y++ {
				if board.Board[y][X] != stateEmpty {
					continue
				}
				board.Board[y][X] = stateShip
			}
			bo = board
			break
		} else {
			X := rand.Int() % 10
			Y := rand.Int() % 10

			if X+size > 10 {
				continue
			}

			for x := X; x < X+size; x++ {
				if board.Board[Y][x] != stateEmpty {
					continue
				}
				board.Board[Y][x] = stateShip
			}
			bo = board
			break
		}
	}

	return bo
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
		board:= [10][10]int{{}, {}}
		fmt.Println(board)
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
