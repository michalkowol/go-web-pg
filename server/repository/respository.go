package repository

import (
	"database/sql"
	"github.com/michalkowol/web-pg/server/domain"
)

type PeopleRepository struct {
	DB *sql.DB
}

func (p PeopleRepository) Find(name string) (*domain.Person, error) {
	var person domain.Person
	err := p.DB.QueryRow("SELECT id, name, age FROM people WHERE name = $1", name).Scan(&person.Id, &person.Name, &person.Age)
	if err != nil {
		return nil, err
	}
	return &person, nil
}

func (p PeopleRepository) List() ([]domain.Person, error) {
	var people []domain.Person
	rows, err := p.DB.Query("SELECT id, name, age FROM people")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var person domain.Person
		err = rows.Scan(&person.Id, &person.Name, &person.Age)
		if err != nil {
			return nil, err
		}
		people = append(people, person)
	}
	return people, err
}



func (p PeopleRepository) ListWithDetails() ([]domain.Person, error) {
	type Row struct {
		PersonId   int64
		PersonName string
		Age        int
		AddressId  sql.NullInt64
		Street     sql.NullString
		CityId     sql.NullInt64
		CityName   sql.NullString
	}

	rows := []Row{}
	dbrow, err := p.DB.Query(`
	SELECT p.id, p.name, p.age, a.id, a.street, c.id, c.name FROM people AS p
	LEFT JOIN addresses_people AS ap ON p.id = ap.person_id
	LEFT JOIN addresses AS a ON ap.address_id = a.id
	LEFT JOIN cities AS c ON a.city_id = c.id
	`)
	if err != nil {
		return nil, err
	}

	for dbrow.Next() {
		var row Row
		err = dbrow.Scan(&row.PersonId, &row.PersonName, &row.Age, &row.AddressId, &row.Street, &row.CityId, &row.CityName)
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}

	addressesMap := make(map[int64][]domain.Address)
	for _, row := range rows {
		addresses := addressesMap[row.PersonId]
		if (row.AddressId.Valid) {
			addressesMap[row.PersonId] = append(addresses, domain.Address{Id: row.AddressId.Int64, Street: row.Street.String, City: row.CityName.String})
		}
	}

	var peopleMap = make(map[int64]domain.Person)
	for _, row := range rows {
		addresses := addressesMap[row.PersonId]
		peopleMap[row.PersonId] = domain.Person{Id: row.PersonId, Name: row.PersonName, Age: row.Age, Addresses: addresses}
	}

	var people []domain.Person
	for  _, person := range peopleMap {
		people = append(people, person)
	}

	return people, nil
}