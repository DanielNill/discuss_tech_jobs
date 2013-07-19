package comments_controller

import(
  "net/http"
  "text/template"
  "github.com/gorilla/sessions"
  "github.com/gorilla/mux"
  "github.com/DanielNill/discuss_tech_jobs/db"
  "github.com/DanielNill/discuss_tech_jobs/models/entry"
)

var store = sessions.NewCookieStore([]byte("auth_token_goes_here"))

func Create(w http.ResponseWriter, r *http.Request){
  session, _ := store.Get(r, "session")
  user_id, _ := session.Values["user_id"]
  if user_id == nil {
    http.Redirect(w, r, "/sign_up_in", http.StatusFound)
  } else {
    conn := db.OpenConnection()
    conn.Exec("INSERT INTO comments (user_id, post_id, text) VALUES ($1, $2, $3)", user_id, r.FormValue("post_id"), r.FormValue("text"))
    http.Redirect(w, r, "/post/" + r.FormValue("post_id"), http.StatusFound )
  }
}

func New(w http.ResponseWriter, r *http.Request){
  vars := mux.Vars(r)
  session, _ := store.Get(r, "session")
  user_id, _ := session.Values["user_id"]
  if user_id == nil {
    http.Redirect(w, r, "/sign_up_in", http.StatusFound)
  } else {
    parent_comment := entry.GetCommentById(vars["id"])
    comments := entry.GetCommentsByParentCommentId(vars["id"])
    context := map[string]interface{}{
      "parent_comment": parent_comment,
      "comments": comments,
    }
    t, _ := template.ParseFiles("templates/new_comment.html")
    t.Execute(w, context)
  }
}
