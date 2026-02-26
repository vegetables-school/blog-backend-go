package services

import (
	"context"
	"time"

	"blog/config"
	"blog/models"
	"blog/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// GetBlogsByTagsWithPagination 分页获取包含所有指定标签的博客文章
func (s *BlogService) GetBlogsByTagsWithPagination(tags []string, page, limit int64) ([]*models.Blog, int64, error) {
	ctx, cancel := utils.NewDefaultContext()
	defer cancel()

	filter := bson.M{}
	if len(tags) > 0 {
		filter["tags"] = bson.M{"$all": tags}
	}

	skip := utils.CalculateSkip(page, limit)

	// 当有筛选条件时，使用精确计数
	var total int64
	var err error
	if len(tags) > 0 {
		total, err = s.collection.CountDocuments(ctx, filter)
	} else {
		total, err = s.collection.EstimatedDocumentCount(ctx)
	}
	if err != nil {
		utils.Error("获取博客总数失败", zap.Error(err))
		return nil, 0, err
	}

	cursor, err := s.collection.Find(ctx, filter, &options.FindOptions{
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
	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
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
	ctx, cancel := utils.NewDefaultContext()
	defer cancel()

	skip := utils.CalculateSkip(page, limit)

	// 使用近似计数，性能更好
	total, err := s.collection.EstimatedDocumentCount(ctx)
	if err != nil {
		utils.Error("获取博客总数失败", zap.Error(err))
		return nil, 0, err
	}

	cursor, err := s.collection.Find(ctx, bson.M{"show": true}, &options.FindOptions{
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
	ctx, cancel := context.WithTimeout(context.Background(), config.QuickQueryTimeout)
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
	ctx, cancel := context.WithTimeout(context.Background(), config.WriteTimeout)
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
	ctx, cancel := context.WithTimeout(context.Background(), config.WriteTimeout)
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

	if tags != nil {
		setFields["tags"] = tags
	}

	if show != nil {
		setFields["show"] = *show
	}

	if views != nil {
		setFields["views"] = *views
	}

	update := bson.M{"$set": setFields}

	// 使用 FindOneAndUpdate 避免二次查询
	var blog models.Blog
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err = s.collection.FindOneAndUpdate(ctx, bson.M{"_id": objID}, update, opts).Decode(&blog)
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

// DeleteBlog 删除博客文章
func (s *BlogService) DeleteBlog(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.WriteTimeout)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = s.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

// SearchBlogs 支持按关键字和标签筛选博客
func (s *BlogService) SearchBlogs(keyword string, tags []string, page, limit int64) ([]*models.Blog, int64, error) {
	ctx, cancel := utils.NewSlowQueryContext()
	defer cancel()

	filter := bson.M{}
	if keyword != "" {
		// 清理关键字，防止正则注入
		safeKeyword := utils.SanitizeSearchKeyword(keyword)
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": safeKeyword, "$options": "i"}},
			{"content": bson.M{"$regex": safeKeyword, "$options": "i"}},
		}
	}
	if len(tags) > 0 {
		filter["tags"] = bson.M{"$all": tags}
	}

	skip := utils.CalculateSkip(page, limit)
	total, err := s.collection.CountDocuments(ctx, filter)
	if err != nil {
		utils.Error("搜索博客计数失败", zap.Error(err))
		return nil, 0, err
	}

	cursor, err := s.collection.Find(ctx, filter, &options.FindOptions{
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

// GetAllTags 获取所有标签（去重）
func (s *BlogService) GetAllTags() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
	defer cancel()

	tags, err := s.collection.Distinct(ctx, "tags", bson.M{})
	if err != nil {
		return nil, err
	}

	var result []string
	for _, t := range tags {
		if str, ok := t.(string); ok {
			result = append(result, str)
		}
	}
	return result, nil
}

// AddTag 给所有包含某标签的博客添加新标签（或单独管理标签集合时用）
// 注意：此方法已弃用，标签请直接在创建或编辑博客时管理
func (s *BlogService) AddTag(tag string) error {
	// 标签管理已集成在博客创建/编辑中，不需要单独的标签集合
	return nil
}

// DeleteTag 删除所有博客中的某个标签
func (s *BlogService) DeleteTag(tag string) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.WriteTimeout)
	defer cancel()

	_, err := s.collection.UpdateMany(ctx, bson.M{"tags": tag}, bson.M{"$pull": bson.M{"tags": tag}})
	return err
}
