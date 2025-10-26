package services

import (
	"context"
	"time"

	"blog/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// PostService 处理博客文章的业务逻辑
type PostService struct {
	collection *mongo.Collection
}

// NewPostService 创建新的PostService实例
func NewPostService(client *mongo.Client, dbName, collectionName string) *PostService {
	collection := client.Database(dbName).Collection(collectionName)
	return &PostService{
		collection: collection,
	}
}

// GetAllPosts 获取所有博客文章
func (s *PostService) GetAllPosts() ([]*models.Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []*models.Blog
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

// GetPostByID 根据ID获取单篇博客文章
func (s *PostService) GetPostByID(id string) (*models.Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var post models.Blog
	err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// CreatePost 创建新博客文章
func (s *PostService) CreatePost(title, content, author string) (*models.Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	post := &models.Blog{
		ID:        primitive.NewObjectID(),
		Title:     title,
		Content:   content,
		Author:    author,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := s.collection.InsertOne(ctx, post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

// UpdatePost 更新博客文章
func (s *PostService) UpdatePost(id string, title, content, author string) (*models.Blog, error) {
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
	var post models.Blog
	err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// DeletePost 删除博客文章
func (s *PostService) DeletePost(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = s.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

// InitializeSampleData 初始化示例数据
func (s *PostService) InitializeSampleData() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 检查是否已有数据
	count, err := s.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // 已有数据，跳过初始化
	}

	post := &models.Blog{
		ID:        primitive.NewObjectID(),
		Title:     "欢迎使用 Go 博客系统",
		Content:   "这是第一篇博客文章。Go 是一门很棒的语言！",
		Author:    "管理员",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = s.collection.InsertOne(ctx, post)
	return err
}
