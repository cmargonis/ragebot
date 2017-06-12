package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
	"strings"
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

var buffer = make([][]byte, 0)

func main() {
	err := loadSound(&buffer)
	if err != nil {
		fmt.Println("Error loading media: ", err)
		return
	}
	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error creating the Bot.")
		return
	}

	// Register ready as a callback for the ready events.
	discord.AddHandler(ready)

	discord.AddHandler(messageCreate)

	// Register guildCreate as a callback for the guildCreate events.
	discord.AddHandler(guildCreate)

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


// This function will be called when the bot receives
// the "ready" event from Discord.
func ready(s *discordgo.Session, event *discordgo.Ready) {
	// Set the playing status
	s.UpdateStatus(0, "!airhorn")

}
// This function will be called (since we register it as a Handler)
// every time a new message is created on any channel that the authenticated bot
// has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("received message: ", m.Content)

	message := strings.ToLower(m.Content)

	// Ignore all messages created by the bot itself
	// This isn't required but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(message, "καλημέρα") {
		s.ChannelMessageSend(m.ChannelID, "Καλημέρα " + m.Author.Mention() + " :sunny:")
	}

	if strings.Contains(message, "help") {
		s.ChannelMessageSend(m.ChannelID, "Μπορώ να κάνω γενικά λίγα πράγματα προς το παρόν :sob:")
	}

	if strings.HasPrefix(message, "!ping") {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if strings.HasPrefix(message, "!pong") {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if strings.HasPrefix(message, "!airhorn") {
		// Find the channel that the message came from.
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			// Could not find channel
			return
		}

		// Find the guild for that channel.
		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			// Could not find guild
			return
		}

		// Look for the message sender in that guild's current voice states.
		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				err = playSound(s, g.ID, vs.ChannelID, buffer)
				if err != nil {
					fmt.Println("Error playing sound: ", err)
				}
				return
			}
		}
	}
}

// This function will be called every time a new guild is joined
func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	if event.Guild.Unavailable {
		return
	}

	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			_, _ = s.ChannelMessageSend(channel.ID, "Airhorn is read! Type !airhorn while in a voice channel to play a sound!")
			return
		}
	}
}