package transport

import (
	"github.com/alifdwt/techtest-unsia-be/internal/db"
	"github.com/alifdwt/techtest-unsia-be/internal/handler"
	"github.com/alifdwt/techtest-unsia-be/internal/service"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, q *db.Queries) {
	quizService := service.NewQuizService(q)
	quizHandler := handler.NewQuizHandler(quizService)

	app.Post("/start", quizHandler.StartQuiz)
}
