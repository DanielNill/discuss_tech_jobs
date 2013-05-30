package models

import(
  _"fmt"
  "github.com/DanielNill/discuss_tech_jobs/db"
)

type Post struct {
  Id int
  Title string
  Points int
  User *User
  CreatedAt string
  UpdatedAt string
}

func GetPostById(id string) *Post {
  conn := db.OpenConnection()
  post := new(Post)
  post.User = new(User)
  row := conn.QueryRow("SELECT p.id, p.title, u.id, u.email FROM posts p, users u WHERE p.user_id = u.id LIMIT 1")
  row.Scan(&post.Id, &post.Title, &post.Id, &post.User.Email)
  return post
}