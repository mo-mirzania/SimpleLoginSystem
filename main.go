package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"text/template"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	Username string
	Fname    string
	Lname    string
	Password string
}

var usersDB = make(map[string]User)
var sessionsDB = make(map[string]string)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("./*.gohtml"))
}

func main() {
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		session, err := r.Cookie("session")
		if err != nil {
			fmt.Println(err.Error())
			log.Fatal()
		}
		un := sessionsDB[session.Value]
		user := usersDB[un]

		fmt.Println(un)
		tpl.ExecuteTemplate(w, "index.gohtml", user)
	}
}

func signup(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session")
	if err != nil {
		if r.Method == http.MethodPost {
			u := User{
				r.FormValue("uname"),
				r.FormValue("fname"),
				r.FormValue("lname"),
				r.FormValue("password"),
			}
			usersDB[u.Username] = u
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	tpl.ExecuteTemplate(w, "signup.gohtml", nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session")
	if err != nil {
		if r.Method == http.MethodPost {
			username := r.FormValue("uname")
			password := r.FormValue("password")
			if _, ok := usersDB[username]; ok {
				if password == usersDB[username].Password {
					c, err := uuid.NewV4()
					if err != nil {
						fmt.Fprintf(w, err.Error())
					}
					cookie := &http.Cookie{
						Name:  "session",
						Value: c.String(),
					}
					http.SetCookie(w, cookie)
					sessionsDB[cookie.Value] = username
					http.Redirect(w, r, "/", http.StatusSeeOther)
				} else {
					tpl.ExecuteTemplate(w, "login.gohtml", nil)
					io.WriteString(w, "User/Pass Wrong!")
				}
			} else {
				tpl.ExecuteTemplate(w, "login.gohtml", nil)
				io.WriteString(w, "User not found! Please signup first!")
				io.WriteString(w, `<br><a href="/signup">signup</a>`)
			}
		} else {
			tpl.ExecuteTemplate(w, "login.gohtml", nil)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		c, _ := r.Cookie("session")
		delete(sessionsDB, c.Value)
		c = &http.Cookie{
			Name:   "session",
			Value:  "",
			MaxAge: -1,
		}
		http.SetCookie(w, c)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
