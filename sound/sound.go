package sound

import (
	"log"

	"github.com/tarm/serial"
)

var (
	c      = &serial.Config{Name: "COM10", Baud: 9600}
	s, err = serial.OpenPort(c)
)

//Play a christmas song from an RTTTL string
func Play(sound string) {
	if s != nil { //if a previous reconnect attempt failed, it may set s to nil
		_, err := s.Write([]byte(sound))
		if err != nil {
			//log.Fatal(err)
			log.Println(err)
			s, err = serial.OpenPort(c)
			if err != nil {
				log.Println("Failed to reconnect: " + err.Error())
			} else { //If it reconnects
				log.Println("Reconnected sound device")
				Play(sound) //try again
			}
		}
	} else { //really, without making it set up the serial device again, it's not recoverable
		log.Println("Sound device disconnected. This is not recoverable.")
	}

}
//Setup a serial port with the specified portname
func Setup(portname string) {
	log.Println("Opening sound device: "+portname)
	c = &serial.Config{Name: portname, Baud: 9600}
}