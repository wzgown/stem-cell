package util

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func CheckInt(t string) (int, bool) {
	n, err := strconv.Atoi(t)
	if err != nil {
		return 0, false
	}
	return n, true
}

func CheckBool(t string) (bool, bool) {
	switch strings.ToLower(t) {
	case "", "y", "yes":
		return true, true
	case "n", "no":
		return false, true
	}
	return false, false
}

func Makedir(err error, path string) error {
	if err != nil {
		return err
	}

	err = os.MkdirAll(path, 0755)
	if err != nil {
		fmt.Println("mkdir failed, detail:", err)
	}
	return err
}

func StripLineSuffix(path string) {
	by, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
		return
	}

	var content []string
	brd := bufio.NewReader(bytes.NewReader(by))
	var line string
	for {
		line, err = brd.ReadString('\n')

		if err != nil && err != io.EOF {
			panic(err)
		}

		newline := strings.TrimRight(line, "\t\n\v\f\r")
		content = append(content, newline)

		if err == io.EOF {
			break
		}
	}

	nf, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
		return
	}
	defer nf.Close()
	for _, l := range content {
		_, err = nf.WriteString(l + "\n")
		if err != nil {
			panic(err)
		}
	}
}

func Camel(in string) (out string) {
	par := strings.Split(in, "_")
	for i, tok := range par {
		par[i] = strings.ToUpper(tok[:1]) + tok[1:len(tok)]
	}
	out = strings.Join(par, "")
	return
}
