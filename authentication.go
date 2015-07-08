package main

import (
  "encoding/json"
  "golang.org/x/oauth2"
  "io/ioutil"
  "net/http"
)

var conf = &oauth2.Config{
  ClientID: "156350471828-1gsbjgc046vtpj216fjhk4vgc4eski5d.apps.googleusercontent.com",
  ClientSecret: "Fiy_KYAW45Ln3txLE2t4M-UG",
  Scopes:[]string{"https://www.googleapis.com/auth/userinfo.profile"},
  RedirectURL: "",
  Endpoint: oauth2.Endpoint{
    AuthURL: "https://accounts.google.com/o/oauth2/auth",
    TokenURL: "https://accounts.google.com/o/oauth2/token",
  },
}

type GoogleProfile struct {
  Name string `json:name`
  GivenName string `json:given_name`
  FamilyName string `json:family_name`
  Id string `json:id`
  Picture string `json:picture`
  Link string `json:link`
  Locale string `json:locale`
  Gender string `json:gender`
}

// Start the authorization process
func googleHandler(w http.ResponseWriter, r *http.Request) {
  conf.RedirectURL = r.Referer() + "oauth2callback"

  //Get the Google URL which shows the Authentication page to the user
  url := conf.AuthCodeURL("state")

  //redirect user to that page
  http.Redirect(w, r, url, http.StatusFound)
}

// Start the authorization process
func facebookHandler(w http.ResponseWriter, r *http.Request) {
  //url := fb_conf.AuthCodeURL("state")
  http.Redirect(w, r, "/", http.StatusFound)
}

// Function that handles the callback from the Google server
const profileInfoURL = "https://www.googleapis.com/oauth2/v1/userinfo?alt=json"
func handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
  //Get the code from the response
  code := r.FormValue("code")
  token, err := conf.Exchange(oauth2.NoContext, code)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  client := conf.Client(oauth2.NoContext, token)
  resp, err := client.Get(profileInfoURL)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  // Reading the body
  raw, err := ioutil.ReadAll(resp.Body)
  defer resp.Body.Close()
  if err != nil {
   http.Error(w, err.Error(), http.StatusInternalServerError)
   return
  }

  // Unmarshalling the JSON of the Profile
  var profile = GoogleProfile{}
  if err := json.Unmarshal(raw, &profile); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  // Add a cookie here to be passed to other pages.
  // TODO(olivier): Probably cookie + session would be better.
  http.SetCookie(w, &http.Cookie{Name: "name", Value: profile.Name})
  http.SetCookie(w, &http.Cookie{Name: "id", Value: profile.Id})

  // Redirect to logged in page
  http.Redirect(w, r, "/client", http.StatusMovedPermanently)
}
