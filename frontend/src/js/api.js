
function getMapFragment(minlon, minlat, maxlon, maxlat) {
    const url = `https://api.openstreetmap.org/api/0.6/map?bbox=${minlon},${minlat},${maxlon},${maxlat}`
    const req = new XMLHttpRequest()
    req.open("GET", url, false)
    req.send(null)
    return req.responseXML
}

function getElevationGrid() {
    const url = `http://localhost:8000/points?min_lat=63.9391&min_lon=50.6874&max_lat=63.9454&max_lon=50.7074`
    const req = new XMLHttpRequest()
    req.responseType = 'json';
    req.open("GET", url, true)
    req.onload  = function() {
        const status = req.status;
        if (status === 200) {
            console.log('SUCCESS\n', req.response)
        } else {
            console.log('ERROR')
        }
    };
    req.send(null)
    return req
}
