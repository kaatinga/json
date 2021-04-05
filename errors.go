package json

import cer "github.com/kaatinga/const-errs"

const (
	ErrInvalidSampleLength cer.Error = "invalid sample length"
	ErrInvalidDataLength   cer.Error = "invalid data length"
	ErrInvalidJSON         cer.Error = "json format is not valid"
	ErrSampleNotSet        cer.Error = "sample must be set first"

	ErrInvalidValue cer.Error = "the value was not read as data format is corrupted"

	WarnNotFound         cer.Warning = "the sample was not found in the provided data"
	WarnNullWasFound     cer.Warning = "null was detected as value"
	WarnBoolWasFound     cer.Warning = "bool value was detected as value"
	WarnUnsupportedArray cer.Warning = "arrays are unsupported"
)
