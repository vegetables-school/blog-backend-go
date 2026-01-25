package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateIndexes 为所有集合创建索引
func CreateIndexes(client *mongo.Client, dbName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db := client.Database(dbName)

	// 创建用户集合索引
	if err := createUserIndexes(ctx, db); err != nil {
		return err
	}

	// 创建博客集合索引
	if err := createBlogIndexes(ctx, db); err != nil {
		return err
	}

	// 创建评论集合索引
	if err := createCommentIndexes(ctx, db); err != nil {
		return err
	}

	// 创建点赞集合索引
	if err := createLikeIndexes(ctx, db); err != nil {
		return err
	}

	log.Println("所有数据库索引创建成功")
	return nil
}

// createUserIndexes 创建用户集合索引
func createUserIndexes(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("users")

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		log.Printf("创建用户索引失败: %v", err)
		return err
	}

	log.Println("用户集合索引创建成功")
	return nil
}

// createBlogIndexes 创建博客集合索引
func createBlogIndexes(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("blogs-dev")

	indexes := []mongo.IndexModel{
		// 创建时间索引（倒序，用于分页）
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
		// 作者索引
		{
			Keys: bson.D{{Key: "author", Value: 1}},
		},
		// 是否显示索引
		{
			Keys: bson.D{{Key: "show", Value: 1}},
		},
		// 标签索引（数组索引）
		{
			Keys: bson.D{{Key: "tags", Value: 1}},
		},
		// 全文搜索索引
		{
			Keys: bson.D{
				{Key: "title", Value: "text"},
				{Key: "content", Value: "text"},
			},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		log.Printf("创建博客索引失败: %v", err)
		return err
	}

	log.Println("博客集合索引创建成功")
	return nil
}

// createCommentIndexes 创建评论集合索引
func createCommentIndexes(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("comments")

	indexes := []mongo.IndexModel{
		// 博客ID索引
		{
			Keys: bson.D{{Key: "blog_id", Value: 1}},
		},
		// 用户ID索引
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
		// 创建时间索引
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		log.Printf("创建评论索引失败: %v", err)
		return err
	}

	log.Println("评论集合索引创建成功")
	return nil
}

// createLikeIndexes 创建点赞集合索引
func createLikeIndexes(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("likes")

	indexes := []mongo.IndexModel{
		// 唯一复合索引：同一用户对同一博客只能点赞一次
		{
			Keys: bson.D{
				{Key: "blog_id", Value: 1},
				{Key: "user_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		// 博客ID索引（用于统计点赞数）
		{
			Keys: bson.D{{Key: "blog_id", Value: 1}},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		log.Printf("创建点赞索引失败: %v", err)
		return err
	}

	log.Println("点赞集合索引创建成功")
	return nil
}
