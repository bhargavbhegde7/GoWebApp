package main
 
import (
    "log"
    "net/http"
	"html/template"
	"fmt"
 
    "github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/securecookie"
)

var cookieHandler = securecookie.New(
    securecookie.GenerateRandomKey(64),
    securecookie.GenerateRandomKey(32))

var store = sessions.NewCookieStore([]byte("something-very-secret"))
 
type User struct {
    UserName string 
}

type ErrorMessage struct {
	Message string
}
 
func GetLoginEndpoint(w http.ResponseWriter, req *http.Request) {

	t, err := template.ParseFiles("login.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, nil)

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


func LoginEndpoint(w http.ResponseWriter, req *http.Request) {

    username := req.FormValue("username")
	passwd := req.FormValue("passwd")
	
	//check if the uname and password match
	if username == "bhargav" && passwd == "bhargav"{
		
		//check if session exists, get the session with uname
		//start a new session if doesn't exist
		setSession(username, w)
		
		//TODO send the user to "/home" with all these incoming data
		http.Redirect(w, req, "/home", http.StatusFound)
			return
	}else{
	
		m := ErrorMessage{Message:"wrong credentials"}
	
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

func clearSession(response http.ResponseWriter) {
    cookie := &http.Cookie{
        Name:   "session",
        Value:  "",
        Path:   "/",
        MaxAge: -1,
    }
    http.SetCookie(response, cookie)
}

func LogoutHandler(response http.ResponseWriter, request *http.Request) {
    clearSession(response)
    http.Redirect(response, request, "/", 302)
}

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

func GetHomeEndpoint(w http.ResponseWriter, req *http.Request) {

	fmt.Println("inside home handler")
	
	username := getUserName(req)
    if username != "" {
		fmt.Println("username exists")
			
		u := User{UserName:username}
		
		t, err := template.ParseFiles("home.html")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = t.Execute(w, u)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		fmt.Println("username doesn't exist")
        http.Redirect(w, req, "/login", 302)
    }
}

func main() {
    router := mux.NewRouter()    
	
	router.HandleFunc("/", GetIndexEndpoint).Methods("GET")
	router.HandleFunc("/home", GetHomeEndpoint).Methods("GET")
    router.HandleFunc("/login", GetLoginEndpoint).Methods("GET")
	router.HandleFunc("/login", LoginEndpoint).Methods("POST")
	router.HandleFunc("/logout", LogoutHandler).Methods("GET")
	//router.HandleFunc("/signup", GetSignupEndpoint).Methods("GET")
	//router.HandleFunc("/signup", SignupEndpoint).Methods("POST")
	
    log.Fatal(http.ListenAndServe(":12345", router))
}