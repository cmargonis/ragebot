package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

// Variables used for command line parameters
var (
	TokenFile string
	Token     string
)

func init() {
	flag.StringVar(&TokenFile, "t", "", "Bot Token file")
	flag.Parse()
	// Read the token from supplied file
	Token = read(TokenFile)
}

func main() {
	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error creating the Bot.")
		return
	}

	discord.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening
	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening connection, ", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the discord session.
	discord.Close()
}

// This function will be called (since we register it as a Handler)
// every time a new message is created on any channel that the authenticated bot
// has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("received message: ", m.Content)

	// Ignore all messages created by the bot itself
	// This isn't required but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
