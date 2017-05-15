package main
 
import (
    "encoding/json"
    "log"
    "net/http"
	"html/template"
	"fmt"
 
    "github.com/gorilla/mux"
)
 
type Person struct {
    ID        string   `json:"id,omitempty"`
    Firstname string   `json:"firstname,omitempty"`
    Lastname  string   `json:"lastname,omitempty"`
    Address   *Address `json:"address,omitempty"`
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
 
func LoginEndpoint(w http.ResponseWriter, req *http.Request) {
    username := req.FormValue("username")
	passwd := req.FormValue("passwd")
	fmt.Println(username+" -- "+passwd)
	
	//TODO send the user to "/home" with all these incoming data
}

func GetIndexEndpoint(w http.ResponseWriter, req *http.Request) {
    http.Redirect(w, req, "/home", http.StatusFound)
		return 
}

func GetHomeEndpoint(w http.ResponseWriter, req *http.Request) {
    http.Redirect(w, req, "/home", http.StatusFound)
		return 
}

func main() {
    router := mux.NewRouter()
    people = append(people, Person{ID: "1", Firstname: "Nic", Lastname: "Raboy", Address: &Address{City: "Dublin", State: "CA"}})
    people = append(people, Person{ID: "2", Firstname: "Maria", Lastname: "Raboy"})
	
	//TODO redirect all traffic at "/" to "/home" and at home handler check if session is set.
	//if the session is set then send to "/home" otherwise send to "/login"
	
	router.HandleFunc("/", GetIndexEndpoint).Methods("GET")
	router.HandleFunc("/home", GetHomeEndpoint).Methods("GET")
    router.HandleFunc("/login", GetLoginEndpoint).Methods("GET")
	router.HandleFunc("/login", LoginEndpoint).Methods("POST")
	router.HandleFunc("/signup", GetSignupEndpoint).Methods("GET")
	router.HandleFunc("/signup", SignupEndpoint).Methods("POST")
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
    router.HandleFunc("/people/{id}", GetPersonEndpoint).Methods("GET")
    router.HandleFunc("/people/{id}", CreatePersonEndpoint).Methods("POST")
    router.HandleFunc("/people/{id}", DeletePersonEndpoint).Methods("DELETE")
    log.Fatal(http.ListenAndServe(":12345", router))
}