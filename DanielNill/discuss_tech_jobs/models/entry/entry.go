package entry

import(
  "fmt"
  "github.com/DanielNill/discuss_tech_jobs/db"
  "github.com/DanielNill/discuss_tech_jobs/models/user"
)

type Post struct {
  Id int
  Title string
  Text string
  Points int
  User *user.User
  Comments *[]Comment
  CreatedAt string
  UpdatedAt string
}

func GetPostById(id string) *Post {
  conn := db.OpenConnection()
  post := new(Post)
  post.User = new(user.User)
  row := conn.QueryRow("SELECT p.id, p.title, p.points, u.id, u.email FROM posts p, users u WHERE p.user_id = u.id AND p.id = $1 LIMIT 1", id)
  row.Scan(&post.Id, &post.Title, &post.Points, &post.User.Id, &post.User.Email)
  return post
}

//**************
// Comments
//**************

type Comment struct {
  Id int
  Text string
  Points int
  User *user.User
  Post *Post
  CreatedAt string
  UpdatedAt string
}

func GetCommentById(id string) *Comment {
  conn := db.OpenConnection()
  comment := new(Comment)
  comment.User = new(user.User)
  comment.Post = new(Post)
  row := conn.QueryRow("SELECT c.id, c.text, c.points, c.post_id, u.email FROM comments c, users u WHERE c.user_id = u.id AND c.id = $1 LIMIT 1", id)
  row.Scan(&comment.Id, &comment.Text, &comment.Points, &comment.Post.Id, &comment.User.Email)
  return comment
}

func GetCommentsByPostId(post_id string) []Comment {
  conn := db.OpenConnection()
  comments := make([]Comment, 0, 1)
  rows, err := conn.Query("SELECT c.id, c.text, c.points, u.email FROM comments c, users u WHERE c.user_id = u.id AND c.post_id = $1", post_id)
  if err != nil {
    fmt.Println(err)
  } else {
    for rows.Next() {
      var comment Comment
      user := new(user.User)
      rows.Scan(&comment.Id, &comment.Text, &comment.Points, &user.Email)
      comment.User = user
      comments = append(comments, comment)
    }
  }
  return comments
}

func GetCommentsByParentCommentId(parent_comment_id string) []Comment {
  conn := db.OpenConnection()
  comments := make([]Comment, 0, 1)
  rows, err := conn.Query("SELECT c.id, c.text, c.points, u.email FROM comments c, users u WHERE c.user_id = u.id AND c.parent_comment_id = $1", parent_comment_id)
  if err != nil {
    fmt.Println(err)
  } else {
    for rows.Next() {
      var comment Comment
      user := new(user.User)
      rows.Scan(&comment.Id, &comment.Text, &comment.Points, &user.Email)
      comment.User = user
      comments = append(comments, comment)
    }
  }
  return comments
}