# scatter demo

A demo involving all the package's function where we
* save 100 random points from the unit cube to a file and
* create the corresponding rotating scatter plot.

```go
package main

import (
    "fmt"
    "github.com/ybeaudoin/go-scatter"
    "log"
    "math/rand"
    "os"
    "time"
)
func main() {
    const(
        dataFile = "demo.dat"
        gifFile  = "demo.gif"
        numPts   = 100
    )
    var opts     scatter.OptSet

    //Pick random points from the unit cube & save to file
    writer, err := os.Create(dataFile)
    if err != nil { log.Fatalln(err) }
    defer writer.Close()

    rand.Seed(time.Now().UnixNano())
    for ptNo := 1; ptNo <= numPts; ptNo++ {
        fmt.Fprintf(writer, "%v\t%v\t%v\n", rand.Float64(), rand.Float64(), rand.Float64())
    }
    if err = writer.Sync();  err != nil { log.Fatalln(err) }
    if err = writer.Close(); err != nil { log.Fatalln(err) }
    fmt.Println("\n- data saved.")
    //Create the scatter plot
    opts.TITLE      = "DEMO"
    opts.XLABEL     = "X-AXIS"
    opts.YLABEL     = "Y-AXIS"
    opts.ZLABEL     = "Z-AXIS"
    opts.XRANGE     = "[0:1]"
    opts.YRANGE     = "[0:1]"
    opts.ZRANGE     = "[0:1]"
    opts.XYPLANE    = 0.
    opts.COLUMNS    = "1:2:3"
    opts.BGCOLOR    = "xe5e5e5"
    opts.PTCOLOR    = "red"
    opts.PLOTDELAY  = 100
    opts.PLOTROT    = 6
    opts.PLOTHEIGHT = 500
    opts.PLOTWIDTH  = 500
    opts.FONT       = "Garamond,10"
    scatter.Plot(opts, dataFile, gifFile)
    fmt.Println("\n- plot created.")
}
```















