package models

import (
	"time"
)

type CommitInfo struct {

	Hash string `json:"hash,omitempty"`

	AuthorName string `json:"author_name,omitempty"`

	AuthorEmail string `json:"author_email,omitempty"`

	Timestamp time.Time `json:"timestamp,omitempty"`

	Message string `json:"message,omitempty"`
}
