package dbmigrate

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UsersSchema() *options.CreateCollectionOptions {
	collectionOptions := options.CreateCollection()
	collectionOptions.SetValidator(
		bson.M{
			"bsonType": "object",
			"required": []string{"username", "password", "role", "timeOfCreation", "version"},
			"properties": bson.M{
				"username": bson.M{
					"bsonType":    "string",
					"description": "the username of the user",
				},
				"password": bson.M{
					"bsonType":    "string",
					"description": "the password of the user",
				},
				"role": bson.M{
					"bsonType":    "string",
					"description": "the role of the user in the system",
				},
				"timeOfCreation": bson.M{
					"bsonType":    "date",
					"description": "the first time when the user have been created",
				},
				"city": bson.M{
					"bsonType":    "string",
					"description": "the city the user lives in",
				},
				"version": bson.M{
					"bsonType":    "int32",
					"description": "the version of the document after updating",
				},
			},
		})

	return collectionOptions
}
