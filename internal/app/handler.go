package app

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

type Hub interface {
	Push(model.Location) error
	Register(string, *websocket.Conn)
	DeRegister(string)
}

func (s *server) CheckHealth(ctx *fiber.Ctx) error {
	//if err := s.db.Ping(); err != nil {
	//	s.logger.Error("Error pinging database", "error", err)
	//	return ctx.Status(fiber.StatusInternalServerError).JSON("Error from db")
	//}

	return ctx.Status(fiber.StatusOK).JSON("OK")
}

func (s *server) PushLocation(ctx *fiber.Ctx) error {
	var location model.Location
	err := ctx.BodyParser(&location)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errMsgFromBodyParse)
	}

	err = s.GetLocationDetails(location)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON("Cannot get location details")
	}

	err = s.hub.Push(location)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON("Failed to push location")
	}

	return ctx.Status(fiber.StatusOK).JSON("The location has been pushed")
}

func (s *server) ServeWs(c *websocket.Conn) {
	clientId := uuid.New().String()
	s.hub.Register(clientId, c)
	defer s.hub.DeRegister(clientId)

	var (
		mt  int
		msg []byte
		err error
	)
	for {
		mt, msg, err = c.ReadMessage()
		if err != nil {
			s.logger.Error("Error from reading message", "error", err)
			break
		}
		s.logger.Info("Received WS message", "clientId", clientId, "info", msg)

		err = c.WriteMessage(mt, msg)
		if err != nil {
			s.logger.Error("Error from writing message", "error", err)
			break
		}
	}
}
