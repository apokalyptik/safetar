package main

import (
	"archive/tar"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var tempDir = ""

func fatalTarErr(h *tar.Header, op string, err error) {
	os.Stderr.WriteString(fmt.Sprintf(
		"\n\nFatal error reading tar from stdin for file: %s (%d): during: %s: %s\n\n",
		h.Name,
		h.Size,
		op,
		err.Error(),
	))
	os.Exit(1)
}

func main() {
	flag.StringVar(&tempDir, "temp-dir", tempDir, "Path to where temp files should be written (otherwise cwd is used)")
	flag.Parse()
	input := tar.NewReader(os.Stdin)
	for {
		header, err := input.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		buffer, err := ioutil.TempFile(tempDir, ".safetar.buffer")
		if err != nil {
			fatalTarErr(header, "temp-file-create", err)
		}
		os.Remove(buffer.Name())
		output := tar.NewWriter(buffer)

		if err := output.WriteHeader(header); err != nil {
			fatalTarErr(header, "write-output-header", err)
		}
		if _, err := io.CopyN(output, input, header.Size); err != nil {
			fatalTarErr(header, "write-output-bytes", err)
		}
		output.Flush()
		buffer.Sync()
		buffer.Seek(0, os.SEEK_SET)
		io.Copy(os.Stdout, buffer)
		buffer.Close()
	}
}
