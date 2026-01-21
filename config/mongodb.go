package config

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
	mongoOnce   sync.Once
)

// getEnv 获取环境变量，若不存在则返回默认值
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// InitMongo 初始化 MongoDB 连接（只执行一次）
func InitMongo() (*mongo.Client, error) {
	var err error
	mongoOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		mongoURI := getEnv("MONGODB_URI", "mongodb://localhost:27017")
		clientOpts := options.Client().ApplyURI(mongoURI)
		mongoClient, err = mongo.Connect(ctx, clientOpts)
		if err != nil {
			return
		}
		err = mongoClient.Ping(ctx, nil)
	})
	if err != nil {
		return nil, fmt.Errorf("MongoDB 初始化失败: %w", err)
	}
	return mongoClient, nil
}

// GetMongoClient 获取 MongoDB 客户端
func GetMongoClient() *mongo.Client {
	return mongoClient
}

// GetMongoDatabase 获取 MongoDB 数据库对象（每次都从环境变量读取数据库名，支持动态切换）
func GetMongoDatabase() *mongo.Database {
	dbName := getEnv("MONGODB_DATABASE", "blog")
	return mongoClient.Database(dbName)
}

// CloseMongo 关闭 MongoDB 连接
func CloseMongo() error {
	if mongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return mongoClient.Disconnect(ctx)
	}
	return nil
}
