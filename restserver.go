package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"

	"io"
	"os"
	"bufio"
	"strings"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

var store = sessions.NewCookieStore([]byte("something-very-secret"))

type User struct {
	UserName string
}

type Image struct{
	Path string
}

type ErrorMessage struct {
	Message string
}

//GetLoginEndpoint . . .
func GetLoginEndpoint(w http.ResponseWriter, req *http.Request) {

	m := ErrorMessage{}

	username := getUserName(req)
	if username != "" {
		fmt.Println("user exists")
		m = ErrorMessage{"Already logged in as " + username}
	}

	t, err := template.ParseFiles("login.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, m)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//GetSignupEndpoint . . .
func GetSignupEndpoint(w http.ResponseWriter, req *http.Request) {

	m := ErrorMessage{}

	username := getUserName(req)
	if username != "" {
		fmt.Println("user exists")
		m = ErrorMessage{"Already logged in as " + username}
	}

	t, err := template.ParseFiles("signup.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, m)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

//this blog right here made it work.
//don't understand the code fully
//	https://mschoebel.info/2014/03/09/snippet-golang-webapp-login-logout/

//LoginEndpoint . . .
func LoginEndpoint(w http.ResponseWriter, req *http.Request) {

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

		m := ErrorMessage{Message: "wrong credentials"}

		t, err := template.ParseFiles("login.html")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = t.Execute(w, m)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

//SignupEndpoint . . .
func SignupEndpoint(w http.ResponseWriter, req *http.Request) {

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

	t, err := template.ParseFiles("signup.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, m)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

//LogoutHandler . . .
func LogoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}

//GetIndexEndpoint . . .
func GetIndexEndpoint(w http.ResponseWriter, req *http.Request) {
	fmt.Println("inside index handler")
	http.Redirect(w, req, "/home", http.StatusFound)
	return
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

//GetHomeEndpoint . . .
func GetHomeEndpoint(w http.ResponseWriter, req *http.Request) {

	fmt.Println("inside home handler")

	username := getUserName(req)
	if username != "" {
		fmt.Println("username exists")

		t, err := template.ParseFiles("home.html")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = t.Execute(w, struct{
			UserStr User
			Paths map[int]string
		}{UserStr:User{UserName:username},Paths:getImageMaps("./images/uploaders.txt")})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		fmt.Println("username doesn't exist")
		http.Redirect(w, req, "/login", 302)
	}
}

func UploadEndpoint(w http.ResponseWriter, r *http.Request) {

	username := getUserName(r)
	if username != "" {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		http.Redirect(w, r, "/home", http.StatusFound)

		f, err := os.OpenFile("./images/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()
		io.Copy(f, file)

		AppendStringToFile("./images/uploaders.txt", username+" > "+"./images/"+handler.Filename)

	}else {
		fmt.Println("username doesn't exist")
		http.Redirect(w, r, "/login", 302)
	}
}

func AppendStringToFile(path, text string) error {
      f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
      if err != nil {
              return err
      }
      defer f.Close()

      _, err = f.WriteString(text+"\n")
      if err != nil {
              return err
      }
      return nil
}

func getImageMaps(path string) map[int]string{

  m := make(map[int]string)

  file, err := os.Open(path)
  if err != nil {
      log.Fatal(err)
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  i := 0
  for scanner.Scan() {
      m[i] = strings.TrimSpace(strings.Split(scanner.Text(), ">")[1])
      i++
  }

  if err := scanner.Err(); err != nil {
      log.Fatal(err)
  }

  return m
}

//TODO in the home page display delete and view. display edit/delete only if it's the logged in user

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", GetIndexEndpoint).Methods("GET")
	router.HandleFunc("/home", GetHomeEndpoint).Methods("GET")
	router.HandleFunc("/login", GetLoginEndpoint).Methods("GET")
	router.HandleFunc("/login", LoginEndpoint).Methods("POST")
	router.HandleFunc("/logout", LogoutHandler).Methods("GET")
	router.HandleFunc("/signup", GetSignupEndpoint).Methods("GET")
	router.HandleFunc("/signup", SignupEndpoint).Methods("POST")
	router.HandleFunc("/upload", UploadEndpoint).Methods("POST")

	router.Handle("/images/{img-path}",
    http.StripPrefix("/images/", http.FileServer(http.Dir("./" + "images/"))))

	log.Fatal(http.ListenAndServe(":12345", router))
}
