package protocol

import (
  "log"
  "time"
  "fmt"
  "github.com/NeilBetham/elements/crc"
  "github.com/NeilBetham/elements/radios"
)


type Hop struct {
  Freq int
  Dwell time.Duration
  HopIndex int
}

func (h Hop)String() string {
  return fmt.Sprintf(
    "Freq: %d, Hop: %2d, Dwell: %3.2f",
    h.Freq,
    h.HopIndex,
    h.Dwell.Seconds(),
  )
}


type ProtocolHandler struct {
  crc.CRC
  stationID int

  hopTime time.Duration
  hopIndex int
  hopPattern []int
  channels []int

  goodPkts int
  badPkts int
  resync bool

  lastPktReceived time.Time
  lastHop time.Time
}

func NewProtocolHandler(stationNumber int) (ph ProtocolHandler){
  ph.CRC = crc.NewCRC("CCITT-16", 0, 0x1021, 0)
  ph.stationID = stationNumber

  ph.hopTime = time.Duration(2562500 + (stationNumber * 62500)) * time.Microsecond

  ph.hopIndex = 0

  ph.channels = []int{
    901862125, 902364460, 902865026, 903367422, 903868415, 904369408,
    904870462, 905372797, 905873790, 906375698, 906876752, 907378172,
    907879653, 908381134, 908883042, 909384950, 909885516, 910387424,
    910888844, 911389898, 911891806, 912393226, 912894280, 913396188,
    913897608, 914399577, 914900570, 915401563, 915903959, 916405379,
    916905945, 917406938, 917909334, 918410815, 918911808, 919413716,
    919915197, 920416617, 920917610, 921418664, 921920572, 922421565,
    922924388, 923424954, 923926435, 924427428, 924929336, 925431244,
    925932725, 926433718, 926935626,
  }

  ph.hopPattern = []int{
    18, 0, 19, 41, 25, 8, 47, 32, 13, 36, 22, 3, 29, 44, 16, 5, 27, 38,
    10, 49, 21, 2, 30, 42, 14, 48, 7, 24, 34, 45, 1, 17, 39, 26, 9, 31,
    50, 37, 12, 20, 33, 4, 43, 28, 15, 35, 6, 40, 11, 23, 46,
  }

  ph.goodPkts = 0
  ph.badPkts = 0
  ph.resync = true

  ph.lastPktReceived = time.Now()
  ph.lastHop = time.Now()
  return
}

func (ph *ProtocolHandler) HandlePacket(pkt radios.Packet, timedout bool) (hop bool, rd Reading){
  for index, data := range pkt.Data {
    pkt.Data[index] = swapBitOrder(data)
  }

  if ph.Checksum(pkt.Data) != 0 || timedout {
    if !timedout{
      log.Printf("Bad: %s", pkt)
    }
    ph.invalidPkt()
    if !timedout && (time.Now().UnixNano() < (ph.lastPktReceived.Add(ph.hopTime).Add(-10 * time.Millisecond)).UnixNano()) {
      hop = false
      return
    }
  } else if int(pkt.Data[0] & 0x0f) != ph.stationID  {
    log.Printf("Wrong Station: %s", pkt)
    hop = false
    return
  } else {
    log.Printf("%s", pkt)

    ph.validPkt(pkt)
    rd = ParsePacket(pkt)
    rd.Valid = true
    hop = true
    return
  }

  if ph.resync {
    if ph.lastHop.Add(ph.hopTime * time.Duration(len(ph.channels))).Before(time.Now()) {
      ph.lastHop = time.Now()
      hop = true
    } else {
      hop = false
    }
  } else {
    ph.lastHop = time.Now()
    hop = true
  }
  return
}

func (ph *ProtocolHandler) invalidPkt(){
  ph.badPkts++
  if ph.badPkts > 5 && !ph.resync {
    log.Printf("Out of sync with transmitter, resyncing...")
    ph.resync = true
    ph.badPkts = 0
  }
}

func (ph *ProtocolHandler) validPkt(pkt radios.Packet) {
  if ph.resync {
    ph.resync = false
  }

  ph.badPkts = 0
  ph.goodPkts++
  ph.lastPktReceived = time.Now()
}

func (ph *ProtocolHandler) NextHop() (hop Hop){
  hop.Freq = ph.channels[ph.hopPattern[ph.hopIndex]]
  hop.HopIndex = ph.hopIndex
  hop.Dwell = ph.hopTime

  ph.hopIndex++
  if ph.hopIndex > 50 {
    ph.hopIndex = 0
  }
  return
}

func (ph *ProtocolHandler) CurrentChannel() (freq int){
  hopIndex := ph.hopIndex
  if hopIndex == 0 {
    hopIndex = 50
  } else {
    hopIndex--
  }

  return ph.channels[ph.hopPattern[hopIndex]]
}

func swapBitOrder(b byte) byte {
  b = ((b & 0xF0) >> 4) | ((b & 0x0F) << 4)
  b = ((b & 0xCC) >> 2) | ((b & 0x33) << 2)
  b = ((b & 0xAA) >> 1) | ((b & 0x55) << 1)
  return b
}
