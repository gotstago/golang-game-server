// Copyright (C) 2010 Alex Gontmakher (gsasha@gmail.com)
// This code is licensed under GPLv3

// This is a room controlling multiple games.
package game

import (
    "fmt"
)

type GameRoom struct {
  // Channels used to pass commands to the event loop.
  connect_cmd chan connectReq
  list_cmd chan listReq
  create_cmd chan createReq
  find_cmd chan findReq
  wait_cmd chan waitReq
  //finished_cmd chan string

  // Data stored in the game room.
  players map[string] *Player
  ids map[string] *Player

  games map[string] *Game
}

type Player struct {
    id string
    name string
    game string
    joined chan string
    rating float
    last_time float
}

func NewGameRoom() *GameRoom {
    room := new(GameRoom)
    room.connect_cmd = make(chan connectReq)
    room.list_cmd = make(chan listReq)
    room.create_cmd = make(chan createReq)
    room.wait_cmd = make(chan waitReq)
    room.find_cmd = make(chan findReq)
    room.players = make(map[string] *Player)
    room.ids = make(map[string] *Player)
    room.games = make(map[string] *Game)

    go room.eventLoop()
    return room
}

func (room *GameRoom) eventLoop() {
    for {
        select {
        case req := <-room.connect_cmd:
            room.connect_stub(req)
        case req := <-room.list_cmd:
            room.list_stub(req)
        case req := <-room.create_cmd:
            room.create_stub(req)
        case req := <-room.find_cmd:
            room.find_stub(req)
        case req := <-room.wait_cmd:
            room.wait_stub(req)
        //case req := <-room.finished_cmd:
        //    room.finished_stub(req)
        }
    }
}

type connectReq struct {
    name string
    reply chan string
}

// Connects to the game room. Returns the id for this name.
// If the name is already taken, returns nil.
func (room *GameRoom) Connect(name string) string {
    req := connectReq{name, make(chan string)}
    room.connect_cmd <- req
    return <-req.reply
}

func (room *GameRoom) connect_stub(req connectReq) {
    req.reply <- room.connect(req.name)
}

// Serial implementation for the room connect.
func (room *GameRoom) connect(name string) string {
    _, found := room.players[name]
    if found {
        return ""
    }
    player := new(Player)
    player.name = name
    player.id = create_random_id(12)
    player.joined = make(chan string, 1)
    player.rating = 0

    room.players[name] = player
    room.ids[player.id] = player
    return player.id
}

type listReq struct {
    reply chan []string
}

// Returns a list of players currently registered to the room.
func (room *GameRoom) List() []string {
    req := listReq{make(chan []string)}
    room.list_cmd <- req
    return <-req.reply
}

func (room *GameRoom) list_stub(req listReq) {
    req.reply <- room.list()
}

// Serial implementation for the room list.
func (room *GameRoom) list() []string {
    result := make([]string, len(room.players))
    i := 0
    for name, _ := range(room.players) {
        result[i] = name
        i++
    }
    return result
}

type createReq struct {
    owner_id string
    guest string
    reply chan string
}

// Creates a new game for given owner and guest
func (room *GameRoom) Create(owner_id, guest string) string {
    req := createReq{owner_id, guest, make(chan string)}
    room.create_cmd <- req
    return <- req.reply
}

func (room *GameRoom) create_stub(req createReq) {
    req.reply <- room.create(req.owner_id, req.guest)
}

func (room *GameRoom) create(owner_id, guest string) string {
    owner_player, ok := room.ids[owner_id]
    if !ok {
        fmt.Println("Owner player not found")
        return ""
    }
    guest_player, ok := room.players[guest]
    if !ok {
        fmt.Println("Guest player not found")
        return ""
    }
    if owner_player == guest_player {
        fmt.Println("Owner player same as guest")
        return ""
    }
    if owner_player.game != "" || guest_player.game != "" {
        fmt.Println("Owner or guest already playing")
        return ""
    }
    game_id := create_random_id(6)
    game := NewGame(owner_player.name)
    room.games[game_id] = game

    owner_player.game = game_id+","+game.key[0]
    guest_player.game = game_id+","+game.key[1]
    guest_player.joined <- guest_player.game
    return owner_player.game
}

type findReq struct {
    game_id string
    reply chan *Game
}

// Creates a new game for given owner and guest
func (room *GameRoom) Find(game_id string) *Game {
    req := findReq{game_id, make(chan *Game)}
    room.find_cmd <- req
    return <- req.reply
}

func (room *GameRoom) find_stub(req findReq) {
    req.reply <- room.find(req.game_id)
}

func (room *GameRoom) find(game_id string) *Game {
    game, ok := room.games[game_id]
    if !ok {
        fmt.Println("Game not found")
        return nil
    }
    return game
}

type waitReq struct {
    owner_id string
    reply chan chan string
}

// Waits for somebody to create a game for this owner.
func (room *GameRoom) Wait(owner_id string) string {
    req := waitReq{owner_id, make(chan chan string)}
    room.wait_cmd <- req
    // The actual waiting has to be done in the caller, otherwise the event
    // loop will get stuck.
    wait_chan := <-req.reply
    return <-wait_chan
}

func (room *GameRoom) wait_stub(req waitReq) {
    req.reply <- room.wait(req.owner_id)
}

func (room *GameRoom) wait(owner_id string) chan string {
    owner_player, ok := room.ids[owner_id]
    if !ok {
        fmt.Println("Owner player not found")
        return nil
    }
    return owner_player.joined
}
