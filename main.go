package main

import (
	"fmt"
	"net/http"

	"./controllers"
	"./models"
	"./views"

	"github.com/gorilla/mux"
)

const (
	host = "localhost"
	port = 5432
	user = "sultan"
	password = "narivo979"
	dbname = "gowebapp_dev"
)

var homeView *views.View

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)
	// cg, err := models.NewCustomerGorm(psqlInfo)
	// if err != nil {
	// 	panic(err)
	// }
	services, err := models.NewServices(
		models.WithGorm("postgres", psqlInfo),
		models.WithCustomer(),
	)
	if err != nil {
		panic(err)
	}
	defer services.Close()

	r := mux.NewRouter()
	homeView = views.NewView("bootstrap", "views/home.gohtml")
	customersC := controllers.NewCustomers(services.Customer, r)

	r.HandleFunc("/", customersC.Index).Methods("GET")
	r.HandleFunc("/customers/create", customersC.Create).Methods("GET")
	r.HandleFunc("/customers/{id:[0-9]+}/edit", customersC.Edit).Methods("GET")
	r.HandleFunc("/customers/{id:[0-9]+}/update", customersC.Update).Methods("POST")
	r.HandleFunc("/customers/store", customersC.Store).Methods("POST")
	r.HandleFunc("/customers/{id:[0-9]+}", customersC.Show).Methods("GET").Name("show_customer")
	http.ListenAndServe(":3000", r)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	homeView.Render(w, nil)
}