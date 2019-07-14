package main

import (
	"flag"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"strings"
)

var Session, _ = discordgo.New()
var commandPrefix string
var debug = true

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

	Session.Token = token
	if Session.Token == "" {
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
	var err error

	Session, err = discordgo.New("Bot " + Session.Token)
	errCheck("error creating discord session", err)

	router := exrouter.New()

	router.On("submit", func(ctx *exrouter.Context) {
		content := strings.Split(ctx.Msg.Content, "!")
		log.Printf("Upserting %s", content[0])
		id := UpsertCategory(content[0])
		if -1 != id {
			if InsertContent(id, content[1], content[2]) {
				ctx.Reply("Accepted as pending.")
			} else {
				ctx.Reply("Some sort of error happened, we recommend screaming and hollering.")
			}
		} else {
			ctx.Reply("Couldn't find the category, sorry.")
		}
	})
	router.On("remove", func(ctx *exrouter.Context) {
		ctx.Reply("remove")
	})
	router.On("pending", func(ctx *exrouter.Context) {
		ctx.Reply("pending")
	})
	router.On("accept", func(ctx *exrouter.Context) {
		ctx.Reply("accept")
	})
	router.On("reject", func(ctx *exrouter.Context) {
		ctx.Reply("reject")
	})
	router.On("categories", func(ctx *exrouter.Context) {
		ctx.Reply("categories")
	})
	router.On("help", func(ctx *exrouter.Context) {
		ctx.Reply("Recognized commands: content, remove, pending, accept, reject, categories and any defined category")
	})
	router.OnMatch("test", matcher, func(ctx *exrouter.Context) {
		request := ctx.Msg.Content
		parts := strings.Split(strings.Split(request, "!")[1], " ")
		if len(parts) >= 1 {
			categoryId := FindCategory(parts[0])
			if len(parts) >= 2 {
				log.Printf("Searching for %s %d %s", parts[0], categoryId, parts[1])
				result := FindContent(categoryId, parts[1])
				log.Printf("length is %d", len(result))
				if len(result) != 0 {
					for _, txt := range result {
						ctx.Reply(txt)
					}
				} else {
					ctx.Reply("Sorry, I couldn't find anything")
				}
			}
		}
	})

	Session.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		router.FindAndExecute(Session, commandPrefix, Session.State.User.ID, m.Message)
	})

	err = Session.Open()
	errCheck("Error opening connection to Discord", err)

	log.Println("Bot is now running")
	<-make(chan struct{})
}

func matcher(param string) bool {
	return -1 != FindCategory(param)
}

func messageLogger(session *discordgo.Session, message *discordgo.MessageCreate) {
	if debug {
		// no need to log our own messages
		if session.State.User.ID == message.Author.ID {
			return
		}

		log.Printf("%s %s %s %s\n", message.GuildID, message.ChannelID, message.Author.Username, message.Content)
	}
}
