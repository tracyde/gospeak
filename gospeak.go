package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type SpeechHandler struct {
	sChan chan string
}

func (h *SpeechHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.sChan <- fmt.Sprintf("%s", r.FormValue("s"))
}

func SpeechServer(c string, sChan chan string) {
	for {
		s := <-sChan
		log.Printf("Command: %s -- Saying: %s\n", c, s)

		cmd := exec.Command(c, s)
		_, err := cmd.Output()
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func testCommand(c string) {
	if c == "" {
		log.Fatalln("Command flag can not be empty")
	}
	if _, err := os.Stat(c); os.IsNotExist(err) {
		log.Fatalln("Command does not exist")
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "gospeak"
	app.Usage = "provides a web interface to the systems text-to-speech system"

	app.Flags = []cli.Flag{
		cli.IntFlag{"port, p", 8080, "port gospeak listens on"},
		cli.StringFlag{"command, c", "", "command to use for speech"},
	}

	app.Action = func(c *cli.Context) {
		command := c.String("command")
		testCommand(command)
		port := fmt.Sprintf(":%d", c.Int("port"))
		var sChan = make(chan string)
		go SpeechServer(command, sChan)
		http.Handle("/speak", &SpeechHandler{sChan})
		log.Fatal(http.ListenAndServe(port, nil))
	}

	app.Run(os.Args)
}
