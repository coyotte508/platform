package v1

import (
	"net/http"

	dataService "github.com/tidepool-org/platform/data/service"
	dataStoreDEPRECATED "github.com/tidepool-org/platform/data/storeDEPRECATED"
	"github.com/tidepool-org/platform/data/types/upload"
	"github.com/tidepool-org/platform/page"
	"github.com/tidepool-org/platform/permission"
	"github.com/tidepool-org/platform/request"
	"github.com/tidepool-org/platform/service"
)

// Array of uploads
type uploads []upload.Upload

// UsersDataSetsGet godoc
// @Summary Get data sets
// @Description Get data sets
// @Description Caller must be a service, the owner, or have the authorizations to do it in behalf of the user.
// @ID platform-data-api-UsersDataSetsGet
// @Produce json
// @Param userId path string true "user ID"
// @Param page query int false "When using pagination, page number, default = 0"
// @Param size query int false "When using pagination, number of elements by page, default = 100 (min 1, max 1000)"
// @Param deleted query bool false "True to return the deleted datasets"
// @Param X-Tidepool-Service-Secret header string false "The platform-data service secret"
// @Param X-Tidepool-Session-Token header string false "A tidepool session token"
// @Param restricted_token header string false "A tidepool restricted token"
// @Param Authorization header string false "A tidepool authorization token"
// @Success 200 {object} uploads "Operation is a success"
// @Failure 400 {object} service.Error "User id is missing or JSON body is malformed"
// @Failure 403 {object} service.Error "Forbiden: caller is not authorized"
// @Failure 500 {object} service.Error "Unable to perform the operation"
// @Router /v1/users/:userId/datasets [get]
func UsersDataSetsGet(dataServiceContext dataService.Context) {
	ctx := dataServiceContext.Request().Context()

	targetUserID := dataServiceContext.Request().PathParam("userId")
	if targetUserID == "" {
		dataServiceContext.RespondWithError(ErrorUserIDMissing())
		return
	}

	if details := request.DetailsFromContext(ctx); !details.IsService() {
		permissions, err := dataServiceContext.PermissionClient().GetUserPermissions(ctx, details.UserID(), targetUserID)
		if err != nil {
			if request.IsErrorUnauthorized(err) {
				dataServiceContext.RespondWithError(service.ErrorUnauthorized())
			} else {
				dataServiceContext.RespondWithInternalServerFailure("Unable to get user permissions", err)
			}
			return
		}
		if _, ok := permissions[permission.Read]; !ok {
			dataServiceContext.RespondWithError(service.ErrorUnauthorized())
			return
		}
	}

	filter := dataStoreDEPRECATED.NewFilter()
	pagination := page.NewPagination()
	if err := request.DecodeRequestQuery(dataServiceContext.Request().Request, filter, pagination); err != nil {
		request.MustNewResponder(dataServiceContext.Response(), dataServiceContext.Request()).Error(http.StatusBadRequest, err)
		return
	}

	dataSets, err := dataServiceContext.DataSession().GetDataSetsForUserByID(ctx, targetUserID, filter, pagination)
	if err != nil {
		dataServiceContext.RespondWithInternalServerFailure("Unable to get data sets for user", err)
		return
	}

	dataServiceContext.RespondWithStatusAndData(http.StatusOK, dataSets)
}
