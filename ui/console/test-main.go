/*
A presentation of the tview package, implemented with tview.

Navigation

The presentation will advance to the next slide when the primitive demonstrated
in the current slide is left (usually by hitting Enter or Escape). Additionally,
the following shortcuts can be used:

  - Ctrl-N: Jump to next slide
  - Ctrl-P: Jump to previous slide
*/
package console

import (
	"fmt"
	"github.com/integrii/flaggy"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/sirupsen/logrus"
	"strconv"
	"sync"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

//const version = "0.3a"

type ConnectionServerConfig struct {
	mailboxServers       string
	ldapServers          []string
	serverTypePreference string
	adminUsername        string
	adminPassword        string
	adminPort            string
	adminProtocol        string
	socksServerString    string
	adminAuthToken       string
	allServers           cmap.ConcurrentMap
	skipSSLChecks        bool
	useSocks5Proxy       bool
	printedSSLSkipNotice bool
}

type RuntimeConfig struct {
}

type DataTracker struct {
}

var log = logrus.New()
var problemAccounts = cmap.New()
var userServerMapping = cmap.New()
var serverTracker = cmap.New()

var serverStats = false
var auditHeaderWritten = false
var mode = "file"

var inputFile = ""
var password = ""
var module = ""
var task = ""
var auditType = ""
var auditValue = ""
var threads = ""

var logTimestamp string
var messageScanLimit string
var scanLimit int
var alreadySeen []string
var email string
var accountsToTest []string
var accountsToSkip []string
var domains []string
var maxThreads int
var wg sync.WaitGroup
var saveToConfig bool
var socks5 bool
var SkipSSL bool
var sshKeyPath string

func init() {

	flaggy.SetName("Zimbra Experience (zmxp)")
	flaggy.SetDescription("A utility to gather information about the quality of the end user Zimbra Experience.")
	flaggy.DefaultParser.ShowHelpOnUnexpected = true
	flaggy.DefaultParser.AdditionalHelpPrepend = "http://github.com/zimbra/zmxp"
	flaggy.String(&threads, "th", "threads", "The Max number of threads to spawn for checking via soap and screenshots.")
	flaggy.String(&inputFile, "f", "file", "Provide an input file with a list of accounts to test (--file=myfile.txt). The file should consist of user email addresses, one per line.")
	flaggy.Bool(&serverStats, "ss", "server-stats", "Gather server stats for success vs failure.")
	flaggy.String(&password, "p", "password", "Manually define a password.")
	flaggy.String(&module, "m", "module", "The module to use.")
	flaggy.String(&task, "t", "task", "The task to to perform.")
	flaggy.String(&auditType, "at", "audit-type", "When 'mail-audit' is specified as the task, this is the audit type.")
	flaggy.String(&auditValue, "av", "audit-value", "When 'mail-audit' is specified as the task, this is the audit value.")
	flaggy.String(&messageScanLimit, "sl", "scan-limit", "When 'mail-audit' is specified as the task, this is the maximum amount of messages which will be scanned per account.")
	flaggy.String(&sshKeyPath, "key", "ssh-key", "When using SSH, you must provide an SSH key.")
	flaggy.SetVersion("zmxp " + version)
	flaggy.Parse()
}
var (
	commit      string
	version     = "unversioned"
	date        string
	buildSource = "unknown"

	configFlag    = false
	debuggingFlag = false
	composeFiles  []string
)
// Slide is a function which returns the slide's main primitive and its title.
// It receives a "nextSlide" function which can be called to advance the
// presentation to the next slide.
type Slide func(nextSlide func()) (title string, content tview.Primitive)

// The application.
var app = tview.NewApplication()

// Starting point for the presentation.
func main() {
	// The presentation slides.
	slides := []Slide{
		Cover,
		Introduction,
		HelloWorld,
		InputField,
		Form,
		TextView1,
		TextView2,
		Table,
		TreeView,
		Flex,
		Grid,
		Colors,
		End,
	}

	// The bottom row has some info on where we are.
	info := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false)
	// The bottom row has some info on where we are.
	topBar := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false)
	topGrid := tview.NewGrid()
	firstStat := tview.NewTex
	topGrid.AddItem()
	topBar.SetText("mtlp00009: 50%\tmtlp00009: 50%\t")
	// Create the pages for all slides.
	currentSlide := 0
	info.Highlight(strconv.Itoa(currentSlide))
	pages := tview.NewPages()
	previousSlide := func() {
		currentSlide = (currentSlide - 1 + len(slides)) % len(slides)
		info.Highlight(strconv.Itoa(currentSlide)).
			ScrollToHighlight()
		pages.SwitchToPage(strconv.Itoa(currentSlide))
	}
	nextSlide := func() {
		currentSlide = (currentSlide + 1) % len(slides)
		info.Highlight(strconv.Itoa(currentSlide)).
			ScrollToHighlight()
		pages.SwitchToPage(strconv.Itoa(currentSlide))
	}
	for index, slide := range slides {
		title, primitive := slide(nextSlide)
		pages.AddPage(strconv.Itoa(index), primitive, true, index == currentSlide)
		fmt.Fprintf(info, `%d ["%d"][darkcyan]%s[white][""]  `, index+1, index, title)
	}

	// Create the main layout.
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(topBar, 1, 1, false).
		AddItem(pages, 0, 1, true).
		AddItem(info, 1, 1, false)

	// Shortcuts to navigate the slides.
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlN {
			nextSlide()
		} else if event.Key() == tcell.KeyCtrlP {
			previousSlide()
		}
		return event
	})

	// Start the application.
	if err := app.SetRoot(layout, true).Run(); err != nil {
		panic(err)
	}
}
