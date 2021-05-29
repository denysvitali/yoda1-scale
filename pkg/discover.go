package yoda1

import (
	"fmt"
	"github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	"sync"
	"time"
)

func Discover(timeout time.Duration) (yodaDevices []YodaDevice, warnings []error, err error) {
	a, err := api.GetDefaultAdapter()
	if err != nil {
		return yodaDevices, warnings, fmt.Errorf("unable to get default bluetooth a: %v", err)
	}

	err = a.StartDiscovery()
	if err != nil {
		return yodaDevices, warnings, fmt.Errorf("unable to start discovery: %v", err)
	}

	_ = a.FlushDevices()

	c, cancel, err := a.OnDeviceDiscovered()
	if err != nil {
		return yodaDevices, warnings, fmt.Errorf("unable to listen for device discovery: %v", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for ev := range c {
			if ev.Type == adapter.DeviceRemoved {
				continue
			}

			dev, err := device.NewDevice1(ev.Path)
			if err != nil {
				warnings = append(warnings, fmt.Errorf("unable to create device: %v", err))
			}

			yodaDevice, err := parseDevice(dev)
			if err != nil {
				warnings = append(warnings, fmt.Errorf("unable to parse device: %v", err))
			}

			if !yodaDevice.IsValid() {
				continue
			}

			yodaDevices = append(yodaDevices, yodaDevice)
		}
		wg.Done()
	}()

	time.Sleep(timeout)
	cancel()
	wg.Wait()
	_ = a.StopDiscovery()
	return yodaDevices, warnings, nil
}

func parseDevice(d *device.Device1) (device YodaDevice, err error) {
	name, err := d.GetName()
	if err != nil {
		return device, fmt.Errorf("unable to get device name: %v", err)
	}
	if name != DeviceBtName {
		// Not our device
		return YodaDevice{isYoda: false}, nil
	}
	// Found our Yoda1 !
	dbusVariant, err := getVariantFromMfData(d.Properties.ManufacturerData)
	if err != nil {
		return device, err
	}

	scaleData, err := getScaleData(dbusVariant)
	if err != nil {
		return device, err
	}

	return YodaDevice{
		MacAddr: d.Properties.Address,
		Rssi:    d.Properties.RSSI,
		Data:    scaleData,
		dev:     d,
		isYoda: true,
	}, nil
}

func getVariantFromMfData(mfData map[uint16]interface{}) (dbus.Variant, error) {
	// Get first value
	var dbusVariant dbus.Variant
	var ok bool

	for _, v := range mfData {
		dbusVariant, ok = v.(dbus.Variant)
		if !ok {
			return dbus.Variant{}, fmt.Errorf("unable to convert value to dbus.Variant")
		}
		break
	}
	return dbusVariant, nil
}

func getScaleData(v dbus.Variant) (ScaleData, error) {
	var dataBytes []uint8
	// Company ID is always random, let's pick the first key of the map
	var ok bool
	data := v.Value()
	dataBytes, ok = data.([]uint8)
	if !ok {
		return ScaleData{}, fmt.Errorf("unable to cast value to []uint8")
	}

	scaleData, err := parseScaleData(dataBytes)
	if err != nil {
		return ScaleData{}, fmt.Errorf("unable to parse scale data: %v", err)
	}
	return scaleData, nil
}
