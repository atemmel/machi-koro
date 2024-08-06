package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"

	"github.com/gorilla/websocket"
)

type RequestOperation int
type ResponseCode int

const (
	JoinRequest RequestOperation = iota
	PongResponse
	LeaveRequest
	StartRequest
	RollRequest
	BuyRequest
)

const (
	JoinAnnouncement ResponseCode = iota
	PingRequest
	LeaveAnnouncement
	StartAnnouncement
	OwnerAssign
	TurnChangeAnnouncement
	RollAnnouncement
	BuyAnnouncement
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
	Operands     []Operand `json:"operands"`
	ResponseCode `json:"responseCode"`
	Player       string `json:"player"`
}

type Client struct {
	Name        string       `json:"name"`
	BoughtCards map[int]Card `json:"cards"`
}

type Game struct {
	Code           Code                       // Code to game
	Owner          *websocket.Conn            // Owner of lobby
	ActivePlayer   *websocket.Conn            // Player whose turn it is
	State          GameState                  // In progess, Lobby
	Phase          string                     // Roll dice, buy cards
	clients        map[*websocket.Conn]Client // All clients
	order          []*websocket.Conn
	availableCards []Card
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
		Phase:   g.Phase,
		Players: players,
	}
}

var (
	gamesMutex = sync.Mutex{}
	games      = map[Code]Game{}
)

func AddPlayer(ws *websocket.Conn) (Code, string, error) {
	if ws == nil {
		panic("Ws was nil")
	}
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
				Operands:     []Operand{Blank},
				ResponseCode: OwnerAssign,
				Player:       msg.Player,
			})
		}

		lookup.clients[ws] = Client{
			Name:        msg.Player,
			BoughtCards: map[int]Card{},
		}
		games[msg.Code] = lookup

		response := &ServerResponse{
			Code:         msg.Code,
			Operands:     []Operand{Blank},
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
						Operands:     []Operand{Blank},
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
		Operands:     []Operand{Blank},
		ResponseCode: LeaveAnnouncement,
		Player:       player.Name,
	})
}

func NewGame() GameResponse {
	g := Game{
		Code:           genRoomCode(),
		State:          Lobby,
		clients:        map[*websocket.Conn]Client{},
		availableCards: CopyOfAllCards(),
	}
	gamesMutex.Lock()
	defer gamesMutex.Unlock()
	games[g.Code] = g
	return *g.Response()
}

func LookupGame(code Code) *GameResponse {
	if !validateRoomCode(code) {
		return nil
	}
	lookup, ok := games[code]
	if ok {
		return lookup.Response()
	}
	return nil
}

func StartGame(ws *websocket.Conn, msg *ClientRequest) {
	gamesMutex.Lock()
	defer gamesMutex.Unlock()

	lookup := games[msg.Code]

	if ws != lookup.Owner {
		// only owner can start
		return
	}

	// just get some guy
	for c := range lookup.clients {
		lookup.ActivePlayer = c
		break
	}

	// make an order of turns
	lookup.order = make([]*websocket.Conn, 0, len(lookup.clients))
	for c := range lookup.clients {
		lookup.order = append(lookup.order, c)
	}

	lookup.State = InProgress
	games[msg.Code] = lookup

	publishMessage(msg.Code, &ServerResponse{
		Player:       msg.Player,
		Code:         msg.Code,
		ResponseCode: StartAnnouncement,
		Operands:     []Operand{Blank},
	})

	activePlayer := lookup.clients[lookup.ActivePlayer]

	publishMessage(msg.Code, &ServerResponse{
		Player:       activePlayer.Name,
		Code:         msg.Code,
		ResponseCode: TurnChangeAnnouncement,
		Operands:     []Operand{Blank},
	})

}

func RollDice(ws *websocket.Conn, msg *ClientRequest) {

	operands := make([]Operand, msg.Operand)

	for i := range operands {
		nBig, err := rand.Int(rand.Reader, big.NewInt(6+1))
		if err != nil {
			panic(err)
		}
		operands[i] = Operand(nBig.Int64())
	}

	gamesMutex.Lock()
	defer gamesMutex.Unlock()
	publishMessage(msg.Code, &ServerResponse{
		Code:         msg.Code,
		Player:       msg.Player,
		ResponseCode: RollAnnouncement,
		Operands:     operands,
	})
}

func BuyCard(ws *websocket.Conn, msg *ClientRequest) {
	gamesMutex.Lock()
	defer gamesMutex.Unlock()

	g, ok := games[msg.Code]
	if !ok {
		return
	}

	c, ok := g.clients[ws]
	if !ok {
		return
	}

	cardIdx := int(msg.Operand)
	card := &g.availableCards[cardIdx]
	if card.Count <= 0 {
		return
	}
	card.Count -= 1

	if existingCard, ok := c.BoughtCards[cardIdx]; ok {
		existingCard.Count++
		c.BoughtCards[cardIdx] = existingCard
	} else {
		newCard := AllCards[cardIdx]
		newCard.Count = 1
		c.BoughtCards[cardIdx] = newCard
	}

	publishMessage(msg.Code, &ServerResponse{
		Player:       msg.Player,
		Code:         msg.Code,
		Operands:     []Operand{Operand(cardIdx)},
		ResponseCode: BuyAnnouncement,
	})

	nextTurn(&g)

	g.Phase = "roll"
	g.clients[ws] = c
	games[msg.Code] = g
}

func nextTurn(g *Game) {
	next := nextPlayer(g.order, g.ActivePlayer)
	g.ActivePlayer = next
	client := g.clients[next]

	publishMessage(g.Code, &ServerResponse{
		Player:       client.Name,
		Code:         g.Code,
		Operands:     []Operand{Blank},
		ResponseCode: TurnChangeAnnouncement,
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

func nextPlayer(order []*websocket.Conn, active *websocket.Conn) *websocket.Conn {
	for i, c := range order {
		if c == active {
			if i+1 == len(order) {
				return order[0]
			}
			return order[i+1]

		}

	}
	return nil
}
