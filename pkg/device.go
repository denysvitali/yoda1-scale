package yoda1

const DeviceBtName = "Yoda1"

type YodaDevice struct {
	MacAddr string
	Rssi    int16
	Data    ScaleData
}
