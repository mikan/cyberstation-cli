package cyberstation

import (
	"io"
	"log"
)

// SafeClose はリソースをクローズし、エラーがあればログに記録します。
func SafeClose(closer io.Closer, tag string) {
	if err := closer.Close(); err != nil {
		log.Printf("failed to close %s: %v", tag, err)
	}
}
