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



func checkErr(err error) { //create an error function.
	if err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w *http.ResponseWriter,r *http.Request) {
	err := rnd.Template(w, http.StatusOK, []string{"static/home/tpl"}, nil)
	checkErr(err)
}
func fetchTodos(w *http.ResponseWriter, r *http.Request){
	todos := []todoModel{}
	if err := db.C(collectionName).Find(bson.M{}.All(&todos)); err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message":"Failed to fetch Todo",
			"error":err,
		})
		return 
	} 
	todoList := []todo{}

	for _, t := range todos{
		todoList = append(todoList, todo{
			ID:         t.ID.Hex(),
			Title:      t.Title,
			Completed:  t.Completed,
			CreatedAt:  t.CreatedAt,
			})
	}
	rnd.JSON(w, http.StatusOK, renderer.M{
		"data": todoList,
		})

}


func main() {

	stopChan := make(chan, os.signal) //stop server gracefully using a go channel.
	signal.Notify(stopChan, os.Interrupt) //making use of the golang os package to recieve and send signals for terminating our golang
	//program gracefully

	r := chi.NewRouter() //initializes a new route handler to handle http requests.

	r.Use(middleware) //help to apply middleware to all registered router.

	r.Get("/", homeHandler) //router for handling GET http request

	r.Mount("/todo", todoHandlers) // mounts group of routes under the /todo prefix.

	srv := &http.Server{
		Addr: port,
		Handler: r,
		ReadTimeout: 60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout: 60 * time.Second,
	}

	go func(){
		log.Println("listening on port", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("%s\n", err)
		}
	}()

	<-stopChan //channel for the stopChan variable
	log.Println("shutting Down server gracefully.")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)
	defer cancel()
	log.Println("server gracefully stopped")

}

func todo.todoHandlers() http.Handler{
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
	r.Get("/" fetchTodos)
	r.Post("/",createTodo)
	r.Put("/{id}", updateTodo)
	r.Delete("/{id}", deleteTodo)


})
return rg
}