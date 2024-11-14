package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
)

type User struct {
	Username string `json:"username"`
}

var users = make(map[string]User)
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
    
    button{    
    font-size: 18px;
border: 2px solid black;
    margitn-top:30px;

   border-radius: 5px;
   background:black;   
   color:white;
  text-align:center;
    }
  
    body{
  text-align:center;
    }
    </style>
</head>
<body>  <div id="cont">
    <button onclick="all_users()">All users</button>
    <div id="all"></div>

    <h2>Add User</h2>
    <input type="text" id="username" placeholder="Enter username">
    <br>
    <button onclick="add()">Add User</button>
    <h2>Delete User</h2>
    <input type="text" id="deleteUsername" placeholder="Enter username to delete">
   <br> <button onclick="del()">Delete User</button>
</div>
    <script>
        function all_users() {
            fetch('/users')
                .then(response => response.json())
                .then(data => {
                    let userList = '<h3 style="width: 100%; text-align: center;">Users:</h3><ul>';
                    data.forEach(user => {
                        userList += '<li>' + user.username + '</li>';
                    });
                    userList += '</ul>';
                    document.getElementById('all').innerHTML = userList;
                })
                .catch(error => console.error('Error:', error));
        }

        function add() {
            const username = document.getElementById('username').value;
            fetch('/users', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username: username })
            })
            .then(response => response.text())
            .then(data => {
                alert(data);
                all_users(); 
            })
            .catch(error => console.error('Error:', error));
        }

        function del() {
            const username = document.getElementById('deleteUsername').value;
            fetch('/users/' + username, {
                method: 'DELETE'
            })
            .then(response => response.text())
            .then(data => {
                alert(data);
                all_users(); 
            })
            .catch(error => console.error('Error:', error));
        }
    </script>
</body>
</html>
`

func getUsers(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userList := make([]User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}
	json.NewEncoder(w).Encode(userList)
}
func add(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if u.Username == "" {
		http.Error(w, "Username cannot be empty", http.StatusBadRequest)
		return
	}
	users[u.Username] = u
	w.Write([]byte("User  added successfully"))
}
func del(w http.ResponseWriter, r *http.Request) {
	u := strings.TrimPrefix(r.URL.Path, "/users/")
	if _, exists := users[u]; exists {
		delete(users, u)
		w.Write([]byte("User  deleted successfully"))
	} else {
		http.Error(w, "User  not found", http.StatusNotFound)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.New("index")
		tmpl, _ = tmpl.Parse(indexHTML)
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getUsers(w, r)
		case http.MethodPost:
			add(w, r)
		default:
			http.Error(w, "choose right method!", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/users/", del)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
