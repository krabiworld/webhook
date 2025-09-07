package api

import (
	"net/http"
	"webhook/internal/db"
	"webhook/internal/dtos"
	"webhook/internal/models"
	"webhook/internal/utils"

	"github.com/rs/zerolog/log"
)

func CreateWebhook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var webhookDto dtos.Webhook
	if err := utils.BindJSON(r, &webhookDto); err != nil {
		log.Error().Err(err).Msg("Error decoding json")
		utils.WriteError(w, http.StatusBadRequest, "Error decoding json: "+err.Error())
		return
	}

	if err := utils.Validate(&webhookDto); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, err)
		return
	}

	webhookModel := models.Webhook{
		ID:    utils.UUID(),
		Name:  webhookDto.Name,
		Token: utils.UUID(),
	}

	err := db.G[models.Webhook]().Create(r.Context(), &webhookModel)
	if err != nil {
		log.Error().Err(err).Msg("Error creating webhook")
		utils.WriteError(w, http.StatusInternalServerError, "Error creating webhook")
		return
	}

	utils.WriteJSON(w, http.StatusCreated, webhookModel)
}

func GetWebhooks(w http.ResponseWriter, r *http.Request) {

}

func GetWebhook(w http.ResponseWriter, r *http.Request) {

}

func PutWebhook(w http.ResponseWriter, r *http.Request) {

}

func DeleteWebhook(w http.ResponseWriter, r *http.Request) {

}
