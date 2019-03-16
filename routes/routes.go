package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/tryOne/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var R *mux.Router

type student struct {
	ID         string    `json:"_id"`
	Name       string    `json:"name"`
	RollNo     string    `json:"rollNo"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
}

var collection = db.DB.Collection("student")

func init() {
	R = mux.NewRouter()

	R.HandleFunc("/", list).Methods("GET")
	R.HandleFunc("/add", addStudent).Methods("POST")
	R.HandleFunc("/update", updateStudent).Methods("PUT")
	R.HandleFunc("/delete", deleteStudent).Methods("DELETE")
}

func list(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.ParseInt(r.FormValue("limit"), 10, 64)
	skip, _ := strconv.ParseInt(r.FormValue("skip"), 10, 64)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	studentCursor, _ := collection.Find(ctx, bson.M{}, options.Find().SetSkip(skip).SetLimit(limit))
	defer studentCursor.Close(ctx)

	var students []student
	for studentCursor.Next(nil) {
		student := student{}
		err := studentCursor.Decode(&student)
		if err != nil {
			log.Fatal("Decode error ", err)
		}
		students = append(students, student)
	}

	jsonRes, _ := json.Marshal(students)
	w.Write(jsonRes)
}

func addStudent(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	decoder := json.NewDecoder(r.Body)
	s := student{}
	decoder.Decode(&s)

	res, err := collection.InsertOne(ctx, bson.D{
		{"name", s.Name},
		{"rollNo", s.RollNo},
		{"createdAt", s.CreatedAt},
		{"modifiedAt", s.ModifiedAt},
	})

	if err != nil {
		fmt.Println("Error while insert" + err.Error())
	}

	w.Write([]byte("Student added successfully" + res.InsertedID.(primitive.ObjectID).Hex()))
}

func updateStudent(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	body := student{}
	inDB := student{}

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&body)

	fmt.Println("#########################")
	fmt.Println(body)

	objectIDS, err := primitive.ObjectIDFromHex(body.ID)
	idDoc := bson.D{{"_id", objectIDS}}

	err = collection.FindOne(ctx, idDoc).Decode(&inDB)

	if err != nil {
		fmt.Errorf("updateTask: couldn't decode task from db: %v", err)
	}

	_, err = collection.UpdateOne(
		ctx,
		idDoc,
		bson.D{
			{"$set", bson.D{
				{"name", body.Name},
				{"rollNo", body.RollNo}},
			},
			{"$currentDate", bson.D{{"modifiedAt", true}}},
		},
	)

	if err != nil {
		fmt.Println("Error while update" + err.Error())
	}

	w.Write([]byte("Student updated successfully"))
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	decoder := json.NewDecoder(r.Body)
	s := student{}
	decoder.Decode(&s)

	objectIDS, err := primitive.ObjectIDFromHex(s.ID)

	if err != nil {
		fmt.Println("deleteTask: couldn't convert student ID from input")
	}

	_, err = collection.DeleteOne(ctx, bson.D{{"_id", objectIDS}})

	if err != nil {
		fmt.Println("deleteTask: couldn't delete student from db" + err.Error())
	}

	w.Write([]byte("Deleted Succesfully"))
}
