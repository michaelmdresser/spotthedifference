# geojson-summary

This should probably be named `geojson-filtering`. Also, pretty much everything in here is deprecated in favor of using PostGIS (see the `postgis` folder, up one directory), which we discovered unfortunately late in the techincal process.

## What does it do?
This is a CLI tool for working with GeoJSON files of the specific type provided for this project (FeatureCollection containing only points). It does one of two things, depending on which bit of code is uncommented (it should just be a flag on the CLI but I didn't get there).
- If using `getBoundingLoop` and `filterGeojsonFilesToRegion`, the code takes a series of points hardcoded in the `getBoundingLoop` function (see [this link](https://godoc.org/github.com/golang/geo/s2#LoopFromPoints) for important info if editing) and filters all provided files to points that are just within the region defined in `getBoundingLoop`. It is slow, but does not run into memory constraints. Usage: `./geojson-summary OUTPUTFILE INPUTFILE [INPUTFILE2 ...]` for example: `./geojson-summary /mnt/d4d/testoutput.geojson /mnt/d4d/geojson/20181101030142_20181113030142_351541N-ACD-BETA/01Nov18S1B101586_02_13Nov18S1B102974_02_ACD_Arrivals.geojson /mnt/d4d/geojson/20181101030142_20181113030142_351541N-ACD-BETA/01Nov18S1B101586_02_13Nov18S1B102974_02_ACD_Arrivals.geojson`
- If using `bucketifyData`, the code turns the entire geographic region (viewed in 2D) into a grid of equal-sized rectangles (the number of which is definable as `latBuckets` and `longBuckets` for each dimension). Total bucket count is approximately equal to `latBuckets * longBuckets`. It then filters the whole dataset into each bucket, producing one GeoJSON file for each bucket. It is extremely slow because a lack of memory to load the entire dataset in means that we have to read each file one time for every bucket. Usage: `./geojson-summary OUTPUTDIR INPUTFILE [INPUTFILE2 ...]` similar to the other usage, just an output directory instead of file.


## How do I run this?

1. Install [golang](https://golang.org/doc/install)
2. `cd` to this directory
3. Run `go get github.com/paulmach/go.geojson`
4. Run `go get github.com/golang/geo`
5. Run `go build`
6. See "Usage" in the above section (make sure to rebuild if changing code)
