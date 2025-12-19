package notes

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotesStats struct {
	Count        int64   `bson:"count" json:"count"`
	AvgContentLn float64 `bson:"avgContentLength" json:"avgContentLength"`
}

type Note struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title"         json:"title"`
	Content   string             `bson:"content"       json:"content"`
	CreatedAt time.Time          `bson:"createdAt"     json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"     json:"updatedAt"`
	ExpiresAt *time.Time         `bson:"expiresAt,omitempty" json:"expiresAt,omitempty"`
}
