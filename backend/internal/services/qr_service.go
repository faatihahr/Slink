package services

import (
	"bytes"
	"log"

	"github.com/skip2/go-qrcode"
)

type QRService struct {
	size int
}

func NewQRService(size int) *QRService {
	return &QRService{size: size}
}

func (s *QRService) GenerateQRCode(url string) ([]byte, error) {
	qrCode, err := qrcode.New(url, qrcode.Medium)
	if err != nil {
		log.Printf("Error generating QR code: %v", err)
		return nil, err
	}

	qrCode.DisableBorder = true

	var buf bytes.Buffer
	err = qrCode.Write(s.size, &buf)
	if err != nil {
		log.Printf("Error writing QR code to buffer: %v", err)
		return nil, err
	}

	return buf.Bytes(), nil
}
