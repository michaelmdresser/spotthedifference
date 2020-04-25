import json
import sys

def count_geojson_features(filename):
    f = open(filename, "rb")
    full = json.loads(f.read())
    f.close()
    return len(full["features"])

if __name__ == "__main__":
    files = sys.argv[1:]

    print("total features:", sum(map(lambda x: count_geojson_features(x), files)))
