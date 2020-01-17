package timechange

import (
	"github.com/tidepool-org/platform/data"
	dataTypesDevice "github.com/tidepool-org/platform/data/types/device"
	"github.com/tidepool-org/platform/structure"
	structureValidator "github.com/tidepool-org/platform/structure/validator"
)

const (
	MethodAutomatic = "automatic"
	MethodManual    = "manual"
	SubType         = "timeChange" // TODO: Rename Type to "device/timeChange"; remove SubType
)

func Methods() []string {
	return []string{
		MethodAutomatic,
		MethodManual,
	}
}

type TimeChange struct {
	dataTypesDevice.Device `bson:",inline"`

	From   *Info   `json:"from,omitempty" bson:"from,omitempty"`
	Method *string `json:"method,omitempty" bson:"method,omitempty"`
	To     *Info   `json:"to,omitempty" bson:"to,omitempty"`

	Change *Change `json:"change,omitempty" bson:"change,omitempty"` // TODO: DEPRECATED
}

func New() *TimeChange {
	return &TimeChange{
		Device: dataTypesDevice.New(SubType),
	}
}

func (t *TimeChange) Parse(parser structure.ObjectParser) {
	if !parser.HasMeta() {
		parser = parser.WithMeta(t.Meta())
	}

	t.Device.Parse(parser)

	t.From = ParseInfo(parser.WithReferenceObjectParser("from"))
	t.Method = parser.String("method")
	t.To = ParseInfo(parser.WithReferenceObjectParser("to"))

	t.Change = ParseChange(parser.WithReferenceObjectParser("change"))
}

func (t *TimeChange) Validate(validator structure.Validator) {
	if !validator.HasMeta() {
		validator = validator.WithMeta(t.Meta())
	}

	t.Device.Validate(validator)

	if t.SubType != "" {
		validator.String("subType", &t.SubType).EqualTo(SubType)
	}

	fromValidator := validator.WithReference("from")
	methodValidator := validator.String("method", t.Method)
	toValidator := validator.WithReference("to")

	changeValidator := validator.WithReference("change")

	if t.From != nil || t.Method != nil || t.To != nil {
		if t.From != nil {
			t.From.Validate(fromValidator)
		}
		methodValidator.OneOf(Methods()...)
		if t.To != nil {
			t.To.Validate(toValidator)
		} else {
			toValidator.ReportError(structureValidator.ErrorValueNotExists())
		}
		if t.Change != nil {
			changeValidator.ReportError(structureValidator.ErrorValueExists())
		}
	} else if t.Change != nil {
		if t.From != nil {
			fromValidator.ReportError(structureValidator.ErrorValueExists())
		}
		methodValidator.NotExists()
		if t.To != nil {
			toValidator.ReportError(structureValidator.ErrorValueExists())
		}
		t.Change.Validate(changeValidator)
	} else {
		validator.ReportError(structureValidator.ErrorValuesNotExistForOne("to", "change"))
	}
}

// IsValid
// returns true if there is no error in the validator
func (t *TimeChange) IsValid(validator structure.Validator) bool {
	return !(validator.HasError())
}

func (t *TimeChange) Normalize(normalizer data.Normalizer) {
	if !normalizer.HasMeta() {
		normalizer = normalizer.WithMeta(t.Meta())
	}

	t.Device.Normalize(normalizer)

	if t.From != nil {
		t.From.Normalize(normalizer.WithReference("from"))
	}
	if t.To != nil {
		t.To.Normalize(normalizer.WithReference("to"))
	}

	if t.Change != nil {
		t.Change.Normalize(normalizer.WithReference("change"))
	}
}
