package sqwiggle

import "time"

// User represents a single person on your organization's team.
type User struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`   // The users full name
	Avatar string `json:"avatar"` // URL to a static avatar for the user

	Role          UserRole   `json:"role"`            // user, owner, manager or banned
	MediaDeviceID string     `json:"media_device_id"` // A string representing current device that media is being received on
	Status        UserStatus `json:"status"`          // busy, available
	Message       string     `json:"message"`         // A status message that other users see, such as “out for lunch”

	Email string `json:"email"` // The users email address

	Snapshot         string    `json:"snapshot"`          // URL to the last snapshot for this user (auto or manual)
	SnapshotInterval int       `json:"snapshot_interval"` // Frequency at which automatic snapshots are taken when the app is open
	Confirmed        bool      `json:"confirmed"`         // The users email confirmation status
	TimeZone         string    `json:"time_zone"`         // Timezone (rails format)
	TimeZoneOffset   float64   `json:"time_zone_offset"`  // Hours offset from UTC, note that this may be a non-integer like 5.5
	CreatedAt        time.Time `json:"created_at"`        // The time this user was created
	LastActiveAt     time.Time `json:"last_active_at"`    // The last time we recorded activity for a user
	LastConnectedAt  time.Time `json:"last_connected_at"` // The time this users current online session started

	// undocumented:
	Type UserType `json:"type,omitempty"`
	// IdleAt
	// Contact
	Support bool `json:"support,omitempty"` // Indicates whether user is support or not
}

// UserRole describes the role of a user, i.e. normal user, owner, manager or banned
type UserRole string

// These UserRole constants define the possible user roles that
// can be returned by the API.
const (
	RoleUser    UserRole = "user"
	RoleOwner            = "owner"
	RoleManager          = "manager"
	RoleBanned           = "banned"
)

func (ur UserRole) String() string {
	return string(ur)
}

// UserStatus describes the current status of the user
type UserStatus string

// The below UserStatus constants define the typical states a user may
// be in, i.e. "busy", "available" or "offline"
const (
	StatusBusy      UserStatus = "busy"
	StatusAvailable            = "available"
	StatusOffline              = "offline"
)

// UserType describes the type of user, i.e. normal user or support
type UserType string

// The below UserType constants define the known UserTypes that may
// be returned by the API.
const (
	TypeUser    UserType = "user"
	TypeClient           = "client"
	TypeSupport          = "support"
)
