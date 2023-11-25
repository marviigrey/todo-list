package main 

import(
	"encoding/json"
	"log"
	"net/http"
	"string"
	"time"
	"context"
	"os"
	"os/signal"
	"github.com/thedevsaddam/renderer"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/middleware"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var rnd *renderer.renderer
var db *mgo.Database

const(
	hostname	   string = "localhost:27017"
	dbName		   string = "demo-todo"
	collectionName string = "todo"
	port		   string = ":9000"
)

type (
	todoModel struct {
		ID			bson.ObjectId	`json:"_id,omitempty" bson
		Title		string	`bson:"title"`
		Completed	bool	`bson:""`
		CreatedAt	time.Time `bson:"createdAt`
	}
	todo struct {
		ID	 string `json:"id"`
		Title string `json:"title"`
		Completed bool `json:"completed"`
		CreatedAt time.Time `json:"created_at"`

	}
)

func init() {
	rnd = renderer.New()
	sess, err := mgo.Dial(hostname) //creating a connection string to mongoDb using monotonic mode.
	
	checkErr(err)
	sess.SetMode(mgo.Monotonic, true)
	db = sess.DB(dbName)
}



func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}