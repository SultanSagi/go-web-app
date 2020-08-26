package controllers

import (
	"github.com/gorilla/mux"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"../models"
	"../views"
)

func NewCustomers(cs models.CustomerService, r *mux.Router) *Customers {
	return &Customers{
		NewView: views.NewView("bootstrap", "views/customers/create.gohtml"),
		ShowView: views.NewView("bootstrap", "views/customers/show.gohtml"),
		EditView: views.NewView("bootstrap", "views/customers/edit.gohtml"),
		IndexView: views.NewView("bootstrap", "views/customers/index.gohtml"),
		CustomerService: cs,
		r: r,
	}
}

type Customers struct {
	NewView *views.View
	ShowView *views.View
	EditView *views.View
	IndexView *views.View
	models.CustomerService
	r *mux.Router
}

type SearchForm struct {
	FirstName    string `schema:"first_name"`
	LastName string `schema:"last_name"`
}

// GET /customers/create
func (c *Customers) Index(w http.ResponseWriter, r *http.Request) {
	// form := SearchForm{}
	// if err := parseForm(r, &form); err != nil {
	// 	panic(err)
	// }
	var customers []models.Customer
	var err error
	FirstName := r.URL.Query().Get("first_name")
	LastName := r.URL.Query().Get("last_name")
	FilterOrderBy := "id"
	FilterSort := "desc"
	var columnList = models.DataList{"id", "first_name", "last_name", "birth_date", "gender", "email", "address"}
	var sortList = models.DataList{"asc", "desc"}
	if r.URL.Query().Get("order_by") != "" && columnList.Has(r.URL.Query().Get("order_by")) {
		FilterOrderBy = r.URL.Query().Get("order_by")
	}
	if r.URL.Query().Get("sort") != "" && sortList.Has(r.URL.Query().Get("sort")) {
		FilterSort = r.URL.Query().Get("sort")
	}
	if FirstName != "" && LastName != "" {
		customers, err = c.CustomerService.Search(FirstName, LastName, FilterOrderBy, FilterSort)
	} else {
		customers, err = c.CustomerService.Fetch(FilterOrderBy, FilterSort)
	}
	if err != nil {
		log.Println(err)
	}
	// fmt.Println(customers)
	var vd views.Data
	vd.Yield = customers
	vd.SearchFirstName = FirstName
	vd.SearchLastName = LastName
	vd.FilterOrderBy = FilterOrderBy
	vd.FilterSort = FilterSort
	c.IndexView.Render(w, vd)
}

// GET /customers/create
func (c *Customers) Create(w http.ResponseWriter, r *http.Request) {
	c.NewView.Render(w, nil)
}

type CustomerForm struct {
	FirstName    string `schema:"first_name"`
	LastName string `schema:"last_name"`
	Email string `schema:"email"`
	Gender string `schema:"gender"`
	Address string `schema:"address"`
	BirthDate string `schema:"birth_date"`
}

// POST /customers
func (c *Customers) Store(w http.ResponseWriter, r *http.Request) {
	form := CustomerForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	timeResult, _ := time.Parse("2006-01-02", form.BirthDate)
	customer := models.Customer{
		FirstName: form.FirstName,
		LastName: form.LastName,
		Email: form.Email,
		Gender: strings.TrimSpace(form.Gender),
		Address: form.Address,
		BirthDate: timeResult,
	}
	if err := c.CustomerService.Create(&customer); err != nil {
		var vd views.Data
		vd.Yield = err
		fmt.Println(err)
		c.NewView.Render(w, vd)
		return
	}
	url, err := c.r.Get("show_customer").URL("id", fmt.Sprintf("%v", customer.ID))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
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
	customer, err := c.CustomerService.ByID(uint(id))
	var vd views.Data
	vd.Yield = customer
	c.ShowView.Render(w, vd)
	// fmt.Fprintln(r, customer)
}

func (c *Customers) Edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
		return
	}
	customer, err := c.CustomerService.ByID(uint(id))
	var vd views.Data
	vd.Yield = customer
	c.EditView.Render(w, vd)
}

func (c *Customers) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
		return
	}
	customer, err := c.CustomerService.ByID(uint(id))
	var vd views.Data
	vd.Yield = customer
	var form CustomerForm
	if err := parseForm(r, &form); err != nil {
		fmt.Println("130")
		fmt.Println(err)
		var vd1 views.Data
		vd1.Yield = err
		c.EditView.Render(w, vd1)
		return
	}
	customer.FirstName = form.FirstName
	customer.LastName = form.LastName
	customer.Email = form.Email
	customer.Gender = form.Gender
	customer.Address = form.Address
	err = c.CustomerService.Update(customer)
	if err != nil {
		fmt.Println("144")
		fmt.Println(err)
		var vd2 views.Data
		vd2.Errors = fmt.Sprintf("%s", err)
		vd2.Yield = customer
		c.EditView.Render(w, vd2)
		return
	}
	c.EditView.Render(w, vd)
}