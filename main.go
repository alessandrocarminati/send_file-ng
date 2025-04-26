package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"github.com/tarm/serial"
	"io"
	"log"
	"os"
	"time"
	"github.com/cheggaaa/pb/v3"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: send_console-ng <input-file> <serial-device>")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	serialPath := os.Args[2]

	c := &serial.Config{
		Name:        serialPath,
		Baud:        115200,
		ReadTimeout: time.Millisecond * 500,
	}
	port, err := serial.OpenPort(c)
	if err != nil {
		log.Fatalf("Failed to open serial port: %v", err)
	}
	defer port.Close()

	reader := bufio.NewReader(port)

	sendLine := func(cmd string) {
		fmt.Printf("-> %s\n", cmd)
		port.Write([]byte(cmd + "\n"))
		time.Sleep(500 * time.Millisecond)
		readSome(reader)
	}

	sendLine("stty -icanon -echo")
	sendLine("stdbuf -o0 cat > output.txt")

	data, err := os.ReadFile(inputPath)
	if err != nil {
		log.Fatalf("Failed to read input file: %v", err)
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write(data)
	gz.Close()

	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	chunkedSend(port, encoded, 128, 25*time.Millisecond)

	port.Write([]byte("\n"))
	port.Write([]byte{3})
	time.Sleep(300 * time.Millisecond)

	sendLine("stty sane")
	fmt.Println("✅ Transfer complete.")
}

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

func readSome(r io.Reader) {

	buf := make([]byte, 512)
	for {
		if port, ok := r.(interface{ SetReadDeadline(time.Time) error }); ok {
			port.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		}

		n, err := r.Read(buf)
		if err != nil {
			break
		}
		if n > 0 {
			fmt.Printf("<- %s", string(buf[:n]))
		}
	}
}
