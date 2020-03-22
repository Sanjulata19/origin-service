package main

import (
	"encoding/json"
	"github.com/nullserve/static-host/cmd"
	"log"
)

type LogMessage struct {
	Error Error `json:"error"`
}

type Error struct {
	Message string `json:"message"`
}

func main() {
	err := cmd.Execute()
	if err != nil {
		if message, err := json.Marshal(LogMessage{
			Error: Error{
				Message: err.Error(),
			},
		}); err == nil {
			log.Fatal(message)
		} else {
			panic(err.Error())
		}
	}
}
