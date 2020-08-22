package controllers

import (
	"fmt"
	"net/http"

	"../models"
	"../views"
)

func NewCustomers(cs models.CustomerService) *Customers {
	return &Customers{
		NewView: views.NewView("bootstrap", "views/customers/create.gohtml"),
		CustomerService: cs,
	}
}

type Customers struct {
	NewView *views.View
	models.CustomerService
}

func (c *Customers) Create(w http.ResponseWriter, r *http.Request) {
	c.NewView.Render(w, nil)
}

type CustomerForm struct {
	FirstName    string `schema:"first_name"`
	LastName string `schema:"last_name"`
	Email string `schema:"email"`
}

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