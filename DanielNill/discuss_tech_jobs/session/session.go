package session

// import(
//   "net/http"
//   "github.com/gorilla/sessions"
//   _ "fmt"
// )

// var store = sessions.NewCookieStore([]byte("auth_token_goes_here"))

// func isLoggedIn(r *http.Request) bool {
//   session, _ := store.Get(r, "session")
//   user_id, _ := session.Values["user_id"]
//   if user_id != nil {
//     return true
//   } else {
//     return false
//   }
// }

// func getCurrentUserId(r *http.Request) int {
//   session, _ := store.Get(r, "session")
//   user_id, _ := session.Values["user_id"]
//   fmt.Println(user_id)
//   return user_id
// }