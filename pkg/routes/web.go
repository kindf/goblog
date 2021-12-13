package routes

import (
    "net/http"
    "github.com/gorilla/mux"
    "goblog/app/http/controllers"
)

func RegisterWebRoutes(r *mux.Router){
    pc := new(controllers.PagesController)
    r.HandleFunc("/", pc.Home).Methods("GET").Name("home")
    r.HandleFunc("/about", pc.About).Methods("GET").Name("about")
    r.NotFoundHandler = http.HandlerFunc(pc.NotFound)
    
    ac := new(controllers.ArticlesController)
    r.HandleFunc("/articles/{id:[0-9]+}", ac.Show).Methods("GET").Name("articles.show")

    r.HandleFunc("/articles", ac.Index).Methods("GET").Name("articles.index")
}
