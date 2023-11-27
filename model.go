package main 

import(
	"encoding/json"
	"log"
	"net/http"
	"strings"
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
		ID			bson.ObjectId	`bson:"_id,omitempty"`
		Title		string	`bson:"title"`
		Completed	bool	`bson:"completed"`
		CreatedAt	time.Time `bson:"createdAt"`
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

func homeHandler(w *http.ResponseWriter,r *http.Request) { //function to handle our home page
	err := rnd.Template(w, http.StatusOK, []string{"static/home/tpl"}, nil)
	checkErr(err)
}

//function to handle retrieved data from database.
func fetchTodos(w http.ResponseWriter, r *http.Request){
	todos := []todoModel{} // declared a variable 
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
//send response back to user frontend.
func createTodo(w, r *http.Request) {
	var t todo
	if err := json.NewDecoder(r.body).Decode(&t); err!=nil {
		rnd.JSON(w, http.StatusProcessing, err)
		return
	}
	if t.Title == ""{
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message":"title is required.",
		})
		return
	}
	td := todoModel{
		ID: bson.NewObjectId(),
		Title: t.Title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	if err := db.C(collectionName).Insert(&tm); err!=nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message":"Error creating todo",
			"error":err,
		})
		return
	}

	rnd.JSON(w, http.StatusCreated, renderer.M{
		"message":"Successfully created todo",
		"todo-id": td.ID.Hex(),
	})


}

func deleteTodo(w http.ResponseWriter, r *Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if !bson.IsObjectIdHex(id){
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message":"Invalid ID format",
		})
		return
	}
	if err := db.C(collectionName).RemoveId(bson.ObjectIdHex(Id)); err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message":"Error deleting todo",
			"error":err,
	})
	return

	}
	rnd.JSON(w, http.StatusOK, renderer.M{
		"message":"Deleted Todo Successfully",
		})


}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))

	if !bson.IsObjectIdHex(id){
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message":"Invalid ID format",
	})
	return
}
var t todo
if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
	rnd.JSON(w, http.StatusProcessing, err)
	return
}
if t.Title == ""{
	rnd.JSON(w, http.StatusBadRequest, renderer.M{
		"Message": "Title field required.",
	})
	return
}

if err := db.C(collectionName).
	Update(
		bson.M{"_id": bson.ObjectIdHex(id)},
		bson.M{"title": t.Title, "completed": t.Completed},
	); err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message":"Error updating todo",
			"error":err,
			})

			return
	}

}


	


func main() {

	stopChan := make(chan os.Signal) //stop server gracefully using a go channel.
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

func todoHandlers() http.Handler{
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
	r.Get("/", fetchTodos)
	r.Post("/", createTodo)
	r.Put("/{id}", updateTodo)
	r.Delete("/{id}", deleteTodo)


})
return rg
}