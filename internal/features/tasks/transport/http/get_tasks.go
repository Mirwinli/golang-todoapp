package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/Mirwinli/golang-todoapp/internal/core/logger"
	"github.com/Mirwinli/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Mirwinli/golang-todoapp/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDTOResponse

// GetTasks godoc
// @Summary Cписок задач
// @Description Перегляд списку задач з опціональною пагінацією та/або фільрацією по ID автора задачі
// @Tags tasks
// @Produce json
// @Param user_id query int false "Фільтрація задач по ID автора"
// @Param limit query int false "Розмір сторінки з задачами"
// @Param offset query int false "Зміщення сторінки з задачами"
// @Success 200 {object} GetTasksResponse "Список задач"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /tasks [get]
func (h *TaskHTTPHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID, limit, offset, err := getUserIDLimitOffsetQueryParam(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID/limit/offset",
		)
		return
	}

	tasks, err := h.taskService.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get tasks",
		)
		return
	}

	response := GetTasksResponse(taskDTOsFromDomains(tasks))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func getUserIDLimitOffsetQueryParam(r *http.Request) (*int, *int, *int, error) {
	const (
		userIDQueryParam = "user_id"
		limitQueryParam  = "limit"
		offsetQueryParam = "offset"
	)

	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParam)
	if err != nil {
		return nil, nil, nil, fmt.Errorf(
			"get 'user_id' query para,: %w",
			err,
		)
	}

	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParam)
	if err != nil {
		return nil, nil, nil, fmt.Errorf(
			"get 'limit' query param: %w", err,
		)
	}
	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParam)
	if err != nil {
		return nil, nil, nil, fmt.Errorf(
			"get 'offset' query param: %w", err,
		)
	}
	return userID, limit, offset, nil
}
