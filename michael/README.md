# Michael's Code


## Folders

#### `geojson-summary`

This folder contains some Golang code that takes the data as provided (in the individual geojson file format) and does one of two things: it either filters the entire dataset to a single geojson file within a certain geographic bounding box or it turns the full geographic region into a grid of rectangles ("buckets") and produces a number of geojson files equal to the number of rectangles with the entire dataset filtered into each region. See the folder's README for a more in-depth explanation

#### `postgis`

This folder contains an explanation for how to get the provided dataset into a PostgreSQL database with PostGIS as well as a simple proof-of-concept notebook that demonstrates filtering the dataset using PostGIS's built-in spatial indexes. We saw massive speed increases for the type of workloads we were interesting in pursuing.


## Files

#### `count_geojson_points.py`

This just opens a list of geojson files provided in command-line args and counts the total number of points in all provided files. This was used for sanity-checking other processing workloads.

#### `merge_geojson_files.py`

This was used to merge a series of geojson files into a single file (it is somewhat specific to the type of geojson files we received). It is moderately useful, but merging the entire dataset creates a geojson file that is much too large to process on any of our computers. This discovery is part of what motivated the pursuit of a database solution like PostGIS.

#### `overtime.ipynb`

This took a geographically-reduced dataset (from `geojson-summary`) and built some relatively simple graphs of different features summary statistics on a monthly basis. This was one of our first explorations of the data and informed some later work.

#### `graphing.ipynb`

This contains a much more mature version of the approach pursued in `overtime.ipynb`, allowing us to build a helpful graph for any feature of our choosing in any filtered region of our choosing. It should be fairly intuitive to follow.


## Potentially useful references

- https://medium.com/@buckhx/unwinding-uber-s-most-efficient-service-406413c5871d
- geojson.io
- https://news.ycombinator.com/item?id=22288276
- geoman.io
