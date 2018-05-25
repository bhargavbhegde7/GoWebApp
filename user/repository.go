package user

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Repository ...
type Repository struct{}

// SERVER the DB server
const SERVER = "localhost:27017"

// DBNAME the name of the DB instance
const DBNAME = "picturestore"

// COLLECTION the name of the collection
const COLLECTION = "users"

// GetUsers returns the list of Users
func (r Repository) GetUsers() Users {
	session, err := mgo.Dial(SERVER)
	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}
	defer session.Close()
	c := session.DB(DBNAME).C(COLLECTION)
	results := Users{}
	if err := c.Find(nil).All(&results); err != nil {
		fmt.Println("Failed to write results:", err)
	}

	return results
}

// AddUser inserts a User in the DB
func (r Repository) AddUser(user User) bool {
	session, err := mgo.Dial(SERVER)
	defer session.Close()

	user.ID = bson.NewObjectId()
	session.DB(DBNAME).C(COLLECTION).Insert(user)

	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}