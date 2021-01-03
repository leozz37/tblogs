package app

import (
	"github.com/ezeoleaf/tblogs/cfg"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type slide func(nextSlide func()) (title string, content tview.Primitive)

// var app = tview.NewApplication()

// App contains the tview application and the layout for the display
type App struct {
	App    *tview.Application
	Layout *tview.Flex
	Config *cfg.Config
}

// Setup returns an instance of the application
func Setup() *App {

	c := cfg.Setup()

	var app *App

	pages, info := app.getPagesInfo()

	// Create the main layout.
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(pages, 0, 1, true).
		AddItem(info, 1, 1, false)

	var appTview *tview.Application

	// Shortcuts to navigate the slides.
	appTview.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch ek := event.Key(); ek {
		case tcell.KeyCtrlH:
			goToSection(helpSection)
		case tcell.KeyCtrlB:
			goToSection(blogsSection)
		case tcell.KeyCtrlT:
			goToSection(homeSection)
		case tcell.KeyCtrlP:
			goToSection(savedPostsSection)
		case tcell.KeyCtrlX:
			goToSection(settingsSection)
		}
		return event
	})

	app.App = appTview
	app.Layout = layout
	app.Config = c

	return app
}

// Start launches the app
func (a *App) Start() {
	if err := a.App.SetRoot(a.Layout, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
