package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/magzhan/geotracker/internal/model"
)

const (
	responseSuccess      = "Success"
	errMsgFromDatabase   = "Error from database"
	errMsgInvalidUUID    = "Invalid UUID input"
	errMsgFromUUIDParse  = "Error while parsing UUID"
	errMsgFromBodyParse  = "Error while parsing body"
	errMsgFromValidation = "Error while validating body"
)

type locationRepository interface {
	GetLocation(key string) (model.Location, error)
	SetLocation(model.Location) error
}

func (s *server) CheckHealth(ctx *fiber.Ctx) error {
	//if err := s.db.Ping(); err != nil {
	//	s.logger.Error("Error pinging database", "error", err)
	//	return ctx.Status(fiber.StatusInternalServerError).JSON("Error from db")
	//}

	return ctx.Status(fiber.StatusOK).JSON("OK")
}

func (s *server) Video(ctx *fiber.Ctx) error {}
