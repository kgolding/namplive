package network

import (
	"testing"
)

func Test_GetInterfaces(t *testing.T) {
	t.Log(GetIPv4NonLocalInterfaces())
}
