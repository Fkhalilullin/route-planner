
function main() {

    const canvas = new MapCanvas()

    let midlon = 73.3137//37.2671//36.9005//0.7115//37.2671//37.26049//9.45874//37.26050
    let midlat = 33.3953//55.4863//55.2004//11.1921//55.4863//55.48261//0.38700//5.48223
    let lonDist = 0.015

    let minlon = midlon - lonDist / 2.
    let maxlon = midlon + lonDist / 2.
    let latDist = (canvas.getHeight() * lonDist) / canvas.getWidth()
    let minlat = midlat - latDist / 2.
    let maxlat = midlat + latDist / 2.

    console.log(minlon, minlat, maxlon, maxlat)
    const responseXml = getMapFragment(minlon, minlat, maxlon, maxlat)
    const mapFragment = new MapFragment(responseXml)
    canvas.drawMap(mapFragment)

    canvas.enableMouseEvents()
    
}
