package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Customer struct {
	gorm.Model
	FirstName string
	LastName string
	Gender string
	Email string `gorm:"not null"`
	Address string
}

type CustomerService interface {
	ByID(id uint) *Customer
	Create(customer *Customer) error
	Update(customer *Customer) error
	Delete(id uint) error
}

func NewCustomerGorm(connectionInfo string) (*CustomerGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	return &CustomerGorm{db}, nil
}

type CustomerGorm struct {
	*gorm.DB
}

func (cg *CustomerGorm) ByID(id uint) *Customer {
	return cg.byQuery(cg.DB.Where("id = ?", id))
}

func (cg *CustomerGorm) byQuery(query *gorm.DB) *Customer {
	ret := &Customer{}
	err := query.First(ret).Error
	switch err {
	case nil:
		return ret
	case gorm.ErrRecordNotFound:
		return nil
	default:
		panic(err)
	}
}

func (cg *CustomerGorm) Create(customer *Customer) error {
	return cg.DB.Create(customer).Error
}

func (cg *CustomerGorm) Update(customer *Customer) error {
	return cg.DB.Save(customer).Error
}

func (cg *CustomerGorm) Delete(id uint) error {
	customer := &Customer{Model: gorm.Model{ID: id}}
	return cg.DB.Delete(customer).Error
}

func (cg *CustomerGorm) DestructiveReset() {
	cg.DropTableIfExists(&Customer{})
	cg.AutoMigrate(&Customer{})
}