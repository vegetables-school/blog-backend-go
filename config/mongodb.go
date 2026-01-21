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

// MongoConfig 保存 MongoDB 配置信息
var (
	MongoURI      = getEnv("MONGODB_URI", "mongodb://localhost:27017")
	MongoDatabase = getEnv("MONGODB_DATABASE", "blog")
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
		clientOpts := options.Client().ApplyURI(MongoURI)
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

// GetMongoDatabase 获取 MongoDB 数据库对象
func GetMongoDatabase() *mongo.Database {
	return mongoClient.Database(MongoDatabase)
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
