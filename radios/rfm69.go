package radios

import (
  "log"
  "reflect"
  "sort"
  "periph.io/x/periph/conn/spi"
  "periph.io/x/periph/conn/spi/spireg"
)

var regAddrs = map[string]uint8{
  "opMode": 0x01,
  "dataModulation": 0x02,
  "bitRateMsb": 0x03,
  "bitRateLsb": 0x04,
  "freqDevMsb": 0x05,
  "freqDevLsb": 0x06,
  "carrierFreqMsb": 0x07,
  "carrierFreqMid": 0x08,
  "carrierFreqLsb": 0x09,
  "rxBwFiltCont": 0x19,
  "afcBwFiltCont": 0x1a,
  "afcFeiContStat": 0x1e,
  "rssiConf": 0x23,
  "rssiValue": 0x24,
  "dioMapping": 0x25,
  "irqFlags1": 0x27,
  "irqFlags2": 0x28,
  "rssiThresh": 0x29,
  "preambleMsb": 0x2c,
  "preambleLsb": 0x2d,
  "syncConfig": 0x2e,
  "syncVal1": 0x2f,
  "syncval2": 0x30,
  "packetConfig1": 0x37,
  "payloadLength": 0x38,
  "packetConfig2": 0x3d,
  "testDagc": 0x6f,
}


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

  // 9.9 kHz RX freq deviation
  r.freqDevMsb = 0x00
  r.freqDevLsb = 0xa1

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
  r.rssiThresh = 228

  // Preambles
  r.preambleMsb = 0
  r.preambleLsb = 4

  // Sync Config
  // 2 sync bytes, allow two errors in sync bytes
  r.syncConfig = (1 << 7) | (1 << 3) | (2 << 0)

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
func NewRFM69(port string) (r RFM69, err error) {
  p, openErr := spireg.Open(port)
  if openErr != nil {
    log.Printf("error: open spi port failed")
    err = openErr
    return
  }

  r.port = p
  r.freq = 915000000
  r.recvBytes = make([]byte, 1)
  r.config = newRFM69Regs()

  err = r.init()
  return
}

// SetFreq Sets the carrier freq of the RFM69
func (r *RFM69) SetFreq(freq uint32) (err error){
  r.config.carrierFreqLsb = uint8(freq / 61)
  r.config.carrierFreqMid = uint8((freq / 61) >> 8)
  r.config.carrierFreqMsb = uint8((freq / 61) >> 16)
  err = r.syncRegs()
  return
}

// StartRead puts the receiver in the RX mode and enables the interrupt
func (r *RFM69) StartRead() (err error){
  r.config.opMode = 0x04 // Put the chip into RX mode
  r.writeReg(regAddrs["opMode"], r.config.opMode)
  return
}

// DumpRegs dumps the current register settings in the RFM69
func (r *RFM69) DumpRegs() (err error){
  reverseMap := make(map[uint8]string)
  for name, addr := range regAddrs {
    reverseMap[addr] = name
  }

  var keys []int
  for addr, _ := range reverseMap {
     keys = append(keys, int(addr))
  }
  sort.Ints(keys)

  for _, addr := range keys {
    val, _ := r.readReg(uint8(addr))
    log.Printf("%30s | 0x%.2x | %.8b\n", reverseMap[uint8(addr)], val, val)
  }
  return
}

func (r *RFM69) init() (err error){
  // Setup the SPI connection with 5MHz baud, CPOL=0, CPHA=0, and 8 bit bytes
  conn, err := r.port.Connect(5000000, spi.Mode0, 8)
  if err != nil {
    return err
  }
  r.conn = conn

  r.syncRegs()
  return
}

func (r *RFM69) syncRegs() (err error) {
  regValRef := reflect.ValueOf(r.config)
  regTypeRef := reflect.TypeOf(r.config)
  for i := 0; i < regValRef.NumField(); i++ {
    regValue := uint8(regValRef.Field(i).Uint())
    regAddr := regAddrs[regTypeRef.Field(i).Name]
    err = r.writeReg(regAddr, regValue)
    if err != nil {
      break
    }
  }
  return
}

func (r *RFM69) writeReg(addr, val uint8) (err error){
  bytesToSend := []byte{addr | 0x80, val}
  bytesReceived := make([]byte, len(bytesToSend))
  if err = r.conn.Tx(bytesToSend, bytesReceived); err != nil{
    return err
  }
  return
}

func (r *RFM69) readReg(addr uint8) (b byte, err error){
  bytesToSend := []byte{addr & 0x7f, 0x00}
  bytesReceived := make([]byte, len(bytesToSend))
  if err = r.conn.Tx(bytesToSend, bytesReceived); err != nil{
    return 0x00, err
  }
  return bytesReceived[1], nil
}

// StartRSSI starts an RSSI reading
func (r *RFM69) StartRSSI() (err error){
  err = r.writeReg(regAddrs["rssiConf"], 0x01)
  return
}
