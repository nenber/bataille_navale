package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"net/url"
	"strings"
	"time"
)

var BOARD = makeBoard()
var LIFE = 0

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

func ParseIPPort(s string) (ip net.IP, port, space string, err error) {
	ip = net.ParseIP(s)
	if ip == nil {
		var host string
		host, port, err = net.SplitHostPort(s)
		if err != nil {
			return
		}
		if port != "" {
			// This check only makes sense if service names are not allowed
			if _, err = strconv.ParseUint(port, 10, 16); err != nil {
				return
			}
		}
		ip = net.ParseIP(host)
	}
	if ip == nil {
		err = errors.New("invalid address format")
	} else {
		space = "IPv6"
		if ip4 := ip.To4(); ip4 != nil {
			space = "IPv4"
			ip = ip4
		}
	}
	return
}

func random(min int, max int) int {
	return rand.Intn(max-min) + min
}
func runServer(port int) {
	go http.HandleFunc("/hit", HitHandler)
	go http.HandleFunc("/board", BoardHandler)
	go http.HandleFunc("/boats", BoatsHandler)
	// fmt.Println("Server started at port 4567")
	addr := ":" + strconv.Itoa(port)
	log.Fatal(http.ListenAndServe(addr, nil))

}

func main() {
	myPort := flag.String("port", ":4567", "Port of server")

	// opponent := flag.String("opp", "default value", "Opponent")
	flag.Parse()

	// ip, port, space, err := ParseIPPort(*opponent)

	// fmt.Println(ip, ":", port, "   ", space)
	number, err := strconv.ParseInt(*myPort, 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	finalIntNum := int(number)
	println(finalIntNum)

	go runServer(finalIntNum)
	rand.Seed(time.Now().UnixNano())

	// grid := [10][10]int{{}, {}}
	// fmt.Println(grid[0])
	// board := makeBoard()
	// for i := 0; i < 10; i++ {
	// 	fmt.Println(board.Board[i])
	// }

	opponents := askOpponentTarget()
	play(opponents)

}
func play(opponents []string) {
	for i := 0; ; i++ {
		println("Quel joueur souhaite-tu attaquer ?")
		for i := 0; i < len(opponents); i++ {
			println("#", i+1, " (", opponents[i], ")")
		}
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		result := scanner.Text()
		if u, err := strconv.Atoi(result); err == nil {
			if len(opponents) >= u {
				println("Quel case ? (A-J:1-10)")
				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan()
				target := strings.Split(scanner.Text(), ":")
				if len(target) < 2 {
					println("Format incorrect. Le format doit respecter \"A-J:1-10\"")
				} else {
					isAlpha := regexp.MustCompile(`^[A-za-z]+$`).MatchString(target[0])
					if j, err := strconv.Atoi(target[1]); err == nil {
						if j >= 1 && j <= 10 && isAlpha == true {
							//on envoi la requete hit
							data := url.Values{
								"hit": {target[0] + ":" + target[1]},
							}
							addr := "http://" + opponents[u-1] + "/hit"
							resp, err := http.PostForm(addr, data)

							if err != nil {
								log.Fatal(err)
							}

							var res map[string]interface{}

							json.NewDecoder(resp.Body).Decode(&res)

							fmt.Println(res["form"])
						} else {
							println("Cette case n'existe pas")
						}
					} else {
						println("Cette case n'existe pas")
					}
				}
			}
		} else {
			println("hey ce joueur n'existe pas, tape son id pour le selectionner !")
		}
	}
}

func askOpponentTarget() []string {

	var opponentsAdrr []string

	for i := 0; ; i++ {
		if i == 0 {
			println("Adresse de l'adversaire #", i+1, " (sous la forme \"127.0.0.1:4567\") : ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			result := scanner.Text()

			ip, port, space, err := ParseIPPort(result)
			_ = space
			if err != nil {
				println("L'adresse n'est pas bonne")
				i--
			} else {
				addr := ip.String() + ":" + port
				opponentsAdrr = append(opponentsAdrr, addr)

			}

		} else {
			println("Adresse de l'adversaire #", i+1, " (sous la forme \"127.0.0.1:4567\") (jouer(j)) : ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			result := scanner.Text()
			if result == "j" || result == "J" {
				break
			}

			ip, port, space, err := ParseIPPort(result)
			_ = space
			if err != nil {
				println("L'adresse n'est pas bonne")
				i--
			} else {
				addr := ip.String() + ":" + port
				opponentsAdrr = append(opponentsAdrr, addr)

			}
		}
	}
	println("Voici les adresses de tes adversaires : ")
	for i := 0; i < len(opponentsAdrr); i++ {
		println("- ", opponentsAdrr[i])
	}
	return opponentsAdrr

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
		get_hit_x := r.Form.Get("hit_x")
		get_hit_y := r.Form.Get("hit_y")
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
		if hit_x > 9 || hit_x < 0 || hit_y > 9 || hit_y < 0 {
			fmt.Fprintf(w, "Valeur invalide")
		}
		// check if it hits the boat on the board
	}
}

// GET /board
func BoardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		fmt.Fprintf(w, "Hello, there\nOnly GET method is allowed")
	} else {

			for i := 0; i < 10; i++ {
				fmt.Println(BOARD.Board[i])
				fmt.Fprintln(w, BOARD.Board[i])
			}
	}
}

// GET /boats
func BoatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		fmt.Fprintf(w, "Hello, there\nOnly GET method is allowed")
	} else {
		for i := 0; i < 10; i++ {
			for j := 0; j < 10; j++ {
				if BOARD.Board[i][j] == 1 {
					LIFE++
				}
			}
		}
		if LIFE == 16 {
			fmt.Fprintln(w, " Votre nombre de vies est de ",LIFE)
		}else{
		fmt.Fprintln(w, "Vous avez perdu ", 16-LIFE, "de vies \n il vous en reste ",LIFE)
		}
	}
}
