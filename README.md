# yoda1-scale

A small utility to get data from the [Yoda1 Bluetooh Body Scale](https://www.aliexpress.com/item/1005002553230163.html).

## Requirements

- Go 1.15+
- Bluetooth scale linked above

## Build instructions

```shell
mkdir -p ./bin
go build -o ./bin/yoda1-scale ./cmd/yoda1-scale
./bin/yoda1-scale
```

## Running

```plain
./bin/yoda1-scale
# Turn on the scale

INFO[0000] Scanning for devices... please step on your scale 
WARN[0015] unable to parse device: unable to get device name: No such property 'Name' 
WARN[0015] unable to parse device: unable to get device name: No such property 'Name' 
INFO[0015] Device detected! Starting listening for events 
INFO[0031] scaleData={23.12}                            
INFO[0052] scaleData={23.12}                            
INFO[0063] scaleData={20.77}                            
INFO[0091] scaleData={20.77}  
```
