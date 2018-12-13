package sound

import (
	"fmt"
	"github.com/tarm/serial"
	"log"
)
var (
	c = &serial.Config{Name: "COM12", Baud: 9600}
	s, err = serial.OpenPort(c)
)

// func main() {
// 	if err != nil {
// 			log.Fatal(err)
// 	}
	
// 	n, err := s.Write([]byte("Jingle Bells:d=4,o=5,b=170:b,b,b,p,b,b,b,p,b,d6,g.,8a,2b.,8p,c6,c6,c6.,8c6,c6,b,b,8b,8b,b,a,a,b,2a,2d6"))
// 	if err != nil {
// 			log.Fatal(err)
// 	}
// 	fmt.Println(n)

// }
//Play plays a christmas song
func Play(sound string) { 
	n, err := s.Write([]byte(sound))
	if err != nil {
			log.Fatal(err)
	}
	fmt.Println(n)
}