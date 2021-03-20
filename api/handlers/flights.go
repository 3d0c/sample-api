package handlers

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/3d0c/sample-api/api/models"
	"github.com/3d0c/sample-api/pkg/helpers"
)

type flightsHandler struct {
	*models.Flight
}

func flights() *flightsHandler {
	return &flightsHandler{Flight: &models.Flight{}}
}

func (f *flightsHandler) create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) (int, error) {
	var (
		err error
	)

	if err = helpers.Decode(r.Body, &f); err != nil {
		return http.StatusInternalServerError, err
	}

	f.ID = 0

	if err = f.Validate(); err != nil {
		return http.StatusBadRequest, err
	}

	if err = f.Create(); err != nil {
		return http.StatusBadRequest, err
	}

	helpers.NewJsonResponder(w).Write(f)

	return http.StatusOK, nil
}

func (f *flightsHandler) remove(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (int, error) {
	var (
		err error
		fid int
	)

	if fid, err = strconv.Atoi(ps.ByName("id")); err != nil {
		return http.StatusBadRequest, err
	}

	if err = f.Delete(fid); err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}

func (f *flightsHandler) update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (int, error) {
	var (
		err error
		fid int
	)

	if fid, err = strconv.Atoi(ps.ByName("id")); err != nil {
		return http.StatusBadRequest, err
	}

	if err = helpers.Decode(r.Body, &f); err != nil {
		return http.StatusInternalServerError, err
	}

	if err = f.Update(fid); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (f *flightsHandler) search(w http.ResponseWriter, r *http.Request, _ httprouter.Params) (int, error) {
	var (
		name        string
		scheduled   string
		departure   string
		destination string
		search      map[string]interface{} = make(map[string]interface{})
	)

	if name = r.URL.Query().Get("flight_name"); name != "" {
		search["name"] = name
	}

	if scheduled = r.URL.Query().Get("scheduled_date"); scheduled != "" {
		search["scheduled"] = scheduled
	}

	if departure = r.URL.Query().Get("departure"); departure != "" {
		search["departure"] = departure
	}

	if destination = r.URL.Query().Get("destination"); destination != "" {
		search["destination"] = destination
	}

	flights, err := f.Find(search)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	helpers.NewJsonResponder(w).Write(flights)

	return http.StatusOK, nil
}
