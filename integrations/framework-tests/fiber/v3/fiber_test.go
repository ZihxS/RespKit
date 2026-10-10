package fiberv3

import (
	"net/http"
	"testing"

	"github.com/ZihxS/RespKit/integrations/framework-tests/internal/testkit"
	fiber "github.com/gofiber/fiber/v3"
)

func TestReal(t *testing.T) {
	app := fiber.New()
	app.Get("/ok", func(c fiber.Ctx) error { return testkit.Execute("ok", c) })
	app.Get("/error", func(c fiber.Ctx) error { return testkit.Execute("error", c) })
	app.Get("/raw", func(c fiber.Ctx) error { return testkit.Execute("raw", c) })
	app.Get("/debug", func(c fiber.Ctx) error { return testkit.Execute("debug", c) })
	app.Post("/created", func(c fiber.Ctx) error { return testkit.Execute("created", c) })
	app.Get("/empty", func(c fiber.Ctx) error { return testkit.Execute("empty", c) })
	testkit.TestResponse(t, "fiber-v3", func(r *http.Request) (*http.Response, error) { return app.Test(r) })
}
