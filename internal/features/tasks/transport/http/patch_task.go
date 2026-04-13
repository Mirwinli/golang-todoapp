package tasks_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Mirwinli/golang-todoapp/internal/core/domain"
	core_logger "github.com/Mirwinli/golang-todoapp/internal/core/logger"
	core_http_request "github.com/Mirwinli/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Mirwinli/golang-todoapp/internal/core/transport/http/response"
	core_http_types "github.com/Mirwinli/golang-todoapp/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title"`
	Description core_http_types.Nullable[string] `json:"description"`
	Completed   core_http_types.Nullable[bool]   `json:"completed"`
}

type PatchTaskResponse struct {
	ID      int
	Version int

	Title       string
	Description *string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time

	AuthorUserID int
}

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("`Title` can't be NULL")
		}
		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("`Title` length must be between 1 and 100")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLen := len([]rune(*r.Description.Value))
			if descriptionLen < 1 || descriptionLen > 1000 {
				return fmt.Errorf("`Description` length must be between 1 and 100")
			}
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("`Completed` can't be NULL")
		}
	}

	return nil
}

// PatchTask godoc
// @Summary Оновити задачу
// @Description Оновлює інформацію вже в існуючій в системі задачі
// @Description ### Логіка оновлення полів (Three-state-logic)
// @Description 1. **Поле не передано**: `description` ігноруєтся, значення в БД не змінюєтся
// @Description 2. **Явно передано значення**: `description`: "Вийти погуляти в 6:30 з собакою" - змінюєьтся поле в БД
// @Descroption 3. **Явно передано null**: `description`: null - очищає поле в БД (set to NULL)
// @Description Обмеження: `title` і `completed` не можуть бути виставденні null
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "ID оновлюваної задачі"
// @Param request body PatchTaskRequest true "PatchTask тіло запиту"
// @Success 200 {object} PatchTaskResponse "Успішно змінена задача"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 404 {object} core_http_response.ErrorResponse "Task not found"
// @Failure 409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /tasks/{id} [patch]
func (h *TaskHTTPHandler) PatchTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	taskID, err := core_http_request.GetIntPathValues(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get taskID path value",
		)
		return
	}

	var req PatchTaskRequest
	if err = core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate request",
		)
		return
	}

	taskPatch := taskPatchFromRequest(req)

	task, err := h.taskService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch task",
		)
		return
	}

	response := PatchTaskResponse{
		ID:           task.ID,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		CreatedAt:    task.CreatedAt,
		CompletedAt:  task.CompletedAt,
		AuthorUserID: task.AuthorUserID,
	}

	responseHandler.JSONResponse(response, http.StatusOK)
}

func taskPatchFromRequest(req PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		req.Title.ToDomain(),
		req.Description.ToDomain(),
		req.Completed.ToDomain(),
	)
}
