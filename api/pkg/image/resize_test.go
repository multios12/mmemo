package image

import (
	"bytes"
	stdimage "image"
	"image/color"
	"image/png"
	"testing"
)

func TestNormalizeImageForStorage_ResizesLargePNG(t *testing.T) {
	img := stdimage.NewRGBA(stdimage.Rect(0, 0, maxStoredImageWidth+1, maxStoredImageHeight+1))
	for y := 0; y < maxStoredImageHeight+1; y++ {
		for x := 0; x < maxStoredImageWidth+1; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x % 255), G: uint8(y % 255), B: 120, A: 255})
		}
	}

	var in bytes.Buffer
	if err := png.Encode(&in, img); err != nil {
		t.Fatalf("png encode failed: %v", err)
	}

	contentType, out, err := NormalizeImageForStorage("image/png", in.Bytes())
	if err != nil {
		t.Fatalf("normalizeImageForStorage failed: %v", err)
	}
	if contentType != "image/png" {
		t.Fatalf("contentType = %s, want image/png", contentType)
	}

	cfg, _, err := stdimage.DecodeConfig(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("DecodeConfig failed: %v", err)
	}
	if cfg.Width > maxStoredImageWidth || cfg.Height > maxStoredImageHeight {
		t.Fatalf("image was not resized enough: %dx%d", cfg.Width, cfg.Height)
	}
	if cfg.Width >= maxStoredImageWidth+1 || cfg.Height >= maxStoredImageHeight+1 {
		t.Fatalf("image size did not shrink: %dx%d", cfg.Width, cfg.Height)
	}
	if len(out) >= len(in.Bytes()) {
		t.Fatalf("expected smaller output, got input=%d output=%d", len(in.Bytes()), len(out))
	}
}

func TestNormalizeImageForStorage_KeepsImageWithinResolutionLimit(t *testing.T) {
	img := stdimage.NewRGBA(stdimage.Rect(0, 0, 640, 360))
	for y := 0; y < 360; y++ {
		for x := 0; x < 640; x++ {
			img.Set(x, y, color.RGBA{R: 10, G: 20, B: 30, A: 255})
		}
	}

	var in bytes.Buffer
	if err := png.Encode(&in, img); err != nil {
		t.Fatalf("png encode failed: %v", err)
	}

	contentType, out, err := NormalizeImageForStorage("image/png", in.Bytes())
	if err != nil {
		t.Fatalf("normalizeImageForStorage failed: %v", err)
	}
	if contentType != "image/png" {
		t.Fatalf("contentType = %s, want image/png", contentType)
	}
	if !bytes.Equal(out, in.Bytes()) {
		t.Fatal("expected image within resolution limit to remain unchanged")
	}
}
