package api

import "math/rand"


func RandomMsg() string {
	msgs := []string{
		"random 1",
		"welcome to the unknown",
		"im a random messsage",
		"KKB on toppp",
		"akjsdhiuandi",
	}

	randomIndex := rand.Intn(len(msgs))
	randomMsg := msgs[randomIndex]

	return randomMsg
}

