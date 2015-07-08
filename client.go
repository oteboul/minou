package main

import (
  "github.com/gorilla/websocket"
)

type client struct {
  socket *websocket.Conn
  User_id string
  im_server *server
  online chan string
  offline chan string
  send chan *Message
}

func newClient(im_server *server, socket *websocket.Conn) *client {
  return &client{socket: socket,
                 User_id: "",
                 im_server: im_server,
                 online: make(chan string),
                 offline: make(chan string),
                 send: make(chan *Message)}
}

func (c *client) readFromSocket() {
  for {
    msg := Message{}
    err := c.socket.ReadJSON(&msg)
    if err != nil {
      return
    }
    
    // Possible types of messages:
    //  1. the client is joining: empty message to no one from 'from'.
    //  2. real message from a client to another
    if len(msg.To) == 0 && len(msg.Text) == 0 {
      c.User_id = msg.From
      c.im_server.join <- c
    } else if len(msg.To) > 0 && len(msg.Text) > 0 {
      msg.print();
      c.im_server.messages <- &msg
    }
  }
  c.socket.Close()
}

func (c *client) writeToSocket() {
  for {
    select {
      case user_id := <-c.online:
        // An empty message with reciepient that the sender is getting online.
        msg := newMessage(user_id, c.User_id, "")
        if err := c.socket.WriteJSON(msg); err != nil {
          break
        }
      case user_id := <-c.offline:
        // And message which only specifies the sender, means that the sender is going offline.
        msg := newMessage(user_id, "", "")
        if err := c.socket.WriteJSON(msg); err != nil {
          break
        }
      case msg := <-c.send:
        if msg == nil {
          break
        }
        if err := c.socket.WriteJSON(msg); err != nil {
          break
        }
      default:
        break
    }
  }
  c.socket.Close()
}