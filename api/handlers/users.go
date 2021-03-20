package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/3d0c/sample-api/api/models"
	"github.com/3d0c/sample-api/pkg/helpers"
)

type usersHandler struct {
	*models.User
}

func users() *usersHandler {
	return &usersHandler{User: &models.User{}}
}

func (u *usersHandler) create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) (int, error) {
	var (
		err error
	)

	if err = helpers.Decode(r.Body, &u); err != nil {
		return http.StatusInternalServerError, err
	}

	if err = u.Validate(); err != nil {
		return http.StatusBadRequest, err
	}

	if err = u.Create(); err != nil {
		return http.StatusBadRequest, err
	}

	// Hide password from output
	u.Password = ""

	helpers.NewJsonResponder(w).Write(u)

	return http.StatusOK, nil
}

func (u *usersHandler) login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) (int, error) {
	var (
		err error
	)

	if err = helpers.Decode(r.Body, &u); err != nil {
		return http.StatusInternalServerError, err
	}

	if err = u.Validate(); err != nil {
		return http.StatusBadRequest, err
	}

	if u.User, err = u.Find(); err != nil {
		return http.StatusNotFound, err
	}

	token, err := u.GenerateJWT()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	helpers.NewJsonResponder(w).Write(token)

	return http.StatusOK, nil
}
