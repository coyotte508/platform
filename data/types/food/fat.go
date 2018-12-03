package food

import (
	"github.com/tidepool-org/platform/data"
	"github.com/tidepool-org/platform/structure"
)

const (
	FatTotalGramsMaximum = 1000.0
	FatTotalGramsMinimum = 0.0
	FatUnitsGrams        = "grams"
)

func FatUnits() []string {
	return []string{
		FatUnitsGrams,
	}
}

type Fat struct {
	Total *float64 `json:"total,omitempty" bson:"total,omitempty"`
	Units *string  `json:"units,omitempty" bson:"units,omitempty"`
}

func ParseFat(parser data.ObjectParser) *Fat {
	if parser.Object() == nil {
		return nil
	}
	datum := NewFat()
	datum.Parse(parser)
	parser.ProcessNotParsed()
	return datum
}

func NewFat() *Fat {
	return &Fat{}
}

func (f *Fat) Parse(parser data.ObjectParser) {
	f.Total = parser.ParseFloat("total")
	f.Units = parser.ParseString("units")
}

func (f *Fat) Validate(validator structure.Validator) {
	validator.Float64("total", f.Total).Exists().InRange(FatTotalGramsMinimum, FatTotalGramsMaximum)
	validator.String("units", f.Units).Exists().OneOf(FatUnits()...)
}

func (f *Fat) Normalize(normalizer data.Normalizer) {}
