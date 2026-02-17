package generator

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"strings"
	"unicode/utf8"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/ean"
	"github.com/innovelabs/microtools-go/internal/models"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
	BarcodeTypeUPCA    = "UPC-A"
	BarcodeTypeEAN13   = "EAN-13"
	BarcodeTypeCode128 = "Code128"

	BarcodeFormatPNG = "png"
	BarcodeFormatSVG = "svg"

	defaultBarcodeWidth  = 300
	defaultBarcodeHeight = 150
	maxBarcodeWidth      = 1024
	maxBarcodeHeight     = 1024
	minBarcodeWidth      = 50
	minBarcodeHeight     = 50

	textPaddingHeight = 20
	maxCode128Length  = 500
)

var (
	ErrInvalidType      = errors.New("invalid barcode type: must be UPC-A, EAN-13, or Code128")
	ErrInvalidFormat    = errors.New("invalid format: must be png or svg")
	ErrInvalidData      = errors.New("invalid data for the specified barcode type")
	ErrChecksumMismatch = errors.New("checksum digit does not match computed value")
)

// BarcodeService defines barcode generation interface
type BarcodeService interface {
	Generate(req models.GenerateRequest) ([]byte, string, error)
}

type defaultBarcodeService struct{}

// NewDefaultBarcodeService creates a new barcode service
func NewDefaultBarcodeService() BarcodeService {
	return &defaultBarcodeService{}
}

// Generate generates a barcode image
func (s *defaultBarcodeService) Generate(req models.GenerateRequest) ([]byte, string, error) {
	applyBarcodeDefaults(&req)

	if err := validateBarcodeRequest(req); err != nil {
		return nil, "", err
	}

	bc, err := encodeBarcode(req.Type, req.Data)
	if err != nil {
		return nil, "", err
	}

	switch req.Format {
	case BarcodeFormatPNG:
		data, err := renderBarcodePNG(bc, req.Width, req.Height, req.IncludeText, req.Data)
		return data, "image/png", err
	case BarcodeFormatSVG:
		data, err := renderBarcodeSVG(bc, req.Width, req.Height, req.IncludeText, req.Data)
		return data, "image/svg+xml", err
	default:
		return nil, "", ErrInvalidFormat
	}
}

func applyBarcodeDefaults(req *models.GenerateRequest) {
	if req.Width == 0 {
		req.Width = defaultBarcodeWidth
	}
	if req.Height == 0 {
		req.Height = defaultBarcodeHeight
	}
}

func validateBarcodeRequest(req models.GenerateRequest) error {
	switch req.Type {
	case BarcodeTypeUPCA, BarcodeTypeEAN13, BarcodeTypeCode128:
	default:
		return ErrInvalidType
	}

	switch req.Format {
	case BarcodeFormatPNG, BarcodeFormatSVG:
	default:
		return ErrInvalidFormat
	}

	if req.Data == "" {
		return fmt.Errorf("%w: data is required", ErrInvalidData)
	}

	if err := validateBarcodeData(req.Type, req.Data); err != nil {
		return err
	}

	if req.Width < minBarcodeWidth || req.Width > maxBarcodeWidth {
		return fmt.Errorf("%w: width must be between %d and %d", ErrInvalidData, minBarcodeWidth, maxBarcodeWidth)
	}
	if req.Height < minBarcodeHeight || req.Height > maxBarcodeHeight {
		return fmt.Errorf("%w: height must be between %d and %d", ErrInvalidData, minBarcodeHeight, maxBarcodeHeight)
	}

	return nil
}

func validateBarcodeData(barcodeType, data string) error {
	switch barcodeType {
	case BarcodeTypeUPCA:
		if !isNumeric(data) {
			return fmt.Errorf("%w: UPC-A data must be numeric", ErrInvalidData)
		}
		n := len(data)
		if n != 11 && n != 12 {
			return fmt.Errorf("%w: UPC-A data must be 11 or 12 digits", ErrInvalidData)
		}
		if n == 12 {
			return validateUPCAChecksum(data)
		}

	case BarcodeTypeEAN13:
		if !isNumeric(data) {
			return fmt.Errorf("%w: EAN-13 data must be numeric", ErrInvalidData)
		}
		n := len(data)
		if n != 12 && n != 13 {
			return fmt.Errorf("%w: EAN-13 data must be 12 or 13 digits", ErrInvalidData)
		}
		if n == 13 {
			return validateEAN13Checksum(data)
		}

	case BarcodeTypeCode128:
		if !utf8.ValidString(data) {
			return fmt.Errorf("%w: Code128 data must be valid UTF-8", ErrInvalidData)
		}
		if len(data) > maxCode128Length {
			return fmt.Errorf("%w: Code128 data exceeds maximum length of %d characters", ErrInvalidData, maxCode128Length)
		}
	}
	return nil
}

func isNumeric(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func computeUPCAChecksum(digits string) int {
	sum := 0
	for i := 0; i < 11; i++ {
		d := int(digits[i] - '0')
		if i%2 == 0 {
			sum += d * 3
		} else {
			sum += d * 1
		}
	}
	return (10 - (sum % 10)) % 10
}

func validateUPCAChecksum(data string) error {
	expected := computeUPCAChecksum(data[:11])
	actual := int(data[11] - '0')
	if expected != actual {
		return fmt.Errorf("%w: expected check digit %d, got %d", ErrChecksumMismatch, expected, actual)
	}
	return nil
}

func computeEAN13Checksum(digits string) int {
	sum := 0
	for i := 0; i < 12; i++ {
		d := int(digits[i] - '0')
		if i%2 == 0 {
			sum += d * 1
		} else {
			sum += d * 3
		}
	}
	return (10 - (sum % 10)) % 10
}

func validateEAN13Checksum(data string) error {
	expected := computeEAN13Checksum(data[:12])
	actual := int(data[12] - '0')
	if expected != actual {
		return fmt.Errorf("%w: expected check digit %d, got %d", ErrChecksumMismatch, expected, actual)
	}
	return nil
}

func encodeBarcode(barcodeType, data string) (barcode.Barcode, error) {
	switch barcodeType {
	case BarcodeTypeUPCA:
		eanData := "0" + data
		bc, err := ean.Encode(eanData)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidData, err)
		}
		return bc, nil

	case BarcodeTypeEAN13:
		bc, err := ean.Encode(data)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidData, err)
		}
		return bc, nil

	case BarcodeTypeCode128:
		bc, err := code128.Encode(data)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidData, err)
		}
		return bc, nil

	default:
		return nil, ErrInvalidType
	}
}

func renderBarcodePNG(bc barcode.Barcode, width, height int, includeText bool, text string) ([]byte, error) {
	barcodeHeight := height
	totalHeight := height
	if includeText {
		totalHeight = height + textPaddingHeight
	}

	scaled, err := barcode.Scale(bc, width, barcodeHeight)
	if err != nil {
		return nil, fmt.Errorf("failed to scale barcode: %w", err)
	}

	canvas := image.NewRGBA(image.Rect(0, 0, width, totalHeight))
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	draw.Draw(canvas, image.Rect(0, 0, width, barcodeHeight), scaled, scaled.Bounds().Min, draw.Over)

	if includeText {
		drawBarcodeTextCentered(canvas, text, height+textPaddingHeight-4, width)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, canvas); err != nil {
		return nil, fmt.Errorf("failed to encode PNG: %w", err)
	}
	return buf.Bytes(), nil
}

func drawBarcodeTextCentered(img *image.RGBA, text string, y int, canvasWidth int) {
	face := basicfont.Face7x13
	textWidth := len(text) * 7
	x := (canvasWidth - textWidth) / 2
	if x < 0 {
		x = 0
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: face,
		Dot: fixed.Point26_6{
			X: fixed.I(x),
			Y: fixed.I(y),
		},
	}
	d.DrawString(text)
}

func renderBarcodeSVG(bc barcode.Barcode, width, height int, includeText bool, text string) ([]byte, error) {
	bounds := bc.Bounds()
	bcWidth := bounds.Max.X - bounds.Min.X

	totalHeight := height
	barcodeHeight := height
	if includeText {
		totalHeight = height + textPaddingHeight
	}

	scaleX := float64(width) / float64(bcWidth)

	var buf bytes.Buffer
	fmt.Fprintf(&buf, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, width, totalHeight, width, totalHeight)
	buf.WriteByte('\n')
	fmt.Fprintf(&buf, `<rect width="%d" height="%d" fill="white"/>`, width, totalHeight)
	buf.WriteByte('\n')

	y := bounds.Min.Y
	startX := -1
	for x := bounds.Min.X; x <= bounds.Max.X; x++ {
		isBar := false
		if x < bounds.Max.X {
			r, _, _, _ := bc.At(x, y).RGBA()
			isBar = (r == 0)
		}
		if isBar && startX == -1 {
			startX = x - bounds.Min.X
		} else if !isBar && startX != -1 {
			svgX := float64(startX) * scaleX
			svgW := float64(x-bounds.Min.X-startX) * scaleX
			fmt.Fprintf(&buf, `<rect x="%.2f" y="0" width="%.2f" height="%d" fill="black"/>`, svgX, svgW, barcodeHeight)
			buf.WriteByte('\n')
			startX = -1
		}
	}

	if includeText {
		textY := barcodeHeight + textPaddingHeight - 2
		fmt.Fprintf(&buf, `<text x="%d" y="%d" text-anchor="middle" font-family="monospace" font-size="12" fill="black">%s</text>`,
			width/2, textY, barcodeSVGEscape(text))
		buf.WriteByte('\n')
	}

	buf.WriteString(`</svg>`)
	return buf.Bytes(), nil
}

func barcodeSVGEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

