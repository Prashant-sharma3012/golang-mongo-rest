package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/tryOne/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, _ := collection.Find(ctx, bson.M{})

	students := make([]student, 1)

	for res.Next(context.Background()) {
		s := student{}
		res.Decode(&s)
		students = append(students, s)
	}

	fmt.Println(students)
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
		fmt.Println("Error while insert")
	}

	w.Write([]byte("Student added successfully" + res.InsertedID.(primitive.ObjectID).Hex()))
}

func updateStudent(w http.ResponseWriter, r *http.Request) {

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
		fmt.Println("deleteTask: couldn't delete student from db")
	}

	w.Write([]byte("Deleted Succesfully"))
}
