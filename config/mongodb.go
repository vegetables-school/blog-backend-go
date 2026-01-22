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

type MongoConfig struct {
	URI      string
	Database string
	Client   *mongo.Client
}

var (
	mongoCfg  *MongoConfig
	mongoOnce sync.Once
)

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// InitMongo 初始化 MongoDB 配置和连接。
// 返回 MongoConfig 结构体（包含连接信息）和 error（初始化失败时返回错误）。
func InitMongo() (*MongoConfig, error) {
	var err error
	mongoOnce.Do(func() {
		uri := getEnv("MONGODB_URI", "mongodb://mongo:27017")
		db := getEnv("MONGODB_DATABASE", "blog")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		clientOpts := options.Client().ApplyURI(uri)
		client, e := mongo.Connect(ctx, clientOpts)
		if e != nil {
			err = e
			return
		}
		if e = client.Ping(ctx, nil); e != nil {
			err = e
			return
		}
		mongoCfg = &MongoConfig{
			URI:      uri,
			Database: db,
			Client:   client,
		}
	})
	if err != nil {
		return nil, fmt.Errorf("MongoDB 初始化失败: %w", err)
	}
	return mongoCfg, nil
}

// GetMongoClient 获取 MongoDB 客户端。
// 返回 *mongo.Client 客户端对象。
func GetMongoClient() *mongo.Client {
	if mongoCfg == nil {
		return nil
	}
	return mongoCfg.Client
}

// GetMongoDatabase 获取 MongoDB 数据库对象。
// 返回 *mongo.Database 数据库对象。
func GetMongoDatabase() *mongo.Database {
	if mongoCfg == nil || mongoCfg.Client == nil {
		return nil
	}
	return mongoCfg.Client.Database(mongoCfg.Database)
}

// CloseMongo 关闭 MongoDB 连接。
// 返回 error，关闭失败时返回错误，否则返回 nil。
func CloseMongo() error {
	if mongoCfg != nil && mongoCfg.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return mongoCfg.Client.Disconnect(ctx)
	}
	return nil
}
