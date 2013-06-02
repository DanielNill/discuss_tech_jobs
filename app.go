package main

import (
  "fmt"
  "html/template"
  "net/http"
  "code.google.com/p/go.crypto/bcrypt"
  "github.com/gorilla/sessions"
  "github.com/gorilla/mux"
  "github.com/DanielNill/discuss_tech_jobs/models"
  "github.com/DanielNill/discuss_tech_jobs/db"
)

//session
var store = sessions.NewCookieStore([]byte("auth_token_goes_here"))

func landing(w http.ResponseWriter, r *http.Request){
  session, _ := store.Get(r, "session")
  user_id, _ := session.Values["user_id"]
  posts := make([]models.Post, 0, 1)
  conn := db.OpenConnection()
  rows, err := conn.Query("SELECT u.email, p.title, p.id FROM posts p, users u WHERE p.user_id = u.id LIMIT 100")
  if err != nil {
    fmt.Println(err)
  } else {
    for rows.Next() {
      var post models.Post
      user := new(models.User)
      rows.Scan(&user.Email, &post.Title, &post.Id)
      post.User = user
      posts = append(posts, post)
    }
  }
  context := map[string]interface{}{
    "user_id": user_id,
    "posts": posts,
  }

  t, _ := template.ParseFiles("home.html")
  t.Execute(w, context)
}

func newPostForm(w http.ResponseWriter, r *http.Request){
  session, _ := store.Get(r, "session")
  user_id, _ := session.Values["user_id"]
  if user_id != nil {
    t, _ := template.ParseFiles("new_post.html")
    t.Execute(w, nil)
  } else {
    http.Redirect(w, r, "/sign_up_in", http.StatusFound)
  }
}

func createNewPost(w http.ResponseWriter, r *http.Request){
  conn := db.OpenConnection()
  session, _ := store.Get(r, "session")
  _, err := conn.Exec("INSERT INTO posts (title, user_id) VALUES ($1, $2)", r.FormValue("title"), session.Values["user_id"])
  if err != nil {
    fmt.Println(err)
  }
  http.Redirect(w, r, "/", http.StatusFound)
}

func signUpInForm(w http.ResponseWriter, r *http.Request){
  t, _ := template.ParseFiles("sign_up_in.html")
  t.Execute(w, nil)
}

func createNewUser(w http.ResponseWriter, r *http.Request){
  conn := db.OpenConnection()
  email := r.FormValue("email")
  password, _ := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
  row := conn.QueryRow("INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id", email, password)
  var user_id int
  row.Scan(&user_id)
  fmt.Println(user_id)
  session, _ := store.Get(r, "session")
  session.Values["user_id"] = user_id
  session.Save(r, w)
  http.Redirect(w, r, "/", http.StatusFound)
}

func signIn(w http.ResponseWriter, r *http.Request){
  conn := db.OpenConnection()
  email := r.FormValue("email")
  row := conn.QueryRow("SELECT id, password_hash FROM users WHERE email = $1", email)
  var user_id int
  var hashed_password string
  row.Scan(&user_id, &hashed_password)
  // need to check password against hash
  fmt.Println(hashed_password)
  if user_id > 0 {
    session, _ := store.Get(r, "session")
    session.Values["user_id"] = user_id
    session.Save(r, w)
    http.Redirect(w, r, "/", http.StatusFound)
  } else {
    http.Redirect(w, r, "/sign_up_in", http.StatusFound)
  }
}

func logout(w http.ResponseWriter, r *http.Request){
  session, _ := store.Get(r, "session")
  session.Values["user_id"] = nil
  session.Save(r, w)
  http.Redirect(w, r, "/", http.StatusFound)
}

func showPost(w http.ResponseWriter, r *http.Request){
  vars := mux.Vars(r)
  post := models.GetPostById(vars["id"])
  context := map[string]interface{}{
    "post": post,
  }
  t, _ := template.ParseFiles("view_post.html")
  t.Execute(w, context)
}

func main(){
  r := mux.NewRouter()
  r.HandleFunc("/", landing)
  r.HandleFunc("/new_post", newPostForm)
  r.HandleFunc("/create_new_post", createNewPost)
  r.HandleFunc("/create_new_user", createNewUser)
  r.HandleFunc("/sign_up_in", signUpInForm)
  r.HandleFunc("/sign_in", signIn)
  r.HandleFunc("/logout", logout)
  r.HandleFunc("/post/{id:[0-9]+}", showPost)
  http.Handle("/", r)
  http.ListenAndServe(":3001", nil)
}