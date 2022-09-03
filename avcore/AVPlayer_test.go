package avcore

import (
	"fmt"
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/dispatch"
	"os"
	"runtime"
	"testing"
	"time"
)

func init() {
	runtime.LockOSThread()
}

func TestMain(m *testing.M) {
	go func() {
		os.Exit(m.Run())
	}()
	app := cocoa.NSApp()
	app.SetActivationPolicy(cocoa.NSApplicationActivationPolicyProhibited)
	app.ActivateIgnoringOtherApps(true)
	app.Run()
}

func TestAsync(t *testing.T) {
	ok := make(chan bool)
	player := AVPlayer_playerWithURL_(core.NSURL_fileURLWithPath_isDirectory_(core.String("/Users/anhoder/Desktop/1.mp3"), false))
	dispatch.Async(dispatch.MainQueue(), func() {
		player.SetAllowsExternalPlayback_(true)
		player.Play()
		go func() {
			ticker := time.Tick(time.Second)
			for {
				<-ticker
				cur := player.CurrentTime()
				total := player.CurrentItem().Duration()
				fmt.Println(cur.Value/int64(cur.Timescale), total.Value/int64(total.Timescale))
				if cur.Value/int64(cur.Timescale) >= total.Value/int64(total.Timescale) {
					break
				}
			}
			ok <- true
		}()
	})
	<-ok
}
