package cyberstation

import (
	"io"
	"log"
)

// SafeClose はリソースをクローズし、エラーがあればログに記録します。
func SafeClose(closer io.Closer, tag string) {
	if err := closer.Close(); err != nil {
		log.Printf("エラー: %s をクローズできませんでした: %v", tag, err)
	}
}
