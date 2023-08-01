package mongomodel

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Username       string             `bson:"username"`
	Password       string             `bson:"password"`
	Role           string             `bson:"role"`
	TimeOfCreation time.Time          `bson:"timeOfCreation"`
	City           string             `bson:"city"`
	Version        int                `bson:"version"`
}
