package utils

import "fmt"

func FormatUniqID(tokenID uint64, denomID string) string {
	return fmt.Sprintf("%d@%s", tokenID, denomID)
}
