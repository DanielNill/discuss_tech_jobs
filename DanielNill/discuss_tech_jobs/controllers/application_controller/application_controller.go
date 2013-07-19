package application_controller

import(
  "net/http"
  "fmt"
  "text/template"
  "github.com/gorilla/sessions"
  "github.com/DanielNill/discuss_tech_jobs/models/entry"
  "github.com/DanielNill/discuss_tech_jobs/models/user"
  "github.com/DanielNill/discuss_tech_jobs/db"
)

var store = sessions.NewCookieStore([]byte("auth_token_goes_here"))

func Index(w http.ResponseWriter, r *http.Request){
  session, _ := store.Get(r, "session")
  user_id, _ := session.Values["user_id"]
  posts := make([]entry.Post, 0, 1)
  conn := db.OpenConnection()
  rows, err := conn.Query("SELECT u.email, p.title, p.id FROM posts p, users u WHERE p.user_id = u.id LIMIT 100")
  if err != nil {
    fmt.Println(err)
  } else {
    for rows.Next() {
      var post entry.Post
      user := new(user.User)
      rows.Scan(&user.Email, &post.Title, &post.Id)
      post.User = user
      posts = append(posts, post)
    }
  }
  context := map[string]interface{}{
    "user_id": user_id,
    "posts": posts,
  }

  t, _ := template.ParseFiles("templates/home.html")
  t.Execute(w, context)
}
