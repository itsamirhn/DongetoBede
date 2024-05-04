package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/itsamirhn/dongetobede/internal/database/entities"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type client struct {
	*mongo.Client
	db *mongo.Database
}

func NewMongoClient(mc *mongo.Client, dbName string) Client {
	return &client{
		Client: mc,
		db:     mc.Database(dbName),
	}
}

func (c *client) AddUser(ctx context.Context, user *entities.User) (int64, error) {
	res, err := c.db.Collection(entities.UserCollectionName).InsertOne(ctx, *user)
	if err != nil {
		return 0, errors.Wrap(err, "failed to insert user")
	}
	return res.InsertedID.(int64), nil
}

func (c *client) GetUserByID(ctx context.Context, id int64) (*entities.User, error) {
	var user entities.User
	err := c.db.Collection(entities.UserCollectionName).FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user by ID")
	}
	return &user, nil
}

func (c *client) GetUsersByIDs(ctx context.Context, ids []int64) ([]*entities.User, error) {
	cursor, err := c.db.Collection(entities.UserCollectionName).Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get users by IDs")
	}
	defer cursor.Close(ctx)
	var users []*entities.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, errors.Wrap(err, "failed to get users by IDs")
	}
	return users, nil
}

func (c *client) UpdateUser(ctx context.Context, user *entities.User) error {
	_, err := c.db.Collection(entities.UserCollectionName).UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": *user})
	if err != nil {
		return errors.Wrap(err, "failed to update user")
	}
	return nil
}

func (c *client) AddDong(ctx context.Context, dong *entities.Dong) (*primitive.ObjectID, error) {
	res, err := c.db.Collection(entities.DongCollectionName).InsertOne(ctx, *dong)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert dong")
	}
	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("failed to cast newly added object ID to objectID")
	}
	return &id, nil
}

func (c *client) GetDongByID(ctx context.Context, id primitive.ObjectID) (*entities.Dong, error) {
	var dong entities.Dong
	err := c.db.Collection(entities.DongCollectionName).FindOne(ctx, bson.M{"_id": id}).Decode(&dong)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get dong by ID")
	}
	return &dong, nil
}

func (c *client) UpdateDong(ctx context.Context, dong *entities.Dong) error {
	_, err := c.db.Collection(entities.DongCollectionName).UpdateOne(ctx, bson.M{"_id": dong.ID}, bson.M{"$set": *dong})
	if err != nil {
		return errors.Wrap(err, "failed to update dong")
	}
	return nil
}
