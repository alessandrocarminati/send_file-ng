package main

import (
	"path/filepath"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"github.com/tarm/serial"
	"io"
	"regexp"
	"log"
	"os"
	"time"
	"github.com/cheggaaa/pb/v3"
)
var Build string
var Version string
var Hash string
var Dirty string
var AppName string

func chunkedSend(w io.Writer, data string, chunkSize int, delay time.Duration) {
	var bar *pb.ProgressBar
	bar = pb.StartNew(len(data)/chunkSize)
	for i := 0; i < len(data); i += chunkSize {
		bar.Increment()
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		w.Write([]byte(data[i:end]))
		time.Sleep(delay)
	}
	bar.Finish()
}

func readOut(r io.Reader) string {

	buf := make([]byte, 512)
	for i:=0; i<10; i++ {
		if port, ok := r.(interface{ SetReadDeadline(time.Time) error }); ok {
			port.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		}

		n, err := r.Read(buf)
		if err != nil {
			break
		}
		if n > 0 {
			return string(buf[:n])
		}
	}
	return ""
}

func versionString() string {
	return fmt.Sprintf("%s %s.%s(%s) - %s\n", AppName, Version, Build, Hash, Dirty)
}

func setupSer(device string, speed int) (*serial.Port, error) {
	c := &serial.Config{
		Name:        device,
		Baud:        speed,
		ReadTimeout: time.Millisecond * 500,
	}
	port, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}
	return port, nil
}

func sendLine(port io.ReadWriter, cmd string) string {
	fmt.Printf("Sending: %s\n", cmd)
	port.Write([]byte(cmd + "\n"))
	time.Sleep(500 * time.Millisecond)
	return readOut(port)
}

func checkCommand(port io.ReadWriter, cmd, regex string) bool {
	fmt.Printf("verifying %s\n", cmd)
	r, _ := regexp.Compile(regex)
	port.Write([]byte(cmd + " --version \n"))
	time.Sleep(500 * time.Millisecond)
	s := readOut(port)
	fmt.Printf("%s --version \n%s\n", cmd, s)
	return r.MatchString(s)

}

func main() {
	var buf bytes.Buffer
	neededCmds1 := []string {"stty", "stdbuf", "cat", "rm", "base64", "gzip"}
	neededCmdsRegex1 := []string {"stty[^0-9]+[0-9\\.]+", "stdbuf[^0-9]+[0-9\\.]+", "cat[^0-9]+[0-9\\.]+", "rm[^0-9]+[0-9\\.]+", "base64[^0-9]+[0-9\\.]+", "gzip[^0-9]+[0-9\\.]+"}

	cl := parseCMDline()

	if cl.Ver {
		fmt.Println(versionString())
		os.Exit(0)
	}

	if cl.Help {
		fmt.Println(versionString())
		fmt.Println(helpText())
		os.Exit(0)
	}

	if cl.Device == "DEFAULT" || cl.Filename == "DEFAULT"  || cl.BaudRate == 0 {
		fmt.Println("Error: Missing or wrong argument\n\n")
		fmt.Println(versionString())
		fmt.Println(helpText())
		os.Exit(1)
	}

	port, err := setupSer(cl.Device, cl.BaudRate)
	defer port.Close()

	for i, cmd := range neededCmds1 {
		if ! checkCommand(port, cmd, neededCmdsRegex1[i]) {
			fmt.Printf("Prerequisite %s is not met\n", cmd)
			os.Exit(1)
		}
	}

	_ = sendLine(port, "stty -icanon -echo")
	_ = sendLine(port, fmt.Sprintf("stdbuf -o0 cat > %s", cl.RemoteTmpFN))

	data, err := os.ReadFile(cl.Filename)
	if err != nil {
		log.Fatalf("Failed to read input file: %v", err)
	}

	gz := gzip.NewWriter(&buf)
	gz.Write(data)
	gz.Close()

	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	chunkedSend(port, encoded, 128, 25*time.Millisecond)

	port.Write([]byte("\n"))
	port.Write([]byte{3})
	time.Sleep(300 * time.Millisecond)

	remoteFilename := filepath.Base(cl.Filename)
	remoteDir := filepath.Dir(cl.RemoteTmpFN)

	_ = sendLine(port, "stty sane")
	_ = sendLine(port, fmt.Sprintf("cat %[1]s | base64 -d | gzip -d >%[2]s/%[3]s && rm -rf  %[1]s", cl.RemoteTmpFN, remoteDir, remoteFilename))
	fmt.Println("Transfer complete.")
}

