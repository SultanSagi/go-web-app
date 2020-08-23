package controllers

import (
	"github.com/gorilla/mux"
	"fmt"
	"net/http"
	"strconv"

	"../models"
	"../views"
)

func NewCustomers(cs models.CustomerService) *Customers {
	return &Customers{
		NewView: views.NewView("bootstrap", "views/customers/create.gohtml"),
		ShowView: views.NewView("bootstrap", "views/customers/show.gohtml"),
		CustomerService: cs,
	}
}

type Customers struct {
	NewView *views.View
	ShowView *views.View
	models.CustomerService
}

// GET /customers/create
func (c *Customers) Create(w http.ResponseWriter, r *http.Request) {
	c.NewView.Render(w, nil)
}

type CustomerForm struct {
	FirstName    string `schema:"first_name"`
	LastName string `schema:"last_name"`
	Email string `schema:"email"`
}

// POST /customers
func (c *Customers) Store(w http.ResponseWriter, r *http.Request) {
	form := CustomerForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	customer := &models.Customer{
		FirstName: form.FirstName,
		LastName: form.LastName,
		Email: form.Email,
	}
	if err := c.CustomerService.Create(customer); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, customer)
}

// GET /customers/:id
func (c *Customers) Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
		return
	}
	customer := c.CustomerService.ByID(uint(id))
	var vd views.Data
	vd.Yield = customer
	c.ShowView.Render(w, vd)
	// fmt.Fprintln(r, customer)
}