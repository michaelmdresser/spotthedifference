import sys
import json

def merge_and_output(output_file, files):
    total_features = 0

    f = open(output_file, "w+")
    header = """{
        "type": "FeatureCollection",
            "features": [
    """
    f.write(header)

    first = True

    for filename in files:
        source_file = open(filename, "rb")
        loaded = json.loads(source_file.read())
        total_features += len(loaded["features"])

        for feature in loaded["features"]:
            if first:
                first = False
            else:
                f.write(",\n")

            f.write(json.dumps(feature))
        source_file.close()

    footer = "]}"
    f.write(footer)
    f.close()

    print("total features:", total_features)

if __name__ == "__main__":
    if len(sys.argv) <= 2:
        print("USAGE: python3 merge_geojson_files.py OUTPUTFILE INPUTFILE1 [INPUTFILE2 ...]")
        sys.exit(1)
    output_file = sys.argv[1]
    files = sys.argv[2:]
    
    merge_and_output(output_file, files)
