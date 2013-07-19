package main

import (
  "net/http"
  "github.com/gorilla/mux"
  "github.com/DanielNill/discuss_tech_jobs/controllers/application_controller"
  "github.com/DanielNill/discuss_tech_jobs/controllers/auth_controller"
  "github.com/DanielNill/discuss_tech_jobs/controllers/posts_controller"
  "github.com/DanielNill/discuss_tech_jobs/controllers/comments_controller"
)

func main(){
  r := mux.NewRouter()
  r.HandleFunc("/", application_controller.Index)
  r.HandleFunc("/new_post", posts_controller.New)
  r.HandleFunc("/create_new_post", posts_controller.Create)
  r.HandleFunc("/create_new_user", auth_controller.Create)
  r.HandleFunc("/sign_up_in", auth_controller.SignUpInForm)
  r.HandleFunc("/sign_in", auth_controller.SignIn)
  r.HandleFunc("/logout", auth_controller.Logout)
  r.HandleFunc("/post/{id:[0-9]+}", posts_controller.Show)
  r.HandleFunc("/create_new_comment", comments_controller.Create)
  r.HandleFunc("/comments/{id:[0-9]+}/new", comments_controller.New)
  http.Handle("/", r)
  http.ListenAndServe(":3001", nil)
}