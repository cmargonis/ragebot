package main

import (
	"strings"
	"github.com/bwmarrin/discordgo"
	"fmt"
)

type Command struct {
	Help        string   // Text used as instructions to the users
	SimpleReply bool     // true if only a simple text is to be send back
	ReplyText   string   // the text to be send back when is simple reply, "" otherwise (?)
	Aliases     []string // different triggers for this command
}

var commands map[string]Command

func init() {
	commands = make(map[string]Command)

	registerCommand("!help", "Αυτό το μήνυμα με οδηγίες", false, "", []string{"!βοήθεια", "!βοηθεια"})
	registerCommand("!ping", "Απαντά με πόνγκ!", true, "Pong!", nil)
	registerCommand("καλημέρα", "Το ευγενικό bot λέει και αυτό καλημέρα", true, "Καλημέρα %s :sunrise:", []string{"καλημερα", "καλημέρες", "καλημερες"})
	registerCommand("καλησπέρα", "Το ευγενικό bot λέει και αυτό καλησπέρα", true, "Καλησπέρα %s :city_sunset:", []string{"καλησπερα", "καλησπέρες", "καλησπερες"})
	registerCommand("!flip", "Όταν είσαι οργισμένος..", true, "(╯°□°）╯︵ ┻━┻", nil)
}

// Used apparently to register a command.
// When isSimpleReply is true, then the replyText is just being passed to the dispatcher
// Otherwise custom logic determines the result -if any-
func registerCommand(operator string, help string, isSimpleReply bool, replyText string, aliases []string) {
	_, ok := commands[operator]

	if !ok {
		a := make([]string, len(aliases))
		a = aliases
		commands[operator] = Command{help, isSimpleReply, replyText, a}
	}
}

func ParseCommand(m *discordgo.MessageCreate) (string) {
	message := strings.ToLower(m.Content)

	for key, value := range commands {
		if (strings.Contains(message, key) || checkAliases(message, &value) ) && value.SimpleReply {
			return simpleReply(m, value)
		} else if (strings.Contains(message, key) || checkAliases(message, &value)) && !value.SimpleReply {
			return complexReply(m, value, key)
		}
	}
	return ""
}

// Checks if the supplied command is an alias for
// a registered command
func checkAliases(message string, cmnd *Command) bool {
	if cmnd.Aliases == nil {
		return false
	}
	for _, alias := range cmnd.Aliases {
		if strings.Contains(message, alias) {
			return true
		}
	}
	return false
}

func simpleReply(m *discordgo.MessageCreate, c Command) (string) {
	addUser := strings.Contains(c.ReplyText, "%s")
	if addUser {
		return fmt.Sprintf(c.ReplyText, m.Author.Username)
	}
	return c.ReplyText
}

func complexReply(m *discordgo.MessageCreate, c Command, operator string) (string) {
	switch operator {
	case "!help":
		return assembleHelpText()
	default:
		return ""
	}
}

func assembleHelpText() string {
	var help string
	help = "Ragetuitter Bot έκδοση: " + Version + "\n\n"
	for key, value := range commands {
		help = help + fmt.Sprintln("`"+key+"`:", value.Help)
	}

	help = help + fmt.Sprintln("Μπορώ να κάνω γενικά λίγα πράγματα προς το παρόν :sob:")
	return help
}
