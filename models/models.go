// Package models defines the different object models used in Project
package models

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var (
	// ContextKeyParsedToken var
	ContextKeyParsedToken = contextKey("ParsedToken")
)

// ParsedToken contains details that are encoded in the received JWT token
type ParsedToken struct {
	Email string
}

// Project defines the project object
type Project struct {
	Name string `bson:"name" json:"name"`
}

// ProjectWithCursor implements the pair of the project with a cursor that determines the
// location of the tennat in the repository.
type ProjectWithCursor struct {
	ProjectID string
	Project   Project
	Cursor    string
}
