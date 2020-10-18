package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	startTime         string
	endTime           string
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

//This function sends json object as a string if the it is present in the collection by searching by _id
func findByIdAndSend(id string) ([]byte, string) {

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
		return nil, id + " is an invalid object id"

	}

	filterCursor, err := collectionMeetings.Find(ctx, bson.M{"_id": objectId})
	if err != nil {
		log.Fatal(err)
	}
	var meeting []bson.M
	if err = filterCursor.All(ctx, &meeting); err != nil {
		log.Fatal(err)
	}
	fmt.Println(meeting)
	jsonResponse, err := json.Marshal(meeting)

	if err != nil {
		return nil, err.Error()
	} else {
		if string(jsonResponse) != "null" {
			fmt.Println(string(jsonResponse))
			return jsonResponse, "passed"
		} else {
			return nil, "No Meeting found corresponding to the given id"
		}
	}

}

func scheduleMeeting(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		r.ParseForm()
		_, participant := r.Form["participant"]
		_, idPresent := r.Form["id"]
		_, startTime := r.Form["startTime"]
		_, stopTime := r.Form["stopTime"]
		if !participant && !startTime && !stopTime {
			idPresent = true
		}
		if idPresent {
			//Send Participant information using given ID
			response, err := findByIdAndSend(r.URL.Path[len("/meetings/"):])

			if err != "passed" {
				fmt.Fprintf(w, "Something went wrong")

			} else {
				w.Header().Set("content-type", "application/json")
				w.Write(response)

			}

		} else {
			if participant {
				fmt.Fprintf(w, "sending information for participant: %v", r.Form["participant"])
			} else {
				if startTime && stopTime {
					fmt.Fprintf(w, "sending information from %v to %v", r.Form["startTime"], r.Form["stopTime"])
				} else {
					fmt.Fprintf(w, "Please specify both the parameters i.e startTime and stopTime")
				}
			}

		}

	case "POST":

		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		var newMeet Meeting
		err = json.Unmarshal(b, &newMeet)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		output, err := json.Marshal(newMeet)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.Write(output)
		// decoder := json.NewDecoder(r.Body)

		// var newMeeting Meeting
		// err := decoder.Decode(&newMeeting)
		// fmt.Println(r.Body)
		// if err != nil {
		// 	fmt.Fprintf(w, err.Error())

		// } else {

		// 	if newMeeting.Title == "" {
		// 		fmt.Fprintf(w, "Please give meeting title")
		// 	} else {
		// 		fmt.Println("Creating new meeting")

		// 		b, err := ioutil.ReadAll(r.Body)
		// 		defer r.Body.Close()

		// 		if err != nil {
		// 			fmt.Fprintf(w, err.Error())
		// 		}

		// 		fmt.Println(b)
		// 		// title := newMeeting.Title
		// 		// // participants := newMeeting.Participants
		// 		// startTime := newMeeting.startTime
		// 		// endTime := newMeeting.endTime
		// 		// creationTimestamp := time.Now()
		// 		// fmt.Println(newMeeting)

		// 		// m1 := &Meeting{Title: title, startTime: startTime, endTime: endTime, creationTimestamp: creationTimestamp}
		// 		// insertResult, err := collectionMeetings.InsertOne(context.TODO(), m1)
		// 		// if err != nil {
		// 		// 	log.Fatal(err)
		// 		// }
		// 		// fmt.Fprintf(w, "Meeting Created with ID as: %v\n", insertResult.InsertedID)
		// 		// fmt.Println("Meeting Created as ")

		// 	}

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

}

func main() {
	mongoInit()

	http.HandleFunc("/meetings/", scheduleMeeting)

	log.Fatal(http.ListenAndServe(":3000", nil))
}
