package bootstrap

import (
	"io"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/kooroshh/fiber-boostrap/app/ws"
	"github.com/kooroshh/fiber-boostrap/pkg/database"
	"github.com/kooroshh/fiber-boostrap/pkg/env"
	"github.com/kooroshh/fiber-boostrap/pkg/router"
	"go.elastic.co/apm"
)

func NewApplication() *fiber.App {
	env.SetupEnvFile()
	setupLogFile()

	database.SetupDatabase()
	database.SetupMongoDB()

	apm.DefaultTracer.Service.Name = "messaging-app"
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Use(recover.New())
	app.Use(logger.New())
	app.Get("/dashboard", monitor.New())

	// setup websocket
	go ws.ServeWSMessaging(app)

	// setup router
	router.InstallRouter(app)

	return app
}

func setupLogFile() {
	logFile, err := os.OpenFile("./logs/messaging_app.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}
