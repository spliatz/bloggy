package cache

import (
	"context"
	"log"
	"sync"

	"github.com/spliatz/bloggy-backend/internal/domain/entity"
	"github.com/spliatz/bloggy-backend/pkg/errors"
)

type UserStorage interface {
  GetAll(ctx context.Context) ([]entity.User, error)
}

type userCache struct {
  m               sync.RWMutex
	userKeyUsername map[string]entity.User
	userKeyId       map[int]entity.User
	queueChan       chan entity.User
	userStorage     UserStorage
}

func NewUserCache(userStorage UserStorage) *userCache {
	return &userCache{
		userKeyUsername: make(map[string]entity.User),
		userKeyId:       make(map[int]entity.User),
		queueChan:       make(chan entity.User),
		userStorage:     userStorage,
	}
}

func (uc *userCache) Init() *userCache {
  return uc.init()	
}

func (uc *userCache) init() *userCache {
	users, err := uc.userStorage.GetAll(context.Background())
	if err != nil {
		log.Fatal(err)
	}
  
  uc.m.Lock()
  defer uc.m.Unlock()
	for _, user := range users {
		uc.userKeyUsername[user.Username] = user
		uc.userKeyId[user.Id] = user
	}

  return uc
}

func (uc *userCache) Set(ctx context.Context, user entity.User) {
  go uc.set(ctx, user)
}

func (uc *userCache) set(ctx context.Context, user entity.User) {
  uc.m.Lock()
  defer uc.m.Unlock()
	uc.userKeyUsername[user.Username] = user
	uc.userKeyId[user.Id] = user
}

func (uc *userCache) GetById(ctx context.Context, id int) (entity.User, error) {
  uc.m.RLock()
  defer uc.m.RUnlock()
	user, ok := uc.userKeyId[id]
	if !ok {
		return entity.User{}, errors.ErrUserNotFound
	}
	return user, nil
}

func (uc *userCache) GetByUsername(ctx context.Context, username string) (entity.User, error) {
  uc.m.RLock()
  defer uc.m.RUnlock()
	user, ok := uc.userKeyUsername[username]
	if !ok {
		return entity.User{}, errors.ErrUserNotFound
	}
	return user, nil
}
