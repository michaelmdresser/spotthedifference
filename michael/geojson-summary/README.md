# geojson-summary
### Could use a different name, currently just does filtering

# What does it do?
CLI utility that takes a single output file and a list of input files. Input files are expected to be in GeoJSON format where each file is a FeatureCollection containing only points. The utility builds up a "master" FeatureCollection of all points in each file that are within a (currently hardcoded) bounded region on the globe. It then writes that to the filename provided.

# How do I run this?

1. Install (golang)[https://golang.org/doc/install]
2. `cd` to this directory
3. Run `go get github.com/paulmach/go.geojson`
4. Run `go get github.com/golang/geo`
5. Run `go build`
6. Run `./geojson-summary OUTPUTFILE INPUTFILE [INPUTFILE2 ...]
