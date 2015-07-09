package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/nklizhe/obscure"
)

var (
	otp           []byte // one-time-pass
	input, output string
	name          string
	pkgName       string
)

func main() {

	flag.StringVar(&input, "i", "stdin", "the input file")
	flag.StringVar(&output, "o", "stdout", "the output file")
	flag.StringVar(&name, "name", "Secret", "name of the variable")
	flag.StringVar(&pkgName, "package", "main", "package name")
	flag.Parse()

	var in io.Reader
	var err error
	if input == "stdin" || input == "" {
		in = os.Stdin
	} else {
		in, err = os.Open(input)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}

	var out io.Writer
	if output == "stdout" || output == "" {
		out = os.Stdout
	} else {
		if path.Ext(output) != "go" {
			output = fmt.Sprintf("%s.go", output)
		}
		out, err = os.Create(output)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}

	defer func() {
		// close input and output
		if input2, ok := in.(io.ReadCloser); ok {
			input2.Close()
		}
		if output2, ok := out.(io.WriteCloser); ok {
			output2.Close()
		}
	}()

	// obscure
	ob := obscure.NewDefaultObscurer(rand.Reader, obscure.DefaultObscurerOptions{
		PkgName: pkgName,
		VarName: name,
	})
	if err := ob.Obscure(in, out); err != nil {
		log.Fatal(err)
	}
}

func decode(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("block size error")
	}
	dst := make([]byte, len(ciphertext[aes.BlockSize:]))
	iv := ciphertext[:aes.BlockSize]
	text := ciphertext[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(dst, text)
	return dst, nil
}
