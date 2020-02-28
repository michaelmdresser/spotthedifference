# Installing PostGIS in PostgreSQL
https://postgis.net/install/

Install database, set up users, `CREATE DATABASE d4d`, then connect to database `\c d4d` and run `CREATE EXTENSION postgis;`

# Getting GeoJSON into PostGIS

https://gis.stackexchange.com/questions/172092/import-geojson-into-postgis
http://www.sarasafavi.com/installing-gdalogr-on-ubuntu.html

```
export GEOJSON_FILE='/mnt/d4d/geojson/20181101030117_20181113030116_365149N-ACD-BETA/01Nov18S1B101586_01_13Nov18S1B102974_01_ACD_Departures.geojson'
ogr2ogr -f "PostgreSQL" PG:"dbname=d4d user=postgres" $GEOJSON_FILE -nln iraq_points -append
```

```
for GEOJSON_FILE in /mnt/d4d/geojson/*/*Departures.geojson /mnt/d4d/geojson/*/*Arrivals.geojson
do
ogr2ogr -f "PostgreSQL" PG:"dbname=d4d user=postgres" $GEOJSON_FILE -nln iraq_points -append
done
```

```
SELECT ST_AsText(ST_FlipCoordinates(wkb_geometry)), timestamp, detect_type FROM iraq_points LIMIT 5;
```

# Querying points in a given polygon
https://gis.stackexchange.com/questions/227892/how-to-find-points-within-a-polygon-in-postgis
https://gis.stackexchange.com/questions/219756/selecting-points-within-a-polygon-in-postgis
https://postgis.net/docs/ST_GeomFromText.html
https://postgis.net/docs/ST_MakePolygon.html
https://postgis.net/docs/ST_Within.html
https://gis.stackexchange.com/questions/68711/postgis-geometry-query-returns-error-operation-on-mixed-srid-geometries-only
https://gis.stackexchange.com/questions/131363/choosing-srid-and-what-is-its-meaning

```
SELECT ST_MakePolygon( ST_GeomFromText('LINESTRING(33.272014 42.447810, 32.361950 40.147858, 29.802965 43.014024, 30.791453 44.904007, 33.272014 42.447810)'));
```

```
SELECT COUNT(*) FROM iraq_points AS ip WHERE ST_Within(
  ST_FlipCoordinates(ip.wkb_geometry),
  ST_MakePolygon( ST_GeomFromText('LINESTRING(33.272014 42.447810, 32.361950 40.147858, 29.802965 43.014024, 30.791453 44.904007, 33.272014 42.447810)', 4326))
);
```

# Backup and Restore
https://postgis.net/workshops/postgis-intro/backup.html

### Backup
```
pg_dump --file=/mnt/d4d/postgis_all_points_2020-02-28.backup --format=c --username=postgres d4d
```

### Restore
on the new system:
Install postgres database, install postgis addon,

In `psql`:
```
CREATE DATABASE d4d;
\c d4d
CREATE EXTENSION postgis;
```

CLI:
```
pg_restore --dbname=d4d --username=postgres postgis_all_points_2020-02-28.backup
```
