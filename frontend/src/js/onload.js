
function main() {

    const canvas = new MapCanvas()

    let midlon = 73.3137
    let midlat = 33.403
    let lonDist = 0.015

    let minlon = midlon - lonDist / 2.
    let maxlon = midlon + lonDist / 2.
    let latDist = (canvas.getHeight() * lonDist) / canvas.getWidth()
    let minlat = midlat - latDist / 2.
    let maxlat = midlat + latDist / 2.

    console.log(minlon, minlat, maxlon, maxlat)
    const responseXml = getMapFragment(minlon, minlat, maxlon, maxlat)
    canvas.lastMapFragment = new MapFragment(responseXml)
    // canvas.drawMap(mapFragment)

    getMesh(minlon, minlat, maxlon, maxlat, canvas)

    // canvas.enableMouseEvents()
    
}
