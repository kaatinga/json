package json

const unicodeMask = 0xf

type Scanner struct {
	position     int
	parsedNumber int64
	sample       []byte
	data         []byte
	parsedData   []byte
	byte
	value      bool // the field is found and we ready to read value
	pass       bool // to pass value indicator
	readText   bool
	parsedBool bool
}

// NewScanner creates new scanner with a sample inside.
func NewScanner(sample []byte) (*Scanner, error) {
	if len(sample) == 0 || len(sample) > maximumSampleLength {
		return nil, ErrInvalidSampleLength
	}

	return &Scanner{sample: sample}, nil
}

// ParsedData returns parsed value.
func (s *Scanner) ParsedData() []byte {
	return s.parsedData
}

// newParameter makes the Scanner ready to parse the next parameter.
func (s *Scanner) newParameter() {
	s.pass = false
	s.value = false
	s.readText = false
}

// SeekIn processes the input data looking for sample value.
// It returns warnings if the value is not string and errors if an error occurred.
func (s *Scanner) SeekIn(data []byte) error {

	s.reset()

	// check if the sample is set
	if s.sample == nil {
		return ErrSampleNotSet
	}

	// check if the length of the data is correct
	if len(data) < minimumJSONLength+len(s.sample) || len(data) > 256 {
		return ErrInvalidDataLength
	}

	// check if the dataset begins with  '{'
	if data[s.position] != ObjectStart {
		return ErrInvalidJSON
	}

	s.position++
	s.data = data

	for ; s.position < len(s.data); s.position++ {

		s.byte = s.data[s.position]

		if s.byte != Comma && s.pass {
			//fmt.println("fast continue")
			continue
		}

		//fmt.printf("we are checking %#v\n", s.byte)
		//fmt.println("position", s.position)

		// passing all whitespaces if it is not value
		if !s.readText && whitespace[s.byte] {
			//fmt.println("passing a whitespace")
			continue
		}

		switch s.byte {
		case QuotationMark:
			if !s.value {
				//fmt.println("checking the token...")
				s.pass = !s.validateToken()
				continue
			}

			//fmt.println("we have to read data now")
			s.readText = true
			continue
		case Colon:
			//fmt.println("waiting for value")
			if s.value {
				//fmt.println("we are waiting for a value already, colon is repeated")
				return ErrInvalidJSON
			}
			s.value = true
			continue
		case ObjectStart:
			//fmt.println("internal objects are not supported")
			s.pass = true
			continue
		case ObjectEnd:
			return WarnNotFound
		case Comma:
			//fmt.println("next sample begins")
			s.newParameter()
			continue
		case True:
			if s.validate(jsonTrue) {
				s.parsedBool = true
				return WarnBoolWasFound
			}
			return ErrInvalidJSON
		case False:
			if s.validate(jsonFalse) {
				return WarnBoolWasFound
			}
			return ErrInvalidJSON
		case Null:
			if s.validate(jsonNull) {
				return WarnNullWasFound
			}
			return ErrInvalidJSON
		case ArrayEnd:
			//fmt.println("array is unsupported, soon we will finish")
			continue
		case ArrayStart:
			//fmt.println("array is unsupported")
			s.pass = true
			continue
		default:
			if s.readText {
				// we found the value start position, the final step is to read data
				return s.readString()
			}
			if s.value {
				//fmt.println("numbers case assumed")
				return s.readNumber()
			}

			//fmt.println("read is not true")
			continue
		}
	}
	return WarnNotFound
}

// validate checks the value and input sample.
func (s *Scanner) validate(sample []byte) bool {
	//fmt.println("checking bool value:", string(sample))
	start := s.position

	end := s.position + len(sample)
	//fmt.println("установили конечную точку сравнения")

	if end > len(s.data) {
		//fmt.println("конечная точка неожиданно дальше максимальной длины данных")
		return false
	}

	for ; s.position != end; s.position++ {
		if s.position < len(s.data) && sample[s.position-start] != s.data[s.position] {
			//fmt.println("один из байтов имени переменной не совпадает или мы вышли за рамки")
			//fmt.printf("comparing %#v and %#v\n", sample[s.position-start], s.data[s.position])
			//fmt.println("position", s.position)
			//fmt.println("end", end)
			return false
		}
	}
	return true
}

func (s *Scanner) reset() {
	s.position = 0
	//s.byte = 0
	//sample     []byte
	//s.data = nil
	s.value = false
	s.pass = false
	s.readText = false
	//s.parsedData = nil
	//s.parsedBool = false
	s.parsedNumber = 0
}

// validateToken compares the sample and the data starting the position.
func (s *Scanner) validateToken() bool {

	if s.data[s.position] == QuotationMark {
		//fmt.println("проверяем название переменной. сдвигаем вперед так как кавычка")
		s.position++
	}

	start := s.position
	end := s.position + len(s.sample)
	//fmt.println("установили конечную точку сравнения")

	if end > len(s.data) {
		//fmt.println("конечная точка неожиданно дальше максимальной длины данных")
		return false
	}

	for ; s.position != end; s.position++ {
		if s.position < len(s.data) && s.sample[s.position-start] != s.data[s.position] {
			//fmt.println("один из байтов имени переменной не совпадает или мы вышли за рамки")
			//fmt.printf("comparing %#v and %#v\n", s.sample[s.position-start], s.data[s.position])
			//fmt.println("position", s.position)
			//fmt.println("end", end)
			return false
		}
	}

	// as we are here, we can check colon that it fallows the last position
	if s.position < len(s.data) && s.data[s.position] != QuotationMark {
		//fmt.println("the data were similar but not the sample :)")
		return false
	}

	return true
}

// readString reads value data starting the position.
func (s *Scanner) readString() error {

	var end int
	var start = s.position
	//fmt.println("value start is set:", start)
	for ; s.position < len(s.data); s.position++ {
		if s.data[s.position] == QuotationMark {
			end = s.position
			//fmt.println("value end is set:", end)
			break
		}
	}

	if end == 0 {
		return ErrInvalidValue
	}

	s.parsedData = s.data[start:end]
	return nil
}

// readNumber reads value data starting the position.
func (s *Scanner) readNumber() error {
	if s.data[s.position] == Minus {
		//fmt.println("the number will be negative")
		s.position++
		return s.checkNegativeNumber()
	}
	return s.checkNumber()
}

// checkNumber checks positive numbers.
func (s *Scanner) checkNumber() error {
	if s.data[s.position] == ObjectEnd || s.data[s.position] == Comma {
		return nil
	}
	if s.data[s.position] == Dot {
		return WarnFloatNotSupported
	}
	if s.data[s.position] < 0x30 || s.data[s.position] > 0x39 {
		//fmt.printf("value %#v\n", s.data[s.position])
		return ErrInvalidJSON
	}
	s.parsedNumber = (s.parsedNumber << 3) + (s.parsedNumber << 1) + int64(s.data[s.position])&unicodeMask
	s.position++
	return s.checkNumber()
}

// checkNumber checks negative numbers.
func (s *Scanner) checkNegativeNumber() error {
	if s.data[s.position] == ObjectEnd || s.data[s.position] == Comma {
		return nil
	}
	if s.data[s.position] == Dot {
		return WarnFloatNotSupported
	}
	if s.data[s.position] < 0x30 || s.data[s.position] > 0x39 {
		return ErrInvalidJSON
	}
	s.parsedNumber = (s.parsedNumber << 3) + (s.parsedNumber << 1) - int64(s.data[s.position])&unicodeMask
	s.position++
	return s.checkNegativeNumber()
}

var whitespace = [256]bool{
	' ':  true,
	'\t': true,
	'\n': true,
	'\r': true,
}

func isWhitespace(c byte) bool {
	return whitespace[c]
}
