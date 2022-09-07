package mediaplayer

import (
	"fmt"
	"github.com/progrium/macdriver/avcore"
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
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

	playingCenter := MPNowPlayingInfoCenter_defaultCenter()
	playingCenter.SetPlaybackState_(MPNowPlayingPlaybackStatePaused)

	item := avcore.AVPlayerItem_playerItemWithURL_(core.NSURL_fileURLWithPath_isDirectory_(core.String("/Users/anhoder/Desktop/1.mp3"), false))
	player := avcore.AVPlayer_playerWithPlayerItem_(item)

	fmt.Println(core.NSKeyValueObservingOptionNew | core.NSKeyValueObservingOptionOld)
	core.AddObserver_forKeyPath_options_context(
		item, core.String("status"), core.NSKeyValueObservingOptionNew|core.NSKeyValueObservingOptionOld, nil,
		func(keyPath core.NSString, ofObject objc.Object, change core.NSDictionary, context objc.Ref) {
			if core.NSNumber_fromRef(change.ObjectForKey(core.String("new"))).IsEqualToNumber_(core.NSNumber_numberWithInt_(1)) {
				playingCenter.SetPlaybackState_(MPNowPlayingPlaybackStatePlaying)
			}
			playingCenter.SetNowPlayingInfo_(nowPlayingInfoOfPlayer(&player))

			fmt.Println(keyPath, ofObject, change, context)
		})

	cls := objc.NewClass("CommandHandler", "NSObject")
	cls.AddMethod("handlePlayCommand:", func(event objc.Object) core.NSInteger {
		fmt.Printf("playing: %#v\n", event)
		playingCenter.SetNowPlayingInfo_(nowPlayingInfoOfPlayer(&player))
		playingCenter.SetPlaybackState_(MPNowPlayingPlaybackStatePlaying)
		player.Play()
		return MPRemoteCommandHandlerStatusSuccess
	})
	cls.AddMethod("handlePausedCommand:", func(event objc.Object) core.NSInteger {
		fmt.Printf("paused: %#v\n", event)
		player.Pause()
		playingCenter.SetPlaybackState_(MPNowPlayingPlaybackStatePaused)
		playingCenter.SetNowPlayingInfo_(nowPlayingInfoOfPlayer(&player))
		return MPRemoteCommandHandlerStatusSuccess
	})
	objc.RegisterClass(cls)
	handler := objc.Get("CommandHandler").Alloc().Init()

	center := MPRemoteCommandCenter_sharedCommandCenter()
	center.SkipBackwardCommand().SetPreferredIntervals_(core.NSArray_arrayWithObject_(core.NSNumber_numberWithFloat_(15.0)))
	center.SkipForwardCommand().SetPreferredIntervals_(core.NSArray_arrayWithObject_(core.NSNumber_numberWithFloat_(15.0)))
	center.PlayCommand().AddTarget_action_(handler, objc.Sel("handlePlayCommand:"))
	center.PauseCommand().AddTarget_action_(handler, objc.Sel("handlePausedCommand:"))
	center.StopCommand().AddTarget_action_(handler, objc.Sel("handlePlayCommand:"))
	center.TogglePlayPauseCommand().AddTarget_action_(handler, objc.Sel("handlePlayCommand:"))
	center.NextTrackCommand().AddTarget_action_(handler, objc.Sel("handlePlayCommand:"))
	center.PreviousTrackCommand().AddTarget_action_(handler, objc.Sel("handlePlayCommand:"))
	center.ChangeRepeatModeCommand().AddTarget_action_(handler, objc.Sel("handlePlayCommand:"))
	center.ChangeShuffleModeCommand().AddTarget_action_(handler, objc.Sel("handlePlayCommand:"))
	center.ChangePlaybackRateCommand().AddTarget_action_(handler, objc.Sel("handlePlayCommand:"))
	center.SeekBackwardCommand().AddTarget_action_(handler, objc.Sel("handlePlayCommand:"))
	center.SeekForwardCommand().AddTarget_action_(handler, objc.Sel("handlePlayCommand:"))
	center.SkipForwardCommand().AddTarget_action_(handler, objc.Sel("handlePlayCommand:"))
	center.SkipBackwardCommand().AddTarget_action_(handler, objc.Sel("handlePlayCommand:"))
	center.ChangePlaybackPositionCommand().AddTarget_action_(handler, objc.Sel("handlePlayCommand:"))
	center.LikeCommand().AddTarget_action_(handler, objc.Sel("handlePlayCommand:"))
	center.DislikeCommand().AddTarget_action_(handler, objc.Sel("handlePlayCommand:"))
	center.BookmarkCommand().AddTarget_action_(handler, objc.Sel("handlePlayCommand:"))
	center.EnableLanguageOptionCommand().AddTarget_action_(handler, objc.Sel("handlePlayCommand:"))
	center.DisableLanguageOptionCommand().AddTarget_action_(handler, objc.Sel("handlePlayCommand:"))

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
	<-ok
}

func nowPlayingInfoOfPlayer(player *avcore.AVPlayer) core.NSDictionary {
	total := player.CurrentItem().Duration().Value / int64(player.CurrentItem().Duration().Timescale)
	ur := player.CurrentTime().Value / int64(player.CurrentTime().Timescale)

	values := core.NSArray_array()
	keys := core.NSArray_array()
	values = values.ArrayByAddingObject_(core.NSNumber_numberWithInt_(int32(total)))
	keys = keys.ArrayByAddingObject_(core.String(MPMediaItemPropertyPlaybackDuration))

	values = values.ArrayByAddingObject_(core.NSNumber_numberWithInt_(int32(ur)))
	keys = keys.ArrayByAddingObject_(core.String(MPNowPlayingInfoPropertyElapsedPlaybackTime))

	values = values.ArrayByAddingObject_(core.NSNumber_numberWithFloat_(player.Rate()))
	keys = keys.ArrayByAddingObject_(core.String(MPNowPlayingInfoPropertyPlaybackRate))

	values = values.ArrayByAddingObject_(core.NSNumber_numberWithFloat_(1.0))
	keys = keys.ArrayByAddingObject_(core.String(MPNowPlayingInfoPropertyDefaultPlaybackRate))
	return core.NSDictionary_dictionaryWithObjects_forKeys_(values, keys)
}