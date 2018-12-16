package sound

import (
	"log"

	"github.com/tarm/serial"
)

var (
	c      = &serial.Config{Name: "COM12", Baud: 9600}
	s, err = serial.OpenPort(c)
)

//Play a christmas song from an RTTTL string
func Play(sound string) string {
	//var reply string
	if s != nil { //if a previous reconnect attempt failed, it may set s to nil
		_, err := s.Write([]byte(sound))
		if err != nil {
			//log.Fatal(err)
			log.Println(err)
			s, err = serial.OpenPort(c)
			if err != nil {
				log.Println("Failed to reconnect: " + err.Error())
				return "err"
			}
			//if it reconnects
			log.Println("Reconnected sound device")
			return Play(sound) //try again
		}
		buf := make([]byte, 1)
		_, err = s.Read(buf)
		if err != nil {
			log.Println(err)
			return "err"
		}
		//log.Printf("Recieved: %q %T", buf, buf)
		if buf[0] == byte('a') {
			log.Printf("Finished")
			return "done"
		}
		return "err"
	} 
	//really, without making it set up the serial device again, it's not recoverable
	log.Println("Sound device disconnected. This is not recoverable.")
	return "err"

}
