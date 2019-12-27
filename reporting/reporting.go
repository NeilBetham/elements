package reporting

import (
  "fmt"
  "time"
  "bytes"
  "errors"
  "net/http"
  "crypto/tls"
  "encoding/json"
  "github.com/NeilBetham/elements/protocol"
  "github.com/NeilBetham/elements/config"
)


type Report struct {
  Reading struct {
    Type string `json:"type"`
    RawValue string `json:"raw_value"`
    DecodedValue string `json:"decoded_value"`
    Timestamp time.Time `json:"timestamp"`
  } `json:"reading"`
}


type Reporter struct {
  Client *http.Client
  Url string
  ApiKey string
}


func NewReporter(c config.Config) (r Reporter) {
  if c.Server.SslVerify == false {
    tr := &http.Transport {
      TLSClientConfig: &tls.Config { InsecureSkipVerify: true },
    }
    r.Client = &http.Client { Transport: tr }
  } else {
    r.Client = &http.Client {}
  }

  protocol := "http"
  if c.Server.Ssl {
    protocol = "https"
  }

  r.Url = fmt.Sprintf("%s://%s:%s/api/stations/%s/reading",
    protocol,
    c.Server.Host,
    c.Server.Port,
    c.Server.StationId,
  )

  r.ApiKey = c.Credentials.ApiKey

  return
}


func (rp *Reporter) ReportReading(r protocol.Reading) (err error) {
  var report Report
  report.Reading.Type = r.SensorName
  report.Reading.RawValue = fmt.Sprintf("%X", r.RawValue)
  report.Reading.DecodedValue = fmt.Sprintf("%f", r.Value)
  report.Reading.Timestamp = time.Now()

  jsonData, err :=  json.Marshal(report)
  if err != nil {
    return
  }

  req, err := http.NewRequest("POST", rp.Url, bytes.NewBuffer(jsonData))
  if err != nil {
    return
  }
  req.Header.Add("Authorization", fmt.Sprintf("Token %s", rp.ApiKey))
  req.Header.Set("Content-Type", "application/json")

  resp, err := rp.Client.Do(req)
  if resp.StatusCode > 300 {
    err = errors.New(fmt.Sprintf("Error posting to API, http code: %i", resp.StatusCode))
  }
  return
}
