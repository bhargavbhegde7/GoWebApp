package main
 
import (
    "encoding/json"
    "log"
    "net/http"
	"html/template"
	"fmt"
 
    "github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))
 
type Person struct {
    ID        string   `json:"id,omitempty"`
    Firstname string   `json:"firstname,omitempty"`
    Lastname  string   `json:"lastname,omitempty"`
    Address   *Address `json:"address,omitempty"`
}

type User struct {
    UserName string //exported field since it begins with a capital letter
	Password string
}
 
type Address struct {
    City  string `json:"city,omitempty"`
    State string `json:"state,omitempty"`
}
 
var people []Person
 
func GetPersonEndpoint(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    for _, item := range people {
        if item.ID == params["id"] {
            json.NewEncoder(w).Encode(item)
            return
        }
    }
    json.NewEncoder(w).Encode(&Person{})
}
 
func GetPeopleEndpoint(w http.ResponseWriter, req *http.Request) {
    json.NewEncoder(w).Encode(people)
}

func CreatePersonEndpoint(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    var person Person
    _ = json.NewDecoder(req.Body).Decode(&person)
    person.ID = params["id"]
    people = append(people, person)
    json.NewEncoder(w).Encode(people)
}
 
func DeletePersonEndpoint(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    for index, item := range people {
        if item.ID == params["id"] {
            people = append(people[:index], people[index+1:]...)
            break
        }
    }
    json.NewEncoder(w).Encode(people)
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
 
func SessionHandler(w http.ResponseWriter, r *http.Request, username string) {
    // Get a session. Get() always returns a session, even if empty.
    session, err := store.Get(r, username)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Set some session values.
    session.Values[username] = username
    //session.Values[42] = 43
    // Save it before we write to the response/return from the handler.
    session.Save(r, w)
}
 
func LoginEndpoint(w http.ResponseWriter, req *http.Request) {
    username := req.FormValue("username")
	passwd := req.FormValue("passwd")
	fmt.Println(username+" -- "+passwd)
	
	//check if the uname and password match
	
	
	//check if session exists, get the session with uname
	//start a new session if doesn't exist
	SessionHandler(w, req, username)
	
	
	
	//TODO send the user to "/home" with all these incoming data
	http.Redirect(w, req, "/home/"+username, http.StatusFound)
		return
	
}

func GetIndexEndpoint(w http.ResponseWriter, req *http.Request) {
    http.Redirect(w, req, "/home", http.StatusFound)
		return 
}

func GetHomeEndpoint(w http.ResponseWriter, req *http.Request) {    
	
	//TODO check if session is valid
	
	params := mux.Vars(req)
    /*var person Person
    _ = json.NewDecoder(req.Body).Decode(&person)*/
    var username = params["username"]
	fmt.Println(username)
	
	//-----------------
	
	t, err := template.ParseFiles("home.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u := User{UserName:username, Password:"abcd"} //define an instance with required field
 
    //t.Execute(os.Stdout, p) //merge template ‘t’ with content of ‘p’
	
	err = t.Execute(w, u)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
    router := mux.NewRouter()
    people = append(people, Person{ID: "1", Firstname: "Nic", Lastname: "Raboy", Address: &Address{City: "Dublin", State: "CA"}})
    people = append(people, Person{ID: "2", Firstname: "Maria", Lastname: "Raboy"})
	
	//TODO redirect all traffic at "/" to "/home" and at home handler check if session is set.
	//if the session is set then send to "/home" otherwise send to "/login"
	
	router.HandleFunc("/", GetIndexEndpoint).Methods("GET")
	router.HandleFunc("/home/{username}", GetHomeEndpoint).Methods("GET")
    router.HandleFunc("/login", GetLoginEndpoint).Methods("GET")
	router.HandleFunc("/login", LoginEndpoint).Methods("POST")
	//router.HandleFunc("/signup", GetSignupEndpoint).Methods("GET")
	//router.HandleFunc("/signup", SignupEndpoint).Methods("POST")
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
    router.HandleFunc("/people/{id}", GetPersonEndpoint).Methods("GET")
    router.HandleFunc("/people/{id}", CreatePersonEndpoint).Methods("POST")
    router.HandleFunc("/people/{id}", DeletePersonEndpoint).Methods("DELETE")
    log.Fatal(http.ListenAndServe(":12345", router))
}