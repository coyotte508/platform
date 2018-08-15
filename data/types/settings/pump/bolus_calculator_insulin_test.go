package pump_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"math"

	dataNormalizer "github.com/tidepool-org/platform/data/normalizer"
	"github.com/tidepool-org/platform/data/types/settings/pump"
	testDataTypes "github.com/tidepool-org/platform/data/types/test"
	testErrors "github.com/tidepool-org/platform/errors/test"
	"github.com/tidepool-org/platform/pointer"
	"github.com/tidepool-org/platform/structure"
	structureValidator "github.com/tidepool-org/platform/structure/validator"
	"github.com/tidepool-org/platform/test"
)

func NewBolusCalculatorInsulin() *pump.BolusCalculatorInsulin {
	units := pointer.FromString(test.RandomStringFromArray(pump.BolusCalculatorInsulinUnits()))
	datum := pump.NewBolusCalculatorInsulin()
	datum.Duration = pointer.FromFloat64(test.RandomFloat64FromRange(pump.BolusCalculatorInsulinDurationRangeForUnits(units)))
	datum.Units = units
	return datum
}

func CloneBolusCalculatorInsulin(datum *pump.BolusCalculatorInsulin) *pump.BolusCalculatorInsulin {
	if datum == nil {
		return nil
	}
	clone := pump.NewBolusCalculatorInsulin()
	clone.Duration = test.CloneFloat64(datum.Duration)
	clone.Units = test.CloneString(datum.Units)
	return clone
}

var _ = Describe("BolusCalculatorInsulin", func() {
	It("BolusCalculatorInsulinDurationHoursMaximum is expected", func() {
		Expect(pump.BolusCalculatorInsulinDurationHoursMaximum).To(Equal(10.0))
	})

	It("BolusCalculatorInsulinDurationHoursMinimum is expected", func() {
		Expect(pump.BolusCalculatorInsulinDurationHoursMinimum).To(Equal(0.0))
	})

	It("BolusCalculatorInsulinUnitsHours is expected", func() {
		Expect(pump.BolusCalculatorInsulinUnitsHours).To(Equal("hours"))
	})

	It("BolusCalculatorInsulinUnits returns expected", func() {
		Expect(pump.BolusCalculatorInsulinUnits()).To(Equal([]string{"hours"}))
	})

	Context("ParseBolusCalculatorInsulin", func() {
		// TODO
	})

	Context("NewBolusCalculatorInsulin", func() {
		It("is successful", func() {
			Expect(pump.NewBolusCalculatorInsulin()).To(Equal(&pump.BolusCalculatorInsulin{}))
		})
	})

	Context("BolusCalculatorInsulin", func() {
		Context("Parse", func() {
			// TODO
		})

		Context("Validate", func() {
			DescribeTable("validates the datum",
				func(mutator func(datum *pump.BolusCalculatorInsulin), expectedErrors ...error) {
					datum := NewBolusCalculatorInsulin()
					mutator(datum)
					testDataTypes.ValidateWithExpectedOrigins(datum, structure.Origins(), expectedErrors...)
				},
				Entry("succeeds",
					func(datum *pump.BolusCalculatorInsulin) {},
				),
				Entry("units missing; duration missing",
					func(datum *pump.BolusCalculatorInsulin) {
						datum.Duration = nil
						datum.Units = nil
					},
					testErrors.WithPointerSource(structureValidator.ErrorValueNotExists(), "/duration"),
					testErrors.WithPointerSource(structureValidator.ErrorValueNotExists(), "/units"),
				),
				Entry("units missing; duration out of range (lower)",
					func(datum *pump.BolusCalculatorInsulin) {
						datum.Duration = pointer.FromFloat64(-0.1)
						datum.Units = nil
					},
					testErrors.WithPointerSource(structureValidator.ErrorValueNotExists(), "/units"),
				),
				Entry("units missing; duration in range (lower)",
					func(datum *pump.BolusCalculatorInsulin) {
						datum.Duration = pointer.FromFloat64(0.0)
						datum.Units = nil
					},
					testErrors.WithPointerSource(structureValidator.ErrorValueNotExists(), "/units"),
				),
				Entry("units missing; duration in range (upper)",
					func(datum *pump.BolusCalculatorInsulin) {
						datum.Duration = pointer.FromFloat64(10.0)
						datum.Units = nil
					},
					testErrors.WithPointerSource(structureValidator.ErrorValueNotExists(), "/units"),
				),
				Entry("units missing; duration out of range (upper)",
					func(datum *pump.BolusCalculatorInsulin) {
						datum.Duration = pointer.FromFloat64(10.1)
						datum.Units = nil
					},
					testErrors.WithPointerSource(structureValidator.ErrorValueNotExists(), "/units"),
				),
				Entry("units invalid; duration missing",
					func(datum *pump.BolusCalculatorInsulin) {
						datum.Duration = nil
						datum.Units = pointer.FromString("invalid")
					},
					testErrors.WithPointerSource(structureValidator.ErrorValueNotExists(), "/duration"),
					testErrors.WithPointerSource(structureValidator.ErrorValueStringNotOneOf("invalid", []string{"hours"}), "/units"),
				),
				Entry("units invalid; duration out of range (lower)",
					func(datum *pump.BolusCalculatorInsulin) {
						datum.Duration = pointer.FromFloat64(-0.1)
						datum.Units = pointer.FromString("invalid")
					},
					testErrors.WithPointerSource(structureValidator.ErrorValueStringNotOneOf("invalid", []string{"hours"}), "/units"),
				),
				Entry("units invalid; duration in range (lower)",
					func(datum *pump.BolusCalculatorInsulin) {
						datum.Duration = pointer.FromFloat64(0.0)
						datum.Units = pointer.FromString("invalid")
					},
					testErrors.WithPointerSource(structureValidator.ErrorValueStringNotOneOf("invalid", []string{"hours"}), "/units"),
				),
				Entry("units invalid; duration in range (upper)",
					func(datum *pump.BolusCalculatorInsulin) {
						datum.Duration = pointer.FromFloat64(10.0)
						datum.Units = pointer.FromString("invalid")
					},
					testErrors.WithPointerSource(structureValidator.ErrorValueStringNotOneOf("invalid", []string{"hours"}), "/units"),
				),
				Entry("units invalid; duration out of range (upper)",
					func(datum *pump.BolusCalculatorInsulin) {
						datum.Duration = pointer.FromFloat64(10.1)
						datum.Units = pointer.FromString("invalid")
					},
					testErrors.WithPointerSource(structureValidator.ErrorValueStringNotOneOf("invalid", []string{"hours"}), "/units"),
				),
				Entry("units hours: duration missing",
					func(datum *pump.BolusCalculatorInsulin) {
						datum.Duration = nil
						datum.Units = pointer.FromString("hours")
					},
					testErrors.WithPointerSource(structureValidator.ErrorValueNotExists(), "/duration"),
				),
				Entry("units hours: duration out of range (lower)",
					func(datum *pump.BolusCalculatorInsulin) {
						datum.Duration = pointer.FromFloat64(-0.1)
						datum.Units = pointer.FromString("hours")
					},
					testErrors.WithPointerSource(structureValidator.ErrorValueNotInRange(-0.1, 0.0, 10.0), "/duration"),
				),
				Entry("units hours: duration in range (lower)",
					func(datum *pump.BolusCalculatorInsulin) {
						datum.Duration = pointer.FromFloat64(0.0)
						datum.Units = pointer.FromString("hours")
					},
				),
				Entry("units hours: duration in range (upper)",
					func(datum *pump.BolusCalculatorInsulin) {
						datum.Duration = pointer.FromFloat64(10.0)
						datum.Units = pointer.FromString("hours")
					},
				),
				Entry("units hours: duration out of range (upper)",
					func(datum *pump.BolusCalculatorInsulin) {
						datum.Duration = pointer.FromFloat64(10.1)
						datum.Units = pointer.FromString("hours")
					},
					testErrors.WithPointerSource(structureValidator.ErrorValueNotInRange(10.1, 0.0, 10.0), "/duration"),
				),
				Entry("units missing",
					func(datum *pump.BolusCalculatorInsulin) { datum.Units = nil },
					testErrors.WithPointerSource(structureValidator.ErrorValueNotExists(), "/units"),
				),
				Entry("units invalid",
					func(datum *pump.BolusCalculatorInsulin) { datum.Units = pointer.FromString("invalid") },
					testErrors.WithPointerSource(structureValidator.ErrorValueStringNotOneOf("invalid", []string{"hours"}), "/units"),
				),
				Entry("units hours",
					func(datum *pump.BolusCalculatorInsulin) {
						datum.Duration = pointer.FromFloat64(0.0)
						datum.Units = pointer.FromString("hours")
					},
				),
				Entry("multiple errors",
					func(datum *pump.BolusCalculatorInsulin) {
						datum.Duration = nil
						datum.Units = nil
					},
					testErrors.WithPointerSource(structureValidator.ErrorValueNotExists(), "/duration"),
					testErrors.WithPointerSource(structureValidator.ErrorValueNotExists(), "/units"),
				),
			)
		})

		Context("Normalize", func() {
			DescribeTable("normalizes the datum",
				func(mutator func(datum *pump.BolusCalculatorInsulin)) {
					for _, origin := range structure.Origins() {
						datum := NewBolusCalculatorInsulin()
						mutator(datum)
						expectedDatum := CloneBolusCalculatorInsulin(datum)
						normalizer := dataNormalizer.New()
						Expect(normalizer).ToNot(BeNil())
						datum.Normalize(normalizer.WithOrigin(origin))
						Expect(normalizer.Error()).To(BeNil())
						Expect(normalizer.Data()).To(BeEmpty())
						Expect(datum).To(Equal(expectedDatum))
					}
				},
				Entry("does not modify the datum",
					func(datum *pump.BolusCalculatorInsulin) {},
				),
				Entry("does not modify the datum; duration missing",
					func(datum *pump.BolusCalculatorInsulin) { datum.Duration = nil },
				),
				Entry("does not modify the datum; units missing",
					func(datum *pump.BolusCalculatorInsulin) { datum.Units = nil },
				),
			)
		})
	})

	Context("BolusCalculatorInsulinDurationRangeForUnits", func() {
		It("returns expected range for units missing", func() {
			minimum, maximum := pump.BolusCalculatorInsulinDurationRangeForUnits(nil)
			Expect(minimum).To(Equal(-math.MaxFloat64))
			Expect(maximum).To(Equal(math.MaxFloat64))
		})

		It("returns expected range for units invalid", func() {
			minimum, maximum := pump.BolusCalculatorInsulinDurationRangeForUnits(pointer.FromString("invalid"))
			Expect(minimum).To(Equal(-math.MaxFloat64))
			Expect(maximum).To(Equal(math.MaxFloat64))
		})

		It("returns expected range for units units", func() {
			minimum, maximum := pump.BolusCalculatorInsulinDurationRangeForUnits(pointer.FromString("hours"))
			Expect(minimum).To(Equal(0.0))
			Expect(maximum).To(Equal(10.0))
		})
	})
})