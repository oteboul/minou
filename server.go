package main
 
import (
  "fmt"
  "gopkg.in/mgo.v2"
  "github.com/gorilla/websocket"
  "net/http"
)

const (
  socketBufferSize = 1024
)

type server struct {
  clients map[string]*client
  join chan *client
  leave chan *client
  messages chan *Message
  mgo_session *mgo.Session
  upgrader *websocket.Upgrader
}

func newServer(session *mgo.Session) *server {
  return &server{
    messages: make(chan *Message),
    join:    make(chan *client),
    leave:   make(chan *client),
    clients: make(map[string]*client),
    mgo_session: session,
    upgrader: &websocket.Upgrader{
      ReadBufferSize: 1024,
      WriteBufferSize: socketBufferSize,
      CheckOrigin: func(r *http.Request) bool { return true },
    },
  }
}

func (s *server) deleteClient(name string) {
  if client, ok := s.clients[name]; ok {
    fmt.Printf("Client %s if leaving us. Bastard!\n", client.User_id)
    delete(s.clients, client.User_id)
    close(client.online)
    close(client.offline)
    close(client.send)
    for _, other := range(s.clients) {
      other.offline <- client.User_id
    }
  }
}

func (s *server) run() {
  for {
    select {
    case client := <-s.join:
      s.clients[client.User_id] = client
      fmt.Printf("Client %s is joining the game!!\n", client.User_id)
      for name, other := range s.clients{
        if name != client.User_id {
          other.online <- client.User_id
          client.online <- other.User_id
        }
      }

    case client := <-s.leave:
      s.deleteClient(client.User_id)
      
    case msg := <-s.messages:
      if len(msg.To) > 0 && len(msg.Text) > 0 {
        if other, ok := s.clients[msg.To]; ok {
          other.send <- msg
          // Save the message in the DB.
        }
      }
    }
  }
}
