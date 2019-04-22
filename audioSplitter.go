package main

import
(
//  "io/ioutil"
  "time"
  "fmt"
  "github.com/go-fingerprint/fingerprint"
  "github.com/go-fingerprint/gochroma"
//  "github.com/micmonay/keybd_event"
  "github.com/gordonklaus/portaudio"
  "os"
//  "os/signal"
  "log"
  "encoding/binary"
  "bytes"
  "text/template"
)

var tmpl = template.Must(template.New("").Parse(
	`{{. | len}} host APIs: {{range .}}
	Name:                   {{.Name}}
	{{if .DefaultInputDevice}}Default input device:   {{.DefaultInputDevice.Name}}{{end}}
	{{if .DefaultOutputDevice}}Default output device:  {{.DefaultOutputDevice.Name}}{{end}}
	Devices: {{range .Devices}}
		Name:                      {{.Name}}
		MaxInputChannels:          {{.MaxInputChannels}}
		MaxOutputChannels:         {{.MaxOutputChannels}}
		DefaultLowInputLatency:    {{.DefaultLowInputLatency}}
		DefaultLowOutputLatency:   {{.DefaultLowOutputLatency}}
		DefaultHighInputLatency:   {{.DefaultHighInputLatency}}
		DefaultHighOutputLatency:  {{.DefaultHighOutputLatency}}
		DefaultSampleRate:         {{.DefaultSampleRate}}
	{{end}}
{{end}}`,
))



func getFingerprint(filename string) string{
  reader, _ := os.Open(filename);
  fpcalc := gochroma.New(gochroma.AlgorithmDefault)
  defer fpcalc.Close()
  fprint, err := fpcalc.Fingerprint(
    fingerprint.RawInfo{
      Src: reader,
      Channels: 2,
      Rate: 44100,
      MaxSeconds: 130,
    })
    if err != nil{
      log.Fatal(err)
      panic(err)
    }
    return fprint
}

func getMicFingerprint() string{
  portaudio.Initialize()
  defer portaudio.Terminate()
  //h, err := portaudio.DefaultHostApi()
  //chk(err)
  //p := portaudio.LowLatencyParameters(h.DefaultInputDevice, nil)
  //p.Input.Channels = 2
  buffer := make([]float32, 44100 * 2)
  audioDevices, _ := portaudio.Devices()
  for i := range audioDevices{
    log.Printf("Id:%d\n   Name: %s\n    MaxInputChannels: %d\n    MaxOutPutChannels: %d\n    DefaultLowInputLatency: %d\n    --------------------------\n\n",i,audioDevices[i].Name,audioDevices[i].MaxInputChannels,audioDevices[i].MaxOutputChannels,audioDevices[i].DefaultLowInputLatency)
  }
  inputDevice := audioDevices[0]
  p := portaudio.LowLatencyParameters(inputDevice,nil)
  p.Input.Channels = 2
  p.Output.Channels = 0
  p.SampleRate = 44100
  p.FramesPerBuffer = len(buffer)
  stream, err := portaudio.OpenStream(p, func(in []float32) {
    for i := range buffer{
      buffer[i] = in[i]
    }
  })
  if err != nil{
    panic(err)
  }
  defer stream.Close()
  buf := new(bytes.Buffer)
  chk(stream.Start())
  for start:= time.Now(); time.Since(start) < time.Millisecond * 20000; {
    binary.Write(buf,binary.BigEndian,&buffer)
  }
  chk(stream.Stop())
  fpcalc := gochroma.New(gochroma.AlgorithmDefault)
  defer fpcalc.Close()
  fprint, err := fpcalc.Fingerprint(
    fingerprint.RawInfo{
      Src: buf,
      Channels: 2,
      Rate: 44100,
      MaxSeconds: 120,
  })
  if err != nil{
    log.Fatal(err)
    panic(err)
  }
  return fprint
}

func resizeFingerprint(fprint1 []int32, fprint2 []int32) []int32{
  nFprint := make([]int32, len(fprint1))
  for k, v := range fprint2{
    nFprint[k] = v
  }
  return nFprint
}

func main() {
  var fprint1 string
  if len(os.Args) < 2 {
		fmt.Println("missing required argument:  <Input audio filename>\nUsing default (lap_extended.wav)")
		fprint1 = getFingerprint("lap_extended.wav")
  }else{
    fprint1 = getFingerprint(os.Args[1])
  }
  for{
    fprint2 := getMicFingerprint()
    fmt.Println(fprint1 , fprint2)
  //  fprint1 = resizeFingerprint(fprint2, fprint1)
  //  s, err := fingerprint.Compare(fprint1, fprint2)
  /*  if err != nil{
      fmt.Printf("Error:%s ",err)
      return
    }
    fmt.Printf("Fingerprint score: %v\n", s)*/
  }
}


func chk(err error) {
	if err != nil {
		panic(err)
	}
}
