package rgbqrcode

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	_ "image/png"
	"os"

	"github.com/skip2/go-qrcode"
)

type (
	GenerateParam struct {
		// Image or logo file path
		// Ex: "name.png"
		LogoPath *os.File

		// Value or data will encoded to QR CODE
		QrValue string
		// Size of QRCODE
		// Default: 256
		QrSize int
	}
	encodeQr struct {
		imageLogo image.Image
		value     string
		size      int
	}
	ResultEncode struct {
		Base64 string
		PNG    bytes.Buffer
	}
)

func New(param GenerateParam) (encodeQr, error) {
	imageLogo, _, err := image.Decode(param.LogoPath)
	if err != nil {
		return encodeQr{}, err
	}

	// validate logo size...
	if validLogoSz := validateLogoSize(
		imageLogo.Bounds().Dx(),
		imageLogo.Bounds().Dy(),
	); !validLogoSz {
		return encodeQr{}, errors.New("Please use logo with size 50 X 50")
	}
	if param.QrSize == 0 {
		param.QrSize = 256
	}

	return encodeQr{
		imageLogo: imageLogo,
		value:     param.QrValue,
		size:      param.QrSize,
	}, nil
}

// validateLogoSize
func validateLogoSize(x, y int) bool {
	if x == 50 && y == 50 {
		return true
	}
	return false
}

// Encode QR Code To base64 and byte
func (eq encodeQr) Encode() (ResultEncode, error) {
	byteQR, err := qrcode.New(eq.value, qrcode.High)
	if err != nil {
		return ResultEncode{}, err
	}

	// RGBA Logo and QRCode
	qrCodeImage := byteQR.Image(eq.size)
	qrcodeRGBA := image.NewNRGBA(qrCodeImage.Bounds())

	eq.overlayRGBLogo(qrCodeImage, eq.imageLogo, qrcodeRGBA)

	// Result...
	resPNG := eq.convertToPNG(qrcodeRGBA)
	resBase64 := eq.convertToBase64(resPNG.Bytes())
	return ResultEncode{
		PNG:    resPNG,
		Base64: resBase64,
	}, nil
}

// convertToPNG
func (eq encodeQr) convertToPNG(qrImage *image.NRGBA) bytes.Buffer {
	var buf bytes.Buffer
	// encode to png image
	if err := png.Encode(&buf, qrImage); err != nil {
		return bytes.Buffer{}
	}
	return buf
}

// convert to base64...
func (eq encodeQr) convertToBase64(data []byte) string {
	return base64.RawStdEncoding.EncodeToString(data)
}

// overlayRGBLogo...
func (eq encodeQr) overlayRGBLogo(qrcode, logo image.Image, rgbQR draw.Image) {
	draw.Draw(rgbQR, qrcode.Bounds(), qrcode, image.Point{}, draw.Over)

	// offset logo
	offset := qrcode.Bounds().Max.X/2 - logo.Bounds().Max.X/2
	for x := 0; x < logo.Bounds().Max.X; x++ {
		for y := 0; y < logo.Bounds().Max.Y; y++ {
			r, g, b, a := logo.At(x, y).RGBA()
			// set logo to qr code with RGBA
			rgbQR.Set(
				x+offset,
				y+offset,
				color.RGBA{
					R: uint8(r),
					G: uint8(g),
					B: uint8(b),
					A: uint8(a),
				},
			)
		}
	}
}
