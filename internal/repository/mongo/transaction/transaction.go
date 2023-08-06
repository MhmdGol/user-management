package transaction

import (
	"context"
	"fmt"
	"time"
	"user-management/internal/model"
	"user-management/internal/repository/mongo/mongomodel"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateTransaction(ctx context.Context, u model.User, collection *mongo.Collection) error {
	filter := bson.M{"username": bson.M{"$eq": u.Username}}

	err := collection.FindOne(ctx, filter).Err()
	if err == nil {
		return err
	}

	user := mongomodel.User{
		Username:       string(u.Username),
		Password:       string(u.Password),
		Role:           string(u.Role),
		TimeOfCreation: time.Now(),
		City:           u.City,
		Version:        1,
	}

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func AllTransaction(ctx context.Context, collection *mongo.Collection) ([]model.User, error) {
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var users []mongomodel.User
	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, err
	}

	result := make([]model.User, len(users))
	for i, u := range users {
		result[i] = model.User{
			ID:             model.ID(u.ID.String()),
			Username:       model.Username(u.Username),
			Password:       model.HashedPass(u.Password),
			Role:           model.Role(u.Role),
			TimeOfCreation: u.TimeOfCreation,
			City:           u.City,
			Version:        u.Version,
		}
	}

	return result, nil
}

func ReadByUsernameTransaction(ctx context.Context, u model.Username, collection *mongo.Collection) (model.User, error) {
	filter := bson.M{"username": bson.M{"$eq": u}}

	var user mongomodel.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return model.User{}, nil
	}

	result := model.User{
		ID:             model.ID(user.ID.String()),
		Username:       model.Username(user.Username),
		Password:       model.HashedPass(user.Password),
		Role:           model.Role(user.Role),
		TimeOfCreation: user.TimeOfCreation,
		City:           user.City,
		Version:        user.Version,
	}

	return result, nil
}

func UpdateByIDTransaction(ctx context.Context, u model.User, collection *mongo.Collection) error {
	objectId, err := primitive.ObjectIDFromHex(string(u.ID))
	if err != nil {
		return err
	}
	filter := bson.M{"_id": bson.M{"$eq": objectId}, "version": bson.M{"$eq": u.Version}}

	var user mongomodel.User
	err = collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return err
	}

	if u.Username != "" {
		user.Username = string(u.Username)
	}
	if u.Password != "" {
		user.Password = string(u.Password)
	}
	if u.Role != "" {
		user.Role = string(u.Role)
	}
	if u.City != "" {
		user.City = u.City
	}
	user.Version = u.Version + 1
	update := bson.M{"$set": user}

	_, err = collection.UpdateOne(ctx, filter, update)

	return err
}

func UpdateByUsernameTransaction(ctx context.Context, u model.User, collection *mongo.Collection) error {
	filter := bson.M{"username": bson.M{"$eq": u.Username}, "version": bson.M{"$eq": u.Version}}

	var user mongomodel.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return err
	}

	if u.Password != "" {
		user.Password = string(u.Password)
	}
	if u.Role != "" {
		user.Role = string(u.Role)
	}
	if u.City != "" {
		user.City = u.City
	}
	user.Version = u.Version + 1
	update := bson.M{"$set": user}

	_, err = collection.UpdateOne(ctx, filter, update)

	return err

}

func DeleteByIDTransaction(ctx context.Context, id model.ID, collection *mongo.Collection) error {
	objectId, _ := primitive.ObjectIDFromHex(string(id))
	filter := bson.M{"_id": bson.M{"$eq": objectId}}
	_, err := collection.DeleteOne(ctx, filter)

	return err
}
