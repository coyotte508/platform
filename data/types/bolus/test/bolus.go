package test

import (
	"github.com/tidepool-org/platform/data/types/bolus"
	dataTypeIobTest "github.com/tidepool-org/platform/data/types/bolus/iob/test"
	dataTypeCommonTest "github.com/tidepool-org/platform/data/types/common/test"
	dataTypesInsulinTest "github.com/tidepool-org/platform/data/types/insulin/test"
	dataTypesTest "github.com/tidepool-org/platform/data/types/test"
)

func NewBolus() *bolus.Bolus {
	datum := &bolus.Bolus{}
	datum.Base = *dataTypesTest.NewBase()
	datum.Type = "bolus"
	datum.SubType = dataTypesTest.NewType()
	datum.InsulinFormulation = dataTypesInsulinTest.NewFormulation(3)
	datum.InsulinOnBoard = dataTypeIobTest.NewIob()
	datum.Prescriptor = dataTypeCommonTest.NewPrescriptor()
	return datum
}

func CloneBolus(datum *bolus.Bolus) *bolus.Bolus {
	if datum == nil {
		return nil
	}
	clone := &bolus.Bolus{}
	clone.Base = *dataTypesTest.CloneBase(&datum.Base)
	clone.SubType = datum.SubType
	clone.InsulinFormulation = dataTypesInsulinTest.CloneFormulation(datum.InsulinFormulation)
	clone.InsulinOnBoard = dataTypeIobTest.CloneIob(datum.InsulinOnBoard)
	clone.Prescriptor = dataTypeCommonTest.ClonePrescriptor(datum.Prescriptor)
	return clone
}
