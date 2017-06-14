package main

import (
	"strings"
	"github.com/bwmarrin/discordgo"
	"fmt"
	"time"
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
	registerCommand("καλημέρα", "Το ευγενικό bot λέει και αυτό καλημέρα", false, "Καλημέρα %s :sunrise:", []string{"καλημερα", "καλημέρες", "καλημερες", "goodmorning"})
	registerCommand("καλησπέρα", "Το ευγενικό bot λέει και αυτό καλησπέρα", false, "Καλησπέρα %s :city_sunset:", []string{"καλησπερα", "καλησπέρες", "καλησπερες"})
	registerCommand("καληνύχτα", "Το ευγενικό bot λέει και αυτό καληνύχτα", false, "Καληνύχτα %s :last_quarter_moon_with_face: :night_with_stars:", []string{"καληνυχτα", "νυχτααα", "καληνύχτες", "goodnight"})
	registerCommand("!flip", "Όταν είσαι οργισμένος..", true, "(╯°□°）╯︵ ┻━┻", nil)
	registerCommand("!review", "Search στα reviews.. επικίνδυνο..", false, "", nil)
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
	case "καλημέρα":
		return assembleGmText(m, c)
	case "καλησπέρα":
		return assembleGEText(m, c)
	case "καληνύχτα":
		return assembleGnText(m, c)
	case "!review":
		return assembleReviewText(m, c)
	default:
		return ""
	}
}

func assembleGnText(m *discordgo.MessageCreate, c Command) string {
	h := time.Now().Hour()
	if h >= 5 && h <= 21 {
		// should I really say goodnight?
		c.ReplyText = "Μη με τρολλάρεις %s :sob:"
	}

	return simpleReply(m, c)
}

func assembleGmText(m *discordgo.MessageCreate, c Command) string {
	h := time.Now().Hour()
	if h >= 12 && h <= 23 {
		// not really morning..
		c.ReplyText = "Καλημέρα τέτοια ώρα;;!!"
	}

	return simpleReply(m, c)
}

func assembleGEText(m *discordgo.MessageCreate, c Command) string {
	h := time.Now().Hour()
	if h >= 0 && h <= 11 {
		// not really evening..
		c.ReplyText = "Καλημέρα λέει ο κόσμος.. τι timezone είσαι %s;"
	}

	return simpleReply(m, c)
}

func assembleHelpText() string {
	var help string
	help = "Ragetuitter Bot έκδοση: " + Version + "\n\n"
	for key, value := range commands {
		help = help + fmt.Sprintln("`"+key+"`:", value.Help)
	}

	help = help + fmt.Sprintf("Μπορώ να κάνω γενικά λίγα πράγματα προς το παρόν :sob:\n\nΜπορείτε να δείτε πώς δουλεύω στο: https://github.com/cmargonis/ragebot\n\nΤο spamming κλειδώνει το bot για %d δευτερόλεπτα", LockFor)
	return help
}

func assembleReviewText(m *discordgo.MessageCreate, c Command) string {
	rqPrefix := "http://ragequit.gr/reviews/item/"
	rqPostfix := "pc-review"
	originalMessage := strings.ToLower(m.Content)
	items := strings.Split(originalMessage, " ")
	urlReview := rqPrefix
	startReading := false
	for i := 0; i < len(items); i++ {
		// fast forward until after the command
		if !startReading && items[i] == "!review" {
			startReading = true
			continue
		}
		if startReading {
			urlReview = urlReview + items[i] + "-"
		}
	}
	urlReview = urlReview + rqPostfix
	c.ReplyText = urlReview
	return simpleReply(m, c)
}
