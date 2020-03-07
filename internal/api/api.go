package api

import "net/http"

const (
	MilkTypeDairy = ":cow:"
	MilkTypeOat   = ":ear_of_rice:"
)

type MilkType string

type Interface interface {
	HandleCommand(w http.ResponseWriter, r *http.Request)
}
