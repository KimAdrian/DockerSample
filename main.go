package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type User struct {
	UserID    int    `bson:user_id`
	Name      string `bson:"name"`
	Email     string `bson:"email`
	Interests string `bson:interests`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	mongoUri := os.Getenv("MONGOURL_LOCAL")

	//connect to mongoDB
	client, err := mongo.Connect(options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatalln(err)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", userFetchHandler)
	mux.HandleFunc("POST /", userSaveHandler)
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	err = http.ListenAndServe(":9000", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func userSaveHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatalln("Unable to parse template file: ", err)
	}
}

func userFetchHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatalln("Unable to parse template file: ", err)
	}

	fetchedUser, err := fetchFromCollection(coll)
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, fetchedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}

func saveToCollection(coll *mongo.Collection, userDetails *User) (*mongo.UpdateResult, error) {
	filter := bson.D{{"user_id", 1}}
	opts := options.UpdateOne().SetUpsert(true)
	result, err := coll.UpdateOne(context.TODO(), filter, userDetails, opts)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func fetchFromCollection(coll *mongo.Collection) (User, error) {
	filter := bson.D{{"user_id", 1}}
	var result User
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Println(err)
		return result, err
	}

	return result, nil
}
