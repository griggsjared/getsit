package qrcode

import (
	"encoding/base64"
	"errors"

	"github.com/skip2/go-qrcode"
)

// ErrValidation is a generic validation error that can be returned when input validation fails
var ErrValidation = errors.New("validation error")

// QRCode is a struct that will hold the body of the QRCode
type QRCode []byte

// NewQRCode will create a new QRCode with an array of bytes as the body
func NewQRCode(body []byte) *QRCode {
	qr := QRCode(body)
	return &qr
}

// String will return the QRCode as a string
func (qr QRCode) String() string {
	return string(qr)
}

// Base64 will return the QRCode as a base64 encoded string
func (qr QRCode) Base64() string {
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(qr)
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

type GenerateInput struct {
	Content string
	Size    int
}

func (s *Service) Generate(input *GenerateInput) (*QRCode, error) {

	var png []byte
	png, err := qrcode.Encode(input.Content, qrcode.High, input.Size)
	if err != nil {
		return nil, err
	}

	return NewQRCode(png), nil
}
