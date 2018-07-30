package util

import (
	"fmt"
)

func FmtModRewrite(sourcePath string, protocol string, hostname string, destPath string) string {
	return fmt.Sprintf("^/%s(.*) %s://%s/%s$1 [R,L]", sourcePath, protocol, hostname, destPath)
}
