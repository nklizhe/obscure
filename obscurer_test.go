package obscure

import (
	"bytes"
	"crypto/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	key = []byte("password")
)

func TestNewDefaultObscurer(t *testing.T) {
	ob := NewDefaultObscurer(rand.Reader, DefaultObscurerOptions{
		Key: key,
	})
	assert.Equal(t, key, ob.key)
	assert.Equal(t, "main", ob.pkgName)

	ob2 := NewDefaultObscurer(rand.Reader, DefaultObscurerOptions{})
	assert.NotEmpty(t, ob2.key)
	assert.NotEqual(t, key, ob2.key)
}

func TestObscure(t *testing.T) {
	ob := NewDefaultObscurer(rand.Reader, DefaultObscurerOptions{})
	text := "clear text"
	input := bytes.NewReader([]byte(text))
	output := &bytes.Buffer{}
	err := ob.Obscure(input, output)
	assert.NoError(t, err)
	code := output.String()
	assert.True(t, strings.Contains(code, ob.pkgName))
	assert.True(t, strings.Contains(code, ob.varName))
	assert.False(t, strings.Contains(code, text))
}

func TestDecode(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "obscure_test")
	os.Mkdir(dir, 0755)
	defer os.RemoveAll(dir)
	// main.go
	mainCode := `package main

import (
	"fmt"
	"os"
)
func main() {
	text, err := decodeSecret()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Print(string(text))
}
`
	main, err := os.Create(filepath.Join(dir, "main.go"))
	assert.NoError(t, err)
	_, err = main.Write([]byte(mainCode))
	assert.NoError(t, err)
	main.Close()

	// secret.go
	output, err := os.Create(filepath.Join(dir, "secret.go"))
	assert.NoError(t, err)

	// generate
	ob := NewDefaultObscurer(rand.Reader, DefaultObscurerOptions{})
	text := "clear text"
	err = ob.Obscure(bytes.NewReader([]byte(text)), output)
	assert.NoError(t, err)
	output.Close()

	// run test program
	os.Chdir(dir)
	t.Logf("executing %s", dir)
	cmd := exec.Command("go", "build")
	_, err = cmd.CombinedOutput()
	assert.NoError(t, err)

	cmd2 := exec.Command("./obscure_test")
	outText, err := cmd2.CombinedOutput()
	assert.NoError(t, err)
	assert.Equal(t, text, string(outText))
}
