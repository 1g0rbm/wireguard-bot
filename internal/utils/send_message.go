package utils

import (
	"log"
)

func SendMessage(render func() ([]byte, error), send func(msg []byte) error) {
	msgBytes, err := render()
	if err != nil {
		log.Fatalf("Render message error.\nerr: %v\n", err)
	}

	err = send(msgBytes)
	if err != nil {
		log.Fatalf("Sending message error.\nerr: %v\n", err)
	}
}
