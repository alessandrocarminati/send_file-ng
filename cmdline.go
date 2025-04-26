package main
import (
	"flag"
	"bytes"
)
type Cmdline struct {
	Device      string
	Filename    string
	RemoteTmpFN string
	Help        bool
	Ver         bool
	BaudRate    int
}

func parseCMDline() *Cmdline {
	DevicePtr      := flag.String ("d", "DEFAULT", "Serial device where send commands")
	FilenamePtr    := flag.String ("f", "DEFAULT", "Local File to send")
	RemoteTmpFNPtr := flag.String ("r", "/tmp/rawtransferred", "remote filename and path for tmp file")
	VerPtr         := flag.Bool   ("v", false, "Returns the version string")
	helpPtr        := flag.Bool   ("help", false, "Show help")
	BaudRatePtr    := flag.Int    ("b", 0, "Serial speed")
	flag.Parse()

	config := Cmdline{
		Device:      *DevicePtr,
		Filename:    *FilenamePtr,
		RemoteTmpFN: *RemoteTmpFNPtr,
		Help:        *helpPtr,
		Ver:         *VerPtr,
		BaudRate:    *BaudRatePtr,
	}
	return &config
}
func helpText() string {
	var buf bytes.Buffer
	flag.CommandLine.SetOutput(&buf)
	flag.Usage()
	return buf.String()
}
