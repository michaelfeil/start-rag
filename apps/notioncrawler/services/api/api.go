package api

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"net/http"
	"notioncrawl/services/api/controller"
	"notioncrawl/services/crawler"
	"notioncrawl/services/state"
)

func Run(state *state.Manager, neo4jOptions crawler.Neo4jOptions, addr string) {
	neo4j, err := neo4j.NewDriverWithContext(neo4jOptions.Address, neo4j.BasicAuth(neo4jOptions.Username, neo4jOptions.Password, ""))
	if err != nil {
		panic(err)
	}

	ctrl := controller.New(neo4j)

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Notion Crawler")
	})

	app.Get("/state", func(c *fiber.Ctx) error {
		s, err := json.Marshal(state.GetState())
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}
		return c.Send(s)
	})

	app.Post("/db/purge", ctrl.PurgeDb)

	log.Fatal(app.Listen(addr))
}
