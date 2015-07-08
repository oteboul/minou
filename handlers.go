// Defines the http handlers to route the different URLs to functionalities.
package main

import (
  "fmt"
  "html/template"
  "net/http"
  "os"
)

//==============================================================================
func IMHandler(s *server) http.Handler {
  fn := func(w http.ResponseWriter, r *http.Request) {
    socket, err := s.upgrader.Upgrade(w, r, nil)
    if err != nil {
      fmt.Println(err)
      return
    }

    // The new socket creates a client which we keep track in the server.
    // At this stage, client has not been identified. We notify the server
    // of the client existence once the client sends its user_id.
    client := newClient(s, socket)
    defer func() { s.leave <- client }()
    go client.writeToSocket()
    client.readFromSocket()
  }
  return http.HandlerFunc(fn)
}

// Just a simple page with a form. When the form is submitted, creates client.
func initHandler(w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("templates/login.html")
  t.Execute(w, r)
}

// When loging in, either login is confirmed and the request is redirected
// to the client instant messaging page. Or it is not and we stay on the page.
// To start, the check is only to see if the name is available or not.
func loginHandler(s *server) http.Handler {
  fn := func(w http.ResponseWriter, r *http.Request) {
    user_name := r.FormValue("user_name")
    if _, ok := s.clients[user_name]; ok {
      http.Redirect(w, r, "/", 302)
    } else {
      // Add a cookie here to be passed to other pages.
      cookie := http.Cookie{Name: "name", Value: user_name}
      http.SetCookie(w, &cookie)
      http.Redirect(w, r, "/client", 302)
    }
  }
  return http.HandlerFunc(fn)
}

// The client has just been created, makes sure there is a get param and serve
// the client page.
func clientHandler(w http.ResponseWriter, r *http.Request) {
  cookie, _ := r.Cookie("name")
  data := struct {
        User_id string
    } {
        cookie.Value,
    }
  t, _ := template.ParseFiles("templates/client.html")
  t.Execute(w, data)
}

//==============================================================================
func main() {
  port := os.Getenv("PORT")

  im_server := newServer()
  go im_server.run()

  // Handlers
  http.HandleFunc("/", initHandler)
  http.Handle("/login", loginHandler(im_server))
  http.HandleFunc("/google", googleHandler)
  http.HandleFunc("/facebook", facebookHandler)
  http.HandleFunc("/oauth2callback", handleOAuth2Callback)
  http.HandleFunc("/client", clientHandler)
  http.Handle("/im", IMHandler(im_server))
  http.Handle(
    "/static/",
    http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

  fmt.Println("Starting Server on port ", port)
  http_err := http.ListenAndServe(port, nil)
  if http_err != nil {
    panic("Error: " + http_err.Error())
  }
}
