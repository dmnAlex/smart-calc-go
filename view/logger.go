package view

import (
	"log"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

func InitLogger() {
	rotation := "hourly"
	path := "logs/logs_%d-%m-%Y-%H-%M-%S"

	var period time.Duration

	switch rotation {
	case "dayly":
		period = time.Hour * 24
	case "monthly":
		period = time.Hour * 24 * 30
	default:
		period = time.Hour
	}

	logger, err := rotatelogs.New(path, rotatelogs.WithRotationTime(period))

	if err != nil {
		log.Printf("failed to create rotatelogs")
		return
	}

	log.SetFlags(log.Ldate | log.Ltime)
	log.SetOutput(logger)
}
