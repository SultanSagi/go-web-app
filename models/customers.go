package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type Customer struct {
	gorm.Model
	FirstName string `gorm:"type:varchar(100);not null"`
	LastName string `gorm:"type:varchar(100);not null"`
	BirthDate time.Time `gorm:"not null"`
	Gender string `gorm:"not null"`
	Email string `gorm:"not null"`
	Address string `gorm:"type:varchar(200)"`
}

type customerValFn func(*Customer) error

type CustomerService interface {
	CustomerDB
}

type CustomerDB interface {
	ByID(id uint) (*Customer, error)
	Fetch(FilterOrderBy string, FilterSort string) ([]Customer, error)
	Search(firstName string, lastName string, FilterOrderBy string, FilterSort string) ([]Customer, error)
	Create(customer *Customer) error
	Update(customer *Customer) error
	Delete(id uint) error
}

func NewCustomerService(db *gorm.DB) CustomerService {
	return &customerService{
		CustomerDB: &customerValidator{&CustomerGorm{db}},
	}
}

type customerService struct {
	CustomerDB
}

type customerValidator struct {
	CustomerDB
}

func (cv *customerValidator) Create(customer *Customer) error {
	err := runCustomerValFuncs(customer,
		cv.emailRequired,
		cv.firstNameRequired,
		cv.lastNameRequired,
		cv.validGender,
		cv.BirthDateRequired,
		cv.DateRange)
	if err != nil {
		return err
	}
	return cv.CustomerDB.Create(customer)
}

func (cv *customerValidator) Update(customer *Customer) error {
	err := runCustomerValFuncs(customer,
		cv.emailRequired,
		cv.firstNameRequired,
		cv.lastNameRequired,
		cv.validGender,
		cv.BirthDateRequired,
		cv.DateRange)
	if err != nil {
		return err
	}
	return cv.CustomerDB.Update(customer)
}

func (cv *customerValidator) emailRequired(c *Customer) error {
	if c.Email == "" {
		return fmt.Errorf("Email is required")
	}
	return nil
}

func (cv *customerValidator) firstNameRequired(c *Customer) error {
	if c.FirstName == "" {
		return fmt.Errorf("First Name is required")
	}
	return nil
}

func (cv *customerValidator) lastNameRequired(c *Customer) error {
	if c.LastName == "" {
		return fmt.Errorf("Last Name is required")
	}
	return nil
}

type DataList []string

func (list DataList) Has(a string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func (cv *customerValidator) validGender(c *Customer) error {
	var gl = DataList{"Male", "Female"}
	if !gl.Has(c.Gender) {
		return fmt.Errorf("Gender should be Male or Female")
	}
	return nil
}

func (cv *customerValidator) BirthDateRequired(c *Customer) error {
	if fmt.Sprintf("%s", c.BirthDate) == "" {
		return fmt.Errorf("BirthDate is required")
	}
	return nil
}

func (cv *customerValidator) DateRange(c *Customer) error {
	// today, _ := time.Parse("2006-01-02", "2001-11-30")
	from := time.Now().Add(-60*365*24*time.Hour)
	to := time.Now().Add(-18*365*24*time.Hour)
	if !(c.BirthDate.After(from) && c.BirthDate.Before(to)) {
		return fmt.Errorf("BirthDate range is 18 til 60 years")
	}
	return nil
}

var _ CustomerDB = &CustomerGorm{}

type CustomerGorm struct {
	db *gorm.DB
}

func NewCustomerGorm(connectionInfo string) (*CustomerGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	return &CustomerGorm{db}, nil
}

func (cg *CustomerGorm) Search(firstName string, lastName string, FilterOrderBy string, FilterSort string) ([]Customer, error) {
	var customers []Customer
	err := cg.db.Where("first_name LIKE ? AND last_name LIKE ?", "%"+firstName+"%", "%"+lastName+"%").Order(FilterOrderBy+" "+FilterSort).Find(&customers).Error
	if err != nil {
		return nil, err
	}
	return customers, nil
}

func (cg *CustomerGorm) Fetch(FilterOrderBy string, FilterSort string) ([]Customer, error) {
	var customers []Customer
	err := cg.db.Order(FilterOrderBy+" "+FilterSort).Find(&customers).Error
	if err != nil {
		return nil, err
	}
	return customers, nil
}

func (cg *CustomerGorm) ByID(id uint) (*Customer, error) {
	var customer Customer
	db := cg.db.Where("id = ?", id)
	err := db.First(&customer).Error
	return &customer, err
	// return cg.byQuery(cg.db.Where("id = ?", id))
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
	return cg.db.Create(customer).Error
}

func (cg *CustomerGorm) Update(customer *Customer) error {
	return cg.db.Save(customer).Error
}

func (cg *CustomerGorm) Delete(id uint) error {
	customer := &Customer{Model: gorm.Model{ID: id}}
	return cg.db.Delete(customer).Error
}

func (cg *CustomerGorm) DestructiveReset() {
	cg.db.DropTableIfExists(&Customer{})
	cg.db.AutoMigrate(&Customer{})
}

type customerValFunc func(*Customer) error

func runCustomerValFuncs(customer *Customer, fns ...customerValFunc) error {
	for _, fn := range fns {
		if err := fn(customer); err != nil {
			return err
		}
	}
	return nil
}