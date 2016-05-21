package trie

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const debugEnabled = false

var logger = log.New(os.Stderr, "\t", log.Lshortfile)

func debugf(format string, args ...interface{}) {
	if debugEnabled {
		msg := strings.Replace(fmt.Sprintf(format, args...), "\n", "\n\t", -1)
		logger.Output(3, msg)
	}
}
