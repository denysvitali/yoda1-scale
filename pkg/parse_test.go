package yoda1

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	input := [][]uint8{
		{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x24, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		{0x2, 0xf3, 0x13, 0x88, 0x0, 0x0, 0x25, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	}

	output := []ScaleData{
		{WeightKG: 0.0},
		{WeightKG: 7.55},
	}
	for i := range input {
		t.Run(fmt.Sprintf("input %d", i), func(t *testing.T) {
			scaleData, err := parseScaleData(input[i])
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, output[i], scaleData)
		})
	}
}
