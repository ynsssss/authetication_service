package repositories

import (
	"authentication_service/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RefreshTokenRepository interface {
	SaveToken(ctx context.Context, token models.RefreshTokenWithHash) error
	GetToken(ctx context.Context, pairToken string) (models.RefreshTokenWithHash, error)
	DeleteToken(ctx context.Context, pairToken string) error
}

type refreshTokenRepository struct {
	db *mongo.Collection
}

func NewRefreshTokenRepository(db *mongo.Collection) RefreshTokenRepository {
	return refreshTokenRepository{
		db: db,
	}
}

func (r refreshTokenRepository) SaveToken(ctx context.Context, token models.RefreshTokenWithHash) error {
	_, err := r.db.InsertOne(ctx, token)
	return err
}

func (r refreshTokenRepository) GetToken(ctx context.Context, pairToken string) (models.RefreshTokenWithHash, error) {
	filter := bson.M{"pairToken": pairToken}

	var result models.RefreshTokenWithHash
	err := r.db.FindOne(ctx, filter).Decode(&result)
	return result, err
}

func (r refreshTokenRepository) DeleteToken(ctx context.Context, pairToken string) error {
	filter := bson.M{"pairToken": pairToken}
	_, err := r.db.DeleteOne(ctx, filter)
	return err
}
