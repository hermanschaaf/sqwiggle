package sqwiggle

import "time"

// AttachmentType represents the type of attachment in a message
type AttachmentType string

// Validation Token Types
const (
	TypeImage         AttachmentType = "image"
	TypeLink                         = "link"
	TypeFile                         = "file"
	TypeTwitterStatus                = "twitter_status"
	TypeTwitterUser                  = "twitter_user"
	TypeVideo                        = "video"
	TypeCode                         = "code"
	TypeGist                         = "gist"
)

func (t AttachmentType) String() string {
	return string(t)
}

// Attachment is a piece of media that belongs to a message.
// It may represent a link, image, video, file upload or more.
// Sqwiggle often adds new attachment types based on demand so it
// is important to filter using the type parameter.
type Attachment struct {
	ID          int            `json:"id"`          // ID of the attachment
	Type        AttachmentType `json:"type"`        // image, link, file, twitter_status, twitter_user, video, code, gist
	URL         string         `json:"url"`         // URL where the attachment content can be accessed
	Title       string         `json:"title"`       // A title for the attachment, for example a filename or webpage title
	Description string         `json:"description"` // A description of the attachment, for example a web page summary
	Image       string         `json:"image"`       // URL of an image representing the attachment, this may not reside on Sqwiggle's servers
	Status      string         `json:"status"`      // If an upload, denotes the uploade status ('pending' or 'uplodaed')
	Animated    bool           `json:"animated"`    // If an image, denotes whether animated
	CreatedAt   time.Time      `json:"created_at"`  // The time that this attachment was created
	UpdatedAt   time.Time      `json:"updated_at"`  // The time that this attachment was last updated or edited

	// undocumented:
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`
}
