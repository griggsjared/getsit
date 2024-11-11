package qrcode

import (
	"crypto/md5"

	"github.com/skip2/go-qrcode"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

type QRCode struct {
	Output []byte
}

type GenerateInput struct {
	Content string
}

func (s *Service) Generate(input *GenerateInput) (*QRCode, error) {

	//we will generate an md5 hash of the content to use as the filename
	md5 := md5.New()
	md5.Write([]byte(input.Content))

	//if that file already exists, return it
	// if exists, err := s.f.Open(string(md5.Sum(nil))); err == nil {
	// 	b := make([]byte, 1024)
	// 	exists.Read(b)
	// 	return &QRCode{
	// 		Output: b,
	// 	}, nil
	// }

	// Generate the qr code
	output, err := generateQRCode(input.Content, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}

	return &QRCode{
		Output: output,
	}, nil
}

func generateQRCode(content string, quality qrcode.RecoveryLevel, size int) ([]byte, error) {
	var png []byte
	png, err := qrcode.Encode(content, quality, size)

	if err != nil {
		return nil, err
	}

	return png, nil
}
