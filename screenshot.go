package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/chromedp"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func screenshot(addr string) {
	// Chromedp-Kontext erstellen (Browser-Interaktionen)
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Timeout-Kontext als Sicherheitsnetz
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// Die URL angeben, die die Screenshot-Interaktion steuert
	url := fmt.Sprintf("http://%s", addr)

	// Wird verwendet, um den Dateinamen des Downloads zu erfassen
	var downloadGUID string

	// Kanal zur Überwachung des Download-Fortschritts
	downloadComplete := make(chan bool)
	chromedp.ListenTarget(ctx, func(v interface{}) {
		if ev, ok := v.(*browser.EventDownloadProgress); ok {
			if ev.State == browser.DownloadProgressStateCompleted {
				downloadGUID = ev.GUID
				close(downloadComplete)
			}
		}
	})

	// Chromedp-Tasks ausführen: Navigation, Warten und Interaktion
	if err := chromedp.Run(ctx, chromedp.Tasks{
		browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllowAndName).
			WithDownloadPath(os.TempDir()).
			WithEventsEnabled(true),

		chromedp.Navigate(url),                             // URL aufrufen
		chromedp.WaitVisible(`#cytoscape-div`),             // Warten, bis das Element sichtbar ist
		chromedp.Click(`#saveGraph`, chromedp.NodeVisible), // "Save Graph"-Button klicken
	}); err != nil && !strings.Contains(err.Error(), "net::ERR_ABORTED") {
		// Ignoriert net::ERR_ABORTED, da Downloads diesen Fehler manchmal auslösen
		log.Fatal(err)
	}

	// Blockiert, bis der Download abgeschlossen ist
	<-downloadComplete

	// Download-Datei verschieben
	e := moveFile(fmt.Sprintf("%v/%v", os.TempDir(), downloadGUID), "./rover.svg")
	if e != nil {
		log.Fatal(e)
	}

	log.Println("Image generation complete.")
}

// Funktion zum Verschieben von Dateien, um plattformspezifische Probleme (z.B. Docker) zu vermeiden
func moveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	// Original-Datei löschen
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}
