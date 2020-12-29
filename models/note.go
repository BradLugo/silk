package models

import "github.com/google/uuid"

type Note struct {
	Id        uuid.UUID `json:"id"`
	Title     *string   `json:"title"`
	Text      *string   `json:"text"`
	Citation  *string   `json:"citation"`
	RelatedTo []*Note   `json:"relatedTo"`
}
