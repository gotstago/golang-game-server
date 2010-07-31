// Copyright (C) 2010 Alex Gontmakher (gsasha@gmail.com)
// This code is licensed under GPLv3

// This is a general multiplayer game implementation
package game

type Game struct {
  name string
  comm [2]chan string
  key [2]string
}

func NewGame(name string) *Game {
  game := new(Game)

  game.name = name
  for i:=0; i<2; i++ {
      game.comm[i] = make(chan string)
      game.key[i] = create_random_id(6)
  }
  return game
}

func (game *Game) Send(key, message string) string {
  if key != game.key[0] && key != game.key[1] {
      return "Bad key"
  }
  for i:=0; i<2; i++ {
    game.comm[i] <- message
  }
  return "OK"
}

func (game *Game) Receive(key string) string {
  for i:=0; i<2; i++ {
      if key==game.key[i] {
          return <-game.comm[i]
      }
  }
  return "bad key"
}

