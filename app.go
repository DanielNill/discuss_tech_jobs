package main

import (
  "fmt"
  "html/template"
  "net/http"
  _ "github.com/lib/pq"
  "database/sql"
  "code.google.com/p/go.crypto/bcrypt"
  "github.com/gorilla/sessions"
  "github.com/DanielNill/discuss_tech_jobs/models"
)

//session
var store = sessions.NewCookieStore([]byte("auth_token_goes_here"))

type Post struct {
  Title string
}

func OpenConnection() *sql.DB {
  conn, err := sql.Open("postgres", "user=Daniel password=*Mrbobn1 dbname=discuss_dev_jobs sslmode=disable")
  if err != nil {
    fmt.Println(err)
  }
  return conn
}


func landing(w http.ResponseWriter, r *http.Request){
  //page := make([]interface{}, 0, 1)
  session, _ := store.Get(r, "session-name")
  user_id, _ := session.Values["user_id"]
  //page = append(page, map[string] interface {"user_id": user_id})
  posts := make([]Post, 0, 1)
  conn := OpenConnection()
  defer conn.Close()
  rows, err := conn.Query("SELECT p.title FROM posts p LIMIT 100")
  if err != nil {
    fmt.Println(err)
  } else {
    for rows.Next() {
      var title string
      rows.Scan(&title)
      posts = append(posts, Post{Title: title})
    }
  }
  //page = append(page, map[String] slice {"pages": pages })
  //fmt.Println(page)
  t, _ := template.ParseFiles("home.html")
  t.Execute(w, posts)
}

func newPostForm(w http.ResponseWriter, r *http.Request){
  t, _ := template.ParseFiles("new_post.html")
  t.Execute(w, nil)
}

func createNewPost(w http.ResponseWriter, r *http.Request){
  conn := OpenConnection()
  _, err := conn.Exec("INSERT INTO posts (title, user_id) VALUES ('" + r.FormValue("title") + "', 1)")
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
  conn := OpenConnection()
  defer conn.Close()
  email := r.FormValue("email")
  password, _ := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
  row := conn.QueryRow("INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING user_id", email, password)
  var user_id int
  row.Scan(&user_id)
  session, _ := store.Get(r, "session-name")
  session.Values["user_id"] = user_id
  session.Save(r, w)
  http.Redirect(w, r, "/", http.StatusFound)
}

func logout(w http.ResponseWriter, r *http.Request){
  session, _ := store.Get(r, "session-name")
  session.Values["user_id"] = nil
  session.Save(r, w)
  http.Redirect(w, r, "/", http.StatusFound)
}

func main(){
  http.HandleFunc("/", landing)
  http.HandleFunc("/new_post", newPostForm)
  http.HandleFunc("/create_new_post", createNewPost)
  http.HandleFunc("/create_new_user", createNewUser)
  http.HandleFunc("/sign_up_in", signUpInForm)
  http.HandleFunc("/logout", logout)
  http.ListenAndServe(":3001", nil)
}