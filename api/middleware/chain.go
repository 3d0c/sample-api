package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/3d0c/sample-api/pkg/helpers"
)

type Middlewares func(res http.ResponseWriter, request *http.Request, p httprouter.Params) (int, error)

func Chain(m ...Middlewares) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var (
			err error
		)

		for _, middleware := range m {
			if _, err = middleware(w, r, p); err != nil {
				break
			}
		}

		if err != nil {
			helpers.NewJsonResponder(w).Write(helpers.Error{Error: err.Error()})
			return
		}
	}
}
