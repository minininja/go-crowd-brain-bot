package main

import (
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

func discordErrorCheck(msg string, err error) {
	if err != nil {
		log.Fatalf("%s %s\n", msg, err)
		panic(err)
	}
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

var session *discordgo.Session

func Discord(token string, commandPrefix string, debug bool) {
	var err error

	session, err = discordgo.New("Bot " + token)
	discordErrorCheck("error creating discord session", err)

	router := exrouter.New()

	router.On("submit", func(ctx *exrouter.Context) {
		raw := ctx.Msg.Content
		log.Printf("%s %d", raw, len(raw))

		raw = raw[1:]
		log.Printf(raw)

		content := strings.Split(raw, " ")

		log.Printf("Upserting %s", content[1])
		id := InsertCategory(content[1])
		if -1 != id {
			if InsertContent(id, content[2], strings.Join(content[3:], " ")) {
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
		raw := ctx.Msg.Content[1:]
		log.Printf("raw: '%s'", raw)

		request := ctx.Msg.Content
		parts := strings.Split(strings.Split(request, "!")[1], " ")
		if len(parts) >= 1 {
			categoryId := FindCategory(parts[0])
			if len(parts) >= 2 {
				log.Printf("Searching for '%s' '%d' '%s'", parts[0], categoryId, parts[1])
				result := FindContent(categoryId, parts[1])
				log.Printf("length is '%d'", len(result))
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

	session.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		router.FindAndExecute(session, commandPrefix, session.State.User.ID, m.Message)
	})

	err = session.Open()
	discordErrorCheck("Error opening connection to Discord", err)
}
