package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/Mirwinli/golang-todoapp/internal/core/logger"
	core_http_utils "github.com/Mirwinli/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Mirwinli/golang-todoapp/internal/core/transport/http/response"
)

type GetUsersResponse []UserDTOResponse

// GetUsers 	godoc
// @Summary 	Список користувачів
// @Description Перегляд списка користувачів з опиціональною пагінацією
// @Tags 		users
// @Produce 	json
// @Param 		limit query int false 						  "Розмір сторінки з користувачами"
// @Param 		offset query int false 						  "Зміщення сторінки з користувачами"
// Success 		200 {object} GetUsersResponse 				  "Успішне получення списка користувачів"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 		/users [get]
func (h *UsersHTTPHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	limit, offset, err := getLimitOffsetQueryParam(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'limit/'offset query param",
		)
		return
	}

	userDomains, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get users",
		)
		return
	}
	response := GetUsersResponse(userDTOFromDomains(userDomains))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func getLimitOffsetQueryParam(r *http.Request) (*int, *int, error) {
	const (
		limitQueryParam  = "limit"
		offsetQueryParam = "offset"
	)

	limit, err := core_http_utils.GetIntQueryParam(r, limitQueryParam)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"get 'limit' query param: %w", err,
		)
	}
	offset, err := core_http_utils.GetIntQueryParam(r, offsetQueryParam)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"get 'offset' query param: %w", err,
		)
	}
	return limit, offset, nil
}
