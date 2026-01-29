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
			return Fail(c, fiber.StatusForbidden, "max attempts reached")
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

type submitRequest struct {
	AttemptID        string  `json:"attempt_id"`
	QuestionID       string  `json:"question_id"`
	SelectedOptionID *string `json:"selected_option_id,omitempty"`
	EssayAnswer      *string `json:"essay_answer,omitempty"`
	Final            bool    `json:"final"`
}

// SubmitAnswer godoc
// @Summary Submit or autosave quiz answer
// @Description Autosave answer per question or finalize quiz submission.
// @Description Multiple choice answers are auto-graded.
// @Description Essay answers will be marked as waiting assessment.
// @Tags Quiz
// @Accept json
// @Produce json
// @Param request body submitRequest true "Submit answer payload"
// @Success 200 {object} handler.APIResponse "Answer saved or quiz submitted"
// @Failure 400 {object} handler.APIResponse "Invalid request or validation error"
// @Failure 403 {object} handler.APIResponse "Attempt expired"
// @Failure 409 {object} handler.APIResponse "Attempt already submitted"
// @Failure 500 {object} handler.APIResponse "Internal server error"
// @Router /submit [post]
func (h *QuizHandler) Submit(c *fiber.Ctx) error {
	var req submitRequest
	if err := c.BodyParser(&req); err != nil {
		return Fail(c, fiber.StatusBadRequest, "invalid request body")
	}

	attemptID, _ := uuid.Parse(req.AttemptID)
	questionID, _ := uuid.Parse(req.QuestionID)

	var optionID *uuid.UUID
	if req.SelectedOptionID != nil {
		id, err := uuid.Parse(*req.SelectedOptionID)
		if err != nil {
			return Fail(c, fiber.StatusBadRequest, "invalid option id")
		}

		optionID = &id
	}

	err := h.service.SubmitAnswer(
		c.Context(),
		attemptID,
		questionID,
		optionID,
		req.EssayAnswer,
		req.Final,
	)

	if err != nil {
		switch err.Error() {
		case "attempt expired":
			return Fail(c, fiber.StatusForbidden, err.Error())
		case "attempt already submitted":
			return Fail(c, fiber.StatusConflict, err.Error())
		default:
			return Fail(c, fiber.StatusBadRequest, err.Error())
		}
	}

	msg := "answer saved"
	if req.Final {
		msg = "quiz submitted"
	}

	return Success(c, msg, nil)
}
