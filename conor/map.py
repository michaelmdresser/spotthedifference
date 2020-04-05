from flask import Flask
import folium

app = Flask(__name__)


@app.route('/')
def index():
    start_coords = (31.634036, 42.556084)
    folium_map = folium.Map(location=start_coords, zoom_start=8)
    return folium_map._repr_html_()


if __name__ == '__main__':
    app.run(debug=True)