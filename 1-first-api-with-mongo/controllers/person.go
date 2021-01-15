package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ernestii/learning-go/1-first-api-with-mongo/models"
)


type PersonController struct {
	personRepo *models.PersonRepository
}

func NewPersonController (pr *models.PersonRepository) PersonController {
	return PersonController{
		personRepo: pr,
	}
}

func (pc *PersonController) GetPeopleEndpoint(res http.ResponseWriter, req *http.Request) {
	people, err := pc.personRepo.GetAll()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(res).Encode(people)
}
