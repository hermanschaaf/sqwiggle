package sqwiggle_test

import (
	"fmt"

	"github.com/hermanschaaf/sqwiggle"
)

// ExampleClient_ListMessages instantiates a client, then calls the
// ListMessages method to return a slice of Messages. If no error occurred, it
// iterates through the messages and prints them out one by one.
func ExampleClient_ListMessages() {
	client := sqwiggle.NewClient("YOUR-API-KEY")

	page, limit := 0, 50
	msgs, err := client.ListMessages(page, limit)

	if err != nil {
		panic(err)
	}

	for _, m := range msgs {
		fmt.Printf("%s: %s\n", m.Author, m.Text)
	}
}
