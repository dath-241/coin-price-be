package repository

import (
	"context"
	"github.com/dath-241/coin-price-be-go/services/admin_service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Find(ctx context.Context, filter bson.M) ([]models.UserDTO, error)
	FindOne(ctx context.Context, filter bson.M) *mongo.SingleResult
	DeleteOne(ctx context.Context, filter bson.M) (*mongo.DeleteResult, error)
	UpdateOne(ctx context.Context, filter bson.M, update bson.M) (*mongo.UpdateResult, error)
	ExistsByFilter(ctx context.Context, filter bson.M) (bool, error) // Mới thêm
    InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) // Mới thêm
}

type MongoUserRepository struct {
	Collection *mongo.Collection
}

func (r *MongoUserRepository) Find(ctx context.Context, filter bson.M) ([]models.UserDTO, error) {
	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.UserDTO
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *MongoUserRepository) FindOne(ctx context.Context, filter bson.M) *mongo.SingleResult {
	return r.Collection.FindOne(ctx, filter)
}

func (r *MongoUserRepository) DeleteOne(ctx context.Context, filter bson.M) (*mongo.DeleteResult, error) {
	return r.Collection.DeleteOne(ctx, filter)
}

// UpdateOne: Cập nhật thông tin người dùng
func (r *MongoUserRepository) UpdateOne(ctx context.Context, filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
    updateOptions := options.Update().SetUpsert(true) // Cập nhật nếu không tồn tại sẽ tạo mới
    return r.Collection.UpdateOne(ctx, filter, update, updateOptions)
}

func (r *MongoUserRepository) ExistsByFilter(ctx context.Context, filter bson.M) (bool, error) {
    var result bson.M
    err := r.Collection.FindOne(ctx, filter).Decode(&result)
    if err == mongo.ErrNoDocuments {
        return false, nil
    }
    if err != nil {
        return false, err
    }
    return true, nil
}

func (r *MongoUserRepository) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
    return r.Collection.InsertOne(ctx, document)
}