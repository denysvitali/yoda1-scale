package yoda1

import (
	"fmt"
	"testing"
	"time"
)

func TestDiscover(t *testing.T){
	devices, warnings, err := Discover(5 * time.Second)
	if err != nil {
		t.Fatal(err)
	}

	if len(warnings) > 0 {
		t.Fatal(warnings)
	}

	fmt.Printf("devices=%v", devices)
}
