package services

import (
	"context"
	"errors"
	"time"

	"blog/models"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// AuthService 处理用户认证的业务逻辑
type AuthService struct {
	collection *mongo.Collection
	jwtSecret  []byte
}

// NewAuthService 创建新的AuthService实例
func NewAuthService(client *mongo.Client, dbName, collectionName string, jwtSecret string) *AuthService {
	collection := client.Database(dbName).Collection(collectionName)
	return &AuthService{
		collection: collection,
		jwtSecret:  []byte(jwtSecret),
	}
}

// Register 用户注册
func (s *AuthService) Register(username, password, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 检查用户名是否已存在
	var existingUser models.User
	err := s.collection.FindOne(ctx, bson.M{"username": username}).Decode(&existingUser)
	if err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	err = s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&existingUser)
	if err == nil {
		return nil, errors.New("邮箱已存在")
	}

	// 哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:  username,
		Password:  string(hashedPassword),
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = result.InsertedID.(primitive.ObjectID)

	return user, nil
}

// Login 用户登录
func (s *AuthService) Login(username, password string) (*models.AuthResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := s.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成 JWT token
	token, err := s.generateToken(user.ID.Hex(), user.Username)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

// ValidateToken 验证 JWT token
func (s *AuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

// GetUserByID 根据ID获取用户
func (s *AuthService) GetUserByID(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// generateToken 生成 JWT token
func (s *AuthService) generateToken(userID, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // 24小时过期
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
