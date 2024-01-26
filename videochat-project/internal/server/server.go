package server

import (
	"flag"

	"github.com/christianferraz/goexpert/videochat-project/internal/handlers"
)

var (
	// addr = flag.String("addr,":"", os.Getenv("PORT"), "")
	cert = flag.String("cert", "", "")
	key  = flag.String("key", "", "")
)

func Run() error {
	flag.Parse()
	if *addr == ":" {
		*addr = ":8080"
	}
	app.Get("/", handlers.Welcome)
	app.Get("/room/create", handlers.RoomCreate)
	app.Get("/room/:uuid", handlers.Room)
	app.Get("/room/:uuid/websocket", handlers.Room)
	app.Get("/room/:uuid/chat", handlers.RoomChat)
	app.Get("/room/:uuid/chat/websocket", websocket.New(handlers.RoomChatWebsocket))
	app.Get("/room/:uuid/viewer/websocket", websocket.New(handlers.RoomViewerWebsocket))
	app.Get("/room/:suuid", handlers.Stream)
	app.Get("/room/:suuid/websocket", websocket.New(handlers.StreamWebsocket))
	app.Get("/room/:suuid/chat/websocket", websocket.New(handlers.StreamWebsocket))
	app.Get("/room/:suuid/viewer/websocket", websocket.New(handlers.StreamWebsocket))
}
