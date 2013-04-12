package main

import (
  "fmt"
  "html/template"
  "net/http"
  _ "github.com/lib/pq"
  "database/sql"
)
func DBConnect() (conn) {
  conn, err := sql.Open("postgres", "user=Daniel password=*Mrbobn1 dbname=discuss_dev_jobs sslmode=disable")
  if err != nil {
    fmt.Println(err)
  }
  return
}

func landing(w http.ResponseWriter, r *http.Request){
  t, _ := template.ParseFiles("home.html")
  t.Execute(w, nil)
}

func newPost(w http.ResponseWriter, r *http.Request){
  t, _ := template.ParseFiles("new.html")
  t.Execute(w, nil)
}

func create(w http.ResponseWriter, r *http.Request){
  conn := DBConnect()
  rows, err := conn.Query("SELECT * FROM names")
  if err != nil {
    fmt.Println(err)
  } else{
    for rows.Next() {
      var name_id int
      var name string
      err = rows.Scan(&name_id, &name)
      fmt.Println(rows)
      fmt.Println(name_id)
      fmt.Println(name)
    }
  }
  t, _ := template.ParseFiles("home.html")
  t.Execute(w, nil)
}

func main(){
  http.HandleFunc("/", landing)
  http.HandleFunc("/new", newPost)
  http.HandleFunc("/create", create)
  http.ListenAndServe(":3001", nil)
}