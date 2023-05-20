package utils

import (
	"fmt"

	models "github.com/dfsantos-source/tangent-backend/models"
)

type BusinessSet struct {
	Businesses []models.Business
	Set        map[string]bool
}

func (s *BusinessSet) Add(element models.Business) {
	business := element
	if !s.Set[business.Id] {
		s.Businesses = append(s.Businesses, business)
		s.Set[business.Id] = true
	} else {
		return
	}
}

func (s *BusinessSet) AddBusinesses(businesses []models.Business) {
	fmt.Println(businesses)
	allBusinesses := businesses
	for i := 0; i < len(allBusinesses); i++ {
		s.Add(allBusinesses[i])
	}
}

func (s *BusinessSet) GetBusinesses() *[]models.Business {
	return &s.Businesses
}
