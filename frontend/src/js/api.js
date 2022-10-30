
function getMapFragment(minlon, minlat, maxlon, maxlat) {
    const url = `https://api.openstreetmap.org/api/0.6/map?bbox=${minlon},${minlat},${maxlon},${maxlat}`
    const req = new XMLHttpRequest()
    req.open("GET", url, false)
    req.send(null)
    return req.responseXML
}