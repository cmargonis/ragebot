package main

import (
	"time"
	"fmt"
)

type SendEvent struct {
	messagesSent int16
	lastSent     int64
}

type UserGreeting struct {
	greetings []*Greeting
}

type Greeting struct {
	greetingType GreetingType
	sent         int64
}

type GreetingType uint8

const (
	GOOD_MORNING GreetingType = iota + 1
	GOOD_EVENING
	GOOD_NIGHT
)

// Stores the last sent event for each guild
var sendevents map[string]*SendEvent
// Stores the greeting used for a user, key is the userid
var userGreetings map[string]*UserGreeting
var eventsPerWindow = int16(2)
var LockFor = int64(30)       // prevent messages for 30 seconds
var windowDuration = int64(5) // only eventsPerWindow allowed for 10 seconds

func init() {
	sendevents = make(map[string]*SendEvent)
	userGreetings = make(map[string]*UserGreeting)
}

// Checks whether it's appropriate to send or not to avoid spamming
func ShouldSend(guildId string) bool {
	lastEvent, ok := sendevents[guildId]

	if !ok {
		// It hasn't sent anything to this guild, so you can send
		sendevents[guildId] = &SendEvent{1, time.Now().Unix()}
		return true
	}

	if lastEvent.messagesSent < eventsPerWindow {
		if (time.Now().Unix() - lastEvent.lastSent) < windowDuration {
			// last sent message was less than lockFor seconds ago
			// increase the counter..
			// otherwise don't
			lastEvent.messagesSent++
		}
		lastEvent.lastSent = time.Now().Unix()
		return true
	}

	// Messages sent equal eventsPerWindow
	// lock it for lockFor milliseconds
	if (time.Now().Unix() - lastEvent.lastSent) > LockFor {
		// lockFor milliseconds have passed since last transmission
		// allow this transmission and set the counter back to 1
		lastEvent.messagesSent = 1
		lastEvent.lastSent = time.Now().Unix()
		return true
	}

	// repeated attempts will renew the lock
	sendevents[guildId].lastSent = time.Now().Unix()
	return false
}

func ShouldSendGreeting(userId string, greetingType GreetingType) bool {
	usergreeting, ok := userGreetings[userId]

	if !ok {
		// User has never interacted with a bot-greeting, let's do it!
		return true
	}

	// the greetings this user has used..
	for _, greeting := range usergreeting.greetings {
		fmt.Printf("iterating greetings %+v\n", greeting)
		if greeting.greetingType == greetingType {
			lastSent := time.Unix(greeting.sent, 0)
			// check if for today the user has interacted with the same greeting...
			if time.Now().Day() == lastSent.Day() && time.Now().Month() == lastSent.Month() {
				// there is a greeting for the given date
				fmt.Println("A greeting already exists")
				return false
			} else {
				return true
			}
		}
	}
	// no specific greeting interaction has been used
	return true
}

func StoreGreeting(userId string, greetingType GreetingType) {
	usergreeting, ok := userGreetings[userId]

	if !ok {
		fmt.Println("No greeting for user")
		greeting := &Greeting{greetingType, time.Now().Unix()}
		greetings := make([]*Greeting, 5)
		greetings[0] = greeting
		userGreetings[userId] = &UserGreeting{greetings}
		fmt.Printf("greetings now: %+v\n", userGreetings)
		return
	}

	for _, greeting := range usergreeting.greetings {
		if greeting.greetingType == greetingType {
			greeting.sent = time.Now().Unix()
			fmt.Printf("Updated greeting: %+v\n", userGreetings)
			return
		}
	}
	fmt.Printf("Greetings before appending: %+v\n", userGreetings)
	usergreeting.greetings = append(usergreeting.greetings, &Greeting{greetingType, time.Now().Unix()})
	fmt.Printf("Inserted new greeting: %+v\n", userGreetings)
}
