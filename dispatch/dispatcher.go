//go:build darwin
// +build darwin

package dispatch

/*
#cgo CFLAGS: -x objective-c -Wno-everything
#include <dispatch/dispatch.h>

void dispatch_async_signal(void *queue);

*/
import "C"

//export dispatcher
func dispatcher() {
	(<-dispatchQueue)()
}
