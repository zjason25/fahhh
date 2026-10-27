package main

/*
#cgo darwin LDFLAGS: -framework ApplicationServices
#include <ApplicationServices/ApplicationServices.h>

static inline int checkAndRequestAccessibility() {
    const void *keys[] = { kAXTrustedCheckOptionPrompt };
    const void *values[] = { kCFBooleanTrue };
    CFDictionaryRef options = CFDictionaryCreate(
        kCFAllocatorDefault,
        keys,
        values,
        1,
        &kCFTypeDictionaryKeyCallBacks,
        &kCFTypeDictionaryValueCallBacks
    );
    Boolean trusted = AXIsProcessTrustedWithOptions(options);
    CFRelease(options);
    return trusted ? 1 : 0;
}
*/
import "C"

import (
	"bytes"
	_ "embed"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	hook "github.com/robotn/gohook"
)

// Bake mp3 of your choice directly into the compiled executable
//
//go:embed fah.mp3
var audioData []byte

func main() {
	// 1. Request macOS Accessibility permission natively
	if C.checkAndRequestAccessibility() == 0 {
		log.Println("Accessibility permission required.")
		log.Println("Please toggle permission ON in System Settings -> Privacy & Security -> Accessibility, then re-open this app.")
		return
	}

	// 2. Decode embedded audio from memory
	streamer, format, err := mp3.Decode(io.NopCloser(bytes.NewReader(audioData)))
	if err != nil {
		log.Fatalf("Failed to decode embedded audio: %v", err)
	}
	defer streamer.Close()

	soundBuffer := beep.NewBuffer(format)
	soundBuffer.Append(streamer)

	// 3. Initialize audio engine
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Millisecond*50))
	if err != nil {
		log.Fatalf("Failed to initialize speaker: %v", err)
	}

	// 4. Start OS keyboard hook
	evChan := hook.Start()
	defer hook.End()

	// 5. Setup clean shutdown listener (Ctrl+C / Kill)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("\nStopping keyboard sound hook...")
		hook.End()
		os.Exit(0)
	}()

	// log.Println("Keyboard sound app active. Type anywhere...")

	// 6. Grace period: Ignore initial keystrokes (like Enter key used to start the app)
	startupGraceEnd := time.Now().Add(300 * time.Millisecond)

	for ev := range evChan {
		if time.Now().Before(startupGraceEnd) {
			continue
		}

		if ev.Kind == hook.KeyDown {
			speaker.Play(soundBuffer.Streamer(0, soundBuffer.Len()))
		}
	}
}
