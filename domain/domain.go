package domain

import (
	"context"
	"rt/data"
	"rt/lib"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Id     string `json:"id"`
	URL    string `json:"URL"`
	Window string `json:"window"`
	Limit  int    `json:"limit"`
}

type Domain interface {
	Request(ctx context.Context, user, URL string, t time.Time) (bool, error)
	AddConfig(ctx context.Context, URL, id string, config Config) error
	GetApiConfig(context.Context, string, string) (Config, error)
	GetApiConfigs(ctx context.Context, URL string) ([]Config, error)
	DeleteConfig(ctx context.Context, URL, id string) error
	GetAllConfigs(ctx context.Context) ([]Config, error)
	UpdateConfig(ctx context.Context, URL, id string, config Config) error
}

type domain struct {
	data data.Data
}

func NewDomain(d data.Data) Domain {
	return &domain{
		data: d,
	}
}

func (s *domain) Request(ctx context.Context, user, URL string, t time.Time) (bool, error) {
	urlpart := strings.Split(URL, "/")
	var configs []Config
	for i := 1; i <= len(urlpart); i++ {
		var cfgs []Config
		URLn := strings.Join(urlpart[:i], "/")
		cfgs, err := s.GetApiConfigs(ctx, URLn)
		if err != nil {
			return false, err
		}
		configs = append(configs, cfgs...)
	}
	if len(configs) == 0 {
		return true, nil
	}
	for _, config := range configs {
		count, err := s.data.Check(ctx, user, config.URL, config.Window, t)
		if err != nil {
			return false, err
		}
		if count >= int64(config.Limit) {
			return false, nil
		}
	}
	for _, config := range configs {
		err := s.data.Request(ctx, user, config.URL, t)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (s *domain) GetApiConfigs(ctx context.Context, URL string) ([]Config, error) {
	results, err := s.data.GetConfigs(ctx, URL)
	var configs []Config
	if err != nil {
		return nil, err
	}
	for _, result := range results {
		var config Config
		config.URL = result["URL"]
		config.Id = result["id"]
		config.Window = result["window"]
		config.Limit, err = strconv.Atoi(result["limit"])
		if err != nil {
			continue
		}
		configs = append(configs, config)

	}
	return configs, nil

}
func (s *domain) AddConfig(ctx context.Context, URL, id string, config Config) error {
	c := lib.STM(config)
	err := s.data.AddConfig(ctx, URL, id, c)
	if err != nil {
		return err
	}
	return nil
}
func (s *domain) UpdateConfig(ctx context.Context, URL, id string, config Config) error {
	c := lib.STM(config)
	err := s.data.UpdateConfig(ctx, URL, id, c)
	if err != nil {
		return err
	}
	return nil
}
func (s *domain) DeleteConfig(ctx context.Context, URL, id string) error {
	err := s.data.DeleteConfig(ctx, URL, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *domain) GetAllConfigs(ctx context.Context) ([]Config, error) {
	results, err := s.data.GetAllConfigs(ctx)
	var configs []Config
	if err != nil {
		return nil, err
	}
	for _, result := range results {
		var config Config
		config.URL = result["URL"]
		config.Id = result["id"]
		config.Window = result["window"]
		config.Limit, err = strconv.Atoi(result["limit"])
		if err != nil {
			continue
		}

		configs = append(configs, config)

	}
	return configs, nil

}

func (s *domain) GetApiConfig(ctx context.Context, URL, id string) (Config, error) {
	var config Config
	result, err := s.data.GetConfig(ctx, URL, id)
	if err != nil {
		return config, err
	}
	config.URL = result["URL"]
	config.Id = result["id"]
	config.Window = result["window"]
	config.Limit, err = strconv.Atoi(result["limit"])
	if err != nil {
		return config, err
	}
	return config, nil

}
