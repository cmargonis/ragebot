package main

import (
	"time"
	_"fmt"
)

type SendEvent struct {
	messagesSent int16
	lastSent     int64
}

// Stores the latest send event for each guild
var sendevents map[string]*SendEvent
var eventsPerWindow = int16(2)
var LockFor = int64(30) // prevent messages for 30 seconds
var windowDuration = int64(5) // only eventsPerWindow allowed for 10 seconds

func init() {
	sendevents = make(map[string]*SendEvent)
}

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
