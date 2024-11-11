package application

import (
	"auther/api/handler/util"
	"auther/internal/database/application"
	applicationSvc "auther/internal/service/application"
	"context"
	"encoding/json"
	"fmt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"log"
	"net/http"
)

func RegisterControllers(mux *http.ServeMux) {
	// TODO: ADD Auth Middleware
	/*
		mux.HandleFunc("POST /api/application", CreateApplicationController)
		mux.HandleFunc("GET /api/application/{id}", ReadApplicationController)
		mux.HandleFunc("PUT /api/application/{id}", UpdateApplicationController)
		//mux.HandleFunc("PATCH /api/application/{id}", UpdateApplicationController)
		mux.HandleFunc("DELETE /api/application/{id}", DeleteApplicationController)

		mux.HandleFunc("POST /api/application/{id}/jwk/generate", GenerateJWKSController)
	*/
}

type applicationClientDTO struct {
	ClientId     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Audiences    []string `json:"audiences"`
}

type applicationDTO struct {
	Id               *string                 `json:"id"`
	ResourceOwnerUrl *string                 `json:"resource_owner_url"`
	Clients          *[]applicationClientDTO `json:"clients"`
	CustomUrls       *[]string               `json:"custom_urls"`
	ActiveKID        *string                 `json:"active_kid"`
	JWKs             *interface{}            `json:"jwks"`
}

func dtoFromEntity(entity *application.Entity) *applicationDTO {
	clientsDTO := make([]applicationClientDTO, len(entity.Clients))
	for i, client := range entity.Clients {
		clientsDTO[i] = applicationClientDTO{
			ClientId:     client.ClientId,
			ClientSecret: client.ClientSecret,
			Audiences:    client.Audiences,
		}
	}
	return &applicationDTO{
		Id:               &entity.Id,
		ResourceOwnerUrl: &entity.ResourceOwnerUrl,
		Clients:          &clientsDTO,
		CustomUrls:       &entity.CustomUrls,
		ActiveKID:        &entity.ActiveKID,
		JWKs:             &entity.JWKs,
	}
}

func dtoToEntity(dto *applicationDTO) *application.Entity {
	clients := make([]application.Client, len(*dto.Clients))
	for i, clientDTO := range *dto.Clients {
		clients[i] = application.Client{
			ClientId:     clientDTO.ClientId,
			ClientSecret: clientDTO.ClientSecret,
			Audiences:    clientDTO.Audiences,
		}
	}
	// convert string to jwk.Set
	if *dto.JWKs != "" {
		marshal, err := json.Marshal(dto.JWKs)
		jwks, err := jwk.Parse(marshal)
		if err != nil {
			log.Printf("Couldn't parse jwks. Here's why: %v\n", err)
		}
		*dto.JWKs = jwks
	}

	return &application.Entity{
		Id:               *dto.Id,
		ResourceOwnerUrl: *dto.ResourceOwnerUrl,
		Clients:          clients,
		CustomUrls:       *dto.CustomUrls,
		ActiveKID:        *dto.ActiveKID,
		JWKs:             *dto.JWKs,
	}
}

func CreateApplicationController(w http.ResponseWriter, r *http.Request) {
	type CreateApplicationRequest struct {
		Id string `json:"id"`
	}
	decoder := json.NewDecoder(r.Body)

	var request CreateApplicationRequest
	err := decoder.Decode(&request)
	if err != nil {
		util.JsonResponse(w, http.StatusBadRequest, fmt.Errorf("error decoding request: %v", err))
		return
	}

	entity, err := applicationSvc.Create(request.Id)
	if err != nil {
		return
	}

	response := dtoFromEntity(entity)

	util.JsonResponseNoCache(w, http.StatusCreated, response)
	return
}

func ReadApplicationController(w http.ResponseWriter, r *http.Request) {
	entity, err := application.Repository.GetByHashKey(context.TODO(), r.PathValue("id"))
	if err != nil {
		util.JsonResponse(w, http.StatusBadRequest, err)
		return
	}
	if entity == nil {
		util.JsonResponse(w, http.StatusNotFound, fmt.Errorf("application not found"))
		return
	}
	response := dtoFromEntity(entity)

	util.JsonResponse(w, http.StatusOK, response)
	return
}

func UpdateApplicationController(w http.ResponseWriter, r *http.Request) {
	// Update Entity
	decoder := json.NewDecoder(r.Body)
	var request applicationDTO
	err := decoder.Decode(&request)
	if err != nil {
		util.JsonResponse(w, http.StatusBadRequest, fmt.Errorf("error decoding request: %v", err))
		return
	}

	entity := dtoToEntity(&request)

	// Save Entity
	err = application.Repository.Save(context.TODO(), entity)
	if err != nil {
		util.JsonResponse(w, http.StatusBadRequest, err)
		return
	}

	// Update JWKS
	err = applicationSvc.UpdateJWKS3(entity)
	if err != nil {
		util.JsonResponse(w, http.StatusBadRequest, err)
		return
	}

	util.JsonResponse(w, http.StatusOK, request)
	return
}

func DeleteApplicationController(w http.ResponseWriter, r *http.Request) {
	err := application.Repository.DeleteByHashKey(context.TODO(), r.PathValue("id"))
	if err != nil {
		util.JsonResponse(w, http.StatusBadRequest, err)
		return
	}
	util.JsonResponse(w, http.StatusNoContent, nil)
	return
}

func GenerateJWKSController(w http.ResponseWriter, r *http.Request) {
	type GenerateJWKSRequest struct {
		Alg jwa.SignatureAlgorithm `json:"alg"`
	}
	decoder := json.NewDecoder(r.Body)
	var request GenerateJWKSRequest
	err := decoder.Decode(&request)
	if err != nil {
		util.JsonResponse(w, http.StatusBadRequest, fmt.Errorf("error decoding request: %v", err))
		return
	}

	entity, err := application.Repository.GetByHashKey(context.TODO(), r.PathValue("id"))
	if err != nil {
		util.JsonResponse(w, http.StatusBadRequest, err)
		return
	}
	if entity == nil {
		util.JsonResponse(w, http.StatusNotFound, fmt.Errorf("application not found"))
		return
	}
	err = applicationSvc.AddGeneratedKey(entity, request.Alg)
	if err != nil {
		util.JsonResponse(w, http.StatusBadRequest, err)
		return
	}
	util.JsonResponse(w, http.StatusOK, dtoFromEntity(entity))
	return
}
