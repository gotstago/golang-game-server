// Copyright (C) 2010 Alex Gontmakher (gsasha@gmail.com)
// This code is licensed under GPLv3

package main

import (
    "flag"
    "fmt"
    "http"
    "log"
    "strings"
    "game/game"
)

var addr = flag.String("addr", ":1718", "http service address") // Q=17, R=18

var room = game.NewGameRoom()

func main() {
    flag.Parse()
    http.HandleFunc("/connect/", Connect)
    http.HandleFunc("/list/", List)
    http.HandleFunc("/wait/", Wait)
    http.HandleFunc("/create/", Create)
    http.HandleFunc("/send/", Send)
    http.HandleFunc("/receive/", Receive)

    fmt.Println("Server starting")
    err := http.ListenAndServe(*addr, nil)
    if err != nil {
        log.Exit("ListenAndServe:", err)
    }
}

func Connect(c *http.Conn, req *http.Request) {
    query_params, _ := http.ParseQuery(req.URL.RawQuery)
    arg_owner, ok_owner := query_params["owner"]
    if ok_owner {
        fmt.Println("Connect owner=", arg_owner)
        owner := arg_owner[0]
        owner_id := room.Connect(owner)
        fmt.Println("Connected as ", owner_id)
        c.Write([]byte(owner_id))
    } else {
        c.Write([]byte("Parameters not specified correctly"))
    }
}

func List(c *http.Conn, req *http.Request) {
    players := room.List()
    c.Write([]byte(strings.Join(players, ",")))
}

func Create(c *http.Conn, req *http.Request) {
    query_params, _ := http.ParseQuery(req.URL.RawQuery)
    arg_owner_id, ok_owner_id := query_params["owner_id"]
    arg_guest, ok_guest := query_params["guest"]
    if ok_owner_id && ok_guest {
        fmt.Println("Create owner_id=", arg_owner_id, ", guest=", arg_guest)
        owner_id := arg_owner_id[0]
        guest := arg_guest[0]
        game_id := room.Create(owner_id, guest)
        fmt.Println("Create ", game_id)
        c.Write([]byte(game_id))
    } else {
        c.Write([]byte("Parameters not specified correctly"))
    }
}

func Wait(c *http.Conn, req *http.Request) {
    query_params, _ := http.ParseQuery(req.URL.RawQuery)
    arg_owner_id, ok_owner_id := query_params["owner_id"]
    if ok_owner_id {
        fmt.Println("Create owner_id=", arg_owner_id)
        owner_id := arg_owner_id[0]
        game_id := room.Wait(owner_id)
        fmt.Println("Create ", game_id)
        c.Write([]byte(game_id))
    } else {
        c.Write([]byte("Parameters not specified correctly"))
    }
}

func Send(c *http.Conn, req *http.Request) {
    query_params, _ := http.ParseQuery(req.URL.RawQuery)
    arg_game, ok_game := query_params["game"]
    arg_msg, ok_msg := query_params["msg"]
    if !ok_game {
      c.Write([]byte("Game parameter not set"))
      return
    }
    if !ok_msg {
      c.Write([]byte("Msg parameter not set"))
      return
    }
    game := arg_game[0]
    msg := arg_msg[0]
    params := strings.Split(game, ",", 2)
    if len(params) != 2 {
      c.Write([]byte("Malformed game parameter"))
      return
    }
    game_id, player_id := params[0], params[1]
    game_instance := room.Find(game_id)
    if room == nil {
      c.Write([]byte("Game not found"))
      return
    }
    result := game_instance.Send(player_id, msg)
    c.Write([]byte(result))
}

func Receive(c *http.Conn, req *http.Request) {
    query_params, _ := http.ParseQuery(req.URL.RawQuery)
    arg_game, ok_game := query_params["game"]
    if !ok_game {
      c.Write([]byte("Game parameter not set"))
      return
    }
    game := arg_game[0]
    params := strings.Split(game, ",", 2)
    if len(params) != 2 {
      c.Write([]byte("Malformed game parameter"))
      return
    }
    game_id, player_id := params[0], params[1]
    game_instance := room.Find(game_id)
    if room == nil {
      c.Write([]byte("Game not found"))
      return
    }
    result := game_instance.Receive(player_id)
    c.Write([]byte(result))
}

