// Package models defines the different object models used in Project
package models

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
