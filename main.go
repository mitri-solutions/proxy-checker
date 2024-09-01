package main

import (
	"fmt"
	proxy "proxy-checker/internal"
	"strconv"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func Worker(input string, url string, wg *sync.WaitGroup, semaphore chan struct{}, updateProgress func(status bool)) {
	defer wg.Done()
	semaphore <- struct{}{}

	status := proxy.CheckProxy(input, url)
	updateProgress(status)

	// Release the semaphore
	<-semaphore
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Proxy Checker v1.0")

	inputTextArea := widget.NewMultiLineEntry()
	inputTextArea.SetPlaceHolder("proxy:port:username:password")

	successTextArea := widget.NewMultiLineEntry()
	successTextArea.SetPlaceHolder("proxy:port:username:password")

	errorTextArea := widget.NewMultiLineEntry()
	errorTextArea.SetPlaceHolder("proxy:port:username:password")

	successTextLabel := widget.NewLabel("Success")
	errorTextLabel := widget.NewLabel("Error")

	urlInput := widget.NewEntry()
	urlInput.SetPlaceHolder("https://www.google.com")
	urlInput.SetText("https://www.google.com")

	copyrightLabel := widget.NewLabel("Proxy Checker by @jamesngdev (t.me/jamesngdev)")
	copyrightLabel.Alignment = fyne.TextAlignCenter

	// concurrencyLabel := widget.NewLabel("Concurrency")
	concurrencyInput := widget.NewEntry()
	concurrencyInput.SetText("3")
	concurrencyInput.SetPlaceHolder("Concurrency")
	concurrencyInput.Validator = func(s string) error {
		if _, err := strconv.Atoi(s); err != nil && s != "" {
			return fmt.Errorf("please enter a valid number")
		}
		return nil
	}

	proxyLabel := widget.NewLabel("Proxy")
	progressBar := widget.NewProgressBar()

	var submitButton *widget.Button // Declare the button variable

	successed := []string{}
	failed := []string{}

	submitButton = widget.NewButton("Check Proxy", func() {
		submitButton.Disable()
		submitButton.SetText("Checking...")

		go func() {
			proxy := inputTextArea.Text
			url := urlInput.Text

			//split the proxy line by line
			proxies := strings.Split(proxy, "\n")

			// Convert to int
			maxConcurrentWorkers, err := strconv.Atoi(concurrencyInput.Text)

			if err != nil {
				myApp.SendNotification(&fyne.Notification{Title: "Proxy Checker", Content: "Please enter a valid number"})
				return
			}

			semaphore := make(chan struct{}, maxConcurrentWorkers)
			var wg sync.WaitGroup
			numTasks := len(proxies)
			proxyLabel.SetText(fmt.Sprintf("Proxy (%d)", numTasks))
			progressBar.SetValue(0)
			wg.Add(numTasks)

			for _, p := range proxies {
				go Worker(p, url, &wg, semaphore, func(status bool) {
					progressBar.SetValue(progressBar.Value + 1.0/float64(numTasks))

					if status {
						successed = append(successed, p)
					} else {
						failed = append(failed, p)
					}

					successTextArea.SetText(strings.Join(successed, "\n"))
					errorTextArea.SetText(strings.Join(failed, "\n"))

					successTextLabel.SetText(fmt.Sprintf("Success (%d)", len(successed)))
					errorTextLabel.SetText(fmt.Sprintf("Error (%d)", len(failed)))
				})
			}

			wg.Wait()
			submitButton.Enable()
			submitButton.SetText("Check Proxy")
			myApp.SendNotification(&fyne.Notification{Title: "Proxy Checker", Content: "All proxies have been checked!"})
		}()
	})

	// Create a function
	myWindow.SetContent(container.NewVBox(
		container.NewGridWithColumns(2,
			container.NewVBox(proxyLabel, inputTextArea),
			container.NewVBox(
				widget.NewLabel("Config"),
				urlInput,
				concurrencyInput,
			),
		),

		submitButton,
		container.NewGridWithColumns(2, container.NewVBox(successTextLabel, successTextArea), container.NewVBox(errorTextLabel, errorTextArea)),
		progressBar,
		copyrightLabel,
	))

	myWindow.Resize(fyne.NewSize(800, 300))
	myWindow.ShowAndRun()
}
