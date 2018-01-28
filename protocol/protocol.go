package protocol

import (
  "log"
  "time"
  "github.com/NeilBetham/elements/crc"
)

// ProtocolHandler handles receipt of packets and channel hopping
type ProtocolHandler struct {
  crc.CRC

  hopTime time.Duration
  hopIndex int
  hopPattern []int
  channels []int

  goodPkts int
  badPkts int
  resync bool
}

// NewProtocolHandler sets up a new protocol handler
func NewProtocolHandler(stationNumber int) (ph ProtocolHandler){
  ph.CRC = crc.NewCRC("CCITT-16", 0, 0x1021, 0)

  ph.hopTime = time.Duration(2562500 + (stationNumber * 62500)) * time.Microsecond

  ph.hopIndex = 0

  ph.channels = []int{
    901862125, 902364460, 902865026, 903367422, 903868415, 904369408, 904870462,
    905372797, 905873790, 906375698, 906876752, 907378172, 907879653, 908381134,
    908883042, 909384950, 909885516, 910387424, 910888844, 911389898, 911891806,
    912393226, 912894280, 913396188, 913897608, 914399577, 914900570, 915401563,
    915903959, 916405379, 916905945, 917406938, 917909334, 918410815, 918911808,
    919413716, 919915197, 920416617, 920917610, 921418664, 921920572, 922421565,
    922924388, 923424954, 923926435, 924427428, 924929336, 925431244, 925932725,
    926433718, 926935626,
  }

  ph.hopPattern = []int{
    1, 30, 21, 11, 41, 15, 46, 26, 5, 34, 18, 48, 38, 8, 24, 44, 14, 31, 0, 2, 39,
    20, 10, 49, 27, 4, 33, 16, 43, 12, 22, 35, 7, 40, 28, 45, 9, 37, 17, 32, 47, 3,
    23, 42, 13, 29, 50, 6, 25, 19, 36,
  }

  ph.goodPkts = 0
  ph.badPkts = 0
  ph.resync = true
  return
}

// HandlePacket handles incomming packets and decides if a hop should happen
func (ph *ProtocolHandler) HandlePacket(pkt []byte) (hop bool){
  // Davis ISS transmits LSB first
  for index, data := range pkt {
    pkt[index] = swapBitOrder(data)
  }

  // If the checksum is valid then hop
  if len(pkt) == 0 || ph.Checksum(pkt) != 0 {
    log.Printf("Bad packet")
    if !ph.resync {
      ph.badPkts++
    }
  } else {
    if ph.resync {
      ph.resync = false
    }
    ph.goodPkts++
  }

  // If we start to accumulate bad packets it's time for a resync
  if ph.badPkts >= 5 {
    log.Printf("Resync needed")
    ph.resync = true
    ph.badPkts = 0
  }

  if ph.resync {
    hop = false
  } else {
    hop = true
  }
  return
}

// NextHop gets the next hopping frequency and increments the hop index
func (ph *ProtocolHandler) NextHop() (freq int){
  freq = ph.channels[ph.hopPattern[ph.hopIndex]]
  ph.hopIndex++
  if ph.hopIndex > 50 {
    ph.hopIndex = 0
  }
  return
}

func swapBitOrder(b byte) byte {
  b = ((b & 0xF0) >> 4) | ((b & 0x0F) << 4)
  b = ((b & 0xCC) >> 2) | ((b & 0x33) << 2)
  b = ((b & 0xAA) >> 1) | ((b & 0x55) << 1)
  return b
}
