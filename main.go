package main

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Code string

type GameState int

const (
	Lobby GameState = iota
	InProgress
)

type Player struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

type Game struct {
	Code    Code      `json:"code"`
	State   GameState `json:"state"`
	Players []Player  `json:"players"`
}

var (
	games = []Game{
		{
			Code:  "a",
			State: Lobby,
			Players: []Player{
				{
					Id:   0,
					Name: "Jocke",
				},
			},
		},
	}

	topicToClients = map[Code]map[*websocket.Conn]bool{}

	upgrader = websocket.Upgrader{} // use default options
)

func queryGame(code Code) *Game {
	//TODO: validate code
	for i, g := range games {
		if g.Code == code {
			return &games[i]
		}
	}
	return nil
}

func makeGame() Game {
	return Game{
		Code:    "a",
		State:   Lobby,
		Players: []Player{},
	}
}

func getGame(c echo.Context) error {
	code := Code(c.Param("code"))
	if g := queryGame(code); g != nil {
		return c.JSON(http.StatusOK, g)
	}
	return c.NoContent(http.StatusNotFound)
}

func postGame(c echo.Context) error {
	g := makeGame()
	return c.JSON(http.StatusOK, &g)
}

func joinWebsocket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {

	}
}

func main() {
	e := echo.New()

	e.GET("/games/:code", getGame)
	e.GET("/ws", joinWebsocket)

	e.Logger.Fatal(e.Start(":1323"))
}
