package domain

import "fmt"

type Person struct {
	Id   int64 `json:"id"`
	Name string `json:"name"`
	Age  int `json:"age"`
	Addresses []Address `json:"addresses,omitempty"`
}

type Address struct {
	Id int64 `json:"id"`
	Street string `json:"street"`
	City string `json:"city"`
}

func (p Person) String() string {
	return fmt.Sprintf("[id=%v, name=%v, age=%v, addresses=%v]", p.Id, p.Name, p.Age, p.Addresses)
}

func (a Address) String() string {
	return fmt.Sprintf("[id=%v, street=%v, city=%v]", a.Id, a.Street, a.City)
}

func (p Person) Smt() int {
	return p.Age * 2
}