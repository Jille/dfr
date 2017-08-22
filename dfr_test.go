package dfr

import (
	"errors"
	"strings"
	"testing"
)

var (
	Err1 = errors.New("Error 1")
	Err2 = errors.New("Error 2")
)

func TestBasic(t *testing.T) {
	err := func() (retErr error) {
		d := D{}
		defer d.Run(&retErr)
		d.AddErr(func() error {
			return Err1
		})
		return
	}()
	if err != Err1 {
		t.Errorf("got %v, want %v", err, Err1)
	}
}

func TestPanicWithoutRetErr(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if r != Err1 {
				t.Errorf("got panic %v, want %v", r, Err1)
			}
		} else {
			t.Errorf("dfr didn't panic but should have")
		}
	}()
	func() {
		d := D{}
		defer d.Run(nil)
		d.AddErr(func() error {
			return Err1
		})
		return
	}()
}

func TestTwoErrors(t *testing.T) {
	err := func() (retErr error) {
		d := D{}
		defer d.Run(&retErr)
		d.AddErr(func() error {
			return Err1
		})
		d.AddErr(func() error {
			return Err2
		})
		return
	}()
	if !strings.Contains(err.Error(), Err1.Error()) {
		t.Errorf("got %v, should have contained %v", err, Err1)
	}
	if !strings.Contains(err.Error(), Err2.Error()) {
		t.Errorf("got %v, should have contained %v", err, Err2)
	}
}

func TestCallbackPanic(t *testing.T) {
	ok := false
	defer func() {
		if r := recover(); r != nil {
			if r != Err2 {
				t.Errorf("got panic %v, want %v", r, Err1)
			}
		} else {
			t.Errorf("dfr didn't panic but should have")
		}
		if !ok {
			t.Errorf("first deferred functions wasn't called when second one paniced")
		}
	}()
	func() (retErr error) {
		d := D{}
		defer d.Run(&retErr)
		d.Add(func() {
			ok = true
		})
		d.Add(func() {
			panic(Err2)
		})
		return
	}()
}

func TestCancel(t *testing.T) {
	err := func() (retErr error) {
		d := D{}
		defer d.Run(&retErr)
		cancel := d.AddErr(func() error {
			return Err1
		})
		d.AddErr(func() error {
			return nil
		})
		cancel(false)
		return
	}()
	if err != nil {
		t.Errorf("callback wasn't cancelled")
	}
}

func TestRunImmediately(t *testing.T) {
	ok := 0
	func() (retErr error) {
		d := D{}
		defer d.Run(&retErr)
		cancel := d.Add(func() {
			ok++
		})
		d.AddErr(func() error {
			return nil
		})
		cancel(true)
		if ok != 1 {
			t.Errorf("cancel(true) didn't immediately execute callback")
		}
		return
	}()
	if ok != 1 {
		t.Errorf("cancel(true) didn't stop callback from running during defer")
	}
}
