package json

import (
	"reflect"
	"testing"
)

// samples
var (
	jsonExample = []byte{ObjectStart, 0x22, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x22, 0x3a, 0x22, 0x31, 0x32, 0x33, 0x34, 0x35, 0x22, ObjectEnd}
	secret      = []byte{0x73, 0x65, 0x63, 0x72, 0x65, 0x74}
)

func TestScanner_validateToken(t *testing.T) {

	tests := []struct {
		name    string
		scanner Scanner
		sample  []byte
		want    bool
	}{
		{"ok", Scanner{
			data: []byte{0x73, 0x65, 0x63, 0x72, 0x65, 0x74},
		}, secret, true},
		{"!ok", Scanner{
			data: []byte{0x73, 0x65, 0x63, 0x72, 0x65, 0x75},
		}, secret, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &tt.scanner
			err := s.SetSample(tt.sample)
			if err != nil {
				t.Errorf("setSample error = %v, want nil", err)
			}
			if got := s.validateToken(); got != tt.want {
				t.Errorf("validateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanner_readData(t *testing.T) {

	tests := []struct {
		name    string
		scanner Scanner
		result  []byte
		wantErr bool
	}{
		{"ok", Scanner{
			position: 11,
			byte:     0,
			sample:   nil,
			data:     jsonExample,
		}, []byte{0x31, 0x32, 0x33, 0x34, 0x35}, false},
		{"ok", Scanner{
			position: 11,
			byte:     0,
			sample:   nil,
			data:     []byte{ObjectStart, 0x22, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x22, 0x3a, 0x22, 0x31, 0x32, 0x33, 0x34, 0x35, 0x21, ObjectEnd},
		}, nil, true},
		{"ok", Scanner{
			position: 11,
			byte:     0,
			sample:   nil,
			data:     []byte{ObjectStart, 0x22, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x22, 0x3a, 0x22, 0x31, 0x32, 0x33, 0x34, 0x35, 0x21},
		}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &tt.scanner
			if err := s.readData(); (err != nil) != tt.wantErr {
				t.Errorf("readData() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(s.parsedData, tt.result) {
				t.Errorf("parsed data is incorrect\n have %v, want %v", s.parsedData, tt.result)
			}
		})
	}
}

func TestScanner_SetSample(t *testing.T) {

	tests := []struct {
		name    string
		scanner Scanner
		sample  []byte
		wantErr bool
	}{
		{"ok", Scanner{}, secret, false},
		{"!ok1", Scanner{}, []byte{}, true},
		{"!ok2", Scanner{}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &tt.scanner
			if err := s.SetSample(tt.sample); (err != nil) != tt.wantErr {
				t.Errorf("SetSample() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestScanner_Seek(t *testing.T) {

	tests := []struct {
		name    string
		scanner Scanner
		sample  []byte
		data    []byte
		wantErr error
	}{
		{"ok1", Scanner{}, secret, jsonExample, nil},
		{"!ok1", Scanner{}, []byte("mmm"), jsonExample, nil},
		{"!ok2", Scanner{}, []byte("mmm"), []byte("mmmmm"), ErrInvalidJSON},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &tt.scanner

			err := s.SetSample(tt.sample)
			if err != nil {
				t.Errorf("SetSample() error, %s", err)
			}

			err = s.SetData(tt.data)
			if err != nil {
				t.Errorf("SetData() error, %s", err)
			}

			err = s.Seek()
			if err != tt.wantErr {
				t.Errorf("Seek() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
