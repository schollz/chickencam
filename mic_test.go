package mic

import (
	"fmt"
	"os"
	"testing"

	"github.com/cocoonlife/goalsa"
)

func TestMic(t *testing.T) {
	dev, err := alsa.NewCaptureDevice("default", 1, alsa.FormatU8, 8000, alsa.BufferParams{})
	if err != nil {
		t.Errorf("Err initializing mic:\n\t%s", err)
		t.FailNow()
	}
	d := 5 //seconds to record
	r := make([]int8, 8000*d)
	fmt.Printf("Recording for %d seconds...\n", d)
	n, err := dev.Read(r)
	if err != nil {
		t.Errorf("Err reading mic:\n\t%s", err)
		t.FailNow()
	}
	if n != len(r) {
		t.Errorf("Didnt read enough: Wanted %d, Got %d", len(r), n)
		t.FailNow()
	}
	//cast to bytes
	b := make([]byte, len(r))
	for i := 0; i < len(b); i++ {
		b[i] = byte(r[i])
	}
	fmt.Printf("Done recording. Saving file.")
	f, err := os.OpenFile("./a.wav", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		t.Errorf("Failed to open file: %s", err)
		t.FailNow()
	}
	_, err = f.Write(b)
	if err != nil {
		t.Errorf("Failed to write to file: %s", err)
		t.FailNow()
	}
	err = f.Close()
	if err != nil {
		t.Errorf("Failed to close file: %s", err)
		t.FailNow()
	}
}
