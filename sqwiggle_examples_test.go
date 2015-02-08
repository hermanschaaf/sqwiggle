package sqwiggle_test

import (
	"fmt"

	"github.com/hermanschaaf/sqwiggle"
)

// The following code instantiates a client, then calls the
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
		fmt.Printf("%s: %s\n", m.Author.Name, m.Text)
	}
}

// The following code instantiates a client, then calls the
// GetMessage method to return a single message. If no error occurred, it
// prints out the author name and text of the message.
func ExampleClient_GetMessage() {
	client := sqwiggle.NewClient("YOUR-API-KEY")

	id := 1000
	m, err := client.GetMessage(id)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s: %s\n", m.Author.Name, m.Text)
}

// The following code instantiates a client, then calls the
// PostMessage method to create a single message. If no error occurred, it
// prints out the author name and text of the new message.
func ExampleClient_PostMessage() {
	client := sqwiggle.NewClient("YOUR-API-KEY")

	streamID := 48914
	text := "Hello from the <b>Go Sqwiggle API Client</b>!"
	options := sqwiggle.PostMessageOptions{
		Format: "html", // allow the use of certain HTML tags in the message
		Parse:  false,  // don't parse rich attachments in the message
	}
	m, err := client.PostMessage(streamID, text, &options) // it is also okay to pass in nil for options
	if err != nil {
		panic(err)
	}

	fmt.Printf("[%d] %s: %s\n", m.ID, m.Author.Name, m.Text)
}

// The following code instantiates a client, then calls the
// UpdateMessage method to update a single message. If no error occurred, it
// prints out the author name and text of the new message.
func ExampleClient_UpdateMessage() {
	client := sqwiggle.NewClient("YOUR-API-KEY")

	messageID := 48914
	text := "Hello again from the <b>Go Sqwiggle API Client</b>!"
	m, err := client.UpdateMessage(messageID, text)
	if err != nil {
		panic(err)
	}

	fmt.Printf("[%d] %s: %s\n", m.ID, m.Author.Name, m.Text)
}

// The following code instantiates a client, then calls the
// DeleteMessage method to delete a single message.
func ExampleClient_DeleteMessage() {
	client := sqwiggle.NewClient("YOUR-API-KEY")

	messageID := 48914
	err := client.DeleteMessage(messageID)
	if err != nil {
		panic(err)
	}
}
