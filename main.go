package main

import(
	"fmt"
	"net/http"
	"log"
	"context"
	"encoding/json"
    "strings"
    //"strconv"

    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// i have used atlas mongodb
// has one database named insta_backend_api
// has 2 collections users and posts
// the host is localhost and port 8081

var clientOptions = options.Client().ApplyURI("mongodb+srv://admin:1357@cluster0.qbbwc.mongodb.net/insta_backend_api?retryWrites=true&w=majority")
var client, err = mongo.Connect(context.TODO(), clientOptions)

type User struct {
	ID        string			 
	Username string             
	Email  string
    Password string
}

type Post struct {
	Userid string			 
	Id string             
	Caption  string
    ImageURL string
    Time string
}

func mainpage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w,"mainpage endpoint hit, get request acknowledged")
}

func add_new_user(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("content-type","application/json")
    var u User
    err := json.NewDecoder(r.Body).Decode(&u) // parsing json to structure
    if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
    }
    collection := client.Database("insta_backend_api").Collection("users")
    insertResult, err1 := collection.InsertOne(context.TODO(), u)
    if err != nil {
             log.Fatal(err1)
	}
    json.NewEncoder(w).Encode(insertResult) // response in json  
}


func add_new_post(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("content-type","application/json")
    var p Post
    err := json.NewDecoder(r.Body).Decode(&p) 
    if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
    }
    collection := client.Database("insta_backend_api").Collection("posts")
    insertResult, err1 := collection.InsertOne(context.TODO(), p)
    if err != nil {
             log.Fatal(err1)
	}
    json.NewEncoder(w).Encode(insertResult)   
}

func find_user(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("content-type","application/json")
    var found_user User
    q := r.URL.String()
	split := strings.Split(q, "/")
    prequery := split[len(split)-1] //to get the user ID from the URL
    query := bson.D{{"id", prequery}}
    collection := client.Database("insta_backend_api").Collection("users")
    err = collection.FindOne(context.TODO(), query).Decode(&found_user)
    if err != nil {
        log.Fatal(err)
    }
    json.NewEncoder(w).Encode(found_user) //response in json
}

func find_post(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("content-type","application/json")
    var found_post Post
    q := r.URL.String()
	split := strings.Split(q, "/")
    prequery := split[len(split)-1] //to get the post ID from the URL
    query := bson.D{{"id", prequery}}
    collection := client.Database("insta_backend_api").Collection("posts")
    err = collection.FindOne(context.TODO(), query).Decode(&found_post)
    if err != nil {
        log.Fatal(err)
    }
    json.NewEncoder(w).Encode(found_post) // response in json
}

func find_user_post(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("content-type","application/json")
    var found_user_post []*Post
    q := r.URL.String()
	split := strings.Split(q, "/")
    prequery := split[len(split)-1] //to get the user ID from the URL
    query := bson.D{{"userid", prequery}}
    //limit, _ := strconv.ParseInt(split[len(split)-2],10,64)
    //offset, _ := strconv.ParseInt(split[len(split)-1],10,64) 
    collection := client.Database("insta_backend_api").Collection("posts")
    //findOptions := options.Find()
    //findOptions.SetLimit(limit)
    //findOptions.SetSkip(offset)
    cursor,err2 := collection.Find(context.TODO(), query)
    if err2 != nil {
        log.Fatal(err2)
    }
    for cursor.Next(context.TODO()) {
    var temp_post Post // iterating cursor and adding decoding each found document into struct from json
    err := cursor.Decode(&temp_post)
        if err != nil {
            log.Fatal(err)
        }
        found_user_post = append(found_user_post, &temp_post)
    }
    json.NewEncoder(w).Encode(found_user_post) // response in json
}



func main() {
    err = client.Ping(context.TODO(), nil)
    if err != nil {
    log.Fatal(err)
    }
    fmt.Println("Connected to MongoDB!")
	http.HandleFunc("/", mainpage) //just for testing the initialization
    http.HandleFunc("/users", add_new_user)
    http.HandleFunc("/users/", find_user)
    http.HandleFunc("/posts", add_new_post)
    http.HandleFunc("/posts/", find_post)
    http.HandleFunc("/posts/users/", find_user_post)
	log.Fatal(http.ListenAndServe(":8081",nil))
}