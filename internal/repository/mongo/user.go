package mongo

import (
	"context"
	"time"
	"user-management/internal/model"
	"user-management/internal/repository"
	"user-management/internal/repository/mongo/mongomodel"

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

	repo := UserRepository{
		db:     db,
		logger: logger,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := repo.Create(ctx, model.User{
		Username: "su",
		Password: "Admin@123",
		Role:     "admin",
	})
	if err != nil {
		logger.Info("su insert failure")
	}
	logger.Info("su inserted to database")

	return &repo
}

func (ur *UserRepository) Create(ctx context.Context, u model.User) error {
	ur.logger.Info("Creating new user")

	filter := bson.M{"username": bson.M{"$eq": u.Username}}

	err := ur.db.Collection("users").FindOne(ctx, filter).Err()
	if err == nil {
		ur.logger.Info("User exists already")
		return err
	}
	ur.logger.Info("Safe to insert")

	user := mongomodel.User{
		Username:       string(u.Username),
		Password:       string(u.Password),
		Role:           string(u.Role),
		TimeOfCreation: time.Now(),
		City:           u.City,
		Version:        1,
	}

	_, err = ur.db.Collection("users").InsertOne(ctx, &user)
	if err != nil {
		ur.logger.Info("User insert failure")
		return err
	}
	ur.logger.Info("User inserted to database")

	return nil
}

func (ur *UserRepository) All(ctx context.Context) ([]model.User, error) {
	ur.logger.Info("Reading all users")

	cursor, err := ur.db.Collection("users").Find(ctx, bson.D{})
	if err != nil {
		ur.logger.Info("Read all users to cursor failure")
		return nil, err
	}
	ur.logger.Info("Read all users to cursor")

	var users []mongomodel.User
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
			Username:       model.Username(u.Username),
			Password:       model.HashedPass(u.Password),
			Role:           model.Role(u.Role),
			TimeOfCreation: u.TimeOfCreation,
			City:           u.City,
			Version:        u.Version,
		}
	}
	ur.logger.Info("All users read")

	return result, nil
}

func (ur *UserRepository) ReadByUsername(ctx context.Context, u model.Username) (model.User, error) {
	ur.logger.Info("Reading an user by username")

	filter := bson.M{"username": bson.M{"$eq": u}}

	var user mongomodel.User
	err := ur.db.Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		ur.logger.Info("Read user by username failure")
		return model.User{}, nil
	}
	ur.logger.Info("User by username read")

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

func (ur *UserRepository) UpdateByID(ctx context.Context, u model.User) error {
	ur.logger.Info("Updating an user")

	objectId, err := primitive.ObjectIDFromHex(string(u.ID))
	if err != nil {
		return err
	}
	filter := bson.M{"_id": bson.M{"$eq": objectId}, "version": bson.M{"$eq": u.Version}}

	var user mongomodel.User
	err = ur.db.Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		ur.logger.Info("User does not exist or version is not valid")
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

	_, err = ur.db.Collection("users").UpdateOne(ctx, filter, update)
	if err != nil {
		ur.logger.Info("User update failure")
		return err
	}
	ur.logger.Info("User updated")

	return nil
}

func (ur *UserRepository) UpdateByUsername(ctx context.Context, u model.User) error {
	ur.logger.Info("Updating an user by username")

	filter := bson.M{"username": bson.M{"$eq": u.Username}, "version": bson.M{"$eq": u.Version}}

	var user mongomodel.User
	err := ur.db.Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		ur.logger.Info("User does not exist or version is not valid")
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

	_, err = ur.db.Collection("users").UpdateOne(ctx, filter, update)
	if err != nil {
		ur.logger.Info("User update by username failure")
		return err
	}
	ur.logger.Info("User updated by username")

	return nil
}

func (ur *UserRepository) DeleteByID(ctx context.Context, id model.ID) error {
	ur.logger.Info("Deleting an user")

	objectId, _ := primitive.ObjectIDFromHex(string(id))
	filter := bson.M{"_id": bson.M{"$eq": objectId}}
	_, err := ur.db.Collection("users").DeleteOne(ctx, filter)
	if err != nil {
		ur.logger.Info("User delete failure")
		return err
	}
	ur.logger.Info("User deleted")

	return nil
}
