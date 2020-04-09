## D4D - Spot the Difference

### TerraDelta - Luke McConnell

#### Coding Deliverable Readme


#### What the program does:

This program takes a month and a calendar day from user input and outputs a map with all arrival and departure points plotted for that day. This make a navigatable heat map to view locations and density of points.

#### What a user would need to know:

The program works by accessing a directory (AKA folder) of geoJSON files. When our team first obtained the data, we separated out the KMZ files from the geoJSONs. What was left, was a directory named 'geojsons' that contains 1159 directories of geoJSON files, most of which contain an arrival.geojson and a departure.geojson pair. This is brought to the users attention because the code is designed to interact with this orientation, as it is mostly a proof of concept. 

#### What a user would need to do before running the program:

The code is written in Python, which accesses directories with a command like: 

>os.chdir('C:\\Users\\luker\\Desktop\\d4d_work\\geojsons')

Therefore, if a user wanted to run the code on their local computer, a correct directory path must be updated in the code for the computer they are working on. This would need to occur at the line shown above, and also the line of code below:

>curr_target_dir_path = 'C:\\Users\\luker\\Desktop\\d4d_work\\geojsons\\' + curr_target_dir


#### How to run the program:

There is a .ipynb Jupyter Notebook file that can be run when Jupyter Notebook is launched.

There is also a .py file that can be run from a terminal/shell using a command like:

>python d4d_final_deliverable.py

#### Program output:

The program outputs a map.html file that will be overwritten each time the program runs. Make sure to save or move a map.html with results that are desired to be seen later without running the Map.




