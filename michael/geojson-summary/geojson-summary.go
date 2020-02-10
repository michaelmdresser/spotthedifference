package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/geo/s2"
	geojson "github.com/paulmach/go.geojson"
)

// getFeatureCollectionFromFile creates a FeatureCollection struct from a given file
func getFeatureCollectionFromFile(file string) (*geojson.FeatureCollection, error) {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %s", file, err)
	}

	fc, err := geojson.UnmarshalFeatureCollection(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal featurecollection from file %s: %s", file, err)
	}

	return fc, nil
}

// getFeatureCollectionsFromFiles creates an array of FeatureCollection structs from a given list of files
func getFeatureCollectionsFromFiles(files []string) ([]*geojson.FeatureCollection, error) {
	unmarshalledGeojsons := make([]*geojson.FeatureCollection, len(files))

	for i, file := range files {
		fc, err := getFeatureCollectionFromFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to get featurecollection from file %s: %s", file, err)
		}

		unmarshalledGeojsons[i] = fc
	}

	return unmarshalledGeojsons, nil
}

// latLngFromGeojsonPoint creates an s2.LatLng struct from a geojson.Feature, expected to be of type Point
func latLngFromGeojsonPoint(f *geojson.Feature) s2.LatLng {
	// TODO: check if point
	long := f.Geometry.Point[0]
	lat := f.Geometry.Point[1]
	// fmt.Printf("lat, long: %f, %f\n", lat, long)

	latlong := s2.LatLngFromDegrees(lat, long)

	return latlong
}

// PointFromGeojsonFeature creates a s2.Point struct from a geojson.Feature, expected to be of type point
func PointFromGeojsonFeature(f *geojson.Feature) s2.Point {
	return s2.PointFromLatLng(latLngFromGeojsonPoint(f))
}

// getBoundingLoop returns a s2.Loop struct for checking membership in a region. It is currently hardcoded because I need to think of a good way to programmatically define the loop in a user-friendly way (also because I just wanted to get it working).
func getBoundingLoop() *s2.Loop {
	p1 := s2.PointFromLatLng(s2.LatLngFromDegrees(36.88, 44.324))
	p2 := s2.PointFromLatLng(s2.LatLngFromDegrees(36.5967, 44.161))
	p3 := s2.PointFromLatLng(s2.LatLngFromDegrees(36.641, 44.687))

	return s2.LoopFromPoints([]s2.Point{p1, p2, p3})
}

// writeFeatureCollectionToRegularFile writes a provided geojson.FeatureCollection to the given filename
func writeFeatureCollectionToRegularFile(fc *geojson.FeatureCollection, file string) {
	bytes, err := fc.MarshalJSON()
	if err != nil {
		// TODO: don't panic
		panic(err)
	}

	err = ioutil.WriteFile(file, bytes, 0644)
	if err != nil {
		// TODO: don't panic
		panic(err)
	}

}

// filterFeatureCollectionToRegion builds a new geojson.FeatureCollection of only Features (expected to be points) in fc that are within the provided region.
// TODO: generalize region to not just be a Loop
func filterFeatureCollectionToRegion(fc *geojson.FeatureCollection, region *s2.Loop) *geojson.FeatureCollection {
	filteredFC := geojson.NewFeatureCollection()
	filteredFC.Type = fc.Type

	for _, feature := range fc.Features {
		point := PointFromGeojsonFeature(feature)

		if region.ContainsPoint(point) {
			filteredFC.Features = append(filteredFC.Features, feature)
		}
	}

	return filteredFC
}

// filterGeojsonFilesToRegion takes a list of files and builds up a "master" geojson.FeatureCollection by merging the results of filterFeatureCollectionToRegion for each file and then outputs it to newFile.
func filterGeojsonFilesToRegion(files []string, newFile string, region *s2.Loop) {
	fullFilteredFC := geojson.NewFeatureCollection()

	for _, file := range files {
		log.Printf("processing %s", file)
		fc, err := getFeatureCollectionFromFile(file)
		if err != nil {
			// TODO: don't panic
			panic(err)
		}

		filteredFC := filterFeatureCollectionToRegion(fc, region)
		log.Printf("%s had %d features before and %d after filtering",
			file,
			len(fc.Features),
			len(filteredFC.Features))

		fullFilteredFC.Features = append(fullFilteredFC.Features, filteredFC.Features...)
	}

	log.Printf("full filtered collection has %d features", len(fullFilteredFC.Features))
	log.Printf("writing full filtered collection to %s", newFile)
	writeFeatureCollectionToRegularFile(fullFilteredFC, newFile)
}

// filterGeojsonFileToRegion takes a filename and outputs the results of filterFeatureCollectionToRegion to newFile.
func filterGeojsonFileToRegion(file, newFile string, region *s2.Loop) {
	log.Printf("processing %s as geojson and outputting to file %s", file, newFile)

	fc, err := getFeatureCollectionFromFile(file)
	if err != nil {
		// TODO: don't panic
		panic(err)
	}
	log.Printf("%s is a feature collection with %d features", file, len(fc.Features))

	filteredFC := filterFeatureCollectionToRegion(fc, getBoundingLoop())

	log.Printf("filtered collection has %d features", len(filteredFC.Features))

	log.Printf("writing filtered feature collection to %s", newFile)
	writeFeatureCollectionToRegularFile(filteredFC, newFile)
}

// generateFilteredFilename is for auto-generating filenames for filtered geojsons if outputting each filtered geojson separates as opposed to building a single master geojson. NOT IN USE.
func generateFilteredFilename(original string) string {
	return original[:len(original)-8] + "_filtered.geojson"
}

func main() {
	args := os.Args[1:]
	output := args[0]
	files := args[1:]

	fmt.Println(generateFilteredFilename(files[0]))

	region := getBoundingLoop()

	filterGeojsonFilesToRegion(files, output, region)
}
