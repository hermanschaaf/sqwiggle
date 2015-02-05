package sqwiggle

// Mention represents a mention as part of a message in any chat stream.
type Mention struct {
	ID          int      `json:"id"`
	MessageID   int      `json:"message_id"`
	Name        string   `json:"name"`
	Indices     []int    `json:"indices"`
	Text        string   `json:"text"`
	SubjectType UserType `json:"subject_type"`
	SubjectID   int      `json:"subject_id"`
}
