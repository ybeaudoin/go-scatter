/*===== Copyright 2016, Webpraxis Consulting Ltd. - ALL RIGHTS RESERVED - Email: webpraxis@gmail.com ===========================
 *  Package scatter:
 *      import "scatter"
 *  Overview:
 *      package for creating an animated GIF of a rotating 3D scatter plot.
 *  Type:
 *      OptSet
 *          Structure for specifying the plot parameters
 *  Function:
 *      Plot(options OptSet, dataFile, gifFile string)
 *          Creates an animated GIF of a rotating 3D scatter plot.
 *  History:
 *      v1.0.0 - November 5, 2016 - Original release.
 *============================================================================================================================*/
package scatter

import(
    "bitbucket.org/binet/go-gnuplot/pkg/gnuplot"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
    "runtime"
    "strings"
    "time"
)
//Exported ---------------------------------------------------------------------------------------------------------------------
type(
    OptSet struct {
        TITLE      string  //centered plot title
        XLABEL     string  //x-axis label
        YLABEL     string  //y-axis label
        ZLABEL     string  //z-axis label (will be rotated 90 degrees)
        XRANGE     string  //x-axis range
        YRANGE     string  //y-axis range
        ZRANGE     string  //z-axis range
        XYPLANE    float64 //position at which the xy plane intersects the z-axis
        COLUMNS    string  //data columns to plot
        BGCOLOR    string  //background color specified as an hex string prefixed with the character "x"
        PTCOLOR    string  //data point color
        PLOTDELAY  int     //delay between display of successive images in units of 1/100 second
        PLOTROT    int     //rotation angle between display of successive images in degrees
        PLOTHEIGHT int     //plot height in pixels
        PLOTWIDTH  int     //plot width in pixels
        FONT       string  //font name & comma-separated optional size
    }
)

func Plot(options OptSet, dataFile, gifFile string) {
/*         Purpose : Creates an animated GIF of a rotating 3D scatter plot.
 *       Arguments : options  = plot options.
 *                   dataFile = filename for the R^3 plot data.
 *                   gifFile  = filename for the resulting GIF plot.
 *         Returns : None.
 * Externals -  In : None.
 * Externals - Out : None.
 *       Functions : halt
 *         Remarks : None.
 *         History : v1.0.0 - November 5, 2016 - Original release.
 */
    if dataFile == "" { halt("the data file was not specified") }
    if gifFile  == "" { halt("the GIF file was not specified") }

    //Create the temporary plot rotation command file
    refTemp, err := ioutil.TempFile("", "scatter_")
    if err != nil { halt("ioutil.TempFile - " + err.Error()) }
    rotCmds := filepath.ToSlash(refTemp.Name())
    writer, err := os.Create(rotCmds)
    if err != nil { halt("os.Create - " + err.Error()) }
    defer writer.Close()
    fmt.Fprintln(writer, "frame_count = frame_count + 1")
    fmt.Fprintf(writer,  "frame_title = sprintf( \"%s\\n( Rotation Angle = %%i%%c )\", zrot, 176 )\n", options.TITLE)
    fmt.Fprintf(writer,  "set title frame_title offset 0,1,0 font \"%s\"\n", options.FONT)
    fmt.Fprintln(writer, "set view xrot,zrot")
    fmt.Fprintln(writer, "replot")
    fmt.Fprintf(writer,  "zrot = ( zrot + %d ) %% 360\n", options.PLOTROT)
    fmt.Fprintf(writer,  "if( frame_count < %d ) reread\n", int(360. / float64(options.PLOTROT)))
    if err = writer.Sync();  err != nil { halt("writer.Sync - " + err.Error()) }
    if err = writer.Close(); err != nil { halt("writer.Close - " + err.Error()) }
    //Create the plot commands
    plotCmds := []string{
                 "unset key",
                 "set hidden3d",
                 `set border 4095 back lt 0 lc rgb "gray40"`,
                 "set grid x y z back",
                 fmt.Sprintf("set xyplane at %f", options.XYPLANE),
                 fmt.Sprintf("set xrange %s noreverse nowriteback", options.XRANGE),
                 fmt.Sprintf("set yrange %s noreverse nowriteback", options.YRANGE),
                 fmt.Sprintf("set zrange %s noreverse nowriteback", options.ZRANGE),
                 "set lmargin 7",
                 fmt.Sprintf(`set xlabel "%s" font "%s"`, options.XLABEL, options.FONT),
                 fmt.Sprintf(`set ylabel "%s" font "%s"`, options.YLABEL, options.FONT),
                 fmt.Sprintf(`set zlabel "%s" font "%s" offset -4,0,0 rotate by 90`, options.ZLABEL, options.FONT),
                 fmt.Sprintf(`set tics font "%s"`, options.FONT),
                 "xrot = 60",
                 fmt.Sprintf("zrot = %d", options.PLOTROT),
                 "set terminal unknown",
                 fmt.Sprintf(`splot "%s" using %s with points lc rgb "%s"`, dataFile, options.COLUMNS, options.PTCOLOR),
                 fmt.Sprintf("set terminal gif size %d,%d animate delay %d loop 0 nooptimize %s",
                             options.PLOTWIDTH, options.PLOTHEIGHT, options.PLOTDELAY, options.BGCOLOR),
                 fmt.Sprintf(`set output "%s"`, gifFile),
                 "frame_count = 0",
                 fmt.Sprintf(`load "%s"`, rotCmds),
                 "quit" }
    //Send the commands to gnuplot
    plotter, err := gnuplot.NewPlotter("", false, false)
    if err != nil { halt("gnuplot.NewPlotter - " + err.Error()) }
    for _, v := range plotCmds { plotter.CheckedCmd("%s", v) }
    plotter.Close()
    //Delete the temp file
    err = os.Remove(rotCmds)
    for err != nil && strings.Contains(err.Error(), "used by another process") {
        time.Sleep(time.Millisecond)
        err = os.Remove(rotCmds)
    }
    if err != nil { halt("os.Remove - " + err.Error()) }
} //end func Plot
//Private ----------------------------------------------------------------------------------------------------------------------
func halt(msg string) {
    pc, _, _, ok := runtime.Caller(1)
    details      := runtime.FuncForPC(pc)
    if ok && details != nil {
        log.Fatalln(fmt.Sprintf("\a%s: %s", details.Name(), msg))
    }
    log.Fatalln("\aoctree: FATAL ERROR!")
} //end func halt
//===== Copyright (c) 2016 Yves Beaudoin - All rights reserved - MIT LICENSE (MIT) - Email: webpraxis@gmail.com ================
//end of Package scatter
