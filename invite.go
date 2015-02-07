package sqwiggle

import "time"

type Invite struct {
	ID        int       `json:"id"`
	FromID    int       `json:"from_id"`    // ID of the user that created the invite
	Email     string    `json:"email"`      // The email address that this invite was sent to
	Avatar    string    `json:"avatar"`     // URL to a static avatar representing the email address
	URL       string    `json:"url"`        // URL to redeem the invite
	CreatedAt time.Time `json:"created_at"` // The time that this invite was created
}
