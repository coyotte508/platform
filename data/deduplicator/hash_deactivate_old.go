package deduplicator

import (
	"strconv"

	"github.com/tidepool-org/platform/app"
	"github.com/tidepool-org/platform/data"
	"github.com/tidepool-org/platform/data/store"
	"github.com/tidepool-org/platform/data/types/upload"
	"github.com/tidepool-org/platform/errors"
	"github.com/tidepool-org/platform/log"
)

type hashDeactivateOldFactory struct {
	*BaseFactory
}

type hashDeactivateOldDeduplicator struct {
	*BaseDeduplicator
}

const _HashDeactivateOldDeduplicatorName = "org.tidepool.hash-deactivate-old"
const _HashDeactivateOldDeduplicatorVersion = "1.1.0"

var _HashDeactivateOldExpectedDeviceManufacturerModels = map[string][]string{
	"Medtronic": {"523", "723", "551", "751", "554", "754", "1510", "1511", "1512", "1710", "1711", "1712", "1715", "1780"},
	"LifeScan":  {"OneTouch Ultra 2", "OneTouch UltraMini"},
	"Abbott":    {"FreeStyle Libre"},
}

func NewHashDeactivateOldFactory() (Factory, error) {
	baseFactory, err := NewBaseFactory(_HashDeactivateOldDeduplicatorName, _HashDeactivateOldDeduplicatorVersion)
	if err != nil {
		return nil, err
	}

	factory := &hashDeactivateOldFactory{
		BaseFactory: baseFactory,
	}
	factory.Factory = factory

	return factory, nil
}

func (h *hashDeactivateOldFactory) CanDeduplicateDataset(dataset *upload.Upload) (bool, error) {
	if can, err := h.BaseFactory.CanDeduplicateDataset(dataset); err != nil || !can {
		return can, err
	}

	if dataset.DeviceID == nil {
		return false, nil
	}
	if *dataset.DeviceID == "" {
		return false, nil
	}
	if dataset.DeviceManufacturers == nil {
		return false, nil
	}
	if dataset.DeviceModel == nil {
		return false, nil
	}

	return allowDeviceManufacturerModel(_HashDeactivateOldExpectedDeviceManufacturerModels, *dataset.DeviceManufacturers, *dataset.DeviceModel), nil
}

func (h *hashDeactivateOldFactory) NewDeduplicatorForDataset(logger log.Logger, dataStoreSession store.Session, dataset *upload.Upload) (data.Deduplicator, error) {
	baseDeduplicator, err := NewBaseDeduplicator(h.name, h.version, logger, dataStoreSession, dataset)
	if err != nil {
		return nil, err
	}

	if dataset.DeviceID == nil {
		return nil, errors.New("deduplicator", "dataset device id is missing")
	}
	if *dataset.DeviceID == "" {
		return nil, errors.New("deduplicator", "dataset device id is empty")
	}
	if dataset.DeviceManufacturers == nil {
		return nil, errors.New("deduplicator", "dataset device manufacturers is missing")
	}
	if dataset.DeviceModel == nil {
		return nil, errors.New("deduplicator", "dataset device model is missing")
	}

	if !allowDeviceManufacturerModel(_HashDeactivateOldExpectedDeviceManufacturerModels, *dataset.DeviceManufacturers, *dataset.DeviceModel) {
		return nil, errors.New("deduplicator", "dataset device manufacturer and model does not contain expected device manufacturers and models")
	}

	return &hashDeactivateOldDeduplicator{
		BaseDeduplicator: baseDeduplicator,
	}, nil
}

func (h *hashDeactivateOldDeduplicator) AddDatasetData(datasetData []data.Datum) error {
	hashes, err := AssignDatasetDataIdentityHashes(datasetData)
	if err != nil {
		return err
	} else if len(hashes) == 0 {
		return nil
	}

	return h.BaseDeduplicator.AddDatasetData(datasetData)
}

func (h *hashDeactivateOldDeduplicator) DeduplicateDataset() error {
	if err := h.dataStoreSession.ArchiveDeviceDataUsingHashesFromDataset(h.dataset); err != nil {
		return errors.Wrapf(err, "deduplicator", "unable to archive device data using hashes from dataset with id %s", strconv.Quote(h.dataset.UploadID))
	}

	return h.BaseDeduplicator.DeduplicateDataset()
}

func (h *hashDeactivateOldDeduplicator) DeleteDataset() error {
	if err := h.dataStoreSession.UnarchiveDeviceDataUsingHashesFromDataset(h.dataset); err != nil {
		return errors.Wrapf(err, "deduplicator", "unable to unarchive device data using hashes from dataset with id %s", strconv.Quote(h.dataset.UploadID))
	}

	return h.BaseDeduplicator.DeleteDataset()
}

func allowDeviceManufacturerModel(allowedDeviceManufacturerModels map[string][]string, deviceManufacturers []string, deviceModel string) bool {
	for _, deviceManufacturer := range deviceManufacturers {
		if deviceModels, found := allowedDeviceManufacturerModels[deviceManufacturer]; found {
			if app.StringsContainsString(deviceModels, deviceModel) {
				return true
			}
		}
	}

	return false
}
