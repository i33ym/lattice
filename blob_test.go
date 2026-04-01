package lattice

import (
	"testing"
)

func TestBlobFromFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		path         string
		wantKey      string
		wantMimeType string
	}{
		{"BlobFromFile/pngFile/setsKeyAndMime", "/data/photo.png", "/data/photo.png", "image/png"},
		{"BlobFromFile/jpgFile/setsImageJpeg", "/data/photo.jpg", "/data/photo.jpg", "image/jpeg"},
		{"BlobFromFile/mp3File/setsAudioMpeg", "/data/song.mp3", "/data/song.mp3", "audio/mpeg"},
		{"BlobFromFile/mp4File/setsVideoMp4", "/data/clip.mp4", "/data/clip.mp4", "video/mp4"},
		{"BlobFromFile/unknownExt/setsOctetStream", "/data/file.xyz", "/data/file.xyz", "application/octet-stream"},
		{"BlobFromFile/wavFile/setsAudioWav", "/data/sound.wav", "/data/sound.wav", "audio/wav"},
		{"BlobFromFile/jsonFile/setsApplicationJson", "/data/config.json", "/data/config.json", "application/json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ref := BlobFromFile(tt.path)
			if ref.Key != tt.wantKey {
				t.Fatalf("BlobFromFile(%q).Key = %q, want %q", tt.path, ref.Key, tt.wantKey)
			}
			if ref.MimeType != tt.wantMimeType {
				t.Fatalf("BlobFromFile(%q).MimeType = %q, want %q", tt.path, ref.MimeType, tt.wantMimeType)
			}
		})
	}
}

func TestBlobFromURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		url     string
		wantKey string
	}{
		{"BlobFromURL/httpURL/setsKey", "https://example.com/file.mp3", "https://example.com/file.mp3"},
		{"BlobFromURL/s3URL/setsKey", "s3://bucket/key", "s3://bucket/key"},
		{"BlobFromURL/emptyURL/setsEmptyKey", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ref := BlobFromURL(tt.url)
			if ref.Key != tt.wantKey {
				t.Fatalf("BlobFromURL(%q).Key = %q, want %q", tt.url, ref.Key, tt.wantKey)
			}
			if ref.MimeType != "" {
				t.Fatalf("BlobFromURL(%q).MimeType = %q, want empty", tt.url, ref.MimeType)
			}
		})
	}
}
