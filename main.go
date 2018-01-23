package main

import (
  "os"
  "time"
  "log"
  "periph.io/x/periph/host"
  "periph.io/x/periph/conn/gpio"
  "periph.io/x/periph/conn/gpio/gpioreg"
  "github.com/NeilBetham/elements/radios"
)

func main() {
  if _, err := host.Init(); err != nil {
		os.Exit(1)
	}

  resetPin := gpioreg.ByName("GPIO4")

  resetPin.Out(gpio.High)
  time.Sleep(5 * time.Millisecond)
  resetPin.Out(gpio.Low)
  time.Sleep(50 * time.Millisecond)

  r, _ := radios.NewRFM69("/dev/spidev0.0")

  r.SetFreq(902355835)
  r.StartRead()

  for{
    time.Sleep(3 * time.Second)
    log.Printf("Dumping Registers...")
    r.DumpRegs()
    r.StartRSSI()
  }

  os.Exit(0)
}
