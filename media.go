package main

import (
	"os"
	"fmt"
	"encoding/binary"
	"io"
	"github.com/bwmarrin/discordgo"
	"time"
)

func loadSound(buffer *[][]byte) (error) {
	file, err := os.Open("assets/airhorn.dca")
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return err
	}

	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file: ", err)
			return err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any EOF errors
		if err != nil {
			fmt.Println("Error reading from dca file: ", err)
			return err
		}

		// Append encoded pcm data to the buffer.
		*buffer = append(*buffer, InBuf)
	}
}

// plays the current buffer to the provided channel.
func playSound(s *discordgo.Session, guildID, channelID string, buffer [][]byte) (err error) {
	// Join the provided channel
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false,true)
	if err != nil {
		return err
	}

	// Sleep for a specified ammount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	// Start speaking
	vc.Speaking(true)

	// Send the buffer data
	for _, buff := range buffer {
		vc.OpusSend <- buff
	}
	// Stop speaking
	vc.Speaking(false)

	// Sleep for a specified amount of time before ending.
	time.Sleep(250 * time.Millisecond)

	// Disconnect from the provided channel.
	vc.Disconnect()

	return nil
}
