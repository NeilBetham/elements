package radios

import (
  "time"
  "fmt"
  "log"
  "reflect"
  "periph.io/x/periph/conn/spi"
  "periph.io/x/periph/conn/spi/spireg"
)

const(
  OpModeReg = 0x01
  DataModulationReg = 0x02
  bitRateMsbReg = 0x03
  bitRateLsbReg = 0x04
  freqDevMsbReg = 0x05
  freqDevLsbReg = 0x06
  carrierFreqMsbReg = 0x07
  carrierFreqMidReg = 0x08
  carrierFreqLsbReg = 0x09
  rxBwFiltContReg = 0x019
  afcBwFiltContReg = 0x1a
  afcFeiContStatReg = 0x1e
  rssiConfReg = 0x23
  dioMappingReg = 0x25
  irqFlags2Reg = 0x28
  rssiThreshReg = 0x29
  preambleMsbReg = 0x2c
  preambleLsbReg = 0x2d
  syncConfigReg = 0x2e
  syncVal1Reg = 0x2f
  syncval2Reg = 0x30
  packetConfig1Reg = 0x37
  payloadLengthReg = 0x38
  packetConfig2Reg = 0x3d
  testDagcReg = 0x6f
)


type rfm69Regs struct {
  opMode uint8
  dataModulation uint8
  bitRateMsb uint8
  bitRateLsb uint8
  freqDevMsb uint8
  freqDevLsb uint8
  carrierFreqMsb uint8
  carrierFreqMid uint8
  carrierFreqLsb uint8
  rxBwFiltCont uint8
  afcBwFiltCont uint8
  afcFeiContStat uint8
  rssiConf uint8
  dioMapping uint8
  irqFlags2 uint8
  rssiThresh uint8
  preambleMsb uint8
  preambleLsb uint8
  syncConfig uint8
  syncVal1 uint8
  syncval2 uint8
  packetConfig1 uint8
  payloadLength uint8
  packetConfig2 uint8
  testDagc uint8
}


func newRFM69Regs() (r rfm69Regs){
  r.opMode = (1 << 2) // Standby mode default
  r.dataModulation = (2 << 0) // Gaussian filter, BT= 0.5

  // 19200 bit rate
  r.bitRateMsb = 0x06
  r.bitRateLsb = 0x83

  // 4.8 kHz RX freq deviation
  r.freqDevMsb = 0x00
  r.freqDevLsb = 0x4e

  // RX Bandwidth
  // DCC Freq = 2, RX BW Mant = 20, RX BW Mant = 4
  r.rxBwFiltCont = (2 << 5) | (1 << 3) | (4 << 0)

  // AFC BW
  // DCC Freq = 2, RX BW Mant = 20, RX BW Mant = 3
  r.afcBwFiltCont = (2 << 5) | (1 << 3) | (3 << 0)

  // AFC/FEI Settings
  // AutoClear on, AutoOn On
  r.afcFeiContStat = (1 << 3) | (1 << 2)

  // DIO Mapping
  // D0 RX Payload Ready interrupt
  r.dioMapping = (1 << 6)

  // IRQ Flags
  // Clears fifo overrun flag
  r.irqFlags2 = (1 << 4)

  // RSSI Threshold
  r.rssiThresh = 170

  // Preambles
  r.preambleMsb = 0
  r.preambleLsb = 4

  // Sync Config
  // 2 sync bytes, allow two errors in sync bytes
  r.syncConfig = (2 << 3) | (2 << 0)

  // Sync Values
  r.syncVal1 = 0xcb
  r.syncval2 = 0x89

  // Packet Config
  // 2 Bit Inter Packet Delay and AutoRxRestartOn On
  r.packetConfig1 = 0
  r.packetConfig2 = (2 << 4) | (1 << 1)

  r.payloadLength = 10

  // Test DAGC Improved No Low Beta On
  r.testDagc = 0x30
  return
}

// RFM69 Handles communication and state for the RFM69 wireless radio
type RFM69 struct {
  port spi.Port
  freq  int32
  conn spi.Conn
  config rfm69Regs

  recvBytes []byte
}

// NewRFM69 sets up a new RFM69 class
func NewRFM69(mosi int, miso int, sclk int, cs int, port string) (r RFM69, err error) {
  p, openErr := spireg.Open(port)
  if openErr {
    log.Printf("error: open spi port failed")
    err = openErr
    return
  }

  r.port = p
  r.freq = 915000000
  r.recv_bytes = make([]byte)
  r.config = newRFM69Regs()

  r.init()
  return
}

func (r *RFM69) init() (err error){
  // Setup the SPI connection with 5MHz baud, CPOL=0, CPHA=0, and 8 bit bytes
  r.Connect(5000000, spi.Mode0, 8)

}

func (r rfm69Regs) syncRegs(regs RFM69Regs) (err error) {
  for i := 0; i < r.Type().NumField(); i++ {

  }
}

func (r *RFM69) writeReg(addr, val uint8) (err error){

}

func (r *RFM69) readReg(addr uint8) (b byte, err error){

}
