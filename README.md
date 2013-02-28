## Terminal progress bar for Go  

Simple progress bar for console programms. Fork of https://github.com/chlu/pb with usage of channels and refresh of time when progress stalls.


### Installation
```
go get github.com/chlu/pb
```   

### Usage   
```Go
package main

import (
	"github.com/chlu/pb"
	"time"
)

func main() {
	count := 100000
	bar := pb.StartNew(count)
	for i := 0; i < count; i++ {
		bar.Increment()
		time.Sleep(time.Millisecond)
	}
	bar.FinishPrint("The End!")
}
```   
Result will be like this:
```
> go run test.go
37158 / 100000 [================>_______________________________] 37.16% 1m11s
```


More functions?  
```Go  
// create bar
bar := pb.New(count)

// refresh info every second (default 200ms)
bar.RefreshRate = time.Second

// show percents (by default already true)
bar.ShowPercent = true

// show bar (by default already true)
bar.ShowBar = true

// no need counters
bar.ShowCounters = false

// show "time left"
bar.ShowTimeLeft = true

// and start
bar.Start()
```    

Not like the looks?
```Go
// insert before usage
pb.BarStart = "<"
pb.BarEnd   = ">"
pb.Empty    = " "
pb.Current  = "-"
pb.CurrentN = "."
```
