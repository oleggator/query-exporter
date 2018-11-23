package main

import (
	"fmt"
	"github.com/icrowley/fake"
)

func GetFakeValueByTypeName(typeName string) string {
	switch typeName {
	case "name":
		return fake.FirstName()
	case "lastname":
		return fake.LastName()
	case "country":
		return fake.Country()
	case "city":
		return fake.City()
	case "timestamp":
		return fmt.Sprintf("%d-%d-%d", fake.Year(1900, 2000), fake.MonthNum(), 10)
	default:
		return fake.Word()
	}
}
