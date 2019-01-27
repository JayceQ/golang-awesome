package utils

import (
	"fmt"
	"testing"
)

func TestNewUUID(t *testing.T) {
	s, _ := NewUUID()
	fmt.Printf("%s", s)
}
