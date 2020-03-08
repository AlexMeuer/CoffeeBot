package api

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	MilkTypeDairy = ":cow:"
	MilkTypeOat   = ":ear_of_rice:"
)

type MilkType string

type Interface interface {
	HandleCommand(w http.ResponseWriter, r *http.Request)
}

type CoffeeRound struct {
	Creator User
	Joiners []User
	Milk    MilkType
	Minutes int
}

func NewCoffeeRound(creator User, milk MilkType, slots int, minutes int) *CoffeeRound {
	return &CoffeeRound{
		Creator: creator,
		Joiners: make([]User, 0, slots),
		Milk:    milk,
		Minutes: minutes,
	}
}

func (cr *CoffeeRound) Join(info User) error {
	if cr.AvailableSlots() > 0 {
		for i := range cr.Joiners {
			if cr.Joiners[i].ID == info.ID {
				return errors.New("you have already joined this round")
			}
		}
		cr.Joiners = append(cr.Joiners, info)
		return nil
	}
	return RoundSlotsFilledError{capacity: cap(cr.Joiners)}
}

func (cr *CoffeeRound) AvailableSlots() int {
	return cap(cr.Joiners) - len(cr.Joiners)
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

type RoundSlotsFilledError struct {
	capacity int
}

func (e RoundSlotsFilledError) Error() string {
	return fmt.Sprintf("all %d slots on this round have been filled", e.capacity)
}
