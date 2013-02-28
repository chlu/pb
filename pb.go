package pb

import (
	"fmt"
	"io"
	"strings"
	"time"
)

var (
	// Default refresh rate - 200ms
	DefaultRefreshRate = time.Millisecond * 200

	BarStart = "["
	BarEnd   = "]"
	Empty    = "_"
	Current  = "="
	CurrentN = ">"
)

// Create new progress bar object
func New(total int) *ProgressBar {
	return &ProgressBar{
		Total:        int64(total),
		RefreshRate:  DefaultRefreshRate,
		ShowPercent:  true,
		ShowCounters: true,
		ShowBar:      true,
		ShowTimeLeft: true,
		increment:    make(chan int64, 1),
		update:       make(chan int64, 1),
		finish:       make(chan bool, 1),
	}
}

// Create new object and start 
func StartNew(total int) (pb *ProgressBar) {
	pb = New(total)
	pb.Start()
	return
}

type Callback func(out string)

type ProgressBar struct {
	Total                                            int64
	RefreshRate                                      time.Duration
	ShowPercent, ShowCounters, ShowBar, ShowTimeLeft bool
	Output                                           io.Writer
	Callback                                         Callback
	NotPrint                                         bool

	increment chan int64
	update    chan int64
	finish    chan bool

	startTime time.Time
}

// Start print
func (pb *ProgressBar) Start() {
	pb.startTime = time.Now()
	go pb.writer()
}

// Increment current value
func (pb *ProgressBar) Increment() {
	pb.Add(1)
}

// Set current value
func (pb *ProgressBar) Set(current int64) {
	pb.update <- current
}

// Add to current value
func (pb *ProgressBar) Add(add int64) {
	pb.increment <- add
}

// End print
func (pb *ProgressBar) Finish() {
	pb.finish <- true

	if !pb.NotPrint {
		fmt.Println()
	}
}

// End print and write string 'str'
func (pb *ProgressBar) FinishPrint(str string) {
	pb.Finish()
	fmt.Println(str)
}

func (pb *ProgressBar) write(current int64) {
	width, _ := terminalWidth()
	var percentBox, countersBox, timeLeftBox, barBox, end, out string

	// percents
	if pb.ShowPercent {
		percent := float64(current) / (float64(pb.Total) / float64(100))
		percentBox = fmt.Sprintf(" %#.02f %% ", percent)
	}

	// counters
	if pb.ShowCounters {
		countersBox = fmt.Sprintf("%d / %d ", current, pb.Total)
	}

	// time left
	if pb.ShowTimeLeft && current > 0 {
		fromStart := time.Now().Sub(pb.startTime)
		perEntry := fromStart / time.Duration(current)
		left := time.Duration(pb.Total-current) * perEntry
		left = (left / time.Second) * time.Second
		if left > 0 {
			timeLeftBox = left.String()
		}
	}

	// bar
	if pb.ShowBar {
		size := width - len(countersBox+BarStart+BarEnd+percentBox+timeLeftBox)
		if size > 0 {
			curCount := int(float64(current) / (float64(pb.Total) / float64(size)))
			emptCount := size - curCount
			barBox = BarStart
			if emptCount < 0 {
				emptCount = 0
			}
			if curCount > size {
				curCount = size
			}
			if emptCount <= 0 {
				barBox += strings.Repeat(Current, curCount)
			} else if curCount > 0 {
				barBox += strings.Repeat(Current, curCount-1) + CurrentN
			}

			barBox += strings.Repeat(Empty, emptCount) + BarEnd
		}
	}

	// check len
	out = countersBox + barBox + percentBox + timeLeftBox
	if len(out) < width {
		end = strings.Repeat(" ", width-len(out))
	}

	out = countersBox + barBox + percentBox + timeLeftBox

	// and print!
	switch {
	case pb.Output != nil:
		fmt.Fprint(pb.Output, out+end)
	case pb.Callback != nil:
		pb.Callback(out + end)
	case !pb.NotPrint:
		fmt.Print("\r" + out + end)
	}
}

func (pb *ProgressBar) writer() {
	var current, i int64

	ticker := time.NewTicker(pb.RefreshRate)
	defer ticker.Stop()

	for {
		select {
		case <-pb.finish:
			return
		case i = <-pb.update:
			current = i
			pb.write(current)
		case i = <-pb.increment:
			current += i
			pb.write(current)
		case <-ticker.C:
			pb.write(current)
		}
	}
}

type window struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}
