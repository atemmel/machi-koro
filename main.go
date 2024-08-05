package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Code string

type GameState int

const (
	Lobby GameState = iota
	InProgress
)

type GameResponse struct {
	Code    Code      `json:"code"`
	State   GameState `json:"state"`
	Phase   string    `json:"phase"`
	Players []Client  `json:"players"`
}

var (
	upgrader = websocket.Upgrader{
		EnableCompression: true,
		CheckOrigin: func(r *http.Request) bool {
			return true //TODO: dev only
		},
	}
)

func getGame(c echo.Context) error {
	code := Code(c.Param("code"))
	if g := LookupGame(code); g != nil {
		return c.JSON(http.StatusOK, g)
	}
	return c.NoContent(http.StatusNotFound)
}

func postGame(c echo.Context) error {
	return c.JSON(http.StatusOK, NewGame())
}

func getCards(c echo.Context) error {
	return c.JSON(http.StatusOK, AllCards)
}

func joinWebsocket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	fmt.Println("New ws conn")
	if err != nil {
		return err
	}
	defer ws.Close()

	code, player, err := AddPlayer(ws)
	if err != nil {
		return err
	}

	// infinite ping-pong loop technology
	go ping(ws, code, player)

	for {
		msg, err := readMessage(ws)
		if err != nil {
			DropPlayer(ws, msg.Code)
			return err
		}

		switch msg.RequestOperation {
		case PongResponse: // ignore
		case LeaveRequest:
			{
				DropPlayer(ws, msg.Code)
				// breaking out of loop will close the connection
				// closing the connection will kill the ping routine aswell
				return nil
			}
		case StartRequest:
			{
				StartGame(ws, msg)
			}
		case RollRequest:
			{
				RollDice(ws, msg)
			}
		case BuyRequest:
			{
				BuyCard(ws, msg)
			}
		}
	}
}

func ping(ws *websocket.Conn, code Code, player string) {
	msg := &ServerResponse{
		Code:         code,
		Operands:     []Operand{Blank},
		ResponseCode: PingRequest,
		Player:       player,
	}

WAIT:
	time.Sleep(1 * time.Second)
	// send ping
	if err := sendMessage(ws, msg); err != nil {
		// stop pinging if error occurs
		DropPlayer(ws, code)
		return
	}
	goto WAIT
}

func main() {
	e := echo.New()
	e.Debug = true

	e.GET("/games/:code", getGame)
	e.POST("/games", postGame)
	e.GET("/cards", getCards)

	e.GET("/ws", joinWebsocket)

	defer DropAllPlayers()
	e.Logger.Fatal(e.Start(":1323"))
}
