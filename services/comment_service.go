package services

import (
	"context"
	"time"

	"blog/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CommentService 评论业务逻辑
type CommentService struct {
	collection *mongo.Collection
}

func NewCommentService(client *mongo.Client, dbName, collectionName string) *CommentService {
	return &CommentService{
		collection: client.Database(dbName).Collection(collectionName),
	}
}

// CreateComment 新增评论
func (s *CommentService) CreateComment(blogID, userID primitive.ObjectID, username, content string) (*models.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	comment := &models.Comment{
		ID:        primitive.NewObjectID(),
		BlogID:    blogID,
		UserID:    userID,
		Username:  username,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err := s.collection.InsertOne(ctx, comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

// DeleteComment 删除评论（只能本人或管理员）
func (s *CommentService) DeleteComment(commentID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := s.collection.DeleteOne(ctx, bson.M{"_id": commentID})
	return err
}

// GetCommentByID 获取评论
func (s *CommentService) GetCommentByID(commentID primitive.ObjectID) (*models.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var comment models.Comment
	err := s.collection.FindOne(ctx, bson.M{"_id": commentID}).Decode(&comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}
