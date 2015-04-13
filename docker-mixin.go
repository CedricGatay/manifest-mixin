package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"errors"
)

func readRemoteDockerFile(url string) (string, error) {
	resp, err := http.Get(url)
	buf := new(bytes.Buffer)
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK{
			buf.ReadFrom(resp.Body)
		}else{
			err = errors.New("Unable to get file")
		}
	}
	return buf.String(), err
}

func cleanupDockerFile(mixinURL string, filecontent string) (cleaned string) {
	cleaned = "##### MIXIN BEGIN (" + mixinURL + ")\n"
	fromRegex, _ := regexp.Compile("^(FROM|MAINTAINER).*")
	lines := strings.Split(filecontent, "\n")
	for _, line := range lines {
		if fromRegex.MatchString(line) {
			cleaned += "# " + line + "\n"
		} else {
			cleaned += line + "\n"
		}
	}
	cleaned += "##### MIXIN END \n"
	return
}

func mixinDockerFile(inpath string, outpath string) {
	lines, err := ReadLines(inpath)
	if err != nil {
		return
	}
	var outLines []string
	mixinRegex, _ := regexp.Compile("^MIXIN\\s*([^\\s]*)")
	for _, line := range lines {
		if mixinRegex.MatchString(line) {
			submatches := mixinRegex.FindStringSubmatch(line)
			for i, v := range submatches{
				if i>0 { //full sentence matched in index 0
					body, err := readRemoteDockerFile(v)
					if err == nil {
						outLines = append(outLines, cleanupDockerFile(v, body))
					} else {
						fmt.Fprintf(os.Stderr, "Unable to get MIXIN content at %s (%s)\n", v, err)
					}
				}
			}
		} else {
			outLines = append(outLines, line)
		}
	}
	WriteLines(outLines, outpath)
}

// Read a whole file into the memory and store it as array of lines
func ReadLines(path string) (lines []string, err error) {
	var (
		file   *os.File
		part   []byte
		prefix bool
	)
	if file, err = os.Open(path); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open %s for reading", path)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			lines = append(lines, buffer.String())
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

func WriteLines(lines []string, path string) (err error) {
	var (
		file *os.File
	)

	if file, err = os.Create(path); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open %s for writing", path)
		return
	}
	defer file.Close()

	for _, item := range lines {
		_, err := file.WriteString(strings.TrimSpace(item) + "\n")
		if err != nil {
			fmt.Println(err)
			break
		}
	}
	return
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "You need to call %s with at least one argument %s (default output file name is Dockerfile)\n", os.Args[0], "infile")
		os.Exit(-1)
	}
	outName := "Dockerfile"
	if len(os.Args) == 3 {
		outName = os.Args[2]
	}
	mixinDockerFile(os.Args[1], outName)
}
