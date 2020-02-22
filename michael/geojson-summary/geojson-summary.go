package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

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

func getFeatureCollectionFromFiles(files []string) (*geojson.FeatureCollection, error) {
	fcs, err := getFeatureCollectionsFromFiles(files)
	if err != nil {
		return nil, err
	}

	masterFc := geojson.NewFeatureCollection()

	for _, fc := range fcs {
		masterFc.Features = append(masterFc.Features, fc.Features...)
	}

	return masterFc, nil
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

	filteredFC := filterFeatureCollectionToRegion(fc, region)

	log.Printf("filtered collection has %d features", len(filteredFC.Features))

	log.Printf("writing filtered feature collection to %s", newFile)
	writeFeatureCollectionToRegularFile(filteredFC, newFile)
}

func min(n1, n2 float64) float64 {
	if n1 < n2 {
		return n1
	}
	return n2
}

func max(n1, n2 float64) float64 {
	if n1 > n2 {
		return n1
	}
	return n2
}

func getMinMaxLatLonBoundsForFeatureCollection(fc *geojson.FeatureCollection) (minLat, minLong, maxLat, maxLong float64) {
	for i, f := range fc.Features {
		featLong := f.Geometry.Point[0]
		featLat := f.Geometry.Point[1]
		if i == 0 {
			minLat = featLat
			minLong = featLong
			maxLat = featLat
			maxLong = featLong
		} else {
			minLat = min(minLat, featLat)
			maxLat = max(maxLat, featLat)
			minLong = min(minLong, featLong)
			maxLong = max(maxLong, featLong)
		}
	}

	return minLat, minLong, maxLat, maxLong
}

func getMinMaxLanLonBoundsMultiFile(files []string) (minLat, minLong, maxLat, maxLong float64) {
	log.Printf("getting min and max latlong pairs from %s", files)

	for i, file := range files {
		fc, err := getFeatureCollectionFromFile(file)
		if err != nil {
			// TODO: don't panic
			panic(err)
		}

		localMinLat, localMinLong, localMaxLat, localMaxLong := getMinMaxLatLonBoundsForFeatureCollection(fc)
		if i == 0 {
			minLat = localMinLat
			minLong = localMinLong
			maxLat = localMaxLat
			maxLong = localMaxLong
		} else {
			minLat = min(minLat, localMinLat)
			minLong = min(minLong, localMinLong)
			maxLat = max(maxLat, localMaxLat)
			maxLong = min(maxLong, localMaxLong)
		}
	}

	return minLat, minLong, maxLat, maxLong
}

func bucketRegionsFromMinMax(minLat, minLong, maxLat, maxLong float64, latBuckets, longBuckets int) []*s2.Loop {
	var regions []*s2.Loop

	latStepSize := (maxLat - minLat) / float64(latBuckets)
	longStepSize := (maxLong - minLong) / float64(longBuckets)

	topLeftLat := minLat
	topLeftLong := minLong
	bottomRightLat := minLat + latStepSize
	bottomRightLong := minLong + longStepSize

	for bottomRightLong <= maxLong {
		for bottomRightLat <= maxLat {
			topRightLat := bottomRightLat
			topRightLong := topLeftLong
			bottomLeftLat := topLeftLat
			bottomLeftLong := bottomRightLong

			topLeftPoint := s2.PointFromLatLng(s2.LatLngFromDegrees(topLeftLat, topLeftLong))
			topRightPoint := s2.PointFromLatLng(s2.LatLngFromDegrees(topRightLat, topRightLong))
			bottomLeftPoint := s2.PointFromLatLng(s2.LatLngFromDegrees(bottomLeftLat, bottomLeftLong))
			bottomRightPoint := s2.PointFromLatLng(s2.LatLngFromDegrees(bottomRightLat, bottomRightLong))

			regions = append(regions, s2.LoopFromPoints([]s2.Point{topLeftPoint, bottomLeftPoint, bottomRightPoint, topRightPoint}))

			topLeftLat += latStepSize
			bottomRightLat += latStepSize
		}
		topLeftLat = minLat
		bottomRightLat = minLat + latStepSize
		topLeftLong += longStepSize
		bottomRightLong += longStepSize
	}

	return regions
}

func bucketifyData(files []string, outputDir string, latBuckets, longBuckets int) {
	if strings.Contains(outputDir[len(outputDir)-1:], "/") {
		outputDir = outputDir[:len(outputDir)-1]
	}
	log.Printf("bucketifying the following files into %d latBuckets and %d longBuckets: %+v\n", latBuckets, longBuckets, files)
	minLat, minLong, maxLat, maxLong := getMinMaxLanLonBoundsMultiFile(files)
	regions := bucketRegionsFromMinMax(minLat, minLong, maxLat, maxLong, latBuckets, longBuckets)
	log.Printf("created %d regions\n", len(regions))

	for i, region := range regions {
		log.Printf("filtering files to region %d\n", i)
		filterGeojsonFilesToRegion(
			files,
			fmt.Sprintf("%s/bucket_%d_lats%d_longs%d.geojson", outputDir, i, latBuckets, longBuckets),
			region,
		)
	}
}

func bucketifyDataMem(files []string, outputDir string, latBuckets, longBuckets int) {
	if strings.Contains(outputDir[len(outputDir)-1:], "/") {
		outputDir = outputDir[:len(outputDir)-1]
	}
	log.Printf("bucketifying the following files into %d latBuckets and %d longBuckets: %+v\n", latBuckets, longBuckets, files)
	minLat, minLong, maxLat, maxLong := getMinMaxLanLonBoundsMultiFile(files)
	regions := bucketRegionsFromMinMax(minLat, minLong, maxLat, maxLong, latBuckets, longBuckets)
	log.Printf("created %d regions\n", len(regions))

	masterFc, err := getFeatureCollectionFromFiles(files)
	if err != nil {
		panic(err)
	}

	for i, region := range regions {
		log.Printf("filtering master fc to region #%d\n", i)
		fc := filterFeatureCollectionToRegion(masterFc, region)
		log.Printf("region #%d has %d points\n", i, len(fc.Features))
		writeFeatureCollectionToRegularFile(fc, fmt.Sprintf("%s/region_%d.geojson", outputDir, i))
	}
}

// generateFilteredFilename is for auto-generating filenames for filtered geojsons if outputting each filtered geojson separates as opposed to building a single master geojson. NOT IN USE.
func generateFilteredFilename(original string) string {
	return original[:len(original)-8] + "_filtered.geojson"
}

func main() {
	args := os.Args[1:]
	// output := args[0]
	// files := args[1:]

	// region := getBoundingLoop()
	// filterGeojsonFilesToRegion(files, output, region)

	outputDir := args[0]
	files := args[1:]
	bucketifyData(files, outputDir, 10, 10)
}
