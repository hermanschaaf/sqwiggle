package sqwiggle

import "time"

// Conversation represents an ephemeral media connection between two or more people.
type Conversation struct {
	ID            int                `json:"id"`
	Status        ConversationStatus `json:"status"`        // open, closed
	Duration      int                `json:"duration"`      // The number of seconds this conversation lasted, or if open has been ongoing.
	CreatedAt     time.Time          `json:"created_at"`    // The time this conversation was started
	Participating []User             `json:"participating"` // A list of User objects that are currently participating in the conversation
	Participated  []User             `json:"participated"`  // A list of User objects that have been and or are currently in the conversation

	// undocumented
	ColorID   int  `json:"color_id"`
	MCU       bool `json:"mcu"`
	MCUServer bool `json:"mcu_server"`
	Locked    bool `json:"locked"`
}

// ConversationStatus describes the status of a conversation
type ConversationStatus string

// These ConversationStatus constants describe the possible statuses
// of a conversation, i.e. open or closed.
const (
	ConversationOpen   ConversationStatus = "open"
	ConversationClosed                    = "closed"
)
