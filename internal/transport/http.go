package transport

import (
	"github.com/alifdwt/techtest-unsia-be/internal/db"
	"github.com/alifdwt/techtest-unsia-be/internal/handler"
	"github.com/alifdwt/techtest-unsia-be/internal/service"
	"github.com/gofiber/fiber/v2"

	_ "github.com/alifdwt/techtest-unsia-be/docs"
	"github.com/gofiber/swagger"
)

func RegisterRoutes(app *fiber.App, q *db.Queries) {
	quizService := service.NewQuizService(q)
	quizHandler := handler.NewQuizHandler(quizService)

	app.Post("/start", quizHandler.StartQuiz)
	app.Post("/submit", quizHandler.Submit)
	app.Get("/result", quizHandler.GetResult)

	app.Get("/swagger/*", swagger.HandlerDefault)
}
