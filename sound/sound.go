package sound

import (
	"github.com/tarm/serial"
	"log"
)
var (
	c = &serial.Config{Name: "COM10", Baud: 9600}
	s, err = serial.OpenPort(c)
)

//Play a christmas song from an RTTTL string
func Play(sound string) { 
	_, err := s.Write([]byte(sound))
	if err != nil {
		//log.Fatal(err)
		log.Println(err)
		s, err = serial.OpenPort(c)
		if err != nil {
			log.Println("Failed to reconnect: "+err.Error())
		}
	}
	
}