package rgbqrcode

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeQrCode(t *testing.T) {
	tests := []struct {
		name     string
		logopath string
		qrsize   int
		QrValue  string
		err      error
	}{
		{
			name:     "Generate qr code with not set logo path",
			logopath: "",
			err:      errors.New("image: unknown format"),
		},
		{
			name:     "encode qr code with logo size 48 x 48",
			logopath: "48_48_logo.png",
			err:      errors.New("Please use logo with size 50 X 50"),
		},
		{
			name:     "encode qr code with logo size 50 X 50 with size 0",
			logopath: "50_50_logo.png",
			qrsize:   0,
			QrValue:  "test",
			err:      nil,
		},
		{
			name:     "encode qr code with logo size 50 X 50 with size 0 and empty value",
			logopath: "50_50_logo.png",
			qrsize:   0,
			QrValue:  "",
			err:      errors.New("no data to encode"),
		},
		{
			name:     "encode qr code with size 50 X 50 with size 256 ",
			logopath: "50_50_logo.png",
			qrsize:   256,
			QrValue:  "test",
			err:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			pwd, _ := os.Getwd()

			pathOfLogo := fmt.Sprintf("%s/%s", pwd, tt.logopath)
			log.Println("PATH OF LOGO => ", pathOfLogo)
			// load logo from path
			file, _ := os.Open(pathOfLogo)
			defer file.Close()

			qr, err := New(GenerateParam{
				LogoPath: file,
				QrSize:   tt.qrsize,
				QrValue:  tt.QrValue,
			})
			if err != nil {
				assert.Error(t, tt.err, err)
			}

			// do encode qr
			_, err = qr.Encode()
			if err != nil {
				assert.Error(t, tt.err, err)
			}
		})
	}
}
