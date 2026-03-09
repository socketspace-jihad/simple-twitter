package storage

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"golang.org/x/image/draw"
)

const MaxImageSize = 2 * 1024 * 1024 // 2MB

// CompressImage reads an image from r, and compresses it to fit within MaxImageSize.
// Returns the compressed bytes and the content type.
func CompressImage(r io.Reader, contentType string) ([]byte, string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read image: %w", err)
	}

	// If already under 2MB, return as-is
	if len(data) <= MaxImageSize {
		return data, contentType, nil
	}

	// Decode the image
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Iteratively reduce quality/size until under 2MB
	// First try JPEG compression at decreasing quality
	for quality := 85; quality >= 10; quality -= 10 {
		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality}); err != nil {
			return nil, "", fmt.Errorf("failed to encode JPEG: %w", err)
		}
		if buf.Len() <= MaxImageSize {
			return buf.Bytes(), "image/jpeg", nil
		}
	}

	// If still too large, resize the image progressively
	bounds := img.Bounds()
	for scale := 0.75; scale >= 0.1; scale -= 0.1 {
		newWidth := int(float64(bounds.Dx()) * scale)
		newHeight := int(float64(bounds.Dy()) * scale)

		resized := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
		draw.BiLinear.Scale(resized, resized.Bounds(), img, bounds, draw.Over, nil)

		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, resized, &jpeg.Options{Quality: 75}); err != nil {
			return nil, "", fmt.Errorf("failed to encode resized JPEG: %w", err)
		}
		if buf.Len() <= MaxImageSize {
			return buf.Bytes(), "image/jpeg", nil
		}
	}

	return nil, "", fmt.Errorf("image too large to compress under 2MB")
}

// IsImageContentType checks if the given content type is a supported image type.
func IsImageContentType(ct string) bool {
	switch ct {
	case "image/jpeg", "image/png", "image/gif", "image/webp":
		return true
	}
	return false
}

func init() {
	// Register PNG decoder (JPEG is registered by default via image/jpeg import)
	image.RegisterFormat("png", "\x89PNG\r\n\x1a\n", png.Decode, png.DecodeConfig)
}
