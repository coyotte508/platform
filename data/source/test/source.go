package test

import (
	"time"

	"github.com/onsi/gomega"
	gomegaGstruct "github.com/onsi/gomega/gstruct"
	gomegaTypes "github.com/onsi/gomega/types"

	authTest "github.com/tidepool-org/platform/auth/test"
	dataSource "github.com/tidepool-org/platform/data/source"
	dataTest "github.com/tidepool-org/platform/data/test"
	errorsTest "github.com/tidepool-org/platform/errors/test"
	"github.com/tidepool-org/platform/pointer"
	requestTest "github.com/tidepool-org/platform/request/test"
	"github.com/tidepool-org/platform/test"
	userTest "github.com/tidepool-org/platform/user/test"
)

func RandomState() string {
	return test.RandomStringFromArray(dataSource.States())
}

func RandomStates() []string {
	return test.RandomStringArrayFromRangeAndArrayWithoutDuplicates(1, len(dataSource.States()), dataSource.States())
}

func RandomFilter() *dataSource.Filter {
	datum := &dataSource.Filter{}
	datum.ProviderType = pointer.FromStringArray(authTest.RandomProviderTypes())
	datum.ProviderName = pointer.FromStringArray(authTest.RandomProviderNames())
	datum.ProviderSessionID = pointer.FromStringArray(authTest.RandomProviderSessionIDs())
	datum.State = pointer.FromStringArray(RandomStates())
	return datum
}

func NewObjectFromFilter(datum *dataSource.Filter, objectFormat test.ObjectFormat) map[string]interface{} {
	if datum == nil {
		return nil
	}
	object := map[string]interface{}{}
	if datum.ProviderType != nil {
		object["providerType"] = test.NewObjectFromStringArray(*datum.ProviderType, objectFormat)
	}
	if datum.ProviderName != nil {
		object["providerName"] = test.NewObjectFromStringArray(*datum.ProviderName, objectFormat)
	}
	if datum.ProviderSessionID != nil {
		object["providerSessionId"] = test.NewObjectFromStringArray(*datum.ProviderSessionID, objectFormat)
	}
	if datum.State != nil {
		object["state"] = test.NewObjectFromStringArray(*datum.State, objectFormat)
	}
	return object
}

func RandomCreate() *dataSource.Create {
	datum := &dataSource.Create{}
	datum.ProviderType = pointer.FromString(authTest.RandomProviderType())
	datum.ProviderName = pointer.FromString(authTest.RandomProviderName())
	datum.ProviderSessionID = pointer.FromString(authTest.RandomProviderSessionID())
	datum.State = pointer.FromString(RandomState())
	return datum
}

func NewObjectFromCreate(datum *dataSource.Create, objectFormat test.ObjectFormat) map[string]interface{} {
	if datum == nil {
		return nil
	}
	object := map[string]interface{}{}
	if datum.ProviderType != nil {
		object["providerType"] = test.NewObjectFromString(*datum.ProviderType, objectFormat)
	}
	if datum.ProviderName != nil {
		object["providerName"] = test.NewObjectFromString(*datum.ProviderName, objectFormat)
	}
	if datum.ProviderSessionID != nil {
		object["providerSessionId"] = test.NewObjectFromString(*datum.ProviderSessionID, objectFormat)
	}
	if datum.State != nil {
		object["state"] = test.NewObjectFromString(*datum.State, objectFormat)
	}
	return object
}

func RandomUpdate() *dataSource.Update {
	state := RandomState()
	datum := &dataSource.Update{}
	if state != dataSource.StateDisconnected {
		datum.ProviderSessionID = pointer.FromString(authTest.RandomProviderSessionID())
	}
	datum.State = pointer.FromString(state)
	datum.Error = errorsTest.RandomSerializable()
	datum.DataSetIDs = pointer.FromStringArray(dataTest.RandomSetIDs())
	datum.EarliestDataTime = pointer.FromTime(test.RandomTimeFromRange(test.RandomTimeMinimum(), time.Now()))
	datum.LatestDataTime = pointer.FromTime(test.RandomTimeFromRange(*datum.EarliestDataTime, time.Now()))
	datum.LastImportTime = pointer.FromTime(test.RandomTimeFromRange(*datum.LatestDataTime, time.Now()))
	return datum
}

func CloneUpdate(datum *dataSource.Update) *dataSource.Update {
	if datum == nil {
		return nil
	}
	clone := &dataSource.Update{}
	clone.ProviderSessionID = pointer.CloneString(datum.ProviderSessionID)
	clone.State = pointer.CloneString(datum.State)
	clone.Error = errorsTest.CloneSerializable(datum.Error)
	clone.DataSetIDs = pointer.CloneStringArray(datum.DataSetIDs)
	clone.EarliestDataTime = pointer.CloneTime(datum.EarliestDataTime)
	clone.LatestDataTime = pointer.CloneTime(datum.LatestDataTime)
	clone.LastImportTime = pointer.CloneTime(datum.LastImportTime)
	return clone
}

func NewObjectFromUpdate(datum *dataSource.Update, objectFormat test.ObjectFormat) map[string]interface{} {
	if datum == nil {
		return nil
	}
	object := map[string]interface{}{}
	if datum.ProviderSessionID != nil {
		object["providerSessionId"] = test.NewObjectFromString(*datum.ProviderSessionID, objectFormat)
	}
	if datum.State != nil {
		object["state"] = test.NewObjectFromString(*datum.State, objectFormat)
	}
	if datum.Error != nil {
		object["error"] = errorsTest.NewObjectFromSerializable(datum.Error, objectFormat)
	}
	if datum.DataSetIDs != nil {
		object["dataSetIds"] = test.NewObjectFromStringArray(*datum.DataSetIDs, objectFormat)
	}
	if datum.EarliestDataTime != nil {
		object["earliestDataTime"] = test.NewObjectFromTime(*datum.EarliestDataTime, objectFormat)
	}
	if datum.LatestDataTime != nil {
		object["latestDataTime"] = test.NewObjectFromTime(*datum.LatestDataTime, objectFormat)
	}
	if datum.LastImportTime != nil {
		object["lastImportTime"] = test.NewObjectFromTime(*datum.LastImportTime, objectFormat)
	}
	return object
}

func MatchUpdate(datum *dataSource.Update) gomegaTypes.GomegaMatcher {
	if datum == nil {
		return gomega.BeNil()
	}
	return gomegaGstruct.PointTo(gomegaGstruct.MatchAllFields(gomegaGstruct.Fields{
		"ProviderSessionID": gomega.Equal(datum.ProviderSessionID),
		"State":             gomega.Equal(datum.State),
		"Error":             gomega.Equal(datum.Error),
		"DataSetIDs":        gomega.Equal(datum.DataSetIDs),
		"EarliestDataTime":  test.MatchTime(datum.EarliestDataTime),
		"LatestDataTime":    test.MatchTime(datum.LatestDataTime),
		"LastImportTime":    test.MatchTime(datum.LastImportTime),
	}))
}

func RandomSource() *dataSource.Source {
	datum := &dataSource.Source{}
	datum.ID = pointer.FromString(RandomID())
	datum.UserID = pointer.FromString(userTest.RandomID())
	datum.ProviderType = pointer.FromString(authTest.RandomProviderType())
	datum.ProviderName = pointer.FromString(authTest.RandomProviderName())
	datum.ProviderSessionID = pointer.FromString(authTest.RandomProviderSessionID())
	datum.State = pointer.FromString(RandomState())
	datum.Error = errorsTest.RandomSerializable()
	datum.DataSetIDs = pointer.FromStringArray(dataTest.RandomSetIDs())
	datum.EarliestDataTime = pointer.FromTime(test.RandomTimeFromRange(test.RandomTimeMinimum(), time.Now()))
	datum.LatestDataTime = pointer.FromTime(test.RandomTimeFromRange(*datum.EarliestDataTime, time.Now()))
	datum.LastImportTime = pointer.FromTime(test.RandomTimeFromRange(*datum.LatestDataTime, time.Now()))
	datum.CreatedTime = pointer.FromTime(test.RandomTimeFromRange(test.RandomTimeMinimum(), time.Now()))
	datum.ModifiedTime = pointer.FromTime(test.RandomTimeFromRange(*datum.CreatedTime, time.Now()))
	datum.Revision = pointer.FromInt(requestTest.RandomRevision())
	return datum
}

func CloneSource(datum *dataSource.Source) *dataSource.Source {
	if datum == nil {
		return nil
	}
	clone := &dataSource.Source{}
	clone.ID = pointer.CloneString(datum.ID)
	clone.UserID = pointer.CloneString(datum.UserID)
	clone.ProviderType = pointer.CloneString(datum.ProviderType)
	clone.ProviderName = pointer.CloneString(datum.ProviderName)
	clone.ProviderSessionID = pointer.CloneString(datum.ProviderSessionID)
	clone.State = pointer.CloneString(datum.State)
	clone.Error = errorsTest.CloneSerializable(datum.Error)
	clone.DataSetIDs = pointer.CloneStringArray(datum.DataSetIDs)
	clone.EarliestDataTime = pointer.CloneTime(datum.EarliestDataTime)
	clone.LatestDataTime = pointer.CloneTime(datum.LatestDataTime)
	clone.LastImportTime = pointer.CloneTime(datum.LastImportTime)
	clone.CreatedTime = pointer.CloneTime(datum.CreatedTime)
	clone.ModifiedTime = pointer.CloneTime(datum.ModifiedTime)
	clone.Revision = pointer.CloneInt(datum.Revision)
	return clone
}

func NewObjectFromSource(datum *dataSource.Source, objectFormat test.ObjectFormat) map[string]interface{} {
	if datum == nil {
		return nil
	}
	object := map[string]interface{}{}
	if datum.ID != nil {
		object["id"] = test.NewObjectFromString(*datum.ID, objectFormat)
	}
	if datum.UserID != nil {
		object["userId"] = test.NewObjectFromString(*datum.UserID, objectFormat)
	}
	if datum.ProviderType != nil {
		object["providerType"] = test.NewObjectFromString(*datum.ProviderType, objectFormat)
	}
	if datum.ProviderName != nil {
		object["providerName"] = test.NewObjectFromString(*datum.ProviderName, objectFormat)
	}
	if datum.ProviderSessionID != nil {
		object["providerSessionId"] = test.NewObjectFromString(*datum.ProviderSessionID, objectFormat)
	}
	if datum.State != nil {
		object["state"] = test.NewObjectFromString(*datum.State, objectFormat)
	}
	if datum.Error != nil {
		object["error"] = errorsTest.NewObjectFromSerializable(datum.Error, objectFormat)
	}
	if datum.DataSetIDs != nil {
		object["dataSetIds"] = test.NewObjectFromStringArray(*datum.DataSetIDs, objectFormat)
	}
	if datum.EarliestDataTime != nil {
		object["earliestDataTime"] = test.NewObjectFromTime(*datum.EarliestDataTime, objectFormat)
	}
	if datum.LatestDataTime != nil {
		object["latestDataTime"] = test.NewObjectFromTime(*datum.LatestDataTime, objectFormat)
	}
	if datum.LastImportTime != nil {
		object["lastImportTime"] = test.NewObjectFromTime(*datum.LastImportTime, objectFormat)
	}
	if datum.CreatedTime != nil {
		object["createdTime"] = test.NewObjectFromTime(*datum.CreatedTime, objectFormat)
	}
	if datum.ModifiedTime != nil {
		object["modifiedTime"] = test.NewObjectFromTime(*datum.ModifiedTime, objectFormat)
	}
	if datum.Revision != nil {
		object["revision"] = test.NewObjectFromInt(*datum.Revision, objectFormat)
	}
	return object
}

func MatchSource(datum *dataSource.Source) gomegaTypes.GomegaMatcher {
	if datum == nil {
		return gomega.BeNil()
	}
	return gomegaGstruct.PointTo(gomegaGstruct.MatchAllFields(gomegaGstruct.Fields{
		"ID":                gomega.Equal(datum.ID),
		"UserID":            gomega.Equal(datum.UserID),
		"ProviderType":      gomega.Equal(datum.ProviderType),
		"ProviderName":      gomega.Equal(datum.ProviderName),
		"ProviderSessionID": gomega.Equal(datum.ProviderSessionID),
		"State":             gomega.Equal(datum.State),
		"Error":             gomega.Equal(datum.Error),
		"DataSetIDs":        gomega.Equal(datum.DataSetIDs),
		"EarliestDataTime":  test.MatchTime(datum.EarliestDataTime),
		"LatestDataTime":    test.MatchTime(datum.LatestDataTime),
		"LastImportTime":    test.MatchTime(datum.LastImportTime),
		"CreatedTime":       test.MatchTime(datum.CreatedTime),
		"ModifiedTime":      test.MatchTime(datum.ModifiedTime),
		"Revision":          gomega.Equal(datum.Revision),
	}))
}

func RandomSources(minimumLength int, maximumLength int) dataSource.Sources {
	datum := make(dataSource.Sources, test.RandomIntFromRange(minimumLength, maximumLength))
	for index := range datum {
		datum[index] = RandomSource()
	}
	return datum
}

func CloneSources(datum dataSource.Sources) dataSource.Sources {
	if len(datum) == 0 {
		return datum
	}
	clone := dataSource.Sources{}
	for _, source := range datum {
		clone = append(clone, CloneSource(source))
	}
	return clone
}

func MatchSources(datum dataSource.Sources) gomegaTypes.GomegaMatcher {
	matchers := []gomegaTypes.GomegaMatcher{}
	for _, d := range datum {
		matchers = append(matchers, MatchSource(d))
	}
	return test.MatchArray(matchers)
}

func RandomID() string {
	return dataSource.NewID()
}
