package validator

/* CHECKLIST
 * [x] Uses interfaces as appropriate
 * [x] Private package variables use underscore prefix
 * [x] All parameters validated
 * [x] All errors handled
 * [x] Reviewed for concurrency safety
 * [x] Code complete
 * [x] Full test coverage
 */

import "github.com/tidepool-org/platform/data"

type StandardFloat struct {
	context   data.Context
	reference interface{}
	value     *float64
}

func NewStandardFloat(context data.Context, reference interface{}, value *float64) *StandardFloat {
	if context == nil {
		return nil
	}

	return &StandardFloat{
		context:   context,
		reference: reference,
		value:     value,
	}
}

func (s *StandardFloat) Exists() data.Float {
	if s.value == nil {
		s.context.AppendError(s.reference, ErrorValueDoesNotExist())
	}
	return s
}

func (s *StandardFloat) EqualTo(value float64) data.Float {
	if s.value != nil {
		if *s.value != value {
			s.context.AppendError(s.reference, ErrorValueNotEqualTo(*s.value, value))
		}
	}
	return s
}

func (s *StandardFloat) NotEqualTo(value float64) data.Float {
	if s.value != nil {
		if *s.value == value {
			s.context.AppendError(s.reference, ErrorValueEqualTo(*s.value, value))
		}
	}
	return s
}

func (s *StandardFloat) LessThan(limit float64) data.Float {
	if s.value != nil {
		if *s.value >= limit {
			s.context.AppendError(s.reference, ErrorValueNotLessThan(*s.value, limit))
		}
	}
	return s
}

func (s *StandardFloat) LessThanOrEqualTo(limit float64) data.Float {
	if s.value != nil {
		if *s.value > limit {
			s.context.AppendError(s.reference, ErrorValueNotLessThanOrEqualTo(*s.value, limit))
		}
	}
	return s
}

func (s *StandardFloat) GreaterThan(limit float64) data.Float {
	if s.value != nil {
		if *s.value <= limit {
			s.context.AppendError(s.reference, ErrorValueNotGreaterThan(*s.value, limit))
		}
	}
	return s
}

func (s *StandardFloat) GreaterThanOrEqualTo(limit float64) data.Float {
	if s.value != nil {
		if *s.value < limit {
			s.context.AppendError(s.reference, ErrorValueNotGreaterThanOrEqualTo(*s.value, limit))
		}
	}
	return s
}

func (s *StandardFloat) InRange(lowerLimit float64, upperLimit float64) data.Float {
	if s.value != nil {
		if *s.value < lowerLimit || *s.value > upperLimit {
			s.context.AppendError(s.reference, ErrorFloatNotInRange(*s.value, lowerLimit, upperLimit))
		}
	}
	return s
}

func (s *StandardFloat) OneOf(allowedValues []float64) data.Float {
	if s.value != nil {
		for _, possibleValue := range allowedValues {
			if possibleValue == *s.value {
				return s
			}
		}
		s.context.AppendError(s.reference, ErrorFloatNotOneOf(*s.value, allowedValues))
	}
	return s
}

func (s *StandardFloat) NotOneOf(disallowedValues []float64) data.Float {
	if s.value != nil {
		for _, possibleValue := range disallowedValues {
			if possibleValue == *s.value {
				s.context.AppendError(s.reference, ErrorFloatOneOf(*s.value, disallowedValues))
				return s
			}
		}
	}
	return s
}
