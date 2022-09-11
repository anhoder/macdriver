//go:build darwin
// +build darwin

package main

import (
	"fmt"
	"os"

	"github.com/progrium/macschema/schema"
)

func main() {
	err := generate("github.com/progrium/macdriver", []pkg{
		// NOTE: the ordering of packages is important since we only expose
		// imports from packages that come earlier in this list. This prevents
		// dependency cycles if packages depend on each other, which appears to
		// happen on occasion, e.g. a type in "Foundation" includes methods that
		// use types from "AppKit".
		{"core", []schemaLoader{
			loadFile("api/quartzcore/calayer.objc.json"),

			loadFile("api/foundation/nsarray.objc.json"),
			loadFile("api/foundation/nsattributedstring.objc.json"),
			loadFile("api/foundation/nsdata.objc.json"),
			loadFile("api/foundation/nsdate.objc.json"),
			loadFile("api/foundation/nsdictionary.objc.json"),
			loadFile("api/foundation/nsnumber.objc.json"),
			loadFile("api/foundation/nsrunloop.objc.json"),
			loadFile("api/foundation/nsstring.objc.json"),
			loadFile("api/foundation/nserrordomain.objc.json"),
			loadFile("api/foundation/nserroruserinfokey.objc.json"),
			loadFile("api/foundation/nserror.objc.json"),
			loadFile("api/foundation/nstimeinterval.objc.json"),
			loadFile("api/foundation/nsthread.objc.json"),
			loadFile("api/foundation/nsurl.objc.json"),
			loadFile("api/foundation/nsurlrequest.objc.json"),
			loadFile("api/foundation/nsnotification.objc.json"),
			loadFile("api/foundation/nsnotificationname.objc.json"),
			loadFile("api/foundation/nsoperationqueue.objc.json"),
			loadFile("api/foundation/nsnotificationcenter.objc.json"),
		}},

		{"cocoa", []schemaLoader{
			loadFile("api/foundation/nsbundle.objc.json"),
			loadFile("api/appkit/nssoundname.objc.json"),
			loadFile("api/appkit/nssoundplaybackdeviceidentifier.objc.json"),
			loadFile("api/appkit/nssound.objc.json"),

			loadFile("api/appkit/nsapplication.objc.json"),
			loadFile("api/appkit/nscontrol.objc.json"),
			loadFile("api/appkit/nsbutton.objc.json"),
			loadFile("api/appkit/nsevent.objc.json"),
			loadFile("api/appkit/nsfont.objc.json"),
			loadFile("api/appkit/nsimage.objc.json"),
			loadFile("api/appkit/nsimageview.objc.json"),
			loadFile("api/appkit/nsnib.objc.json"),
			loadFile("api/appkit/nspasteboard.objc.json"),
			loadFile("api/appkit/nslayoutmanager.objc.json"),
			loadFile("api/appkit/nsmenu.objc.json"),
			loadFile("api/appkit/nspopover.objc.json"),
			loadFile("api/appkit/nsmenuitem.objc.json"),
			loadFile("api/appkit/nsscreen.objc.json"),
			loadFile("api/appkit/nsstatusbar.objc.json"),
			loadFile("api/appkit/nsstatusbarbutton.objc.json"),
			loadFile("api/appkit/nsstatusitem.objc.json"),
			loadFile("api/appkit/nstext.objc.json"),
			loadFile("api/appkit/nstextcontainer.objc.json"),
			loadFile("api/appkit/nsviewcontroller.objc.json"),
			loadFile("api/appkit/nsvisualeffectview.objc.json"),
			loadFile("api/appkit/nswindow.objc.json"),

			loadFile("api/appkit/nscolor.objc.json").Then(func(s *schema.Schema) error {
				s.Class.TypeMethods = append(s.Class.TypeMethods, schema.Method{
					Identifier: schema.Identifier{
						Name: "colorWithRed:green:blue:alpha:",
					},
					Return: schema.DataType{Name: "NSColor", IsPtr: true},
					Args: []schema.Arg{
						{Name: "red", Type: schema.DataType{Name: "CGFloat"}},
						{Name: "green", Type: schema.DataType{Name: "CGFloat"}},
						{Name: "blue", Type: schema.DataType{Name: "CGFloat"}},
						{Name: "alpha", Type: schema.DataType{Name: "CGFloat"}},
					},
				})
				// TODO would be nice to make this fast to just extend with custom
				// props and methods, using macschema to parse them from a
				// declaration?
				// Though this one exists in the docs, but just nested under another
				// level where macschema doesn't detect it
				s.Class.TypeProperties = append(s.Class.TypeProperties, schema.Property{
					Identifier: schema.Identifier{
						Name:        "clearColor",
						Description: "Returns a color object whose grayscale and alpha values are both 0.0.",
						Declaration: "@property(class, strong, readonly) NSColor *clearColor;",
						Deprecated:  false,
						TopicURL:    "https://developer.apple.com/documentation/appkit/nscolor/1527217-clearcolor?language=objc",
					},
					Type: schema.DataType{
						Name:  "NSColor",
						IsPtr: true,
					},
					Attrs: map[string]interface{}{
						"class":    true,
						"readonly": true,
						"strong":   true,
					},
				})
				return nil
			}),
			loadFile("api/appkit/nstextview.objc.json").Then(func(s *schema.Schema) error {
				s.Class.InstanceProperties = append(s.Class.InstanceProperties, schema.Property{
					Identifier: schema.Identifier{
						Name: "font",
					},
					Type: schema.DataType{Name: "NSFont", IsPtr: true},
				})
				return nil
			}),
			loadFile("api/appkit/nsview.objc.json").Then(filterProps(func(p schema.Property) bool {
				// only available on macOS 11+, causing build errors on GitHub
				return p.Name != "safeAreaRect"
			})).Then(func(s *schema.Schema) error {
				s.Class.InstanceProperties = append(s.Class.InstanceProperties, schema.Property{
					Identifier: schema.Identifier{
						Name: "backgroundColor",
					},
					Type: schema.DataType{Name: "NSColor", IsPtr: true},
				})
				return nil
			}),
		}},

		{"avcore", []schemaLoader{
			loadFile("api/avfoundation/avplayeractionatitemend.objc.json"),
			loadFile("api/avfoundation/avplayerstatus.objc.json"),
			loadFile("api/avfoundation/avasset.objc.json"),
			loadFile("api/avfoundation/avplayeritem.objc.json"),
			loadFile("api/avfoundation/avplayer.objc.json"),
			loadFile("api/avfoundation/avqueueplayer.objc.json"),
		}},

		{"mediaplayer", []schemaLoader{
			loadFile("api/mediaplayer/mpnowplayinginfocenter.objc.json"),
			loadFile("api/mediaplayer/mpremotecommandevent.objc.json"),
			loadFile("api/mediaplayer/mpchangeplaybackpositioncommandevent.objc.json"),
			loadFile("api/mediaplayer/mpchangelanguageoptioncommandevent.objc.json"),
			loadFile("api/mediaplayer/mpchangeplaybackratecommandevent.objc.json"),
			loadFile("api/mediaplayer/mpchangerepeatmodecommandevent.objc.json"),
			loadFile("api/mediaplayer/mpchangeshufflemodecommandevent.objc.json"),
			loadFile("api/mediaplayer/mpseekcommandevent.objc.json"),
			loadFile("api/mediaplayer/mpskipintervalcommandevent.objc.json"),

			loadFile("api/mediaplayer/mpremotecommand.objc.json"),
			loadFile("api/mediaplayer/mpchangeshufflemodecommand.objc.json"),
			loadFile("api/mediaplayer/mpchangerepeatmodecommand.objc.json"),
			loadFile("api/mediaplayer/mpchangeplaybackpositioncommand.objc.json"),
			loadFile("api/mediaplayer/mpchangeplaybackratecommand.objc.json"),
			loadFile("api/mediaplayer/mpskipintervalcommand.objc.json"),
			loadFile("api/mediaplayer/mpfeedbackcommand.objc.json"),
			loadFile("api/mediaplayer/mpratingcommand.objc.json"),
			loadFile("api/mediaplayer/mpremotecommandcenter.objc.json"),
		}},

		{"webkit", []schemaLoader{
			loadFile("api/webkit/wknavigation.objc.json"),
			loadFile("api/webkit/wkuserscript.objc.json"),
			loadFile("api/webkit/wkwebview.objc.json").Then(filterProps(func(p schema.Property) bool {
				return p.Name != "pageZoom"
			})),
			loadFile("api/webkit/wkwebviewconfiguration.objc.json"),
			loadFile("api/webkit/wkpreferences.objc.json").Then(func(s *schema.Schema) error {
				s.Class.InstanceMethods = append(s.Class.InstanceMethods, schema.Method{
					Identifier: schema.Identifier{
						Name: "setValue:forKey:",
					},
					Return: schema.DataType{Name: "void"},
					Args: []schema.Arg{
						{Name: "value", Type: schema.DataType{Name: "id"}},
						{Name: "key", Type: schema.DataType{Name: "NSString", IsPtr: true}},
					},
				})
				return nil
			}),
		}},
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
