package yoda1

import (
	"fmt"
	"testing"
)

func TestDiscover(t *testing.T){
	devices, err := Discover()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("devices=%v", devices)
}
