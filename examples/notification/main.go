package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -lobjc -framework AppKit
#include <AppKit/AppKit.h>
*/
import "C"

import (
	"fmt"
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

type NSUserNotification struct {
	objc.Object
}

var NSUserNotification_ = objc.Get("NSUserNotification")

type NSUserNotificationCenter struct {
	objc.Object
}

var NSUserNotificationCenter_ = objc.Get("NSUserNotificationCenter")

func main() {
	handlerCls := objc.NewClass("EventHandler", "NSObject")
	handlerCls.AddMethod("outputDeviceChanged:", func(notification objc.Object) {
		fmt.Println(notification)
	})
	objc.RegisterClass(handlerCls)
	handler := objc.Get("EventHandler").Alloc().Init()

	app := cocoa.NSApp_WithDidLaunch(func(app objc.Object) {
		notification := NSUserNotification{NSUserNotification_.Alloc().Init()}
		notification.Set("title:", core.String("Hello, world!"))
		notification.Set("informativeText:", core.String("More text"))

		sel := objc.Sel("outputDeviceChanged:")
		fmt.Println(sel.SelectorAddress())
		//fmt.Println(objc.Get("AVAudioSession"))
		fmt.Println(objc.Get("AVAudioSession").Get("sharedInstance"))

		center := core.NSNotificationCenter_defaultCenter()
		center.AddObserver_selector_name_object_(handler, objc.Sel("outputDeviceChanged:"), core.String("AVAudioSessionRouteChangeReasonNewDeviceAvailable"), objc.Get("AVAudioSession").Get("sharedInstance"))

		fmt.Println(handler)
	})

	app.SetActivationPolicy(cocoa.NSApplicationActivationPolicyRegular)
	app.ActivateIgnoringOtherApps(true)
	app.Run()
}
