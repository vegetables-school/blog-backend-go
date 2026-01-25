package services

import (
	"context"
	"time"

	"blog/config"
	"blog/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// LikeService 点赞业务逻辑
type LikeService struct {
	collection *mongo.Collection
}

func NewLikeService(client *mongo.Client, dbName, collectionName string) *LikeService {
	return &LikeService{
		collection: client.Database(dbName).Collection(collectionName),
	}
}

// AddLike 新增点赞
func (s *LikeService) AddLike(blogID, userID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.WriteTimeout)
	defer cancel()
	like := &models.Like{
		ID:        primitive.NewObjectID(),
		BlogID:    blogID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}
	_, err := s.collection.InsertOne(ctx, like)
	return err
}

// RemoveLike 取消点赞（只能本人或管理员）
func (s *LikeService) RemoveLike(blogID, userID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.WriteTimeout)
	defer cancel()
	_, err := s.collection.DeleteOne(ctx, bson.M{"blog_id": blogID, "user_id": userID})
	return err
}

// GetLikeByBlogAndUser 获取点赞
func (s *LikeService) GetLikeByBlogAndUser(blogID, userID primitive.ObjectID) (*models.Like, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.QuickQueryTimeout)
	defer cancel()
	var like models.Like
	err := s.collection.FindOne(ctx, bson.M{"blog_id": blogID, "user_id": userID}).Decode(&like)
	if err != nil {
		return nil, err
	}
	return &like, nil
}
