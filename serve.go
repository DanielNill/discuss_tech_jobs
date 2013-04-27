package main

import (
  "fmt"
  "html/template"
  "net/http"
  _ "github.com/lib/pq"
  "database/sql"
  "code.google.com/p/go.crypto/bcrypt"
)

func OpenConnection() *sql.DB {
  conn, err := sql.Open("postgres", "user=Daniel password=*Mrbobn1 dbname=discuss_dev_jobs sslmode=disable")
  if err != nil {
    fmt.Println(err)
  }
  return conn
}

func landing(w http.ResponseWriter, r *http.Request){
  conn := OpenConnection()
  defer conn.Close()
  rows, err := conn.Query("SELECT * FROM posts p JOIN users u ON p.user_id = u.user_id LIMIT 100")
  if err != nil {
    fmt.Println(err)
  } else {
    //columns, _ := rows.Columns()
    //values := make([]interface{}, len(columns))
    //scanArgs := make([]interface{}, len(values))
    //for i := range values {
    //  scanArgs[i] = &values[i]
    //}
    for rows.Next() {
      var post_title string
      var user_name string
      var post_id int
      rows.Scan(&post_title, &user_name, &post_id)
    }
  }
  t, _ := template.ParseFiles("home.html")
  t.Execute(w, nil)
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
  t, _ := template.ParseFiles("home.html")
  t.Execute(w, nil)
}

func newUserForm(w http.ResponseWriter, r *http.Request){
  t, _ := template.ParseFiles("new_user.html")
  t.Execute(w, nil)
}

func createNewUser(w http.ResponseWriter, r *http.Request){
  conn := OpenConnection()
  email := r.FormValue("email")
  password, _ := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), 0)
  conn.Exec("INSERT INTO users (email, password_hash) VALUES (?, ?)", email, password)
}

func main(){
  http.HandleFunc("/", landing)
  http.HandleFunc("/new_post", newPostForm)
  http.HandleFunc("/create_new_post", createNewPost)
  http.HandleFunc("/new_user", newUserForm)
  http.HandleFunc("/create_new_user", createNewUser)
  http.ListenAndServe(":3001", nil)
}