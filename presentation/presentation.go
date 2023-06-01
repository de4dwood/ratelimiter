package presentation

import (
	"fmt"
	"rt/domain"
	"rt/lib"

	"github.com/gofiber/fiber/v2"
)

type httpApi struct {
	engine *fiber.App
	app    domain.Domain
}

func NewHttp(s domain.Domain) *httpApi {
	return &httpApi{
		engine: fiber.New(),
		app:    s,
	}
}

func (H *httpApi) Start(listen string) {
	H.engine.Get("/request", H.Request)
	H.engine.Post("/config", H.AddConfig)
	H.engine.Delete("/config", H.DeleteConfig)
	H.engine.Get("/config", H.GetConfig)
	H.engine.Listen(listen)

}
func (H *httpApi) Request(c *fiber.Ctx) error {
	URL, err := lib.UrlParser(c.Query("url"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "url is not correct"})
	}
	user := c.Query("user")
	if URL == "" {
		return c.Status(400).JSON(fiber.Map{"error": "plz send URL address"})
	}
	if user == "" {
		return c.Status(400).JSON(fiber.Map{"error": "plz send user"})
	}
	fmt.Println(&c.Request().Header)
	permit, err := H.app.Request(c.Context(), user, URL, c.Context().Time())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"there is an error ": err.Error()})
	}
	if permit {
		return c.SendStatus(200)
	} else {
		return c.SendStatus(429)

	}
}

func (H *httpApi) AddConfig(c *fiber.Ctx) error {
	URL, err := lib.UrlParser(c.Query("url"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "url is not correct"})
	}
	if URL == "" {
		return c.Status(400).JSON(fiber.Map{"error": "plz send URL address"})
	}
	id := c.Query("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "plz send config id"})
	}
	limit := c.QueryInt("limit")
	if limit == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "plz send limit count"})
	}
	window := c.Query("window")
	if window == "" {
		return c.Status(400).JSON(fiber.Map{"error": "plz send window size"})
	}
	config := domain.Config{
		Limit:  limit,
		URL:    URL,
		Window: window,
		Id:     id,
	}
	err = H.app.AddConfig(c.Context(), URL, id, config)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"there is an error ": err.Error()})
	}
	return c.SendStatus(200)

}
func (H *httpApi) DeleteConfig(c *fiber.Ctx) error {
	URL, err := lib.UrlParser(c.Query("url"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "url is not correct"})
	}
	if URL == "" {
		return c.Status(400).JSON(fiber.Map{"error": "plz send URL address"})
	}
	id := c.Query("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "plz send config id"})
	}

	err = H.app.DeleteConfig(c.Context(), URL, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"there is an error ": err.Error()})
	}
	return c.SendStatus(200)

}

func (H *httpApi) GetConfig(c *fiber.Ctx) error {
	var outs []map[string]interface{}
	URL, err := lib.UrlParser(c.Query("url"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "url is not correct"})
	}
	id := c.Query("id")
	if URL != "" && id != "" {
		config, err := H.app.GetApiConfig(c.Context(), URL, id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"there is an error ": err.Error()})
		}
		return c.Status(200).JSON(fiber.Map{"results ": lib.STM(config)})

	} else if URL != "" {
		configs, err := H.app.GetApiConfigs(c.Context(), URL)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"there is an error ": err.Error()})
		}
		for _, config := range configs {
			outs = append(outs, lib.STM(config))
		}

	} else {
		configs, err := H.app.GetAllConfigs(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"there is an error ": err.Error()})
		}
		for _, config := range configs {
			outs = append(outs, lib.STM(config))
		}

	}
	return c.Status(200).JSON(fiber.Map{"results ": outs})

}
