package selfmonitored_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/tidepool-org/platform/data/context"
	"github.com/tidepool-org/platform/data/normalizer"
	"github.com/tidepool-org/platform/data/types/base"
	"github.com/tidepool-org/platform/data/types/base/selfmonitored"
	"github.com/tidepool-org/platform/data/types/base/testing"
	"github.com/tidepool-org/platform/data/types/common/bloodglucose"
	"github.com/tidepool-org/platform/data/validator"
	"github.com/tidepool-org/platform/service"
)

func NewRawObject() map[string]interface{} {
	rawObject := testing.RawBaseObject()
	rawObject["type"] = "smbg"
	rawObject["units"] = "mmol/L"
	rawObject["subType"] = "manual"
	rawObject["value"] = 5
	return rawObject
}

func NewMeta() interface{} {
	return &base.Meta{
		Type: "smbg",
	}
}

var _ = Describe("SelfMonitored", func() {
	Context("units", func() {
		DescribeTable("invalid when", testing.ExpectFieldNotValid,
			Entry("is empty", NewRawObject(), "units", "",
				[]*service.Error{testing.ComposeError(validator.ErrorStringNotOneOf("", []string{bloodglucose.Mmoll, bloodglucose.MmolL, bloodglucose.Mgdl, bloodglucose.MgdL}), "/units", NewMeta())},
			),
			Entry("is not one of the predefined values", NewRawObject(), "units", "wrong",
				[]*service.Error{testing.ComposeError(validator.ErrorStringNotOneOf("wrong", []string{bloodglucose.Mmoll, bloodglucose.MmolL, bloodglucose.Mgdl, bloodglucose.MgdL}), "/units", NewMeta())},
			),
		)

		DescribeTable("valid when", testing.ExpectFieldIsValid,
			Entry("is mmol/l", NewRawObject(), "units", "mmol/l"),
			Entry("is mmol/L", NewRawObject(), "units", "mmol/L"),
			Entry("is mg/dl", NewRawObject(), "units", "mg/dl"),
			Entry("is mg/dL", NewRawObject(), "units", "mg/dL"),
		)
	})

	Context("subType", func() {
		DescribeTable("invalid when", testing.ExpectFieldNotValid,
			Entry("is not one of the predefined values", NewRawObject(), "subType", "wrong",
				[]*service.Error{testing.ComposeError(validator.ErrorStringNotOneOf("wrong", []string{"manual", "linked"}), "/subType", NewMeta())},
			),
		)

		DescribeTable("valid when", testing.ExpectFieldIsValid,
			Entry("is manual", NewRawObject(), "subType", "manual"),
			Entry("is linked", NewRawObject(), "subType", "linked"),
		)
	})

	Context("value", func() {
		DescribeTable("value when", testing.ExpectFieldNotValid,
			Entry("is less than 0", NewRawObject(), "value", -0.1,
				[]*service.Error{testing.ComposeError(validator.ErrorFloatNotInRange(-0.1, bloodglucose.MgdLFromValue, bloodglucose.MgdLToValue), "/value", NewMeta())},
			),
			Entry("is greater than 1000", NewRawObject(), "value", 1000.1,
				[]*service.Error{testing.ComposeError(validator.ErrorFloatNotInRange(1000.1, bloodglucose.MgdLFromValue, bloodglucose.MgdLToValue), "/value", NewMeta())},
			),
		)

		DescribeTable("valid when", testing.ExpectFieldIsValid,
			Entry("is above 0", NewRawObject(), "value", 0.1),
			Entry("is below 1000", NewRawObject(), "value", 990.85745),
			Entry("is an integer", NewRawObject(), "value", 12),
		)
	})

	Context("normalized when mmol/L", func() {
		DescribeTable("normalization", func(val, expected float64) {
			selfMonitoredBg, err := selfmonitored.New()
			Expect(err).To(BeNil())
			selfMonitoredBg.Value = &val
			selfMonitoredBg.Units = &bloodglucose.Mmoll

			testContext := context.NewStandard()
			standardNormalizer, err := normalizer.NewStandard(testContext)
			Expect(err).To(BeNil())
			selfMonitoredBg.Normalize(standardNormalizer)
			Expect(selfMonitoredBg.Units).To(Equal(&bloodglucose.MmolL))
			Expect(selfMonitoredBg.Value).To(Equal(&expected))
		},
			Entry("is expected lower bg value", 3.7, 3.7),
			Entry("is below max", 54.99, 54.99),
			Entry("is expected upper bg value", 23.0, 23.0),
		)
	})

	Context("normalized when mg/dL", func() {
		DescribeTable("normalization", func(val, expected float64) {
			selfMonitoredBg, err := selfmonitored.New()
			selfMonitoredBg.Value = &val
			selfMonitoredBg.Units = &bloodglucose.Mgdl

			testContext := context.NewStandard()
			standardNormalizer, err := normalizer.NewStandard(testContext)
			Expect(err).To(BeNil())
			selfMonitoredBg.Normalize(standardNormalizer)
			Expect(selfMonitoredBg.Units).To(Equal(&bloodglucose.MmolL))
			Expect(selfMonitoredBg.Value).To(Equal(&expected))
		},
			Entry("is expected lower bg value", 60.0, 3.33044879462732),
			Entry("is below max", 990.85745, 55.0),
			Entry("is expected upper bg value", 400.0, 22.202991964182132),
		)
	})
})
