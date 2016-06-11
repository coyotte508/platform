package v1

/* CHECKLIST
 * [ ] Uses interfaces as appropriate
 * [ ] Private package variables use underscore prefix
 * [ ] All parameters validated
 * [ ] All errors handled
 * [ ] Reviewed for concurrency safety
 * [ ] Code complete
 * [ ] Full test coverage
 */

import (
	"net/http"

	"github.com/tidepool-org/platform/data/deduplicator/root"
	"github.com/tidepool-org/platform/data/store"
	"github.com/tidepool-org/platform/data/types/base/upload"
	"github.com/tidepool-org/platform/dataservices/server"
	"github.com/tidepool-org/platform/userservices/client"
)

func DatasetsUpdate(context server.Context) {
	datasetID := context.Request().PathParam("datasetid")
	if datasetID == "" {
		context.RespondWithError(ErrorDatasetIDMissing())
		return
	}

	// TODO: Improve context.DataStoreSession() Find - more specific
	var datasetUpload upload.Upload
	if err := context.DataStoreSession().Find(store.Query{"type": "upload", "uploadId": datasetID}, &datasetUpload); err != nil {
		context.RespondWithError(ErrorDatasetIDNotFound(datasetID))
		return
	}

	// TODO: Validate
	targetUserID := datasetUpload.UserID

	err := context.UserServicesClient().ValidateTargetUserPermissions(context, context.RequestUserID(), targetUserID, client.UploadPermissions)
	if err != nil {
		if client.IsUnauthorizedError(err) {
			context.RespondWithError(ErrorUnauthorized())
		} else {
			context.RespondWithInternalServerFailure("Unable to validate target user permissions", err)
		}
		return
	}

	if datasetUpload.DataState != "open" {
		context.RespondWithError(ErrorDatasetClosed(datasetID))
		return
	}

	datasetUpload.SetDataState("closed")

	if err = context.DataStoreSession().Update(map[string]interface{}{"type": "upload", "uploadId": datasetID}, datasetUpload); err != nil {
		context.RespondWithInternalServerFailure("Unable to insert dataset", err)
		return
	}

	deduplicator, err := root.NewFactory().NewDeduplicator(context.Logger(), context.DataStoreSession(), &datasetUpload)
	if err != nil {
		context.RespondWithInternalServerFailure("No duplicator found matching dataset", err)
		return
	}

	if err = deduplicator.FinalizeDataset(); err != nil {
		context.RespondWithInternalServerFailure("Unable to finalize dataset", err)
		return
	}

	// TODO: Filter datasetUpload to only "public" fields
	context.Response().WriteHeader(http.StatusOK)
	context.Response().WriteJson(datasetUpload)
}
