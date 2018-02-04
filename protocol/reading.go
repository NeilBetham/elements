package protocol

import (
  "log"
  "fmt"
  "github.com/NeilBetham/elements/radios"
)

type Reading struct {
  StationID int
  Sensor Sensor
  SensorName string
  Value float64
  Valid bool

  WindSpeed float64
  WindDir float64

  StationBatLow bool
}

func (r Reading) String() string {
  if r.StationBatLow {
    batLow := "yes"
  } else {
    batLow := "no"
  }

  return fmt.Sprintf(
    "Reading for %s, station: %d, wind speed %2.0f, wind direction %3.0f, value %f, battery low %s",
    r.SensorName,
    r.StationID,
    r.WindSpeed,
    r.WindDir,
    r.Value,
    batLow,
  )
}

func ParsePacket(pkt radios.Packet) (rd Reading){
  rd.StationID = int((pkt.Data[0] & 0x0f) + 1)
  rd.Sensor = Sensor((pkt.Data[0] & 0xf0) >> 4)
  rd.SensorName = fmt.Sprintf("%s", rd.Sensor)
  rd.StationBatLow = (pkt.Data[0] & 0x08) > 0

  switch rd.Sensor {
  case SuperCapVoltage:
    rd.Value = convertSuperCapVoltage(pkt.Data[3:6])
  case UVIndex:
    rd.Value = 0
  case RainRate:
    rd.Value = convertRainRate(pkt.Data[3:6])
  case SolarRadiation:
    rd.Value = 0
  case Light:
    rd.Value = convertLight(pkt.Data[3:6])
  case Temperature:
    rd.Value = convertTemperature(pkt.Data[3:6])
  case WindGustSpeed:
    rd.Value = convertGustSpeed(pkt.Data[3:6])
  case Humidity:
    rd.Value = convertHumidity(pkt.Data[3:6])
  case RainClicks:
    rd.Value = convertRainClicks(pkt.Data[3:6])
  default:
    fmt.Sprintf("Unknown Reading Type: %0x", rd.Sensor)
  }

  rd.WindSpeed = float64(pkt.Data[1])
  rd.WindDir = ((float64(pkt.Data[2]) * 360) / 255)
  return
}

func convertSuperCapVoltage(data []byte) float64 {
  return float64((uint(data[0]) * 4 + ((uint(data[1]) & 0xC0) / 64)) / 100)
}

func convertRainRate(data []byte) float64 {
  if data[0] == 0xff{
    return 0
  } else if data[1] & 0x40 == 0 {
    return ((float64(uint(data[1]) & 0x30) / 16 * 250) + float64(data[0]))
  } else if data[1] & 0x40 == 0x40 {
    return ((float64(uint(data[1]) & 0x30) / 16 * 250) + float64(data[0])) / 16
  }
  log.Printf("Unkown Rain Rate State")
  return 0
}

func convertLight(data []byte) float64 {
  return float64(uint(data[0]) * 4) + (float64(uint(data[1]) & 0xc0) / 64)
}

func convertTemperature(data []byte) float64 {
  return float64(uint(data[0]) * 256 + uint(data[1])) / 160
}

func convertGustSpeed(data []byte) float64 {
  return float64(data[0])
}

func convertHumidity(data []byte) float64 {
  return float64(((uint(data[1]) & 0xf0) * 16) + uint(data[0])) / 10
}

func convertRainClicks(data []byte) float64 {
  return float64(uint(data[0]) & 0x7f)
}

type Sensor byte

const (
  SuperCapVoltage Sensor = 2
  UVIndex         Sensor = 4
  RainRate        Sensor = 5
  SolarRadiation  Sensor = 6
  Light           Sensor = 7
  Temperature     Sensor = 8
  WindGustSpeed   Sensor = 9
  Humidity        Sensor = 0xA
  RainClicks      Sensor = 0xE
)

func (r Sensor) String() string {
  switch r {
  case SuperCapVoltage:
    return "SuperCapVoltage"
  case UVIndex:
    return "UVIndex"
  case RainRate:
    return "RainRate"
  case SolarRadiation:
    return "SolarRadiation"
  case Light:
    return "Light"
  case Temperature:
    return "Temperature"
  case WindGustSpeed:
    return "WindGustSpeed"
  case Humidity:
    return "Humidity"
  case RainClicks:
    return "RainClicks"
  default:
    return fmt.Sprintf("Unknown Reading Type: %0x", uint(r))
  }
}
