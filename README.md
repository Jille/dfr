# dfr

dfr is a library to make defers more powerful. Usage is simple:

```
func() (retErr error) {
	d := dfr.D{}
	defer d.Run(&retErr)
	d.AddError(func() error {
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
