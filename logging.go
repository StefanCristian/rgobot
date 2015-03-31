package main

import (
	"fmt"
	"io"
	"os"
	"time"
)


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
