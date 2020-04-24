# Conor's Code


## Files

#### `clustering_example.geojson`

This is an example output, using Michael's `geojson-summary` code, of all of the geojsons over the whole time period of the given dataset within the following bounding box: `[(33.272014, 42.447810), (32.361950, 40.147858), (29.802965, 43.014024), (30.791453, 44.904007)]`. This is an area focused around the Arabian desert in southern Iraq because it is a remote area.

#### `Understanding_Clustering.ipynb`

This is a Jupyter Notebook that demonstrates how 10km clusters are formed from **a singlular** timestamp in `clustering_example.geojson`: March 16th, 2019. The notebook is annotated to help explain how clustering is done. Refer to our write-up for a conceptual understanding of what clusters are. The final output of this notebook is an HTML file displaying the clusters (represented as circles) on a Leaflet map. The notebook will automatically open this HTML file in a web browser. Smaller clusters are represented in blue and larger clusters are represented in red. Click on a cluster for a popup on the cluster number and size.

Note that any local file paths will have to be changed to run in a different environment.

#### `sample_map.html`

This is the output file from the `Understanding_Clustering.ipynb` notebook described above. You may open this file in a browser to view the output without running the notebook.

#### `Cluster_All.ipynb`

Here is the most exciting piece of my work. This is a Jupyter Notebook that demonstrates how 10km clusters are formed from **every timestamp/collect** in `clustering_example.geojson`. Refer to our write-up for a conceptual understanding of what clusters are. The final output of this notebook is a folder of many HTML files displaying the clusters (represented as circles) on a Leaflet map. The notebook will automatically open the first timestamp of the dataset in a web browser. Smaller clusters are represented in blue and larger clusters are represented in red. Click on a cluster for a popup on the cluster number and size. **There are "next" and "previous" buttons on the side of the Leaflet map that will cycle through the timestamps/collects.** 

Note that any local file paths will have to be changed to run in a different environment.


## Folders

#### `map_outputs`

This is a folder containing all of the outputted HTML files from the `Cluster_All.ipynb` notebook described above. You may open `map1.html` in a browser to start at the first timestamp and cycle through each collect without running the notebook.


## Noteworthy Technical Details on Clustering

* The capability of the bounding box being analyzed is restricted by file size; processing too many points at once can cause the program to crash, so bounding boxes may need to be broken down into smaller units / boxes.
* Our implementation of clustering used `geojson-summary` to create a bounding box that is small enough for the program to run.
* The default cluster distance threshold in our program is 10 km
* The radius of the circles that represent clusters use the average distance from the centroid to each point in the cluster by default; however, maximum distance is also calculated in the program, and this can be used for a circular representation that will capture every point in a cluster with 100% accuracy
