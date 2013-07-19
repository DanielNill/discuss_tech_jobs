package posts_controller

import(
  "net/http"
  "fmt"
  "text/template"
  "github.com/gorilla/sessions"
  "github.com/gorilla/mux"
  "github.com/DanielNill/discuss_tech_jobs/db"
  "github.com/DanielNill/discuss_tech_jobs/models/entry"
)

var store = sessions.NewCookieStore([]byte("auth_token_goes_here"))

func New(w http.ResponseWriter, r *http.Request){
  session, _ := store.Get(r, "session")
  user_id, _ := session.Values["user_id"]
  if user_id != nil {
    t, _ := template.ParseFiles("templates/new_post.html")
    t.Execute(w, nil)
  } else {
    http.Redirect(w, r, "/sign_up_in", http.StatusFound)
  }
}

func Create(w http.ResponseWriter, r *http.Request){
  session, _ := store.Get(r, "session")
  user_id, _ := session.Values["user_id"]
  if user_id != nil {
    conn := db.OpenConnection()
    _, err := conn.Exec("INSERT INTO posts (title, user_id) VALUES ($1, $2)", r.FormValue("title"), user_id)
    if err != nil {
      fmt.Println(err)
    }
    http.Redirect(w, r, "/", http.StatusFound)
  } else {
    http.Redirect(w, r, "/sign_up_in", http.StatusFound)
  }
}

func Show(w http.ResponseWriter, r *http.Request){
  vars := mux.Vars(r)
  post := entry.GetPostById(vars["id"])
  comments := entry.GetCommentsByPostId(vars["id"])
  context := map[string]interface{}{
    "post": post,
    "comments": comments,
  }
  t, _ := template.ParseFiles("templates/view_post.html")
  t.Execute(w, context)
}