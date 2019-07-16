package main

import (
	"flag"
	"log"
	"os"
)

var discordToken string
var commandPrefix string
var debug = false

func init() {
	var token string
	flag.StringVar(&token, "t", "", "Discord Auth Token")
	flag.StringVar(&commandPrefix, "cp", "!", "Discord command prefix")
	flag.BoolVar(&debug, "debug", false, "Enable debug message logger mode")

	flag.Parse()

	// fall back to environment variables
	if token == "" {
		log.Printf("Looking in environment for token")
		token = os.Getenv("DG_TOKEN")
	}
	if commandPrefix == "" {
		log.Printf("Looking in environment for command prefix")
		commandPrefix = os.Getenv("DG_COMMAND_PREFIX")
	}

	log.Printf("Using %s as command prefix", commandPrefix)
	if debug {
		log.Printf("Message logging enabled")
	}

	discordToken = token
	if discordToken == "" {
		log.Fatal("A discord token must be provided")
		return
	}
}

func errCheck(msg string, err error) {
	if err != nil {
		log.Fatalf("%s %s\n", msg, err)
		panic(err)
	}
}

func main() {
	// start the discord side
	log.Printf("Starting discord client")
	Discord(discordToken, commandPrefix, debug)

	// start the rest side
	log.Printf("Starting rest server")
	Rest()

	//<-make(chan struct{})
}
