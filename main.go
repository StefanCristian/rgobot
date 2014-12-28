/*
GoBot

An IRC bot written in Go.

Copyright (C) 2014  Brian C. Tomlinson

Contact: brian.tomlinson@linux.com

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License along
with this program; if not, write to the Free Software Foundation, Inc.,
51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
*/
package main

import (
	"encoding/json"
	"fmt"
	"github.com/thoj/go-ircevent"
	"time"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Config struct {
	Server, Channel, BotUser, BotNick, LogDir string
}

// ParseCmds takes PRIVMSG strings containing a preceding bang "!"
// and attempts to turn them into an ACTION that makes sense.
// Returns a msg string.
func ParseCmds(cmdMsg string) string {
	cmdArray := strings.SplitAfterN(cmdMsg, "!", 2)
	msgArray := strings.SplitN(cmdArray[1], " ", 2)
	cmd := fmt.Sprintf("%vs", msgArray[0])

	// This should give us something like:
	//     "Snuffles slaps $USER, FOR SCIENCE!"
	// If given the command:
	//     "!slap $USER"
	msg := fmt.Sprintf("\x01"+"ACTION %v %v, FOR SCIENCE!\x01", cmd, msgArray[1])
	return msg
}
// Creates Log Directory
func LogDir(CreateDir string) {

        //Check if the LogDir Exists. And if not Create it.
        if _, err := os.Stat(CreateDir); os.IsNotExist(err) {
                fmt.Printf("No such file or directory: %s", CreateDir)
                os.Mkdir(CreateDir, 0777)
        } else {
                fmt.Printf("Its There: %s", CreateDir)
        }
}

func LogFile(CreateFile string) {
        //Check if the Log File for the Channel(s) Exists if not create it
        if _, err := os.Stat(CreateFile + ".log"); os.IsNotExist(err) {
                fmt.Printf("Log File " + CreateFile + ".log Doesn't Exist.\n")
		os.Create(CreateFile + ".log")
		fmt.Printf("Created the log file\n")
        } else {
                fmt.Printf("Log File Exists.\n")
        }
}

// UrlTitle attempts to extract the title of the page that a
// pasted URL points to.
// Returns a string message with the title and URL on success, or a string
// with an error message on failure.
func UrlTitle(msg string) string {
	var (
		newMsg, url, title, word string
	)

	regex, _ := regexp.Compile(`<title[^>]*>([^<]+)<\/title>`)

	msgArray := strings.Split(msg, " ")

	for _, word = range msgArray {
		if strings.Contains(word, "http") || strings.Contains(word, "www") {
			url = word
			break
		}
	}

	resp, err := http.Get(word)

	if err != nil {
		newMsg = fmt.Sprintf("Could not resolve URL %v, beware...\n", word)
		return newMsg
	}

	defer resp.Body.Close()

	rawBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		newMsg = fmt.Sprintf("Could not read response Body of %v ...", word)
		return newMsg
	}

	body := string(rawBody)
	title = regex.FindString(body)
	newMsg = fmt.Sprintf("[ %v ]->( %v )", title, url)

	return newMsg
}

func ChannelLogger(LogDir string, UserNick string, message string) {
	STime := time.Now().Format(time.ANSIC)
	logLoc := fmt.Sprintf("%d-%s-%d", time.Now().Day(), time.Now().Month(), time.Now().Year())

	f, err := os.OpenFile(LogDir + logLoc + ".log", os.O_CREATE|os.O_RDWR|os.O_APPEND|os.O_SYNC, 0666)
	if err != nil {
		fmt.Println(f, err)
	}

	n, err := io.WriteString(f, STime + " > " + UserNick + message + "\n")
	if err != nil {
		fmt.Println(n, err)
	}
	f.Close()
}

// AddCallbacks is a single function that does what it says.
// It's merely a way of decluttering the main function.
func AddCallbacks(conn *irc.Connection, config *Config) {
	conn.AddCallback("001", func(e *irc.Event) {
		conn.Join(config.Channel)
	})
	conn.AddCallback("JOIN", func(e *irc.Event) {
		if e.Nick == config.BotNick {
			dateLog := fmt.Sprintf("%d-%s-%d", time.Now().Day(), time.Now().Month(), time.Now().Year())
			LogDir(config.LogDir)
			LogFile(config.LogDir+dateLog)
		}
		message := " has joined"
		ChannelLogger(config.LogDir, e.Nick, message)
	})

        conn.AddCallback("PART", func (e *irc.Event) {
		pmessage := " parted "
		message := e.Message()
		spaceAround := "@"
		spaceAccoladeOne := "( "
		spaceAccoladeZwei := " )"
                ChannelLogger(config.LogDir, fmt.Sprintf("%v%v%v", e.Nick, spaceAround, e.Host), fmt.Sprintf("%v%v%v%v", pmessage, spaceAccoladeOne, message, spaceAccoladeZwei))
        })

        conn.AddCallback("QUIT", func (e *irc.Event) {
                qmessage := " has quit "
                message := e.Message()
		spaceAround := "@"
		spaceAccoladeOne := "( "
		spaceAccoladeZwei := " )"
                ChannelLogger(config.LogDir, fmt.Sprintf("%v%v%v", e.Nick, spaceAround, e.Host), fmt.Sprintf("%v%v%v%v", qmessage, spaceAccoladeOne, message, spaceAccoladeZwei))
        })

	conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		var response string
		message := e.Message()

		if e.Host == "unaffiliated/blacknoxis" && strings.Contains(message, "#meriacas") && strings.Index(message, "#meriacas") == 0 {
			// This is intentionally borken. Further functions will follow
			response = ParseCmds(message)
		}

		/* if strings.Contains(e.Message, "http") || strings.Contains(e.Message, "www") {
			response = UrlTitle(e.Message)
		} */

		if strings.Contains(message, "#sursa") || strings.Contains(message, "#surse") || strings.Contains(message, "#sources") {
			conn.Privmsg(config.Channel, "http://github.com/BlackNoxis http://github.com/Rogentos")
		}

		if strings.Contains(message, "#wiki") {
			conn.Privmsg(config.Channel, "http://wiki.rogentos.ro/")
		}

                if strings.Contains(message, "#logs") {
			conn.Privmsg(config.Channel, "http://bpr.bluepink.ro/~rogentos/logs")
		}

		if len(response) > 0 {
			conn.Privmsg(config.Channel, response)
		}
		if len(message) > 0 {
			if e.Arguments[0] != config.BotNick {
				spacePoint := ":"
				ChannelLogger(config.LogDir, e.Nick, fmt.Sprintf("%v %v", spacePoint, message))
			} else {
				// Someone is trying to speak to the bot
				conn.Privmsg(e.Nick, "There is no function implemented for private messages")
			}

		}
	})
}

// Connect tries up to three times to get a connection to the server
// and channel, hopefully with a nil err value at some point.
// Returns error
func Connect(conn *irc.Connection, config *Config) error {
	var err error

	for attempt := 1; attempt <= 3; attempt++ {
		if err = conn.Connect(config.Server); err != nil {
			fmt.Println("Connection attempt %v failed, trying again...", attempt)
			continue
		} else {
			break
		}
	}
	return err
}

func main() {

	// Read the config file and populate our Config struct.
	file, err := os.Open("config.json")

	if err != nil {
		fmt.Println("Couldn't read config file, dying...")
		panic(err)
	}

	decoder := json.NewDecoder(file)
	config := &Config{}
	decoder.Decode(&config)

	conn := irc.IRC(config.BotNick, config.BotUser)
	err = Connect(conn, config)

	if err != nil {
		fmt.Println("Failed to connect.")
		// Without a connection, we're useless, panic and die.
		panic(err)
	}

	AddCallbacks(conn, config)
	conn.Loop()
}
