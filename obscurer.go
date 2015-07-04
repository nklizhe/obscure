package obscure

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"io/ioutil"
	"text/template"
)

// Obscurer provides an interface for obscuring clear text
type Obscurer interface {
	// Obscure reads all clear text from input and writes obscured text to output
	Obscure(input io.Reader, output io.Writer) error
}

// DefaultObscurer obscures clear text using AES-256 algorithm
type DefaultObscurer struct {
	rand     io.Reader
	key      []byte
	template string
	pkgName  string
	varName  string
}

// DefaultObscurerOptions provides options for creating a DefaultObscurer
type DefaultObscurerOptions struct {
	// Key contains the encryption key for the obscurer.
	// If Key is nil, a random 32-bytes key will be generated automatically.
	Key []byte
	// PkgName contains the package name of the generated go code.
	// If PkgName is not provided, default package "main" will be used.
	PkgName string
	// VarName contains the variable name of the generated go code.
	// If VarName is empty, default name "Secret" will be used.
	VarName string
}

// NewDefaultObscurer creates a new DefaultObscurer
func NewDefaultObscurer(rand io.Reader, options DefaultObscurerOptions) *DefaultObscurer {
	if options.Key == nil {
		buf := make([]byte, 32)
		rand.Read(buf)
		options.Key = buf
	}
	if options.PkgName == "" {
		options.PkgName = "main"
	}
	if options.VarName == "" {
		options.VarName = "Secret"
	}
	return &DefaultObscurer{
		rand:     rand,
		key:      options.Key,
		pkgName:  options.PkgName,
		varName:  options.VarName,
		template: defaultTemplate,
	}
}

// Obscure implements the Obscurer interface
// It reads all clear text from input, encrypt it using AES-256 with the key,
// then generates go code and write to output
func (ob *DefaultObscurer) Obscure(input io.Reader, output io.Writer) error {
	buf, err := ioutil.ReadAll(input)
	if err != nil {
		return err
	}

	// encrypt text
	block, err := aes.NewCipher(ob.key)
	if err != nil {
		return err
	}
	ciphertext := make([]byte, aes.BlockSize+len(buf))
	iv := ciphertext[:aes.BlockSize]
	if _, err := ob.rand.Read(iv); err != nil {
		return err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(buf))

	// render go code
	return ob.render(output, ciphertext)
}

func (ob *DefaultObscurer) render(output io.Writer, data []byte) error {
	outPass := &bytes.Buffer{}
	fmt.Fprintf(outPass, "[]byte{")
	for i, c := range []byte(ob.key) {
		if i == 0 {
			fmt.Fprintf(outPass, "\n\t\t")
		}
		fmt.Fprintf(outPass, "%d, ", c)
	}
	fmt.Fprintf(outPass, "\n\t}")

	outData := &bytes.Buffer{}
	fmt.Fprintf(outData, "[]byte{")
	for i, c := range []byte(data) {
		if i == 0 {
			fmt.Fprintf(outData, "\n\t\t")
		}
		fmt.Fprintf(outData, "%d, ", c)
	}
	fmt.Fprintf(outData, "\n\t}")

	type templateData struct {
		Package string
		Name    string
		Key     string
		Data    string
	}
	tmpl, err := template.New("template").Parse(ob.template)
	if err != nil {
		return err
	}
	return tmpl.Execute(output, templateData{
		Package: ob.pkgName,
		Name:    ob.varName,
		Key:     outPass.String(),
		Data:    outData.String(),
	})
}

var (
	defaultTemplate = `package {{.Package}}

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

var (
	data{{.Name}} = {{.Data}}
	key{{.Name}} = {{.Key}}
)

// decode{{.Name}} decrypts and returns the original text
func decode{{.Name}}() ([]byte, error){
	block, err := aes.NewCipher(key{{.Name}})
	if err != nil {
		return nil, err
	}
	if len(data{{.Name}}) < aes.BlockSize {
		return nil, errors.New("error block size")
	}
	iv := data{{.Name}}[:aes.BlockSize]
	text := data{{.Name}}[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	return text, nil
}`
)
