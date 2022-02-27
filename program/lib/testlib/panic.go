package testlib

// DoesPanic runs f and checks if it panicked or not.
// When f does panic, returns the recovered value.
func DoesPanic(f func()) (panicked bool, recovered interface{}) {

	// In principle this function could just return the value of recover.
	// However that wouldn't allow to tell the difference between f calling panic(nil) and not panicking at all.

	defer func() {
		recovered = recover()
	}()

	panicked = true
	f()
	panicked = false

	return
}
