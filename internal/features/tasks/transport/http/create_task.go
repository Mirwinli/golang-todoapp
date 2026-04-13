package tasks_transport_http

import (
	"net/http"

	"github.com/Mirwinli/golang-todoapp/internal/core/domain"
	core_logger "github.com/Mirwinli/golang-todoapp/internal/core/logger"
	core_http_request "github.com/Mirwinli/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Mirwinli/golang-todoapp/internal/core/transport/http/response"
)

type CreateTaskRequest struct {
	Title        string  `json:"title" validate:"required,min=1,max=100"         example:"Hello Wo"`
	Description  *string `json:"description" validate:"omitempty,min=1,max=1000" example:"Сробити домашку до 16 квітня""`
	AuthorUserID int     `json:"author_user_id" validate:"required"              example:"1"`
}

type CreateTaskResponse TaskDTOResponse

// CreateTask 	godoc
// @Summary 	Створити задачу
// @Description Створити нову задачу в системі
// @Tags 		tasks
// @Accept 		json
// @Produce 	json
// @Param 		request body CreateTaskRequest true 		  "CreateTask тіло запиту"
// @Success 	201 {object} CreateTaskResponse 			  "Успішно створена задача"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 	404 {object} core_http_response.ErrorResponse "Author not found"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 		/tasks [post]
func (h *TaskHTTPHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(
		log,
		w,
	)

	var request CreateTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate request",
		)
		return
	}

	taskDomain := domain.NewTaskInitialized(
		request.Title,
		request.Description,

		request.AuthorUserID,
	)

	taskDomain, err := h.taskService.CreateTask(ctx, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create task",
		)

		return
	}
	response := CreateTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)
}
