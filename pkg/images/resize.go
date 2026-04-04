package images

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"math"
	"net/http"
)

const (
	maxStoredImageWidth  = 1920
	maxStoredImageHeight = 1080
	jpegQuality          = 85
)

func normalizeImageForStorage(contentType string, data []byte) (string, []byte, error) {
	if len(data) == 0 {
		return contentType, data, nil
	}

	detectedType := http.DetectContentType(data)
	if detectedType != "" {
		contentType = detectedType
	}

	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return contentType, data, nil
	}
	if !needsResize(cfg.Width, cfg.Height) {
		return normalizeContentType(contentType, format), data, nil
	}
	if format != "jpeg" && format != "png" {
		return normalizeContentType(contentType, format), data, nil
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return "", nil, fmt.Errorf("画像を読み込めません: %w", err)
	}

	bounds := img.Bounds()
	nextWidth, nextHeight := scaledSize(bounds.Dx(), bounds.Dy())
	if nextWidth >= bounds.Dx() && nextHeight >= bounds.Dy() {
		return normalizeContentType(contentType, format), data, nil
	}

	resized := resizeImage(img, nextWidth, nextHeight)
	out, err := encodeImage(format, resized)
	if err != nil {
		return "", nil, err
	}

	return normalizeContentType(contentType, format), out, nil
}

func needsResize(width int, height int) bool {
	return width >= maxStoredImageWidth || height >= maxStoredImageHeight
}

func scaledSize(width int, height int) (int, int) {
	if !needsResize(width, height) {
		return width, height
	}

	scaleW := float64(maxStoredImageWidth) / float64(width)
	scaleH := float64(maxStoredImageHeight) / float64(height)
	scale := math.Min(scaleW, scaleH)
	if scale >= 1 {
		scale = math.Nextafter(1, 0)
	}

	nextWidth := max(1, int(math.Floor(float64(width)*scale)))
	nextHeight := max(1, int(math.Floor(float64(height)*scale)))
	if nextWidth >= width && width > 1 {
		nextWidth = width - 1
	}
	if nextHeight >= height && height > 1 {
		nextHeight = height - 1
	}
	return nextWidth, nextHeight
}

func resizeImage(src image.Image, width int, height int) image.Image {
	srcBounds := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		srcY := srcBounds.Min.Y + (y*srcBounds.Dy())/height
		for x := 0; x < width; x++ {
			srcX := srcBounds.Min.X + (x*srcBounds.Dx())/width
			dst.Set(x, y, src.At(srcX, srcY))
		}
	}

	return dst
}

func encodeImage(format string, img image.Image) ([]byte, error) {
	var out bytes.Buffer
	switch format {
	case "jpeg":
		if err := jpeg.Encode(&out, img, &jpeg.Options{Quality: jpegQuality}); err != nil {
			return nil, fmt.Errorf("JPEG画像を保存できません: %w", err)
		}
	case "png":
		if err := png.Encode(&out, img); err != nil {
			return nil, fmt.Errorf("PNG画像を保存できません: %w", err)
		}
	default:
		return nil, fmt.Errorf("未対応の画像形式です: %s", format)
	}
	return out.Bytes(), nil
}

func normalizeContentType(contentType string, format string) string {
	switch format {
	case "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	default:
		return contentType
	}
}
