
function getMapFragment(minlon, minlat, maxlon, maxlat) {
    const url = `https://api.openstreetmap.org/api/0.6/map?bbox=${minlon},${minlat},${maxlon},${maxlat}`
    const req = new XMLHttpRequest()
    req.open("GET", url, false)
    req.send(null)
    console.log(req.response)
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

function getRoute(topLeftPoint, botRightPoint, beginPoint, endPoint, self) {
    const url = `http://localhost:8000/route`
    const req = new XMLHttpRequest()
    req.responseType = 'json';
    req.open("POST", url, true)

    req.onload = function() {
        const status = req.status;
        if (status === 200) {
            console.log(req.response)
            console.log(req.response.length)
            console.log('SUCCESS\n', req.response)
            self.drawMap(self.lastMapFragment)

            let route = req.response
            let coords = new Float32Array(route.length * 2)
            for (let i = 0; i < route.length; ++i) {
                console.log(route[i].lon, route[i].lat)
                coords[i * 2    ] = route[i].lon;
                coords[i * 2 + 1] = route[i].lat;
                self.drawPoint(route[i].lon, route[i].lat, "#ff0000")
            }
            console.log(coords)
            let canvasCoords = self._convertCoordsToCanvas(coords)
            self.mapWebGLContext.drawPolyline(canvasCoords, "#ff0000")

            self.drawPoint(self.startRoute.x, self.startRoute.y, "#ff9090")
            self.drawPoint(self.endRoute.x, self.endRoute.y, "#ff0000")

        } else {
            console.log('ERROR')
        }
    }

    let body = JSON.stringify(new RouteRq(
        new MapNode(topLeftPoint.lat, topLeftPoint.lon),
        new MapNode(botRightPoint.lat, botRightPoint.lon),
        new MapNode(beginPoint.lat, beginPoint.lon),
        new MapNode(endPoint.lat, endPoint.lon)
    ));

    console.log(body)

    req.send(body)
    return req.response
}
