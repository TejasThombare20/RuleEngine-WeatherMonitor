package services

import (
	"log"

	"github.com/TejasThombare20/weather-engine/repositories"
)

type UserSerivce struct {
	userRepo *repositories.UserRepository
}

func NewUserSerivce(userRepo *repositories.UserRepository) *UserSerivce {
	return &UserSerivce{userRepo: userRepo}
}

func (s *UserSerivce) AddUserwithCityThreashold(email string, temperature_unit string, thre_temperautes map[string]float64, consecutiveAlerts int) error {

	user, err := s.userRepo.CreateUser(email, temperature_unit)

	if err != nil {
		log.Println("error creating user")
		return err
	}

	for cityName, maxTemp := range thre_temperautes {

		err := s.userRepo.SetCityThreshold(user.ID, cityName, maxTemp, consecutiveAlerts)
		if err != nil {
			log.Printf("error setting city threshold for %s: %v\n", cityName, err)
			// return err
			continue
		}
	}
	return nil
}
