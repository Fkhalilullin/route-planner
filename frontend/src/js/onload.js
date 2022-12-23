
function main() {

    const canvas = new MapCanvas()

    let midLon = 73.3137
    let midLat = 33.403
    let lonDist = 0.015

    let minLon = midLon - lonDist / 2.
    let maxLon = midLon + lonDist / 2.
    let latDist = (canvas.getHeight() * lonDist) / canvas.getWidth()
    let minLat = midLat - latDist / 2.
    let maxLat = midLat + latDist / 2.

    console.log(`minLon=${minLon}, minLat=${minLat}, maxLon=${maxLon}, maxLat=${maxLat}`)
    // const responseXml = getMapFragment(minLon, minLat, maxLon, maxLat)
    // canvas.lastMapFragment = new MapFragment(responseXml)
    // canvas.drawMap(mapFragment)

    canvas.lastMapFragment = new MapFragment().withCoords(minLon, minLat, maxLon, maxLat)

    getMesh(minLon, minLat, maxLon, maxLat, canvas)
    // sleep(5000)
    // const responseXml = getMapFragment(minLon, minLat, maxLon, maxLat)
    // canvas.lastMapFragment = new MapFragment(responseXml)
    // canvas.drawMap(canvas.lastMapFragment)


    canvas.enableMouseEvents()
    
}

function sleep(sleepDuration){
    let now = new Date().getTime();
    while(new Date().getTime() < now + sleepDuration){
        /* Do nothing */
    }
}
