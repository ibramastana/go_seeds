package main

import (
	"context"
	"encoding/json"
	"fmt"
	_ "log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName = "go_seeds"
)

type QuizPrize struct {
	Icon_prize string
	Price      int64
}

type Play struct {
	Questions          int
	Played             int
	Durations          int
	Quiz_periode_start time.Time
	Quiz_periode_end   time.Time
	Term_conditions    []string
	Quiz_prize         []QuizPrize
	Sponsors           []string
	Comunities         []string
}

func main() {
	http.HandleFunc("/api/play", requestHandler)
	http.ListenAndServe(":8080", nil)
}

func requestHandler(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{}

	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/"))

	if err != nil {
		fmt.Println(err.Error())
	}

	collection := client.Database(dbName).Collection("play")

	data := map[string]interface{}{}

	err = json.NewDecoder(req.Body).Decode(&data)

	if err != nil {
		fmt.Println(err.Error())
	}

	collection.Drop(ctx)

	term_conditions := []string{"Single Competition, participans compete for prizes", "Start with 50 million virtual capital", "Winner based on highest equuity score", "Participans must follow instagram @seed_finance", "Ticket fee: 100.000/entry (no promo code)"}
	quiz_prize := []QuizPrize{
		{
			Icon_prize: "1",
			Price:      3000000,
		},
		{
			Icon_prize: "2",
			Price:      2000000,
		},
		{
			Icon_prize: "3",
			Price:      1000000,
		},
	}
	sponsor := []string{"Starbuck Icon"}
	community := []string{"Community Icon"}

	insertRes, err := collection.InsertOne(ctx, Play{
		Questions:          15,
		Played:             100,
		Durations:          19,
		Quiz_periode_start: time.Now(),
		Quiz_periode_end:   time.Now(),
		Term_conditions:    term_conditions,
		Quiz_prize:         quiz_prize,
		Sponsors:           sponsor,
		Comunities:         community,
	})

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(insertRes.InsertedID)

	switch req.Method {
	case "POST":
		response, err = createRecord(ctx, data)

		if err != nil {
			response = map[string]interface{}{"error": err.Error()}
		}

		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")

		if err := enc.Encode(response); err != nil {
			fmt.Println(err.Error())
		}
	case "GET":
		response, err = getRecords(collection, ctx)

		if err != nil {
			response = map[string]interface{}{"error": err.Error()}
		}

		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")

		if err := enc.Encode(response); err != nil {
			fmt.Println(err.Error())
		}
	}
}
