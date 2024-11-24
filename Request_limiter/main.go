package main

import (
	"encoding/json"
	"html/template"
	"net/http"

	//"fmt"
	"github.com/go-chi/chi/v5"
	"golang.org/x/time/rate"
)

type Users struct {
	Username []string `json:"username"`
}

var users Users

var indexHTML = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Users</title>
    <style>
    #cont{
        background:cyan;
        border-radius: 50px;
        display: inline-block;
        width: 40%;
        border: 5px solid red;
        padding: 20px;
    }
    input{
        background:white;    
        font-size: 18px;
        border: 2px solid blue;
        border-radius: 5px;
    }
    *{
        margin-top:30px;
    }
    button{    
        font-size: 18px;
        border: 2px solid black;
        border-radius: 5px;
        background:black;   
        color:white;
        text-align:center;
    }
    body{
        text-align:center;
    }
    </style>
    <script>
        function all_users() {
            fetch('/users')
                .then(response => response.json())
                .then(data => {
                    document.getElementById('all').innerHTML = data.username.join(' ');
                });
        }

        function add() {
            const username = document.getElementById('username').value;
            fetch('/users', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ username: username }),
            }).then(() => {
                document.getElementById('username').value = '';
                all_users();
            });
        }

        function del() {
            const username = document.getElementById('deleteUsername').value;
            fetch('/users/' + username, {
                method: 'DELETE',
            }).then(() => {
                document.getElementById('deleteUsername').value = '';
                all_users();
            });
        }
    </script>
</head>
<body>
    <div id="cont">
        <button onclick="all_users()">All users</button>
        <div id="all"></div>

        <h2>Add User</h2>
        <input type="text" id="username" placeholder="Enter username">
        <br>
        <button onclick="add()">Add User</button>
        <h2>Delete User</h2>
        <input type="text" id="deleteUsername" placeholder="Enter username to delete">
        <br>
        <button onclick="del()">Delete User</button>
    </div>
</body>
</html>
`

var limiter = rate.NewLimiter(1, 3)

func rateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func basicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || username != "1" || password != "1" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func home(w http.ResponseWriter, _ *http.Request) {
	tmp, err := template.New("home").Parse(indexHTML)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmp.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func get(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func add(w http.ResponseWriter, r *http.Request) {
	var new struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&new); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	users.Username = append(users.Username, new.Username)
	w.WriteHeader(http.StatusCreated)
}

func delete(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	for i, user := range users.Username {
		if user == username {
			users.Username = append(users.Username[:i], users.Username[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "this user doesn't exist", http.StatusNotFound)
}

func main() {
	r := chi.NewRouter()
	r.Use(rateLimit)
	r.Use(basicAuth)
	r.Get("/", home)
	r.Get("/users", get)
	r.Post("/users", add)
	r.Delete("/users/{username}", delete)

	http.ListenAndServe(":8080", r)
}
