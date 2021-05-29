package yoda1

import "encoding/binary"

type ScaleData struct {
	WeightKG float32
}

func parseScaleData(data []uint8) (ScaleData, error){
	weight := float32(binary.BigEndian.Uint16(data[0:2]))/100 // 55.43 is sent as [0x15, 0xA7] (=5543)

	return ScaleData{
		weight,
	}, nil
}
