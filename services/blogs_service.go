package services

import (
	"context"
	"time"

	"blog/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// GetBlogsWithPagination 分页获取博客文章
func (s *BlogService) GetBlogsWithPagination(page, limit int64) ([]*models.Blog, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 计算跳过的文档数
	skip := (page - 1) * limit

	// 获取总数
	total, err := s.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	// 分页查询，按创建时间倒序
	cursor, err := s.collection.Find(ctx, bson.M{}, &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
		Sort:  bson.M{"created_at": -1},
	})
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var blogs []*models.Blog
	if err = cursor.All(ctx, &blogs); err != nil {
		return nil, 0, err
	}
	return blogs, total, nil
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
func (s *BlogService) CreateBlog(title, content, author string, tags []string, show bool) (*models.Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	blog := &models.Blog{
		ID:        primitive.NewObjectID(),
		Title:     title,
		Content:   content,
		Author:    author,
		Tags:      tags,
		Views:     0,
		Show:      show,
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
func (s *BlogService) UpdateBlog(id string, title, content, author *string, tags []string, show *bool, views *int64) (*models.Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	setFields := bson.M{
		"updated_at": time.Now(),
	}

	if title != nil {
		setFields["title"] = *title
	}
	if content != nil {
		setFields["content"] = *content
	}
	if author != nil {
		setFields["author"] = *author
	}

	// tags == nil -> don't change tags; empty slice -> clear tags
	if tags != nil {
		setFields["tags"] = tags
	}

	if show != nil {
		setFields["show"] = *show
	}

	// views == nil -> don't change; non-nil -> set to provided value
	if views != nil {
		setFields["views"] = *views
	}

	update := bson.M{"$set": setFields}

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
