package services

import (
	"sync"
	"time"

	"blog/models"
)

// PostService 处理博客文章的业务逻辑
type PostService struct {
	posts   map[int]*models.Blog
	nextID  int
	postsMu sync.RWMutex
}

// NewPostService 创建新的PostService实例
func NewPostService() *PostService {
	return &PostService{
		posts:  make(map[int]*models.Blog),
		nextID: 1,
	}
}

// GetAllPosts 获取所有博客文章
func (s *PostService) GetAllPosts() []*models.Blog {
	s.postsMu.RLock()
	defer s.postsMu.RUnlock()

	postList := make([]*models.Blog, 0, len(s.posts))
	for _, post := range s.posts {
		postList = append(postList, post)
	}
	return postList
}

// GetPostByID 根据ID获取单篇博客文章
func (s *PostService) GetPostByID(id int) (*models.Blog, bool) {
	s.postsMu.RLock()
	defer s.postsMu.RUnlock()

	post, exists := s.posts[id]
	return post, exists
}

// CreatePost 创建新博客文章
func (s *PostService) CreatePost(title, content, author string) *models.Blog {
	s.postsMu.Lock()
	defer s.postsMu.Unlock()

	post := &models.Blog{
		ID:        s.nextID,
		Title:     title,
		Content:   content,
		Author:    author,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	s.posts[s.nextID] = post
	s.nextID++
	return post
}

// UpdatePost 更新博客文章
func (s *PostService) UpdatePost(id int, title, content, author string) (*models.Blog, bool) {
	s.postsMu.Lock()
	defer s.postsMu.Unlock()

	post, exists := s.posts[id]
	if !exists {
		return nil, false
	}

	post.Title = title
	post.Content = content
	post.Author = author
	post.UpdatedAt = time.Now()
	return post, true
}

// DeletePost 删除博客文章
func (s *PostService) DeletePost(id int) bool {
	s.postsMu.Lock()
	defer s.postsMu.Unlock()

	_, exists := s.posts[id]
	if !exists {
		return false
	}

	delete(s.posts, id)
	return true
}

// InitializeSampleData 初始化示例数据
func (s *PostService) InitializeSampleData() {
	s.posts[1] = &models.Blog{
		ID:        1,
		Title:     "欢迎使用 Go 博客系统",
		Content:   "这是第一篇博客文章。Go 是一门很棒的语言！",
		Author:    "管理员",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	s.nextID = 2
}