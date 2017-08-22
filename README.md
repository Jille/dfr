# dfr

[![Build Status](https://travis-ci.org/Jille/dfr.png)](https://travis-ci.org/Jille/dfr)

dfr is a library to make defers more powerful. Usage is simple:

```
func() (retErr error) {
	d := dfr.D{}
	defer d.Run(&retErr)
	d.AddErr(func() error {
		// I'll be run during defers
		return nil
	})
	cancel := d.Add(func() {
		panic("I'll never be run")
	})
	cancel(false)
	what := "crap"
	cancel = d.Add(func() {
		what = "awesome"
	})
	// Giving true to the cancel function will execute it immediately.
	cancel(true)
	fmt.Printf("This library is %s\n", what)
}
```

# But why?

Here are some code samples where your code would benefit from `dfr`.

## Avoid forgetting a mtx.Unlock()

```
func() {
	mtx.Lock()
	if x == 5 {
		mtx.Unlock()
		return
	}
	if y == 6 {
		mtx.Unlock()
		return
	}
	mtx.Unlock()
	doALotOfWorkOutsideTheLock()
}
```

becomes:

```
func() {
	d := dfr.D{}
	defer d.Run(nil)
	mtx.Lock()
	unlock := d.Add(mtx.Unlock)
	if x == 5 {
		return
	}
	if y == 6 {
		return
	}
	unlock(true)
	doALotOfWorkOutsideTheLock()
}
```

## Deferring in a loop

```
func() {
	for {
		mtx.Lock()
		doSomethingThatMightPanic()
		mtx.Unlock()
	}
}
```

becomes:

```
func() {
	d := dfr.D{}
	defer d.Run(nil)
	for {
		mtx.Lock()
		unlock := d.Add(mtx.Unlock)
		doSomethingThatMightPanic()
		unlock(true)
	}
}
```

## Ignoring error codes from deferred calls

```
func() error {
	fh, err := os.Create("test.txt")
	if err != nil {
		return err
	}
	// return error from fh.Close() gets silently discard!
	defer fh.Close()
}
```

becomes:

```
func() (retErr error) {
	d := dfr.D{}
	defer d.Run(&retErr)
	fh, err := os.Create("test.txt")
	if err != nil {
		return err
	}
	d.AddErr(fh.Close)
}
```
