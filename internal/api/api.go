package api

import "net/http"

type Interface interface {
	HandleCommand(w http.ResponseWriter, r *http.Request)
}
