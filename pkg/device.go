package yoda1

import (
	"fmt"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	"github.com/sirupsen/logrus"
)

const DeviceBtName = "Yoda1"

type YodaDevice struct {
	MacAddr string
	Rssi    int16
	Data    ScaleData
	dev     *device.Device1

	isYoda bool
}

func (y YodaDevice) IsValid() bool {
	return y.isYoda
}

func (y YodaDevice) WatchEvents() (<- chan *ScaleData, error) {
	outChannel := make(chan *ScaleData, 1)
	inChannel, err := y.dev.WatchProperties()
	if err != nil {
		return outChannel, fmt.Errorf("unable to get properties watch channel: %v\n", err)
	}

	go func() {
		for v := range inChannel {
			if v.Name != "ManufacturerData" {
				// We are only interested in Manufacturer Data changes
				continue
			}

			properties, err := y.dev.GetProperties()
			if err != nil {
				logrus.Errorf("unable to get device properties: %v", err)
				continue
			}

			dbusVariant, err := getVariantFromMfData(properties.ManufacturerData)
			if err != nil {
				logrus.Errorf("unable to get variant: %v", err)
				continue
			}

			scaleData, err := getScaleData(dbusVariant)
			if err != nil {
				logrus.Errorf("unable to get scale data: %v", err)
				continue
			}

			outChannel <- &scaleData
		}
	}()
	return outChannel, nil
}