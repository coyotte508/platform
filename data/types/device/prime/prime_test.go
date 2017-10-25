package prime_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"

	testData "github.com/tidepool-org/platform/data/test"
	"github.com/tidepool-org/platform/data/types/device"
	"github.com/tidepool-org/platform/service"
)

func NewRawObject() map[string]interface{} {
	rawObject := testData.RawBaseObject()
	rawObject["type"] = "deviceEvent"
	rawObject["subType"] = "prime"
	rawObject["volume"] = 0.0
	return rawObject
}

func NewRawObjectWithCannula() map[string]interface{} {
	rawObject := NewRawObject()
	rawObject["primeTarget"] = "cannula"
	return rawObject
}

func NewRawObjectWithTubing() map[string]interface{} {
	rawObject := NewRawObject()
	rawObject["primeTarget"] = "tubing"
	return rawObject
}

func NewMeta() interface{} {
	return &device.Meta{
		Type:    "deviceEvent",
		SubType: "prime",
	}
}

var _ = Describe("Prime", func() {
	Context("primeTarget", func() {
		DescribeTable("invalid when", testData.ExpectFieldNotValid,
			Entry("is empty", NewRawObject(), "primeTarget", "",
				[]*service.Error{testData.ComposeError(service.ErrorValueStringNotOneOf("", []string{"cannula", "tubing"}), "/primeTarget", NewMeta())},
			),
			Entry("is not one of the predefined types", NewRawObject(), "primeTarget", "bad",
				[]*service.Error{testData.ComposeError(service.ErrorValueStringNotOneOf("bad", []string{"cannula", "tubing"}), "/primeTarget", NewMeta())},
			),
		)

		DescribeTable("valid when", testData.ExpectFieldIsValid,
			Entry("is cannula type", NewRawObject(), "primeTarget", "cannula"),
			Entry("is tubing type", NewRawObject(), "primeTarget", "tubing"),
		)
	})

	Context("cannula volume", func() {
		DescribeTable("invalid when", testData.ExpectFieldNotValid,
			Entry("is less than 0", NewRawObjectWithCannula(), "volume", -0.1,
				[]*service.Error{testData.ComposeError(service.ErrorValueNotInRange(-0.1, 0.0, 10.0), "/volume", NewMeta())},
			),
			Entry("is more than 10", NewRawObjectWithCannula(), "volume", 10.1,
				[]*service.Error{testData.ComposeError(service.ErrorValueNotInRange(10.1, 0.0, 10.0), "/volume", NewMeta())},
			),
		)

		DescribeTable("valid when", testData.ExpectFieldIsValid,
			Entry("is 0", NewRawObjectWithCannula(), "volume", 0.0),
			Entry("is 10.0", NewRawObjectWithCannula(), "volume", 10.0),
			Entry("is no decimal", NewRawObjectWithCannula(), "volume", 2),
		)
	})

	Context("tubing volume", func() {
		DescribeTable("invalid when", testData.ExpectFieldNotValid,
			Entry("is less than 0", NewRawObjectWithTubing(), "volume", -0.1,
				[]*service.Error{testData.ComposeError(service.ErrorValueNotInRange(-0.1, 0.0, 100.0), "/volume", NewMeta())},
			),
			Entry("is more than 100", NewRawObjectWithTubing(), "volume", 100.1,
				[]*service.Error{testData.ComposeError(service.ErrorValueNotInRange(100.1, 0.0, 100.0), "/volume", NewMeta())},
			),
		)

		DescribeTable("valid when", testData.ExpectFieldIsValid,
			Entry("is 0", NewRawObjectWithTubing(), "volume", 0.0),
			Entry("is 100.0", NewRawObjectWithTubing(), "volume", 100.0),
			Entry("is no decimal", NewRawObjectWithTubing(), "volume", 55),
		)
	})
})
