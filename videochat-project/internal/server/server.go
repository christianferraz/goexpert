package server

import (
	"flag"
	"time"

	"os"

	"github.com/christianferraz/goexpert/videochat-project/internal/handlers"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/template/html/v2"
	"gorm.io/gorm/logger"
)

var (
	addr = flag.String("addr", "", os.Getenv("PORT"))
	cert = flag.String("cert", "", "")
	key  = flag.String("key", "", "")
)

func Run() error {
	flag.Parse()
	if *addr == ":" {
		*addr = ":8080"
	}
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Use(logger.New())
	app.Use(cors.New())
	app.Get("/", handlers.Welcome)
	app.Get("/room/create", handlers.RoomCreate)
	app.Get("/room/:uuid", handlers.Room)
	app.Get("/room/:uuid/websocket", websocket.New(handlers.RoomWebsocket), websocket.Config{
		HandshakeTimeout: 20 * time.Second,
	})
	app.Get("/room/:uuid/chat", handlers.RoomChat)
	app.Get("/room/:uuid/chat/websocket", websocket.New(handlers.RoomChatWebsocket))
	app.Get("/room/:uuid/viewer/websocket", websocket.New(handlers.RoomViewerWebsocket))
	app.Get("/room/:suuid", handlers.Stream)
	app.Get("/room/:suuid/websocket", websocket.New(handlers.StreamWebsocket))
	app.Get("/room/:suuid/chat/websocket", websocket.New(handlers.StreamWebsocket))
	app.Get("/room/:suuid/viewer/websocket", websocket.New(handlers.StreamWebsocket))
}
