package users_transport_http

import (
	"net/http"

	core_logger "github.com/Mirwinli/golang-todoapp/internal/core/logger"
	core_http_utils "github.com/Mirwinli/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Mirwinli/golang-todoapp/internal/core/transport/http/response"
)

type GetUserResponse UserDTOResponse

func (h *UsersHTTPHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID, err := core_http_utils.GetIntPathValues(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)
		return
	}
	user, err := h.usersService.GetUser(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get user",
		)
	}

	response := GetUserResponse(userDTOFromDomain(user))

	responseHandler.JSONResponse(response, http.StatusOK)
}
