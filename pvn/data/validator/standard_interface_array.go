package validator

/* CHECKLIST
 * [ ] Uses interfaces as appropriate
 * [ ] Private package variables use underscore prefix
 * [ ] All parameters validated
 * [ ] All errors handled
 * [ ] Reviewed for concurrency safety
 * [ ] Code complete
 * [ ] Full test coverage
 */

import "github.com/tidepool-org/platform/pvn/data"

type StandardInterfaceArray struct {
	context   data.Context
	reference interface{}
	value     *[]interface{}
}

func NewStandardInterfaceArray(context data.Context, reference interface{}, value *[]interface{}) *StandardInterfaceArray {
	if context == nil {
		return nil
	}

	return &StandardInterfaceArray{
		context:   context,
		reference: reference,
		value:     value,
	}
}

func (s *StandardInterfaceArray) Exists() data.InterfaceArray {
	if s.value == nil {
		s.context.AppendError(s.reference, ErrorValueDoesNotExist())
	}
	return s
}

func (s *StandardInterfaceArray) LengthEqualTo(limit int) data.InterfaceArray {
	if s.value != nil {
		if length := len(*s.value); length != limit {
			s.context.AppendError(s.reference, ErrorLengthNotEqualTo(length, limit))
		}
	}
	return s
}

func (s *StandardInterfaceArray) LengthNotEqualTo(limit int) data.InterfaceArray {
	if s.value != nil {
		if length := len(*s.value); length == limit {
			s.context.AppendError(s.reference, ErrorLengthEqualTo(length, limit))
		}
	}
	return s
}

func (s *StandardInterfaceArray) LengthLessThan(limit int) data.InterfaceArray {
	if s.value != nil {
		if length := len(*s.value); length >= limit {
			s.context.AppendError(s.reference, ErrorLengthNotLessThan(length, limit))
		}
	}
	return s
}

func (s *StandardInterfaceArray) LengthLessThanOrEqualTo(limit int) data.InterfaceArray {
	if s.value != nil {
		if length := len(*s.value); length > limit {
			s.context.AppendError(s.reference, ErrorLengthNotLessThanOrEqualTo(length, limit))
		}
	}
	return s
}

func (s *StandardInterfaceArray) LengthGreaterThan(limit int) data.InterfaceArray {
	if s.value != nil {
		if length := len(*s.value); length <= limit {
			s.context.AppendError(s.reference, ErrorLengthNotGreaterThan(length, limit))
		}
	}
	return s
}

func (s *StandardInterfaceArray) LengthGreaterThanOrEqualTo(limit int) data.InterfaceArray {
	if s.value != nil {
		if length := len(*s.value); length < limit {
			s.context.AppendError(s.reference, ErrorLengthNotGreaterThanOrEqualTo(length, limit))
		}
	}
	return s
}

func (s *StandardInterfaceArray) LengthInRange(lowerlimit int, upperLimit int) data.InterfaceArray {
	if s.value != nil {
		if length := len(*s.value); length < lowerlimit || length > upperLimit {
			s.context.AppendError(s.reference, ErrorLengthNotInRange(length, lowerlimit, upperLimit))
		}
	}
	return s
}