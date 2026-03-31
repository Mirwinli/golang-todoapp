package users_transport_http

import (
	"net/http"

	core_logger "github.com/Mirwinli/golang-todoapp/internal/core/logger"
	core_http_response "github.com/Mirwinli/golang-todoapp/internal/core/transport/http/response"
	core_http_utils "github.com/Mirwinli/golang-todoapp/internal/core/transport/http/utils"
)

func (h *UsersHTTPHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID, err := core_http_utils.GetIntPathValues(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID from path parameter",
		)
		return
	}

	if err = h.usersService.DeleteUser(ctx, userID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete user",
		)
		return
	}

	responseHandler.NoContentResponse()
}
