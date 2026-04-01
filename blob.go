package lattice

import "path/filepath"

// BlobRef references a binary object in a BlobStore.
type BlobRef struct {
	Key      string
	MimeType string
	Size     int64
	Meta     map[string]string
}

// BlobMeta holds metadata about a blob.
type BlobMeta struct {
	Key      string
	MimeType string
	Size     int64
}

// BlobFromFile creates a BlobRef from a local file path.
func BlobFromFile(path string) BlobRef {
	return BlobRef{
		Key:      path,
		MimeType: mimeFromExt(filepath.Ext(path)),
	}
}

// BlobFromURL creates a BlobRef from a URL.
func BlobFromURL(url string) BlobRef {
	return BlobRef{
		Key: url,
	}
}

func mimeFromExt(ext string) string {
	switch ext {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".mp4":
		return "video/mp4"
	case ".json":
		return "application/json"
	default:
		return "application/octet-stream"
	}
}
