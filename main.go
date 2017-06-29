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
	Version   string
	OwnName   string
	Debug     bool
	DevSrvID  string
)

// TODO move stuff to proper config file
func init() {
	flag.StringVar(&TokenFile, "t", "", "Bot Token file")
	flag.Parse()
	// Read the token from supplied file
	Token = read(TokenFile)
	DevSrvID = read("config") // TODO sloppy..
	Version = "0.2.7"
	Debug = true

	if Debug {
		OwnName = "Ragequitter - Test"
	} else {
		OwnName = "Ragequitter"
	}

	// set error logging file
	if !Debug {
		logFile, _ := os.OpenFile("./ragebot.err", os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0755)
		syscall.Dup2(int(logFile.Fd()), 1)
		syscall.Dup2(int(logFile.Fd()), 2)
	}
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
	discord.AddHandler(messageReactionAdd)

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
	s.UpdateStatus(0, "Σας παρακολουθώ")
}

// This function will be called (since we register it as a Handler)
// every time a new message is created on any channel that the authenticated bot
// has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	c, err := s.Channel(m.ChannelID)
	if err != nil {
		return
	}
	if !isOnDevelopmentServer(c.GuildID) && Debug {
		// do not send messages when in development mode
		return
	}

	simpleReplyText := ParseCommand(m)

	if simpleReplyText != "" && ShouldSend(c.GuildID) {
		s.ChannelMessageSend(m.ChannelID, simpleReplyText)
	}
}

//case func(*Session, *Disconnect):
//return disconnectEventHandler(v)
func messageReactionAdd(s *discordgo.Session, d *discordgo.MessageReactionAdd) {
	if d.Emoji.Name == "🖕" {
		message, err := s.ChannelMessage(d.ChannelID, d.MessageID)
		if err != nil {
			return
		}
		if message.Author.ID != s.State.User.ID {
			// not for me!
			return
		}
		user, usrerr := s.User(d.UserID)
		// TODO refactor to shared code for channel cheking
		// TODO add functionality for private messages
		c, err := s.Channel(d.ChannelID)
		if err != nil || usrerr != nil {
			return
		}

		if !isOnDevelopmentServer(c.GuildID) && Debug {
			// do not send messages when in development mode
			return
		}
		if ShouldSend(c.GuildID) {
			msg := "Χμ όχι και πολύ ευγενικό αυτό " + user.Mention() + " ε; :poop:"
			s.ChannelMessageSend(d.ChannelID, msg)
		}
	}

}

// This function will be called every time a new guild is joined
func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	if event.Guild.Unavailable {
		return
	}

	if !isOnDevelopmentServer(event.Guild.ID) && Debug {
		// Avoid spamming other servers when testing
		return
	}
	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			_, _ = s.ChannelMessageSend(channel.ID, "Χαίρετε :smile:")
			s.GuildMemberNickname(event.Guild.ID, "@me", OwnName)
			return
		}
	}
}

func isOnDevelopmentServer(guildId string) bool {
	return DevSrvID == guildId
}
