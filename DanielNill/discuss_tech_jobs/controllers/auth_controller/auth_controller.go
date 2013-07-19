package auth_controller

import(
  "net/http"
  "fmt"
  "github.com/DanielNill/discuss_tech_jobs/db"
  "github.com/gorilla/sessions"
  "text/template"
  //"crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("auth_token_goes_here"))

func SignUpInForm(w http.ResponseWriter, r *http.Request){
  t, _ := template.ParseFiles("templates/sign_up_in.html")
  t.Execute(w, nil)
}

func Create(w http.ResponseWriter, r *http.Request){
  conn := db.OpenConnection()
  email := r.FormValue("email")
  //password, _ := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
  row := conn.QueryRow("INSERT INTO users (email, password_hash) VALUES ($1) RETURNING id", email)
  var user_id int
  row.Scan(&user_id)
  fmt.Println(user_id)
  session, _ := store.Get(r, "session")
  session.Values["user_id"] = user_id
  session.Save(r, w)
  http.Redirect(w, r, "/", http.StatusFound)
}

func SignIn(w http.ResponseWriter, r *http.Request){
  conn := db.OpenConnection()
  email := r.FormValue("email")
  row := conn.QueryRow("SELECT id, password_hash FROM users WHERE email = $1", email)
  var user_id int
  var hashed_password string
  row.Scan(&user_id, &hashed_password)
  // need to check password against hash
  fmt.Println(hashed_password)
  if user_id > 0 {
    session, _ := store.Get(r, "session")
    session.Values["user_id"] = user_id
    session.Save(r, w)
    http.Redirect(w, r, "/", http.StatusFound)
  } else {
    http.Redirect(w, r, "/sign_up_in", http.StatusFound)
  }
}

func Logout(w http.ResponseWriter, r *http.Request){
  session, _ := store.Get(r, "session")
  session.Values["user_id"] = nil
  session.Save(r, w)
  http.Redirect(w, r, "/", http.StatusFound)
}