package models

import "time"

type Category struct {
	ID        int       `json:"id" bun:",pk,autoincrement"`
	Name      string    `json:"name"`
	ParentID  int       `json:"parentID" bun:",nullzero"`
	CreatedAt time.Time `json:"-" bun:"default:current_timestamp"`
	UpdatedAt time.Time `json:"-" bun:"default:current_timestamp"`
}
