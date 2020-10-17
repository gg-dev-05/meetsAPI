package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Defining Participant
type Participant struct {
	Name  string
	Email string
	RSVP  string
}

//Defining Meeting
type Meeting struct {
	Title string
	// Participants      []Participant //TODO: Add participants
	startTime         time.Time
	endTime           time.Time
	creationTimestamp time.Time
}

var collectionMeetings *mongo.Collection
var collectionParticipants *mongo.Collection

var ctx = context.TODO()

//connect to mongoDB database
func mongoInit() {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collectionMeetings = client.Database("appointy").Collection("meetings")
	collectionParticipants = client.Database("appointy").Collection("participants")

}

func scheduleMeeting(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		r.ParseForm()
		_, participant := r.Form["participant"]

		if participant {
			query := r.URL.Query()
			fmt.Fprintf(w, "Sending meeting details of the participant with email: %v\n", query["participant"][0])
		} else {
			fmt.Fprintf(w, "This participant does not exist in the database")

		}

	case "POST":
		decoder := json.NewDecoder(r.Body)

		var newMeeting Meeting
		err := decoder.Decode(&newMeeting)

		if err != nil {
			fmt.Fprintf(w, err.Error())

		} else {

			if newMeeting.Title == "" {
				fmt.Fprintf(w, "Please give meeting title")
			} else {
				fmt.Println("Here")
				title := newMeeting.Title
				// participants := newMeeting.Participants
				startTime := newMeeting.startTime
				endTime := newMeeting.endTime
				creationTimestamp := time.Now()
				m1 := &Meeting{Title: title, startTime: startTime, endTime: endTime, creationTimestamp: creationTimestamp}
				insertResult, err := collectionMeetings.InsertOne(context.TODO(), m1)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Fprintf(w, "Meeting Created with ID as: %v\n", insertResult.InsertedID)
				fmt.Println("Meeting Created with ID as: ", insertResult.InsertedID)

			}
		}

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

}

func main() {
	mongoInit()

	http.HandleFunc("/meetings", scheduleMeeting)

	log.Fatal(http.ListenAndServe(":3000", nil))
}
