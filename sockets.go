package main

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/gorilla/websocket"
)

type RequestOperation int
type ResponseCode int

const (
	JoinRequest RequestOperation = iota
	PongResponse
	LeaveRequest
	StartRequest
)

const (
	JoinAnnouncement ResponseCode = iota
	PingRequest
	LeaveAnnouncement
	StartAnnouncement
	OwnerAssign
	TurnChangeAnnouncement
)

type Operand int

const (
	Blank Operand = 0
)

type ClientRequest struct {
	Code
	Operand
	RequestOperation
	Player string
}

type ServerResponse struct {
	Code         `json:"code"`
	Operand      `json:"operand"`
	ResponseCode `json:"responseCode"`
	Player       string `json:"player"`
}

type Client struct {
	Name string `json:"name"`
}

type Game struct {
	Code         Code                       // Code to game
	Owner        *websocket.Conn            // Owner of lobby
	ActivePlayer *websocket.Conn            // Player whose turn it is
	State        GameState                  // In progess, Lobby
	clients      map[*websocket.Conn]Client // All clients
	order        []*websocket.Conn
}

func (g *Game) Response() *GameResponse {
	players := make([]Client, 0, len(g.clients))
	for _, v := range g.clients {
		if v.Name != "" {
			players = append(players, v)
		}
	}
	return &GameResponse{
		Code:    g.Code,
		State:   g.State,
		Players: players,
	}
}

var (
	gamesMutex = sync.Mutex{}
	games      = map[Code]Game{}
)

func AddPlayer(ws *websocket.Conn) (Code, string, error) {
	msg, err := readMessage(ws)
	if err != nil {
		return "", "", err
	}

	fmt.Println("joined with", msg)

	{
		gamesMutex.Lock()
		defer gamesMutex.Unlock()

		lookup := games[msg.Code]
		if len(lookup.clients) == 0 {
			lookup.clients = map[*websocket.Conn]Client{}
			// first to join becomes owner
			fmt.Println("owner assigned")
			lookup.Owner = ws
			sendMessage(ws, &ServerResponse{
				Code:         msg.Code,
				Operand:      Blank,
				ResponseCode: OwnerAssign,
				Player:       msg.Player,
			})
		}

		lookup.clients[ws] = Client{
			Name: msg.Player,
		}
		games[msg.Code] = lookup

		response := &ServerResponse{
			Code:         msg.Code,
			Operand:      Blank,
			ResponseCode: JoinAnnouncement,
			Player:       msg.Player,
		}

		for client := range lookup.clients {
			sendMessage(client, response)
		}
	}

	return msg.Code, msg.Player, nil
}

func DropPlayer(ws *websocket.Conn, code Code) {
	gamesMutex.Lock()
	defer gamesMutex.Unlock()
	lookup, ok := games[code]
	if !ok {
		return
	}

	player := lookup.clients[ws]

	if lookup.Owner == ws {
		if len(lookup.clients) <= 1 {
			// last player, just toss the game
			delete(games, code)
			return
		} else {
			// make someone else the owner
			for c := range lookup.clients {
				if c != ws {
					lookup.Owner = c
					sendMessage(ws, &ServerResponse{
						Code:         code,
						Operand:      Blank,
						ResponseCode: OwnerAssign,
						Player:       lookup.clients[c].Name,
					})
					break
				}
			}
			delete(lookup.clients, ws)
			lookup.order = remove(lookup.order, ws)
			games[code] = lookup
		}
	}

	publishMessage(code, &ServerResponse{
		Code:         code,
		Operand:      Blank,
		ResponseCode: LeaveAnnouncement,
		Player:       player.Name,
	})
}

func StartGame(ws *websocket.Conn, msg *ClientRequest) {
	gamesMutex.Lock()
	defer gamesMutex.Unlock()

	lookup := games[msg.Code]

	if ws != lookup.Owner {
		fmt.Println("was not owner", unsafe.Pointer(ws), lookup.Owner)
		// only owner can start
		return
	}

	// just get some guy
	for c := range lookup.clients {
		lookup.ActivePlayer = c
		break
	}

	// make an order of turns
	lookup.order = make([]*websocket.Conn, len(lookup.clients))
	for c := range lookup.clients {
		lookup.order = append(lookup.order, c)
	}

	lookup.State = InProgress
	games[msg.Code] = lookup

	publishMessage(msg.Code, &ServerResponse{
		Player:       msg.Player,
		Code:         msg.Code,
		ResponseCode: StartAnnouncement,
		Operand:      Blank,
	})

	activePlayer := lookup.clients[lookup.ActivePlayer]

	publishMessage(msg.Code, &ServerResponse{
		Player:       activePlayer.Name,
		Code:         msg.Code,
		ResponseCode: TurnChangeAnnouncement,
		Operand:      Blank,
	})
}

func DropAllPlayers() {
	gamesMutex.Lock()
	defer gamesMutex.Unlock()
	for _, g := range games {
		for c := range g.clients {
			c.Close()
		}
	}
}

// assumes gamesMutex to be locked
func publishMessage(code Code, msg *ServerResponse) {
	lookup := games[code]

	for c := range lookup.clients {
		_ = sendMessage(c, msg)
	}
}

func sendMessage(ws *websocket.Conn, msg *ServerResponse) error {
	return ws.WriteJSON(msg)
}

func readMessage(ws *websocket.Conn) (*ClientRequest, error) {
	type UntypedClientRequest struct {
		Code             string `json:"code"`
		Operand          `json:"operand"`
		RequestOperation `json:"requestOperation"`
		Player           string `json:"player"`
	}
	msg := &UntypedClientRequest{}
	err := ws.ReadJSON(msg)
	return &ClientRequest{
		Operand:          msg.Operand,
		Code:             Code(msg.Code),
		Player:           msg.Player,
		RequestOperation: msg.RequestOperation,
	}, err
}

func remove(slice []*websocket.Conn, item *websocket.Conn) []*websocket.Conn {
	idx := 0
	for ; idx < len(slice); idx++ {
		if slice[idx] == item {
			break
		}
	}
	if idx >= len(slice) {
		return slice
	}
	return append(slice[:idx], slice[idx+1:]...)
}
