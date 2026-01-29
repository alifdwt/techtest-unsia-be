package handler

import (
	"github.com/alifdwt/techtest-unsia-be/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type QuizHandler struct {
	service *service.QuizService
}

func NewQuizHandler(s *service.QuizService) *QuizHandler {
	return &QuizHandler{
		service: s,
	}
}

type startQuizRequest struct {
	QuizID string `json:"quiz_id"`
	UserID string `json:"user_id"`
}

// StartQuiz godoc
// @Summary Start quiz attempt
// @Description Initialize quiz attempt and fetch questions
// @Tags Quiz
// @Accept json
// @Produce json
// @Param request body startQuizRequest true "Start quiz request"
// @Success 200 {object} handler.APIResponse
// @Failure 400 {object} handler.APIResponse
// @Failure 403 {object} handler.APIResponse
// @Failure 500 {object} handler.APIResponse
// @Router /start [post]
func (h *QuizHandler) StartQuiz(c *fiber.Ctx) error {
	var req startQuizRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "invalid request body")
	}

	quizID, err := uuid.Parse(req.QuizID)
	if err != nil {
		return Fail(c, fiber.StatusBadRequest, "invalid quiz id")
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return Fail(c, fiber.StatusBadRequest, "invalid user id")
	}

	attempt, quiz, deadline, questions, err := h.service.StartQuiz(c.Context(), quizID, userID)
	if err != nil {
		if err.Error() == "max attempts exceeded" {
			return Fail(c, fiber.StatusForbidden, "maximum attempts reached")
		}
		return Fail(c, fiber.StatusInternalServerError, "failed to start quiz")
	}

	return Success(c, "quiz started", map[string]interface{}{
		"attempt_id": attempt.ID,
		"started_at": attempt.StartedAt.Time,
		"deadline":   deadline,
		"duration":   quiz.DurationSeconds,
		"questions":  questions,
	})
}
