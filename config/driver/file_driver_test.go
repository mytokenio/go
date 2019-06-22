package driver

import (
	"github.com/mytokenio/go/log"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestFile(t *testing.T) {
	content := []byte("test value")
	tmpfile, err := ioutil.TempFile("", "config-test")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	d := NewFileDriver(Path(tmpfile.Name()))
	val, err := d.Get("")
	assert(t, val.String(), "test value")
	assert(t, err, nil)

	_, err = NewFileDriver().Get("xxx")
	b := strings.Contains(err.Error(), "no such file or directory")
	assert(t, b, true)
}

