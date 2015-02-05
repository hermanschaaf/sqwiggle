package sqwiggle

import "time"

// Message represents an item in any chat stream.
// They may be created by humans, bots or the Sqwiggle servers.
type Message struct {
	ID          int          `json:"id"`          // ID of the message
	StreamID    int          `json:"stream_id"`   // ID of the chat stream that this message belongs to
	Text        string       `json:"text"`        // The plain text content of the message, HTML will be escaped
	Author      User         `json:"author"`      // An object representing the user or API client that created the message
	Attachments []Attachment `json:"attachments"` // A list of Attachment objects to be displayed with this message
	Mentions    []Mention    `json:"mentions"`    // A list of users tagged / mentioned in this message
	CreatedAt   time.Time    `json:"created_at"`  // The time that this message was created
	UpdatedAt   time.Time    `json:"updated_at"`  // The time that this message was last updated or edited

	// Undocumented
	ConversationID *int `json:"conversation_id,omitempty"`
}
