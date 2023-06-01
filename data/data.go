package data

import (
	"context"
	"time"
)

type Data interface {
	Check(ctx context.Context, user, api, td string, t time.Time) (int64, error)
	Request(ctx context.Context, user, api string, t time.Time) error
	AddConfig(ctx context.Context, api, id string, c map[string]interface{}) error
	UpdateConfig(ctx context.Context, api, id string, c map[string]interface{}) error
	GetConfig(ctx context.Context, api, id string) (map[string]string, error)
	GetConfigs(ctx context.Context, api string) ([]map[string]string, error)
	GetAllConfigs(ctx context.Context) ([]map[string]string, error)
	DeleteConfig(ctx context.Context, api, id string) error
}
