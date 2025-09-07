package api

import (
	"errors"
	"io"
	"net/http"
	"webhook/internal/db"
	"webhook/internal/dtos"
	"webhook/internal/models"
	"webhook/internal/utils"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func CreateWebhook(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("Error closing body")
		}
	}(r.Body)

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
	webhookModels, err := db.G[models.Webhook]().Find(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Error getting webhooks")
		utils.WriteError(w, http.StatusInternalServerError, "Error getting webhooks")
		return
	}

	utils.WriteJSON(w, http.StatusOK, webhookModels)
}

func GetWebhook(w http.ResponseWriter, r *http.Request) {
	webhookModel, err := db.G[models.Webhook]().Where("id = ?", r.PathValue("id")).Find(r.Context())
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.WriteError(w, http.StatusNotFound, "Webhook not found")
		return
	} else if err != nil {
		log.Error().Err(err).Msg("Error getting webhook")
		utils.WriteError(w, http.StatusInternalServerError, "Error getting webhook")
		return
	}

	utils.WriteJSON(w, http.StatusOK, webhookModel)
}

func PutWebhook(w http.ResponseWriter, r *http.Request) {

}

func DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	rowsAffected, err := db.G[models.Webhook]().Where("id = ?", r.PathValue("id")).Delete(r.Context())
	if rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, "Webhook not found")
		return
	} else if err != nil {
		log.Error().Err(err).Msg("Error deleting webhook")
		utils.WriteError(w, http.StatusInternalServerError, "Error deleting webhook")
		return
	}

	utils.Write(w, http.StatusNoContent)
}
