package mongo

import (
	"context"
	"user-management/internal/model"
	"user-management/internal/repository"
	"user-management/internal/repository/mongo/transaction"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type UserRepository struct {
	cl     *mongo.Client
	db     *mongo.Database
	logger *zap.Logger
}

var _ repository.UserRepository = (*UserRepository)(nil)

func NewUserRepo(db *mongo.Database, cl *mongo.Client, logger *zap.Logger) *UserRepository {
	logger.Info("Creating new User repo")

	repo := UserRepository{
		db:     db,
		cl:     cl,
		logger: logger,
	}

	return &repo
}

func (ur *UserRepository) Create(ctx context.Context, u model.User) error {
	ur.logger.Info("creating new user")

	// session, err := ur.cl.StartSession()
	// if err != nil {
	// 	ur.logger.Info("user creation failed")
	// 	return err
	// }
	// defer session.EndSession(ctx)

	// transactionCtx := mongo.NewSessionContext(ctx, session)

	// if err := session.StartTransaction(); err != nil {
	// 	ur.logger.Info("user creation failed")
	// 	return err
	// }

	// if err := transaction.CreateTransaction(transactionCtx, u, ur.db.Collection("users")); err != nil {
	// 	ur.logger.Info("transaction failed")
	// 	session.AbortTransaction(ctx)
	// 	fmt.Println(err)
	// } else {
	// 	ur.logger.Info("transaction completed")
	// 	session.CommitTransaction(ctx)
	// }

	err := transaction.CreateTransaction(ctx, u, ur.db.Collection("users"))
	if err != nil {
		ur.logger.Info("user not inserted")
	}
	return err
}

func (ur *UserRepository) All(ctx context.Context) ([]model.User, error) {
	ur.logger.Info("Reading all users")

	result, err := transaction.AllTransaction(ctx, ur.db.Collection("users"))

	return result, err
}

func (ur *UserRepository) ReadByUsername(ctx context.Context, u model.Username) (model.User, error) {
	ur.logger.Info("Reading an user by username")

	result, err := transaction.ReadByUsernameTransaction(ctx, u, ur.db.Collection("users"))

	return result, err
}

func (ur *UserRepository) UpdateByID(ctx context.Context, u model.User) error {
	ur.logger.Info("Updating an user")

	err := transaction.UpdateByIDTransaction(ctx, u, ur.db.Collection("users"))

	return err
}

func (ur *UserRepository) UpdateByUsername(ctx context.Context, u model.User) error {
	ur.logger.Info("Updating an user by username")

	err := transaction.UpdateByUsernameTransaction(ctx, u, ur.db.Collection("users"))

	return err
}

func (ur *UserRepository) DeleteByID(ctx context.Context, id model.ID) error {
	ur.logger.Info("Deleting an user")

	err := transaction.DeleteByIDTransaction(ctx, id, ur.db.Collection("users"))

	return err
}
