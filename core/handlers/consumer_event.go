package handlers

import (
	"encoding/json"
	"errors"
	"structure-golang/common/logs"
	"structure-golang/core/models"
	"structure-golang/core/services"
)

type consumerHand struct {
	log     logs.AppLog
	userSrv services.UserService
}

func NewConsumerHandler(log logs.AppLog, jobSrv services.UserService) consumerHand {
	return consumerHand{log, jobSrv}
}

func (h *consumerHand) UpdateData(topic string, data []byte) error {
	if topic != "example_topic" {
		h.log.Error(models.ErrUnexpected)
		return errors.New(models.ErrUnexpected)
	}

	body := models.HandUpdateDataModel{}
	err := json.Unmarshal(data, &body)
	if err != nil {
		h.log.Error(models.ErrUnexpected)
		return errors.New(models.ErrUnexpected)
	}

	// Update ticket type service

	return nil
}
