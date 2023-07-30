package nosql

import (
	"context"
	"time"
	"user-management/internal/model"
	"user-management/internal/repository"
	"user-management/internal/repository/nosql/nosqlmodel"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type UserRepository struct {
	db     *mongo.Database
	logger *zap.Logger
}

var _ repository.UserRepository = (*UserRepository)(nil)

func NewUserRepo(db *mongo.Database, logger *zap.Logger) *UserRepository {
	logger.Info("Creating new User repo")
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (ur *UserRepository) Create(u model.User) error {
	ur.logger.Info("Creating new user")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := ur.db.Collection("users").FindOne(ctx, bson.M{"username": u.Username}).Err()
	if err == nil {
		ur.logger.Info("User exists already")
		return err
	}

	user := nosqlmodel.User{
		Username:       u.Username,
		Password:       u.Password,
		TimeOfCreation: time.Now(),
		City:           u.City,
	}

	_, err = ur.db.Collection("users").InsertOne(ctx, &user)
	if err != nil {
		ur.logger.Info("User insert failure")
		return err
	}
	ur.logger.Info("User inserted to database")

	return nil
}

func (ur *UserRepository) All() ([]model.User, error) {
	ur.logger.Info("Reading all users")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cursor, err := ur.db.Collection("users").Find(ctx, bson.D{})
	if err != nil {
		ur.logger.Info("Read all users to cursor failure")
		return nil, err
	}
	ur.logger.Info("Read all users to cursor")

	var users []nosqlmodel.User
	err = cursor.All(ctx, &users)
	if err != nil {
		ur.logger.Info("Read all users from cursor failure")
		return nil, err
	}
	ur.logger.Info("Read all users from cursor")

	result := make([]model.User, len(users))
	for i, u := range users {
		result[i] = model.User{
			ID:             model.ID(u.ID.String()),
			Username:       u.Username,
			Password:       u.Password,
			TimeOfCreation: u.TimeOfCreation,
			City:           u.City,
		}
	}
	ur.logger.Info("All users read")

	return result, nil
}

func (ur *UserRepository) ReadByUsername(u model.User) (model.User, error) {
	ur.logger.Info("Reading an user by username")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var user nosqlmodel.User
	err := ur.db.Collection("users").FindOne(ctx, bson.M{"username": u.Username}).Decode(&user)
	if err != nil {
		ur.logger.Info("Read user by username failure")
		return model.User{}, nil
	}
	ur.logger.Info("User by username read")

	var result model.User
	result.ID = model.ID(user.ID.String())
	result.Username = user.Username
	result.Password = user.Password
	result.TimeOfCreation = user.TimeOfCreation
	result.City = user.City

	return result, nil
}

func (ur *UserRepository) UpdateByID(u model.User) error {
	ur.logger.Info("Updating an user")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var user nosqlmodel.User
	err := ur.db.Collection("users").FindOne(ctx, bson.M{"_id": string(u.ID)}).Decode(&user)
	if err != nil {
		ur.logger.Info("User doesn't exists")
		return err
	}

	user.Username = u.Username
	user.Password = u.Password
	user.City = u.City

	_, err = ur.db.Collection("users").UpdateOne(ctx, bson.M{"_id": u.ID}, &user)
	if err != nil {
		ur.logger.Info("User update failure")
		return err
	}
	ur.logger.Info("User updated")

	return nil
}

func (ur *UserRepository) UpdateByUsername(u model.User) error {
	ur.logger.Info("Updating an user by username")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var user nosqlmodel.User
	err := ur.db.Collection("users").FindOne(ctx, bson.M{"username": u.Username}).Decode(&user)
	if err != nil {
		ur.logger.Info("User doesn't exists")
		return err
	}

	user.Password = u.Password
	user.City = u.City

	_, err = ur.db.Collection("users").UpdateOne(ctx, bson.M{"username": u.Username}, &user)
	if err != nil {
		ur.logger.Info("User update by username failure")
	}
	ur.logger.Info("User updated by username")

	return nil
}

func (ur *UserRepository) DeleteByID(id model.ID) error {
	ur.logger.Info("Deleting an user")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	mId, _ := primitive.ObjectIDFromHex(string(id))
	_, err := ur.db.Collection("users").DeleteOne(ctx, bson.M{"_id": mId})
	if err != nil {
		ur.logger.Info("User delete failure")
		return err
	}
	ur.logger.Info("User deleted")

	return nil
}
