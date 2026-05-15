package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"Prak_9/internal/cache"
	"Prak_9/internal/config"
	"Prak_9/internal/task"
	"github.com/redis/go-redis/v9"
)

var listKey = "tasks:list"

type TaskService struct {
	repo  *task.Repo
	redis *redis.Client
	cfg   config.Config
}

func NewTaskService(repo *task.Repo, redisClient *redis.Client, cfg config.Config) *TaskService {
	return &TaskService{repo: repo, redis: redisClient, cfg: cfg}
}

func (s *TaskService) GetTasks(ctx context.Context) ([]task.Task, error) {
	if s.redis != nil {
		cached, err := s.redis.Get(ctx, listKey).Result()
		if err == nil {
			var t []task.Task
			if err := json.Unmarshal([]byte(cached), &t); err == nil {
				log.Println("cache hit:", listKey)
				return t, nil
			}
			log.Println("cache decode error:", err)
		} else if !errors.Is(err, redis.Nil) {
			log.Println("redis read error:", err)
		} else {
			log.Println("cache miss:", listKey)
		}
	}

	t, err := s.repo.GetAllTasks()
	if err != nil {
		return []task.Task{}, err
	}

	if s.redis != nil {
		bytes, err := json.Marshal(t)
		if err != nil {
			log.Println("cache encode error:", err)
			return t, nil
		}

		ttl := cache.TTLWithJitter(s.cfg.CacheTTL, s.cfg.CacheTTLJitter)
		if err := s.redis.Set(ctx, listKey, bytes, ttl).Err(); err != nil {
			log.Println("redis write error:", err)
		} else {
			log.Println("cache set:", listKey, "ttl:", ttl)
		}
	}

	return t, nil
}

func (s *TaskService) GetTaskByID(ctx context.Context, id int64) (task.Task, error) {
	key := cache.TaskByIDKey(id)

	if s.redis != nil {
		cached, err := s.redis.Get(ctx, key).Result()
		if err == nil {
			var t task.Task
			if err := json.Unmarshal([]byte(cached), &t); err == nil {
				log.Println("cache hit:", key)
				return t, nil
			}
			log.Println("cache decode error:", err)
		} else if !errors.Is(err, redis.Nil) {
			log.Println("redis read error:", err)
		} else {
			log.Println("cache miss:", key)
		}
	}

	t, err := s.repo.GetByID(id)
	if err != nil {
		return task.Task{}, err
	}

	if s.redis != nil {
		bytes, err := json.Marshal(t)
		if err != nil {
			log.Println("cache encode error:", err)
			return t, nil
		}

		ttl := cache.TTLWithJitter(s.cfg.CacheTTL, s.cfg.CacheTTLJitter)
		if err := s.redis.Set(ctx, key, bytes, ttl).Err(); err != nil {
			log.Println("redis write error:", err)
		} else {
			log.Println("cache set:", key, "ttl:", ttl)
		}
	}

	return t, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, t task.Task) error {
	if err := s.repo.Update(t); err != nil {
		return err
	}

	if s.redis != nil {
		key := cache.TaskByIDKey(t.ID)
		if err := s.redis.Del(ctx, key).Err(); err != nil {
			log.Println("redis delete error:", err)
		} else {
			log.Println("cache invalidated:", key)
		}
	}

	return nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id int64) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	if s.redis != nil {
		key := cache.TaskByIDKey(id)
		if err := s.redis.Del(ctx, key).Err(); err != nil {
			log.Println("redis delete error:", err)
		} else {
			log.Println("cache invalidated:", key)
		}
	}

	return nil
}
