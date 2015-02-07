package sqwiggle

import "time"

// Organization represents a single organization that every user must belong to.
type Organization struct {
	ID        int         `json:"id"`
	Name      string      `json:"name"`       // The organizations name
	CreatedAt time.Time   `json:"created_at"` // The time that this company was created
	Path      string      `json:"path"`       // The url path to access company on app, eg sqwiggle.com/:path
	Billing   OrgBilling  `json:"billing"`
	Security  OrgSecurity `json:"security"`

	// Undocumented
	InviteURL                   string `json:"invite_url"`
	UserCount                   int    `json:"user_count"`
	MaxConversationParticipants int    `json:"max_conversation_participants"`
}

// OrgBilling represents the billing information for an organization.
type OrgBilling struct {
	Plan       string `json:"plan"`        // undocumented
	Status     string `json:"status"`      // undocumented, but can be "trial"
	ActiveCard bool   `json:"active_card"` // whether the organization has an active credit card
	Receipts   bool   `json:"receipts"`    // whether the organization wants receipts
	Email      string `json:"email"`       // the organization billing email address
}

// OrgSecurity represents the known fields returned in the security field
// of Organizations. These fields are not explained in the documentation,
// so it's up to you to figure out what they mean (or PR if you know).
type OrgSecurity struct {
	MediaAccept     bool `json:"media_accept"`
	DomainRestrict  bool `json:"domain_restrict"`
	DomainSignup    bool `json:"domain_signup"`
	OpenInvites     bool `json:"open_invites"`
	UploadsDisabled bool `json:"uploads_disabled"`
	ManualDisabled  bool `json:"manual_disabled"`
}
