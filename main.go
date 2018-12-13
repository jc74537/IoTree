package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"./sound"
	"github.com/gorilla/mux"
)

type song struct {
	name  string
	rtttl string
}
type command struct {
	pin int
	command bool
}
type pattern struct {
	commands []command
}



var (
	library []song
)

func main() {
	library = []song{
		{"Jingle Bells", "Jingle Bells:d=4,o=5,b=170:b,b,b,p,b,b,b,p,b,d6,g.,8a,2b.,8p,c6,c6,c6.,8c6,c6,b,b,8b,8b,b,a,a,b,2a,2d6"},
		{"Santa Claus is coming to town", "Santa Clause is Coming Tonight:d=4,o=5,b=180:g,8e,8f,g,g.,8g,8a,8b,c6,2c6,8e,8f,g,g,g,8a,8g,f,2f,e,g,c,e,d,2f,b4,1c,p,g,8e,8f,g,g.,8g,8a,8b,c6,2c6,8e,8f,g,g,g,8a,8g,f,f,e,g,c,e,d,2f,b4,1c,p,c6,d6,c6,b,c6,a,2a,c6,d6,c6,b,c6,2a.,d6,e6,d6,c#6,d6,b,b,b,8b,8c6,d6,c6,b,a,g,p,g.,8g,8e,8f,g,g.,8g,8a,8b,c6,2c6,8e,8f,g,g,g,8a,8g,8f,2f,e,g,c,e,d,2f,d6,1c6."},
		{"We Wish you a Merry Christmas", "We Wish you a Merry Christmas:d=8,o=5,b=140:4d,4g,g,a,g,f#,4e,4c,4e,4a,a,b,a,g,4f#,4d,4f#,4b,b,c6,b,a,4g,4e,4d,4e,4a,4f#,2g"},
		{"12 days of Christmas", "On the 12th Day of christmas:d=8,o=5,b=150:d,d,4g,g,g,4g,g,g,a,b,c6,a,4b.,p,4d6,a,b,c6,a,d6,d6,a,b,c6,a,4d6,4e6,4d.6,p,d6,c6,b,a,4g,a,b,4c6,4e,4e,4d,g,a,b,c6,4b,4a,2g."},
		{"We Wish You A Merry Christmas", "We Wish you a Merry Christmas:d=8,o=5,b=140:4d,4g,g,a,g,f#,4e,4c,4e,4a,a,b,a,g,4f#,4d,4f#,4b,b,c6,b,a,4g,4e,4d,4e,4a,4f#,2g"},
	}
	playSong(library[1])
	startServer()
}

func playSong(s song) {
	sound.Play(s.rtttl)
}

func startServer() {
	portPtr := flag.Int("p", 8081, "Port number to run the server on")
	flag.Parse()
	port := *portPtr
	mr := mux.NewRouter()
	apiRouter := mr.PathPrefix("/api").Subrouter()
	//Setup a static router for HTML/CSS/JS
	mr.PathPrefix("/client/").Handler(http.StripPrefix("/client/", http.FileServer(http.Dir("./resources")))) //test for directory traversal!
	//CRUD API routes for pastes
	pasteRouter := apiRouter.PathPrefix("/song").Subrouter()
	/*Play A Song*/ pasteRouter.HandleFunc("/{id}", songHandler).Methods("POST")
	fmt.Println("Listening for requests")
	http.ListenAndServe(fmt.Sprintf(":%v", port), mr)
}

func songHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Println("Playing song: " + id)
	idNo, err := strconv.Atoi(id)
	if err != nil || idNo > len(library) {
		w.WriteHeader(http.StatusNotFound)
	} else  {
		playSong(library[idNo])
		w.WriteHeader(http.StatusOK)
	}
	
}
