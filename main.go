package main

func main() {

	r := chi.NewRouter()
	r.Use(middleware)
	r.Get("/", homeHandler) 
	r.Mount("/todo", todoHandlers)

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
	}
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