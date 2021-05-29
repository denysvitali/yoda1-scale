package main

import (
	yoda1 "github.com/denysvitali/yoda1-scale/pkg"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"sync"
	"time"
)

var args struct {
	Address string `arg:"-a" help:"BT Mac Address of the scale (for when you're dealing with multiple scales)'"`
}

func main(){
	log := logrus.New()

	log.Infof("Scanning for devices... please step on your scale")

	d, warns, err := yoda1.Discover(8 * time.Second)
	if err != nil {
		log.Fatalf("unable to discover devices: %v", err)
	}

	for _, w := range warns {
		log.Warn(w)
	}

	if len(d) == 0 {
		log.Fatalf("Yoda1 not found")
	}

	var device *yoda1.YodaDevice
	if args.Address != "" {
		for _, yDevice := range d {
			if strings.ToLower(yDevice.MacAddr) == strings.ToLower(args.Address) {
				device = &yDevice
			}
		}

		if device == nil {
			log.Errorf("the device you specified could not be found")
			printDevices(d, log)
			log.Fatal()
		}
	}

	if len(d) > 1 {
		printDevices(d, log)
		log.Fatalf("multiple scales found but no selector specified, please use the -a argument")
	} else if len(d) > 1 {
		printDevices(d, log)
		if device == nil {
			log.Fatalf("unable to find the device with the specified Mac Address")
		}
	}

	device = &d[0]

	log.Infof("Device detected! Starting listening for events")

	c, err := device.WatchEvents()
	if err != nil {
		log.Fatalf("unable to get properties watch channel: %v\n", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for scaleData := range c {
			logrus.Printf("scaleData=%v", scaleData)
		}
		wg.Done()
	}()

	go func() {
		osSigChan := make(chan os.Signal, 1)
		for s := range osSigChan {
			if s == os.Interrupt {
				wg.Done()
			}
		}
	}()
	wg.Wait()
	log.Infof("Bye bye")
}

func printDevices(d []yoda1.YodaDevice, log *logrus.Logger) {
	for _, yDevice := range d {
		log.Printf("%s\t%d", yDevice.MacAddr, yDevice.Rssi)
	}
}