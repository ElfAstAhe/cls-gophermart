package utils

import (
	"fmt"
	"io"
)

func CloseOnly(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		fmt.Printf("Error closing io.Closer: [%s]\n", err)
	}
}
