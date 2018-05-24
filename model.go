package main

import "gopkg.in/mgo.v2/bson"

//User represents a user of the site
type User struct {
	ID       bson.ObjectId `bson:"_id"`
	Email    string        `json:"email"`
	Password string        `json:"password"`
}

//Users is an array of User
type Users []User
