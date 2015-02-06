package sqwiggle

import "time"

// Stream represents a stream object in the Sqwiggle API, which is similar
// to the concept of a Room.
type Stream struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`    // Id of the user that created this stream
	Name       string    `json:"name"`       // The full stream name
	Path       string    `json:"path"`       // The path to access stream in the web app, eg app.sqwiggle.com/:path
	Icon       string    `json:"icon"`       // An icon representing the stream
	IconColor  string    `json:"icon_color"` // A color representing the stream in hex format (eg: #121212)
	Subscribed bool      `json:"subscribed"` // Whether the user receives notifications for this stream
	CreatedAt  time.Time `json:"created_at"` // The time that this stream was created

	// Undocumented:
	Status      StreamStatus `json:"status"`
	Type        StreamType   `json:"type"`
	Description string       `json:"description"`
	// Position int         `json:"position"`
}

// StreamStatus represents the status of a stream
type StreamStatus string

// The StreamStatus constants below represent the possible stream
// states - this is undocumented in the API docs, so take it with
// a grain of salt.
const (
	StreamStatusActive   StreamStatus = "active"
	StreamStatusInactive              = "inactive"
)

// StreamType represents the type of a stream
type StreamType string

// The StreamType constants below represent the possible stream
// types - this is undocumented in the API docs, so take it with
// a grain of salt.
const (
	StreamTypeStandard StreamType = "standard"
	StreamTypeSupport             = "support"
)
