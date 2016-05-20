package pump_test

import (
	"github.com/tidepool-org/platform/pvn/data/context"
	"github.com/tidepool-org/platform/pvn/data/normalizer"
	"github.com/tidepool-org/platform/pvn/data/types/base/pump"
	"github.com/tidepool-org/platform/pvn/data/types/base/testing"
	"github.com/tidepool-org/platform/pvn/data/types/common/bloodglucose"
	"github.com/tidepool-org/platform/pvn/data/validator"
	"github.com/tidepool-org/platform/service"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pump Settings", func() {

	var rawObject = testing.RawBaseObject()

	BeforeEach(func() {

		rawObject["type"] = "pumpSettings"
		rawObject["activeSchedule"] = "standard"

		rawObject["units"] = map[string]interface{}{
			"carb": "grams",
			"bg":   "mmol/L",
		}

		rawObject["carbRatio"] = []interface{}{
			map[string]interface{}{"amount": 12, "start": 0},
			map[string]interface{}{"amount": 10, "start": 21600000},
		}

		rawObject["bgTarget"] = []interface{}{
			map[string]interface{}{"low": 5.5, "high": 6.7, "start": 0},
			map[string]interface{}{"low": 5.0, "high": 6.1, "start": 18000000},
		}

		rawObject["insulinSensitivity"] = []interface{}{
			map[string]interface{}{"amount": 3.6, "start": 0},
			map[string]interface{}{"amount": 2.5, "start": 18000000},
		}

		rawObject["basalSchedules"] = []interface{}{
			map[string][]interface{}{
				"standard": {
					map[string]interface{}{"rate": 0.8, "start": 0},
					map[string]interface{}{"rate": 0.75, "start": 3600000},
				}},
		}

	})

	Context("activeSchedule", func() {

		DescribeTable("invalid when", testing.ExpectFieldNotValid,
			Entry("empty", rawObject, "activeSchedule", "", []*service.Error{validator.ErrorValueNotTrue()}),
		)

		DescribeTable("valid when", testing.ExpectFieldIsValid,
			Entry("more than 1 characters", rawObject, "activeSchedule", "A"),
			Entry("freetext", rawObject, "activeSchedule", "standard"),
		)

	})

	Context("units", func() {

		DescribeTable("invalid when", testing.ExpectFieldNotValid,
			Entry("bg empty", rawObject, "units", map[string]interface{}{"carb": "grams", "bg": ""}, []*service.Error{validator.ErrorValueNotTrue()}),
			Entry("bg not predefined type", rawObject, "units", map[string]interface{}{"carb": "grams", "bg": "na"}, []*service.Error{validator.ErrorValueNotTrue()}),
			Entry("carb empty", rawObject, "units", map[string]interface{}{"carb": "", "bg": "mmol/L"}, []*service.Error{validator.ErrorValueNotTrue()}),
		)

		DescribeTable("valid when", testing.ExpectFieldIsValid,
			Entry("carbs set and bg set as mmol/L", rawObject, "units", map[string]interface{}{"carb": "grams", "bg": "mmol/L"}),
			Entry("carbs set and bg set as mg/dl", rawObject, "units", map[string]interface{}{"carb": "grams", "bg": "mg/dl"}),
		)

	})

	Context("carbRatio", func() {

		DescribeTable("invalid when", testing.ExpectFieldNotValid,
			Entry("start negative", rawObject, "carbRatio",
				[]interface{}{map[string]interface{}{"amount": 12, "start": -1}},
				[]*service.Error{validator.ErrorValueNotTrue()},
			),
			Entry("start greater than 86400000", rawObject, "carbRatio",
				[]interface{}{map[string]interface{}{"amount": 12, "start": 86400001}},
				[]*service.Error{validator.ErrorValueNotTrue()},
			),
			Entry("amount negative", rawObject, "carbRatio",
				[]interface{}{map[string]interface{}{"amount": -1, "start": 21600000}},
				[]*service.Error{validator.ErrorValueNotTrue()},
			),
			Entry("amount greater than 250", rawObject, "carbRatio",
				[]interface{}{map[string]interface{}{"amount": 251, "start": 21600000}},
				[]*service.Error{validator.ErrorValueNotTrue()},
			),
		)

		DescribeTable("valid when", testing.ExpectFieldIsValid,
			Entry("start and amount within bounds", rawObject, "carbRatio",
				[]interface{}{map[string]interface{}{"amount": 12.0, "start": 0}},
			),
		)

	})

	Context("insulinSensitivity", func() {

		DescribeTable("invalid when", testing.ExpectFieldNotValid,
			Entry("start negative", rawObject, "insulinSensitivity",
				[]interface{}{map[string]interface{}{"amount": 12, "start": -1}},
				[]*service.Error{validator.ErrorValueNotTrue()},
			),
			Entry("start greater than 86400000", rawObject, "insulinSensitivity",
				[]interface{}{map[string]interface{}{"amount": 12, "start": 86400001}},
				[]*service.Error{validator.ErrorValueNotTrue()},
			),
			Entry("amount negative", rawObject, "insulinSensitivity",
				[]interface{}{map[string]interface{}{"amount": -0.1, "start": 21600000}},
				[]*service.Error{validator.ErrorValueNotTrue()},
			),
			Entry("amount greater than 1000.0", rawObject, "insulinSensitivity",
				[]interface{}{map[string]interface{}{"amount": 1000.1, "start": 21600000}},
				[]*service.Error{validator.ErrorValueNotTrue()},
			),
		)

		DescribeTable("valid when", testing.ExpectFieldIsValid,
			Entry("start and amount within bounds", rawObject, "insulinSensitivity",
				[]interface{}{map[string]interface{}{"amount": 12, "start": 0}},
			),
		)

	})

	Context("bgTarget", func() {

		Context("start, target, range", func() {

			DescribeTable("invalid when", testing.ExpectFieldNotValid,
				Entry("start negative", rawObject, "bgTarget",
					[]interface{}{map[string]interface{}{"start": -1, "target": 99.0, "range": 15}},
					[]*service.Error{validator.ErrorValueNotTrue()},
				),
				Entry("start greater than 86400000", rawObject, "bgTarget",
					[]interface{}{map[string]interface{}{"start": 86400001, "target": 99.0, "range": 15}},
					[]*service.Error{validator.ErrorValueNotTrue()},
				),
				Entry("target negative", rawObject, "bgTarget",
					[]interface{}{map[string]interface{}{"start": 21600000, "target": -0.1, "range": 15}},
					[]*service.Error{validator.ErrorValueNotTrue()},
				),
				Entry("target greater than 1000.0", rawObject, "bgTarget",
					[]interface{}{map[string]interface{}{"start": 21600000, "target": 1000.1, "range": 15}},
					[]*service.Error{validator.ErrorValueNotTrue()},
				),
				Entry("range negative", rawObject, "bgTarget",
					[]interface{}{map[string]interface{}{"start": 21600000, "target": 99.0, "range": -1}},
					[]*service.Error{validator.ErrorValueNotTrue()},
				),
				Entry("range greater than 51", rawObject, "bgTarget",
					[]interface{}{map[string]interface{}{"start": 21600000, "target": 199.0, "range": 51}},
					[]*service.Error{validator.ErrorValueNotTrue()},
				),
			)

			DescribeTable("valid when", testing.ExpectFieldIsValid,
				Entry("within bounds", rawObject, "bgTarget",
					[]interface{}{map[string]interface{}{"start": 21600000, "target": 99.9, "range": 10}},
				),
			)
		})

		Context("start, target, high", func() {

			DescribeTable("invalid when", testing.ExpectFieldNotValid,
				Entry("start negative", rawObject, "bgTarget",
					[]interface{}{map[string]interface{}{"start": -1, "target": 99.0, "high": 180.0}},
					[]*service.Error{validator.ErrorValueNotTrue()},
				),
				Entry("start greater than 86400000", rawObject, "bgTarget",
					[]interface{}{map[string]interface{}{"start": 86400001, "target": 99.0, "high": 180.0}},
					[]*service.Error{validator.ErrorValueNotTrue()},
				),
				Entry("target negative", rawObject, "bgTarget",
					[]interface{}{map[string]interface{}{"start": 21600000, "target": -0.1, "high": 180.0}},
					[]*service.Error{validator.ErrorValueNotTrue()},
				),
				Entry("target greater than 1000.0", rawObject, "bgTarget",
					[]interface{}{map[string]interface{}{"start": 21600000, "target": 1000.1, "high": 180.0}},
					[]*service.Error{validator.ErrorValueNotTrue(), validator.ErrorValueNotTrue()},
				),
				Entry("high less than target", rawObject, "bgTarget",
					[]interface{}{map[string]interface{}{"start": 21600000, "target": 90.0, "high": 80.0}},
					[]*service.Error{validator.ErrorValueNotTrue()},
				),
				Entry("high greater than 1000.0", rawObject, "bgTarget",
					[]interface{}{map[string]interface{}{"start": 21600000, "target": 90.0, "high": 1000.1}},
					[]*service.Error{validator.ErrorValueNotTrue()},
				),
			)

			DescribeTable("valid when", testing.ExpectFieldIsValid,
				Entry("within bounds", rawObject, "bgTarget",
					[]interface{}{map[string]interface{}{"start": 21600000, "target": 99.9, "high": 180.0}},
				),
			)
		})

	})

	Context("normalized", func() {

		DescribeTable("normalization when mmol/L", func(lowVal, lowExpected, highVal, highExpected, targetVal, targetExpected float64) {
			pumpSettings := pump.New()
			pumpSettings.Units = &pump.Units{BloodGlucose: &common.MmolL}

			pumpSettings.BloodGlucoseTargets = &[]*pump.BloodGlucoseTarget{
				&pump.BloodGlucoseTarget{High: &highVal, Low: &lowVal, Target: &targetVal},
				&pump.BloodGlucoseTarget{High: &highVal, Low: &lowVal, Target: &targetVal},
			}

			testContext := context.NewStandard()
			standardNormalizer, err := normalizer.NewStandard(testContext)
			Expect(err).To(BeNil())
			pumpSettings.Normalize(standardNormalizer)
			Expect(pumpSettings.Units.BloodGlucose).To(Equal(&common.MmolL))

			for _, bgTarget := range *pumpSettings.BloodGlucoseTargets {
				Expect(bgTarget.High).To(Equal(&highExpected))
				Expect(bgTarget.Low).To(Equal(&lowExpected))
				Expect(bgTarget.Target).To(Equal(&targetExpected))
			}
		},
			Entry("expected lower bg value", 2.1, 2.1, 3.1, 3.1, 2.5, 2.5),
			Entry("below max", 54.0, 54.0, 55.0, 55.0, 54.5, 54.5),
			Entry("expected upper bg value", 4.0, 4.0, 12.0, 12.0, 8.0, 8.0),
		)

		DescribeTable("normalization when mg/dL", func(lowVal, lowExpected, highVal, highExpected, targetVal, targetExpected float64) {
			pumpSettings := pump.New()
			pumpSettings.Units = &pump.Units{BloodGlucose: &common.MgdL}

			pumpSettings.BloodGlucoseTargets = &[]*pump.BloodGlucoseTarget{
				&pump.BloodGlucoseTarget{High: &highVal, Low: &lowVal, Target: &targetVal},
				&pump.BloodGlucoseTarget{High: &highVal, Low: &lowVal, Target: &targetVal},
			}

			testContext := context.NewStandard()
			standardNormalizer, err := normalizer.NewStandard(testContext)
			Expect(err).To(BeNil())
			pumpSettings.Normalize(standardNormalizer)
			Expect(pumpSettings.Units.BloodGlucose).To(Equal(&common.MmolL))

			for _, bgTarget := range *pumpSettings.BloodGlucoseTargets {
				Expect(bgTarget.High).To(Equal(&highExpected))
				Expect(bgTarget.Low).To(Equal(&lowExpected))
				Expect(bgTarget.Target).To(Equal(&targetExpected))
			}
		},
			Entry("expected lower bg value", 50.0, 2.7753739955227665, 55.0, 3.0529113950750433, 52.0, 2.8863889553436772),
			Entry("below max", 970.0, 53.84225551314167, 990.0, 54.95240511135078, 980.0, 54.397330312246226),
			Entry("expected upper bg value", 70.0, 3.8855235937318735, 180.0, 9.991346383881961, 99.0, 5.495240511135078),
		)
	})

})
