package services

import (
	"context"
	"time"

	"blog/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// BlogService 处理博客文章的业务逻辑
type BlogService struct {
	collection *mongo.Collection
}

// NewBlogService 创建新的BlogService实例
func NewBlogService(client *mongo.Client, dbName, collectionName string) *BlogService {
	collection := client.Database(dbName).Collection(collectionName)
	return &BlogService{
		collection: collection,
	}
}

// GetAllBlogs 获取所有博客文章
func (s *BlogService) GetAllBlogs() ([]*models.Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var blogs []*models.Blog
	if err = cursor.All(ctx, &blogs); err != nil {
		return nil, err
	}
	return blogs, nil
}

// GetBlogByID 根据ID获取单篇博客文章
func (s *BlogService) GetBlogByID(id string) (*models.Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var blog models.Blog
	err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&blog)
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

// CreateBlog 创建新博客文章
func (s *BlogService) CreateBlog(title, content, author string) (*models.Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	blog := &models.Blog{
		ID:        primitive.NewObjectID(),
		Title:     title,
		Content:   content,
		Author:    author,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := s.collection.InsertOne(ctx, blog)
	if err != nil {
		return nil, err
	}
	return blog, nil
}

// UpdateBlog 更新博客文章
func (s *BlogService) UpdateBlog(id string, title, content, author string) (*models.Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"title":      title,
			"content":    content,
			"author":     author,
			"updated_at": time.Now(),
		},
	}

	_, err = s.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return nil, err
	}

	// 获取更新后的文档
	var blog models.Blog
	err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&blog)
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

// DeleteBlog 删除博客文章
func (s *BlogService) DeleteBlog(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = s.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}
