package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"./sound"
	"github.com/gorilla/mux"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

type song struct {
	name  string
	rtttl string
}

//SongReply to be used when generating JSON replies for songs
type SongReply struct {
	Name  string `json:"SongName"`
	Index int    `json:"SongID"`
}

type command struct {
	pin         *gpio.DirectPinDriver //LED Pin or Nil, if nil signifies wait
	instruction byte                  //PWM level, or 255 for fully on, or wait duration in ms
	pwm         bool                  //Is it a PWM command?
}

//PatternReply to be used when generating JSON replies for patterns
type PatternReply struct {
	Name  string `json:"PatternName"`
	Index int    `json:"PatternID"`
}

type pattern struct {
	commands []command
	name     string
}

var (
	library = []song{
		{"Jingle Bells", "Jingle Bells:d=4,o=5,b=170:b,b,b,p,b,b,b,p,b,d6,g.,8a,2b.,8p,c6,c6,c6.,8c6,c6,b,b,8b,8b,b,a,a,b,2a,2d6"},
		{"Santa Claus is coming to town", "Santa Clause is Coming Tonight:d=4,o=5,b=180:g,8e,8f,g,g.,8g,8a,8b,c6,2c6,8e,8f,g,g,g,8a,8g,f,2f,e,g,c,e,d,2f,b4,1c,p,g,8e,8f,g,g.,8g,8a,8b,c6,2c6,8e,8f,g,g,g,8a,8g,f,f,e,g,c,e,d,2f,b4,1c,p,c6,d6,c6,b,c6,a,2a,c6,d6,c6,b,c6,2a.,d6,e6,d6,c#6,d6,b,b,b,8b,8c6,d6,c6,b,a,g,p,g.,8g,8e,8f,g,g.,8g,8a,8b,c6,2c6,8e,8f,g,g,g,8a,8g,8f,2f,e,g,c,e,d,2f,d6,1c6."},
		{"We Wish you a Merry Christmas", "We Wish you a Merry Christmas:d=8,o=5,b=160:4d,4g,g,a,g,f#,4e,4c,4e,4a,a,b,a,g,4f#,4d,4f#,4b,b,c6,b,a,4g,4e,4d,4e,4a,4f#,2g"},
		{"12 days of Christmas", "On the 12th Day of christmas:d=8,o=5,b=150:d,d,4g,g,g,4g,g,g,a,b,c6,a,4b.,p,4d6,a,b,c6,a,d6,d6,a,b,c6,a,4d6,4e6,4d.6,p,d6,c6,b,a,4g,a,b,4c6,4e,4e,4d,g,a,b,c6,4b,4a,2g."},
		{"Happy New Year", "Happy New Year:d=4,o=5,b=125:a5,d.,8d,d,f#,e.,8d,e,8f#,8e,d.,8d,f#,a,2b.,b,a.,8f#,f#,d,e.,8d,e,8f#,8e,d.,8b5,b5,a5,2d,16p"},
		{"Hark The Herald Angels", "Hark The Herald Angels:d=4,o=5,b=100:10a.,19p,40a.,36p,20a.,20p,26g.,5p,24e.,15p,6f.,40p,e.,96p,23e.,72f.,30p,240,g.,19p,c.,40p,3c.,17p,3d.,13p,1c."},
		{"Have Yourself A Merry Little Christmas", "Have Yourself A Merry Little Christmas:d=4,o=5,b=100:16b,8c,8e.,g.,16b,8c6,8g,f,16p,8e,16d,c,16c,d.,8p,16b,8c,8e,g.,16b,8c.6,2g,p,e,g,c6,16d6,8e.6,8d.6,8c6,8b,8a,g,16e,f,2e,16p,8f,16g,16d,16d,8c,8d,16a,c.,e,g,16b,8c.6,8g,f.,16p,8e,16d,8c,16d,16e,8d.,8e,8p,16c,8e,g.,16b,8c6,2g,p,e,g,c6,d6,8d.6,16c6,16p,8b,16a,g#.,8b.,16b,2c.6,16p"},
		{"I Saw Mommy Kissing Santa Claus", "I Saw Mommy Kissing Santa Claus:d=4,o=5,b=170:2c.,d,e,g,a,8c6,2b.,g,1e,a,g,e,c,a,g,e,8c.,1b.,f,e,d,d.,8c#.,2d.,8a.,16b.,8a.,g.,f#.,8a.,2g.,e,d,e,f#,g,a,g#,a,a#,b,a,f,8e,2d.,g,2c.,d,e,g,a,8c6,1b,g,1e,a,g,2e,16c,a,g,e,8c.,1a,16g#,8a.,16b,c6,8c.6,c6,8d.6,1b,f#,8a.,16a,g,f,8e,2d,8e.,2f,g,a,8c.6,a,2c6,2d6,1c6"},
		//{"I Wish It Could Be Christmas Everyday", "I Wish It Could Be Christmas Everyday:d=4,o=5,b=90:16g,a,16g,b.,c.7,b.,a,g.,4a.,16g,f,4a.,16g,1e,g,16g,f#.,g,f#.,16f#,e,16e,d.,e,4d,c.,e.,f#,2d,16d5,d5,d.5,16g,a,16g,b.,c.7,b.,a,g.,4a.,16g,f,4a.,16g,1e,g,16a,2b,16p,b.,2a,16p,b.,a.,g.,b.5,a.5,g.5"}, //BROKEN
		{"I ll Be Home For Christmas", "I ll Be Home For Christmas:d=4,o=5,b=160:2g.,f#,2a,g,2d,1d.,2e.,d,2f.,e,2a,e,2g.,f#,d,2a.,b,2d.,c,2b.,d,2g.,f#,2a,a,2f#.,f#,1e.,d.,2g.,f#,2a,g,2d,1d.,2e.,d,2f.,e,2c.,d,e,b,a,g,2a.,g,2a.,g,1d,2e.,e,2a,2b,2g,2a,2g."},
		{"Jingle Bell Rock", "Jingle Bell Rock:d=4,o=5,b=100:8a.,8c.6,16d.6,16d6,8c.6,8g#,16g#,16c6,16p,d.6,8p,8c6,16a,8b,16a,8g.,2c.6,8p,8c6,16c6,8c.6,16p,8b,16b,8b.,8a,16b,8a,e.,8p,8a,16b,8a,8e.,8g,16p,8a,16b,8a,f.,8p,8d,8e.,16f,16g,8p,8a,16g,8d,16e,16f,16p,2g"},
		{"Jingle Bells (Long)", "Jingle Bells:d=4,o=5,b=125:8g,8e6,8d6,8c6,2g,8g,8e6,8d6,8c6,2a,8a,8f6,8e6,8d6,8b,8g,8b,8d6,8g.6,16g6,8f6,8d6,2e6,8g,8e6,8d6,8c6,2g,16f#,8g,8e6,8d6,8c6,2a,8a,8f6,8e6,8d6,8g6,16g6,16f#6,16g6,16f#6,16g6,16g#6,8a.6,16g6,8e6,8d6,c6,g6,8e6,8e6,8e.6,16d#6,8e6,8e6,8e.6,16d#6,8e6,8g6,8c.6,16d6,2e6,8f6,8f6,8f.6,16f6,8f6,8e6,8e6,16e6,16e6,8e6,8d6,8d6,8e6,2d6"},
		{"Jolly Old St Nick", "Jolly Old St Nick:d=4,o=5,b=112:8d6,8d6,8d6,8d6,8c6,8c6,c6,8a#,8a#,8a#,8a#,2d6,8g,8g,8g,8g,8f,8f,a#,8a,8a#,8c6,8d6,2c6,8d6,8d6,8d6,8d6,8c6,8c6,c6,8a#,8a#,8a#,8a#,2d6,8g,8g,8g,8g,8f,8f,a#,8c6,8a#,8c6,8d6,a#,p,8e6,8e6,8e6,8e6,8d6,8d6,d6,8c6,8c6,8c6,8c6,2e6,8a,8a,8a,8a,8g,8g,c6,8b,8c6,8d6,8e6,2d6,8e6,8e6,8e6,8e6,8d6,8d6,d6,8c6,8c6,8c6,8c6,2e6,8a,8a,8a,8a,8g,8g,c6,8d6,8c6,8d6,8e6,c6,p"},
		{"Joy To The World", "Joy To The World:d=4,o=5,b=112:d6,8c#.6,16b,a.,8g,f#,e,d,8p,8a,b,8p,8b,c#6,8p,8c#6,2d.6,8p,8d6,8d6,8c#6,8b,8a,8a.,16g,8f#,8d6,8d6,8c#6,8b,8a,8a.,16g,8f#,8f#,8f#,8f#,8f#,16f#,16g,a.,16g,16f#,8e,8e,8e,16e,16f#,g,8p,16f#,16e,8d,8d6,8p,8b,8a.,16g,8f#,8g,f#,e,2d"},
		{"Last Christmas", "Last Christmas:d=4,o=5,b=112:16d6,e6,8p,e6,8d6,8p,8a,8e6,8e6,8f#6,d.6,8b,8b,8e6,8e6,f#6,d.6,8b,8c#6,8d6,8c#6,2b.,16e6,f#6,8p,e.6,8p,8b,8f#6,8g6,8f#6,2e6,8d6,8c#6,8d6,8c#6,c#6,8d6,8p,8c#6,8p,2a,16d6,e6,8p,e6,8d6,8p,8a,8e6,8e6,8f#6,d.6,8b,8b,8e6,8e6,f#6,d.6,8b,8c#6,8d6,8c#6,2b.,16e6,f#6,8p,e.6,8p,8b,8f#6,8g6,8f#6,2e6,8d6,8c#6,8d6,8c#6,c#6,8d6,8p,8c#6,8p,a"},
		{"God Rest Ye Gentleman", "God Rest Ye Gentleman:d=4,o=5,b=112:d.,d.,a.,a.,g.,f.,e.,d.,c.,d.,e.,f.,g.,2a,p,d.,d.,a.,a.,g.,f.,e.,d.,c.,d.,e.,f.,g.,2a,p,a.,a#.,g.,a.,a#.,c.7,d.7,a.,g.,f.,d.,e.,f.,4g.,f.,g.,4a.,a#.,a.,a.,g.,f.,e.,4d.,16f.,16e.,d.,4g.,f.,g.,a.,a#.,c.7,d.7,a.,g.,f.,e.,2d."},
		{"Let It Snow", "Let It Snow:d=4,o=5,b=125:8c,8c,8c6,8c6,a#,a,g,f,2c,8c,16c,g.,8f,g.,8f,e,2c,d,8d6,8d6,c6,a#,a,2g.,8e.6,16d6,c6,8c.6,16a#,a,8a#.,16a,2f.,c,8c6,8c6,a#,a,g,f,2c,8c.,16c,g.,8f,g.,8f,e,2c,d,8d6,8d6,c6,a#,a,2g.,8e.6,16d6,c6,8c.6,16a#,a,8a.,16g,2f. "},
		{"Little Drummer Boy", "Little Drummer Boy:d=4,o=5,b=140:2c,4p,d,p,e,10p,e,p,e,p,e,p,f,e,16f,5p,2e,4p,d,p,e,p,f,p,g,p,g,p,g,p,a,p,g,f,16e,5p,4d,5p,e,32p,d,e,10p,4c"},
		{"Mary's Boy Child", "Mary s Boy Child:d=4,o=5,b=140:8g,8g,g,32p,c6,c6,32p,a,8f,d,32p,8a,8a,g,8b,a,f,e,16p,8g,8g,32p,e6,d6,8c6,8p,a,8f,d,8p,8a,g,c6,8b,d6,c6"},
		{"O Christmas Tree", "O Christmas Tree:d=4,o=5,b=140:c,8f.,16f,f,g,8a.,16a,a.,8p,8a,8g,8a,a#,e,g,f"},
		{"O Come All Ye Faithful", "O Come All Ye Faithful:d=4,o=5,b=100:9g.,9p,3g.,12p,7d.,18p,7g.,17p,4a.,p,3d.,10p,b.,17p,9a.,15p,7b.,25p,c.6,11p,90b.,80a.,2p,p,7a.,14p,g.,12p,3g.,9p,7f#.,20p,e.,14p,f#.,18p,7g.,12p,7a.,16p,b.,15p,3f#.,9p,4e.,17p,23d.,17p,2d. "},
		{"O Little Town of Bethlehem", "O Little Town of Bethlehem:d=4,o=5,b=125:8d.,15p,9d.,10p,6d.,p,6c#.,p,6d.,p,6f.,p,6d#.,153,20p,22p,0p,12p,153,24p,21p,0p,24p,153,20p,21p,0p,12p,153,18p,23p,0p,32p,4p,8d.,15p,8d.,15p,8d.,15p,6g.,p,6f.,p,6f.,p,8d#.,15p,6g.5,153,18p,21p,0p,19p,p,16a.5,p,16a#.5,p,4"},
		{"Rudolph the Red Nose Reindeer", "Rudolph the Red Nose Reindeer:d=4,o=5,b=125:8g.,16a,8p,16g.,e,c6,a,2g.,8g.,16a,8g.,16a,g,c6,1b,8f.,16g,8p,16f.,d,b,a,2g.,8g.,16a,8g.,16a,g,a,1e,8g.,16a,8p,16g.,e,c6,a,1g,8g.,16a,8g.,16a,g,c6,1b,8f.,16g,8p,16f.,d,b,a,1g,8g.,16a,8g.,16a,g,d6,1c6,a,a,c6,a,g,8e.,2g,f,a,g,8f.,1e,d,e,g,a,b,8b.,1b,c6,8p,16c6,b,a,g,8f.,2d"},
		{"Silent Night", "Silent Night:d=4,o=5,b=90:a.,16b,a,4f#,p,a.,16b,a,4f#,p,e.6,16d#6,e6,4c#6,p,d.6,16c#6,d6,4a"},
		{"Silver Bells", "Silver Bells:d=4,o=5,b=140:9f.7,d.7,4c7,4a,9f.7,d.7,4c7,4a,9a.7,g.7,3f.7,f.,a#.,c.7,f.7,3f.7,12p,9e.7,f.7,4g7,4e.7,c.7,b.,4a#,4c.7,12p,a#.,4a#,3a.,12p,a.5,c.,f.,c.,f.,c.,f.,c.7,f.,9a.,g.,4f. "},
		{"Sleigh Bells", "Sleigh Bells:d=4,o=5,b=112:8d6,8d6,8d6,8d6,8e6,16d6,16b,8g,8a,8b,16a,16g,8e,2d,8p,16d,16e,16f#,16g,16a,16b,8d6,8e6,16d6,16b,16a,16g,8a,16a,16b,16a,16g,8e,2g,p,16d#,16c#,8b,16d#,16c#,8b,16d#,16c#,8b,8c#6,2a#,8b,16f#,16f#,8d#,2g#"},
		{"So This Is Christmas", "So This Is Christmas:d=4,o=5,b=56:a,b,c#6,a,e,4p,e,a,b,c#6,4b,5p,16f#,b,c#6,d6,c#6,4b,16e,16e,c#6,e6,16c#6,16b,4a"},
		{"The First Noel", "The First Noel:d=4,o=5,b=63:50f#.,153,50e.,153,16d.,153,50e.,153,50f#.,153,50g.,153,16a.,32p,50b.,153,50c#.6,153,d.6,153,c#.6,153,b.,153,16a.,32p,50b.,153,50c#.6,153,d.6,153,j#.6,153,b.,153,a.,153,b.,153,c#.6,153,d.6,1536"},
		{"The Holly and The Ivy", "The Holly and The Ivy:d=4,o=5,b=90:8c.,c.,c.,8c.,8a.,8g.,4e,32p,c.,c.,c.,8c.,8a.,4g.,g.,f.,e.,d.,8c.,e.,e.,a.6,a.6,8g.6,c.,d.,e.,f.,8e.,8d.,4c"},
		{"The Snowman", "The Snowman:d=4,o=5,b=90:a,d6,d6,c6,c6,1a.,a,d6,d6,c6,c6,4a.,f,2g.,g,a#,a#,a,a,4g.,d,f,f,e,e,2d."},
		{"We Three Kings", "We Three Kings:d=4,o=5,b=140:4g,f,4d#,c,d,d#,d,4c.,4g,f,4d#,c,d,d#,d,4c.,4d#,d#,4f,f,4g,g,a#,g#,g,f,g,f,4d#,d,2c."},
		{"We Wish You A Merry Christmas", "We Wish You A Merry Christmas:d=4,o=5,b=200:d,g,8g,8a,8g,8f#,e,e,e,a,8a,8b,8a,8g,f#,d,d,b,8b,8c6,8b,8a,g,e,d,e,a,f#,2g,d,g,8g,8a,8g,8f#,e,e,e,a,8a,8b,8a,8g,f#,d,d,b,8b,8c6,8b,8a,g,e,d,e,a,f#,1g,d,g,g,g,2f#,f#,g,f#,e,2d,a,b,8a,8a,8g,8g,d6,d,d,e,a,f#,2g"},
		{"While Shepherds Watched", "While Shepherds Watched:d=4,o=5,b=90:g.,4b,16b.,a.,g.,c.6,c.6,b.,a.,b.,d.6,d.6,c#.6,2d6,d.6,4e6,16d.6,c.6,b.,a.,g.,f#.,b.,a.,g.,g.,f#.,4g."},
		{"Winter Wonderland", "Winter Wonderland:d=4,o=5,b=140:8a#.,16a#,2a#.,8a#.,16a#,g,2a#,8a#.,16a#,2a#.,8a#.,16a#,g#,2a#,8p,16a#,8d.6,16d6,8d.6,c.6,8p,16c6,8a#.,16a#,8a#.,g#.,8p,16g#,8g.,16g,8g.,16g,8f.,16f,8f.,16f,2d#,p,8a#.,16a#,2a#.,8a#.,16a#,g,2a#,8a#.,16a#,2a#.,8a#.,16a#,g#,2a#,8p,16a#,8d.6,16d6,8d.6,c.6,8p,16c6,8a#.,16a#,8a#.,g#.,8p,16g#,8g.,16g,8g.,16g,8f.,16f,8f.,16f,2d#,p,8d.,16d,8b.,16b,8e.,16e,8c.6,16c6,b,2g,p,8d.,16d,8b.,16b,8e.,16e,8c.6,16c6,2b.,p"},
		// {"Monty Python", "Monty Python:d=8,o=5,b=180:d#6,d6,4c6,b,4a#,a,4g#,g,f,g,g#,4g,f,2a#,p,a#,g,p,g,g,f#,g,d#6,p,a#,a#,p,g,g#,p,g#,g#,p,a#,2c6,p,g#,f,p,f,f,e,f,d6,p,c6,c6,p,g#,g,p,g,g,p,g#,2a#,p,a#,g,p,g,g,f#,g,g6,p,d#6,d#6,p,a#,a,p,f6,f6,p,f6,2f6,p,d#6,4d6,f6,f6,e6,f6,4c6,f6,f6,e6,f6,a#,p,a,a#,p,a,2a#"}, //BROKEN
	}
	firmataAdaptor = firmata.NewAdaptor("COM3")
	commandChan    chan command
	led            = []*gpio.DirectPinDriver{
		gpio.NewDirectPinDriver(firmataAdaptor, "5"),  //Red channel 1
		gpio.NewDirectPinDriver(firmataAdaptor, "6"),  //Red channel 2
		gpio.NewDirectPinDriver(firmataAdaptor, "10"), //Green channel 1
		gpio.NewDirectPinDriver(firmataAdaptor, "11"), //Green channel 2
	}
	patterns = []pattern{
		pattern{
			name: "PyGen",
			commands: []command{
				command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[2], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[3], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(235)}, command{pin: led[0], pwm: true, instruction: byte(127)}, command{pin: nil, pwm: true, instruction: byte(60)}, command{pin: led[0], pwm: true, instruction: byte(50)}, command{pin: nil, pwm: true, instruction: byte(60)}, command{pin: led[0], pwm: true, instruction: byte(25)}, command{pin: nil, pwm: true, instruction: byte(60)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(127)}, command{pin: nil, pwm: true, instruction: byte(60)}, command{pin: led[1], pwm: true, instruction: byte(50)}, command{pin: nil, pwm: true, instruction: byte(60)}, command{pin: led[1], pwm: true, instruction: byte(25)}, command{pin: nil, pwm: true, instruction: byte(60)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[2], pwm: true, instruction: byte(127)}, command{pin: nil, pwm: true, instruction: byte(60)}, command{pin: led[2], pwm: true, instruction: byte(50)}, command{pin: nil, pwm: true, instruction: byte(60)}, command{pin: led[2], pwm: true, instruction: byte(25)}, command{pin: nil, pwm: true, instruction: byte(60)}, command{pin: led[2], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[3], pwm: true, instruction: byte(127)}, command{pin: nil, pwm: true, instruction: byte(60)}, command{pin: led[3], pwm: true, instruction: byte(50)}, command{pin: nil, pwm: true, instruction: byte(60)}, command{pin: led[3], pwm: true, instruction: byte(25)}, command{pin: nil, pwm: true, instruction: byte(60)}, command{pin: led[3], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)},
			},
		},
		pattern{
			name: "One pulse",
			commands: []command{
				command{pin: led[0], pwm: true, instruction: byte(12)}, command{pin: nil, pwm: true, instruction: byte(50)}, command{pin: led[0], pwm: true, instruction: byte(25)}, command{pin: nil, pwm: true, instruction: byte(50)}, command{pin: led[0], pwm: true, instruction: byte(50)}, command{pin: led[1], pwm: true, instruction: byte(12)}, command{pin: nil, pwm: true, instruction: byte(50)}, command{pin: led[0], pwm: true, instruction: byte(100)}, command{pin: led[1], pwm: true, instruction: byte(25)}, command{pin: nil, pwm: true, instruction: byte(50)}, command{pin: led[0], pwm: true, instruction: byte(200)}, command{pin: led[1], pwm: true, instruction: byte(50)}, command{pin: led[2], pwm: true, instruction: byte(12)}, command{pin: nil, pwm: true, instruction: byte(50)}, command{pin: led[1], pwm: true, instruction: byte(100)}, command{pin: led[2], pwm: true, instruction: byte(25)}, command{pin: nil, pwm: true, instruction: byte(50)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[1], pwm: true, instruction: byte(200)}, command{pin: led[2], pwm: true, instruction: byte(50)}, command{pin: led[3], pwm: true, instruction: byte(12)}, command{pin: nil, pwm: true, instruction: byte(50)}, command{pin: led[2], pwm: true, instruction: byte(100)}, command{pin: led[3], pwm: true, instruction: byte(25)}, command{pin: nil, pwm: true, instruction: byte(50)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: led[2], pwm: true, instruction: byte(200)}, command{pin: led[3], pwm: true, instruction: byte(50)}, command{pin: nil, pwm: true, instruction: byte(50)}, command{pin: led[3], pwm: true, instruction: byte(100)}, command{pin: nil, pwm: true, instruction: byte(50)}, command{pin: led[3], pwm: true, instruction: byte(200)}, command{pin: led[2], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[3], pwm: true, instruction: byte(0)},
			},
		},
		pattern{
			name: "Blink pairs",
			commands: []command{
				command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: led[2], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: led[3], pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[2], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: led[3], pwm: true, instruction: byte(0)}, command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: led[2], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: led[3], pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[2], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: led[3], pwm: true, instruction: byte(0)}, command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: led[2], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: led[3], pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[2], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: led[3], pwm: true, instruction: byte(0)}, command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: led[2], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: led[3], pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[2], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: led[3], pwm: true, instruction: byte(0)}, command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: led[2], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: led[3], pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[2], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: led[3], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)},
			},
		},
		pattern{
			name: "Blink pairs long",
			commands: []command{
				command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: led[2], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: led[3], pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[2], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: led[3], pwm: true, instruction: byte(0)}, command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: led[2], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: led[3], pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[2], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: led[3], pwm: true, instruction: byte(0)}, command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: led[2], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: led[3], pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[2], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: led[3], pwm: true, instruction: byte(0)}, command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: led[2], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: led[3], pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[2], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: led[3], pwm: true, instruction: byte(0)}, command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: led[2], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: led[3], pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[2], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: led[3], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: led[2], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: led[3], pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[2], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: led[3], pwm: true, instruction: byte(0)}, command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: led[2], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: led[3], pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[2], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: led[3], pwm: true, instruction: byte(0)}, command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: led[2], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: led[3], pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[2], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: led[3], pwm: true, instruction: byte(0)}, command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: led[2], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: led[3], pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[2], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: led[3], pwm: true, instruction: byte(0)}, command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: led[2], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: led[3], pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[2], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: led[3], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: led[2], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: led[3], pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[2], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: led[3], pwm: true, instruction: byte(0)}, command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: led[2], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: led[3], pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[2], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: led[3], pwm: true, instruction: byte(0)}, command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: led[2], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: led[3], pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[2], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: led[3], pwm: true, instruction: byte(0)}, command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: led[2], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: led[3], pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[2], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: led[3], pwm: true, instruction: byte(0)}, command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: led[2], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: led[3], pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[2], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: led[3], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)},
			},
		},
		pattern{
			name: "Crossover",
			commands: []command{
				command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: led[2], pwm: true, instruction: byte(255)}, command{pin: led[3], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(50)}, command{pin: led[0], pwm: true, instruction: byte(220)}, command{pin: led[1], pwm: true, instruction: byte(35)}, command{pin: led[2], pwm: true, instruction: byte(220)}, command{pin: led[3], pwm: true, instruction: byte(35)}, command{pin: nil, pwm: true, instruction: byte(50)}, command{pin: led[0], pwm: true, instruction: byte(200)}, command{pin: led[1], pwm: true, instruction: byte(55)}, command{pin: led[2], pwm: true, instruction: byte(200)}, command{pin: led[3], pwm: true, instruction: byte(55)}, command{pin: nil, pwm: true, instruction: byte(50)}, command{pin: led[0], pwm: true, instruction: byte(180)}, command{pin: led[1], pwm: true, instruction: byte(75)}, command{pin: led[2], pwm: true, instruction: byte(180)}, command{pin: led[3], pwm: true, instruction: byte(75)}, command{pin: nil, pwm: true, instruction: byte(50)}, command{pin: led[0], pwm: true, instruction: byte(150)}, command{pin: led[1], pwm: true, instruction: byte(105)}, command{pin: led[2], pwm: true, instruction: byte(150)}, command{pin: led[3], pwm: true, instruction: byte(105)}, command{pin: nil, pwm: true, instruction: byte(50)}, command{pin: led[0], pwm: true, instruction: byte(120)}, command{pin: led[1], pwm: true, instruction: byte(135)}, command{pin: led[2], pwm: true, instruction: byte(120)}, command{pin: led[3], pwm: true, instruction: byte(135)}, command{pin: nil, pwm: true, instruction: byte(50)}, command{pin: led[0], pwm: true, instruction: byte(100)}, command{pin: led[1], pwm: true, instruction: byte(155)}, command{pin: led[2], pwm: true, instruction: byte(100)}, command{pin: led[3], pwm: true, instruction: byte(155)}, command{pin: nil, pwm: true, instruction: byte(50)}, command{pin: led[0], pwm: true, instruction: byte(80)}, command{pin: led[1], pwm: true, instruction: byte(175)}, command{pin: led[2], pwm: true, instruction: byte(80)}, command{pin: led[3], pwm: true, instruction: byte(175)}, command{pin: nil, pwm: true, instruction: byte(50)}, command{pin: led[0], pwm: true, instruction: byte(60)}, command{pin: led[1], pwm: true, instruction: byte(195)}, command{pin: led[2], pwm: true, instruction: byte(60)}, command{pin: led[3], pwm: true, instruction: byte(195)}, command{pin: nil, pwm: true, instruction: byte(50)}, command{pin: led[0], pwm: true, instruction: byte(30)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: led[2], pwm: true, instruction: byte(30)}, command{pin: led[3], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(60)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: led[2], pwm: true, instruction: byte(0)}, command{pin: led[3], pwm: true, instruction: byte(0)},
			},
		},
		pattern{
			name: "Binary",
			commands: []command{
				command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: led[2], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: led[2], pwm: true, instruction: byte(0)}, command{pin: led[3], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: led[2], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: led[2], pwm: true, instruction: byte(0)}, command{pin: led[3], pwm: true, instruction: byte(0)},
			},
		},
		pattern{
			name: "Sequence",
			commands: []command{
				command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: led[2], pwm: true, instruction: byte(0)}, command{pin: led[3], pwm: true, instruction: byte(0)}, command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[2], pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[3], pwm: true, instruction: byte(255)}, command{pin: led[2], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[2], pwm: true, instruction: byte(255)}, command{pin: led[3], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(255)}, command{pin: led[2], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(255)}, command{pin: led[1], pwm: true, instruction: byte(0)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: nil, pwm: true, instruction: byte(255)}, command{pin: led[0], pwm: true, instruction: byte(0)},
			},
		},
	}
)

func main() {
	commandChan = make(chan command, 20)
	playSong(song{"startup", "startup:d=4,o=4,b=200:c,p,c"})
	go startServer()
	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		work,
	)

	robot.Start()
}

func work() {
	gobot.Every(time.Millisecond*10, func() {
		c := <-commandChan
		runCommand(c)
	})
}
func playSong(s song) string {
	return (sound.Play(s.rtttl))
}

func runCommand(c command) {
	if c.pin != nil {
		if c.pwm {
			c.pin.PwmWrite(c.instruction)
		} else {
			c.pin.DigitalWrite(c.instruction)
		}
	} else {
		time.Sleep(time.Millisecond * time.Duration(int(c.instruction)))
	}
}

func startServer() {
	portPtr := flag.Int("p", 8081, "Port number to run the server on")
	flag.Parse()
	port := *portPtr
	mr := mux.NewRouter()
	apiRouter := mr.PathPrefix("/api").Subrouter()
	//Setup a static router for HTML/CSS/JS
	mr.PathPrefix("/client/").Handler(http.StripPrefix("/client/", http.FileServer(http.Dir("./resources")))) //test for directory traversal!
	//CRUD API routes for songs
	songRouter := apiRouter.PathPrefix("/song").Subrouter()
	/*Play A Song   */ songRouter.HandleFunc("/{id}", songHandler).Methods("POST")
	/*List all songs*/ songRouter.HandleFunc("/list", songLister).Methods("GET")
	//API routes for lights
	lightRouter := apiRouter.PathPrefix("/lights").Subrouter()
	/*Play A Pattern*/ lightRouter.HandleFunc("/{id}", lightHandler).Methods("POST")
	/*List patterns */ lightRouter.HandleFunc("/list", patternLister).Methods("GET")
	log.Println("Listening for requests")
	http.ListenAndServe(fmt.Sprintf(":%v", port), mr)
}
func lightHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	idNo, err := strconv.Atoi(id)
	if err != nil || idNo >= len(patterns) {
		log.Println("Failed to run pattern: " + id)
		w.WriteHeader(http.StatusNotFound)
	} else {
		log.Println("Running pattern: " + id)
		w.WriteHeader(http.StatusOK)
		if len(commandChan) < 3 {
			for _, c := range patterns[idNo].commands {
				commandChan <- c
			}
		}
	}
}
func patternLister(w http.ResponseWriter, r *http.Request) {
	var tempPatterns []PatternReply
	w.Header().Set("Content-Type", "application/json")
	for index, p := range patterns {
		tempPatterns = append(tempPatterns, PatternReply{Index: index, Name: p.name})
	}
	json.NewEncoder(w).Encode(tempPatterns)
}
func songLister(w http.ResponseWriter, r *http.Request) {
	var songs []SongReply
	w.Header().Set("Content-Type", "application/json")
	for index, s := range library {
		songs = append(songs, SongReply{s.name, index})
	}
	json.NewEncoder(w).Encode(songs)
}

func songHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	idNo, err := strconv.Atoi(id)
	if err != nil || idNo >= len(library) {
		log.Println("Failed to play song: " + id)
		w.WriteHeader(http.StatusNotFound)
	} else {
		log.Println("Playing song: " + id)
		log.Println(playSong(library[idNo]))
		w.WriteHeader(http.StatusOK)
	}
}
