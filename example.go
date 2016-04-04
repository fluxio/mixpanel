package main

import (
	"fmt"

	"./mixpanel"
)

func main() {
	m := mixpanel.Client{
		ApiKey:    "YOUR KEY",
		ApiSecret: "YOUR SECRET",
	}
	res, err := m.Request([]string{"events"}, map[string](interface{}){
		"event":    []string{"pages"},
		"unit":     "hour",
		"interval": 24,
		"type":     "general",
	}, "POST", "")
	fmt.Printf("Response %v\nError %v\n", string(res), err)
}
