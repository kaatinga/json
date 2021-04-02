package json

const (
	maximumSampleLength = 65535
	minimumJSONLength   = 2

	// samples
	//  { Object start
	//  [ Array start
	//  } Object end
	//  ] Array End
	//  , Literal comma
	//  : Literal colon
	//  t JSON true
	//  f JSON false
	//  n JSON null
	//  " A string, possibly containing backslash escaped entities.
	//  -, 0-9 A number

	ObjectStart   byte = '{'
	ObjectEnd     byte = '}'
	Colon         byte = ':'
	Comma         byte = ','
	ArrayStart    byte = '['
	ArrayEnd      byte = ']'
	QuotationMark byte = '"'
	True          byte = 't'
	False         byte = 'f'
	Null          byte = 'n'
)

var (
	jsonTrue  = []byte("true")
	jsonFalse = []byte("false")
	jsonNull  = []byte("null")
)
