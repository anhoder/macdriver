package dispatch

import (
	"errors"
	"fmt"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
	"os"
	"runtime"
	"testing"

	"github.com/progrium/macdriver/cocoa"
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
	Async(MainQueue(), func() {
		//data, err := os.ReadFile("/Users/anhoder/Desktop/1.mp3")
		//fmt.Println(err)

		cls := objc.NewClass("SoundDelegate", "NSObject")
		cls.AddMethod("sound:didFinishPlaying:", func(sound objc.Object, didFinishPlaying bool) {
			fmt.Println("finish playing: ", didFinishPlaying)
			if didFinishPlaying {
				ok <- true
			}
		})
		objc.RegisterClass(cls)

		delegate := objc.Get("SoundDelegate").Alloc().Init()

		url := core.NSURL_fileURLWithPath_isDirectory_(core.NSString_FromString("/Users/anhoder/Desktop/1.mp3"), false)
		s := cocoa.NSSound_InitWithURL(url)
		s.Set("delegate:", delegate)
		//d := core.NSData_WithBytes(data, uint64(len(data)))
		//s := cocoa.NSSound_InitWithData(d)
		s.Play()
	})
	<-ok

	//player := objc.Get("AVAudioPlayer").Alloc().Init()

	//player = player.Send("initWithContentsOfURL:", url)

}

func TestSync(t *testing.T) {
	var ok bool
	Sync(MainQueue(), func() {
		ok = true
	})
	if !ok {
		t.Fatal("ok not set to true")
	}
}

func TestDo(t *testing.T) {
	err := errors.New("test")
	d := Do(MainQueue(), func() error {
		return err
	})
	ret := d.Wait()
	if ret != err {
		t.Fatal("unexpected return value from Wait:", ret)
	}
}
