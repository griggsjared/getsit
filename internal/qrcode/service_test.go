package qrcode_test

import (
	"testing"

	"github.com/griggsjared/getsit/internal/qrcode"
)

func longContent(length int) string {
	var content string
	for i := 0; i < length; i++ {
		content += "a"
	}
	return content
}

func TestQRCode_NewQRCode(t *testing.T) {
	body := []byte("test")
	qr := qrcode.NewQRCode(body)
	if qr == nil {
		t.Errorf("NewQRCode() = %v, want a new QRCode", qr)
	}
}

func TestQRCode_String(t *testing.T) {
	body := []byte("test")
	qr := qrcode.NewQRCode(body)
	if qr.String() != "test" {
		t.Errorf("String() = %v, want %v", qr.String(), "test")
	}
}

func TestQRCode_Base64(t *testing.T) {
	body := []byte("test")
	qr := qrcode.NewQRCode(body)
	if qr.Base64() != "data:image/png;base64,dGVzdA==" {
		t.Errorf("Base64() = %v, want %v", qr.Base64(), "data:image/png;base64,dGVzdA==")
	}
}

func TestService_NewService(t *testing.T) {
	s := qrcode.NewService()
	if s == nil {
		t.Errorf("NewService() = %v, want a new Service", s)
	}
}

func TestService_Generate(t *testing.T) {

	s := qrcode.NewService()

	tests := []struct {
		name    string
		input   *qrcode.GenerateInput
		wantErr bool
	}{
		{
			name: "valid input",
			input: &qrcode.GenerateInput{
				Content: "https://example.com",
				Size:    256,
			},
			wantErr: false,
		},
		{
			name: "empty content",
			input: &qrcode.GenerateInput{
				Content: "",
				Size:    256,
			},
			wantErr: true,
		},
		{
			name: "content too large",
			input: &qrcode.GenerateInput{
				Content: longContent(1664), //apparently 1664 is the max length for a QR code
				Size:    256,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.Generate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
