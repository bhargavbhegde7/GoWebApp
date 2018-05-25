package user

import (
	"net/http"
	"fmt"
	"text/template"
	"github.com/gorilla/securecookie"
)

//Controller ...
type Controller struct {
	Repository Repository
}

type ErrorMessage struct {
	Message string
}

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

// Index GET /
func (c *Controller) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("inside index handler")
	http.Redirect(w, r, "/home", http.StatusFound)
	return
}

// Entry GET /
func (c *Controller) Entry(w http.ResponseWriter, req *http.Request) {
	m := ErrorMessage{}

	username := getUserName(req)
	if username != "" {
		fmt.Println("user exists")
		m = ErrorMessage{"Already logged in as " + username}
	}

	t, err := template.ParseFiles("./static/entry.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, m)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Login POST /
func (c *Controller) Login(w http.ResponseWriter, req *http.Request) {
	username := req.FormValue("username")
	passwd := req.FormValue("passwd")

	//check if the uname and password match
	if username == "bhargav" && passwd == "bhargav" {

		//check if session exists, get the session with uname
		//start a new session if doesn't exist

		setSession(username, w)

		//TODO send the user to "/home" with all these incoming data
		http.Redirect(w, req, "/home", http.StatusFound)
		return
	} else if username == "bhargav2" && passwd == "bhargav2" {

		//check if session exists, get the session with uname
		//start a new session if doesn't exist
		setSession(username, w)

		//TODO send the user to "/home" with all these incoming data
		http.Redirect(w, req, "/home", http.StatusFound)
		return
	} else {

		//TODO URL stays /login in browser when login fail happens. Need to redirect to /entry with the message.
		// how to do - redirect the user to /entry?error=loginFailed

		m := ErrorMessage{Message: "wrong credentials"}

		t, err := template.ParseFiles("./static/entry.html")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = t.Execute(w, m)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		//http.Redirect(w, req, "/home", http.StatusFound)
	}
}

// Signup POST /
func (c *Controller) Signup(w http.ResponseWriter, req *http.Request) {
	var m ErrorMessage

	username := req.FormValue("username")
	passwd1 := req.FormValue("passwd1")
	passwd2 := req.FormValue("passwd2")

	//TODO create a validation method
	if username != "" && passwd1 != "" && passwd1 == passwd2 {
		fmt.Println("username : " + username + " , password : " + passwd1)
		http.Redirect(w, req, "/home", http.StatusFound)
		return
	}

	m = ErrorMessage{Message: "error : enter proper username and password"}

	t, err := template.ParseFiles("./static/entry.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, m)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Home GET /
func (c *Controller) Home(w http.ResponseWriter, req *http.Request) {
	fmt.Println("inside home handler")

	username := getUserName(req)
	if username != "" {
		fmt.Println("username exists")

		t, err := template.ParseFiles("./static/home.html")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		m := ErrorMessage{Message: username}
		err = t.Execute(w, m)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		fmt.Println("username doesn't exist")
		http.Redirect(w, req, "/entry", 302)
	}
}

// Logout GET /
func (c *Controller) Logout(w http.ResponseWriter, req *http.Request) {
	clearSession(w)
	http.Redirect(w, req, "/", 302)
}

func getUserName(request *http.Request) (username string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			username = cookieValue["username"]
		}
	}
	return username
}

func setSession(username string, response http.ResponseWriter) {
	value := map[string]string{
		"username": username,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}