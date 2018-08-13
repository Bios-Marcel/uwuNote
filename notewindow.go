package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	uuid "github.com/satori/go.uuid"
)

func createWindowForNote(file string, x int, y int) {
	//Error variable to be reused
	var gtkError error

	// Create a new toplevel window and connect it to the
	// "destroy" signal to exit the GTK main loop when it is destroyed.
	win, gtkError := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	panicOnError(gtkError)

	// The app isn't killable for now.
	/*win.Connect("destroy", func() {
		gtk.MainQuit()
	})*/

	win.SetTitle(file)

	newButton, gtkError := gtk.ButtonNew()
	panicOnError(gtkError)

	newButton.SetLabel("New")
	newButton.Connect("clicked", func() {
		fileName := uuid.Must(uuid.NewV4())
		newNotePath := notePath + string(os.PathSeparator) + fileName.String() + ".md"
		os.Create(newNotePath)
		createWindowForNote(newNotePath, x+20, y+20)
	})
	newButton.SetHExpand(false)

	deleteButton, gtkError := gtk.ButtonNew()
	panicOnError(gtkError)

	deleteButton.SetLabel("Delete")
	deleteButton.Connect("clicked", func() {
		os.Remove(file)
		win.Destroy()
	})
	deleteButton.SetHExpand(false)
	deleteButton.SetHAlign(gtk.ALIGN_END)

	topBar, gtkError := gtk.HeaderBarNew()
	panicOnError(gtkError)

	topBar.PackStart(newButton)
	topBar.PackEnd(deleteButton)

	var hAdjustment, vAdjustment *gtk.Adjustment
	textViewScrollPane, gtkError := gtk.ScrolledWindowNew(hAdjustment, vAdjustment)
	panicOnError(gtkError)

	textView, gtkError := gtk.TextViewNew()
	panicOnError(gtkError)

	textView.SetVExpand(true)
	textView.SetHExpand(true)

	//TODO Currently saving is triggered manualy by pressing Ctrl + S, but later on it is supposedto be saving automatically.

	textView.Connect("key_release_event", func(widget *gtk.TextView, event *gdk.Event) {
		//Subtract default modifiers according to:
		//https://developer.gnome.org/gtk3/stable/checklist-modifiers.html
		//modifiers := gtk.AcceleratorGetDefaultModMask()

		keyEvent := gdk.EventKeyNewFromEvent(event)
		if keyEvent.KeyVal() == gdk.KEY_s {
			if (keyEvent.State() & gdk.GDK_CONTROL_MASK) == gdk.GDK_CONTROL_MASK {
				//TODO Following errors on which i now panic should probably inform the user about not being able to save.

				currentNoteBuffer, bufferError := textView.GetBuffer()
				if bufferError != nil {
					panic(bufferError)
				}

				iterStart, iterEnd := currentNoteBuffer.GetBounds()
				//TODO Check if I need the "Hidden chars"
				textToSave, textError := currentNoteBuffer.GetText(iterStart, iterEnd, false)
				if textError != nil {
					panic(textError)
				}
				writeError := ioutil.WriteFile(file, []byte(textToSave), os.ModeType)
				if writeError != nil {
					panic(writeError)
				} else {
					fmt.Println("Successfully saved content")
				}
			}
		}
	})

	//Wrapping the textView in a scrollpane, otherwise the window will expand instead
	textViewScrollPane.Add(textView)

	buffer, gtkError := textView.GetBuffer()
	panicOnError(gtkError)

	fileContent, _ := ioutil.ReadFile(file)
	buffer.SetText(string(fileContent))

	nodeLayout, gtkError := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	panicOnError(gtkError)

	nodeLayout.Add(textViewScrollPane)
	nodeLayout.SetVExpand(true)

	win.SetTitlebar(topBar)
	win.Add(nodeLayout)

	win.SetSkipTaskbarHint(true)
	win.SetSkipPagerHint(true)
	win.SetKeepBelow(true)
	win.Stick()

	win.Connect("window-state-event", func(window *gtk.Window, event *gdk.Event) {
		windowEvent := gdk.EventWindowStateNewFromEvent(event)
		newWindowState := windowEvent.NewWindowState()

		if (newWindowState & gdk.WINDOW_STATE_ICONIFIED) == gdk.WINDOW_STATE_ICONIFIED {
			window.Present()
		}
	})

	win.Move(x, y)
	win.SetDefaultSize(300, 350)

	// Recursively show all widgets contained in this window.
	win.ShowAll()
}