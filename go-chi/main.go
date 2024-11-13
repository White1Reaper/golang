package main
import (
    "encoding/json"
    "net/http"
    "github.com/go-chi/chi/v5"
)
type User struct {
    Text string `json:"text"`
}

// GET /users
// POST /users {'name':"fdfsd"}
// DELETE /users/fdfsd

// 7. - implement using http package
// 9 - implement using go-chi middleware

func Hello(w http.ResponseWriter, r *http.Request) {
    m := Mes{Text: "hello world"}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(m)
}
func main() {
    r := chi.NewRouter()
    r.Get("/hello",Hello)
    http.ListenAndServe(":8080", r)
}

