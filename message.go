// Defines a message between two clients and convenience functions to manipulate them.
package main

import (
  "fmt"
)

type Message struct {
  From string
  To string 
  Text string
}

func newMessage(from string, to string, text string) *Message {
  return &Message{From: from, To: to, Text: text}
}

func newEmptyMessage() *Message {
  return newMessage("", "", "")
}

func (m *Message) empty() bool {
  return len(m.From) == 0 && len(m.To) == 0 && len(m.Text) == 0;
}

func (m *Message) print() {
  fmt.Printf("\tMessage [%s] from {%s} to {%s}\n", m.Text, m.From, m.To)
}