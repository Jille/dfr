package dfr

import (
	"github.com/Jille/errchain"
)

type D struct{
	defers []func() error
}

func (d *D) AddErr(cb func() error) func(bool) error {
	d.defers = append(d.defers, cb)
	idx := len(d.defers)-1
	return func(ex bool) error {
		defer func() {
			d.defers[idx] = nil
		}()
		if ex {
			return d.defers[idx]()
		}
		return nil
	}
}

func (d *D) Add(cb func()) func(bool) {
	canceler := d.AddErr(func() error {
		cb()
		return nil
	})
	return func(ex bool) {
		canceler(ex)
	}
}

func (d *D) Run(retErr *error) {
	d.run(retErr, 0)
}

func (d *D) run(retErr *error, offset int) {
	n := len(d.defers)
	defer func() {
		// If more defers were added during callbacks, do another pass for them.
		if len(d.defers) > n {
			d.run(retErr, n)
		}
	}()
	for i := offset; n > i; i++ {
		i := i
		defer func() {
			if cb := d.defers[i]; cb != nil {
				if err := cb(); err != nil {
					if retErr == nil {
						panic(err)
					}
					errchain.Append(retErr, err)
				}
			}
		}()
	}
}
