package db

import(
  "fmt"
  _ "github.com/lib/pq"
  "database/sql"
)

func OpenConnection() *sql.DB {
  conn, err := sql.Open("postgres", "user=Daniel password=*Mrbobn1 dbname=discuss_dev_jobs sslmode=disable")
  if err != nil {
    fmt.Println(err)
  }
  //defer conn.Close()
  return conn
}