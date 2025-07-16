package testutil_test

import (
	"fmt"

	"go.tkw01536.de/ggman/internal/testutil"
)

func ExampleRecordingT() {
	var r1 testutil.RecordingT
	r1.Helper()
	r1.Errorf("hello %s", "world")

	fmt.Println(r1.HelperCalled)
	fmt.Println(r1.Message)

	var r2 testutil.RecordingT
	r2.Errorf("1 + 2 = %d", 1+2)

	fmt.Println(r2.HelperCalled)
	fmt.Println(r2.Message)

	// Output: true
	// hello world
	// false
	// 1 + 2 = 3
}
