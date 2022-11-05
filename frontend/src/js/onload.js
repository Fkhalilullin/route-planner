
function main() {
    const canvas = new MapCanvas()

    console.log(ColorUtils.getGradientColor('#000000', '#ffff00', 0, 255, 128));
    console.log(ColorUtils.getGradientColor('#000000', '#ffff00', 0, 255, 255));
    console.log(ColorUtils.getGradientColor('#000000', '#ffff00', 0, 255, 0));
    console.log(ColorUtils.getGradientColor('#000000', '#ffff00', 0, 255, 64));

    let midlon = 37.2671//37.26049//9.45874//37.26050
    let midlat = 55.4863//55.48261//0.38700//5.48223
    let lonDist = 0.009

    let minlon = midlon - lonDist / 2.
    let maxlon = midlon + lonDist / 2.
    let latDist = (canvas.getHeight() * lonDist) / canvas.getWidth()
    let minlat = midlat - latDist / 2.
    let maxlat = midlat + latDist / 2.

    const responseXml = getMapFragment(minlon, minlat, maxlon, maxlat)
    const mapFragment = new MapFragment(responseXml)
    canvas.drawMap(mapFragment)

    let elevationGrid = [
        new ElevationGrid(30., new Location(-0.5, 0.5)),
        new ElevationGrid(40., new Location(0.5, 0.5)),
        new ElevationGrid(20., new Location(-0.5, -0.5)),
        new ElevationGrid(80., new Location(0.5, -0.5)),
        // new ElevationGrid(40., new Location(-0.8, 0.5)),
        // new ElevationGrid(55., new Location(0.8, 0.5)),
        // new ElevationGrid(65., new Location(-0.8, -0.5)),
        // new ElevationGrid(50., new Location(0.55, -0.4))
    ]
    let elevationRs = new ElevationRs(20., 80., elevationGrid)
    canvas.drawReliefMap(elevationRs)

    canvas.enableMouseEvents()

}
