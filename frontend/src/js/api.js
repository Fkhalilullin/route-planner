
function getMapFragment(minlon, minlat, maxlon, maxlat) {
    const url = `https://api.openstreetmap.org/api/0.6/map?bbox=${minlon},${minlat},${maxlon},${maxlat}`
    const req = new XMLHttpRequest()
    req.open("GET", url, false)
    req.send(null)
    return req.responseXML
}

function getMesh(minLon, minLat, maxLon, maxLat, mapCanvas) {
    const url = `http://localhost:8000/mesh`
    const req = new XMLHttpRequest()
    req.responseType = 'json'
    req.open("POST", url, true)

    req.onload = function() {
        if (req.status === 200) {
            console.log('Response received (getting mesh). SUCCESS.\n', req.response)

            let mesh = req.response

            mapCanvas.drawElevations(mesh)
            const responseXml = getMapFragment(minLon, minLat, maxLon, maxLat)
            mapCanvas.lastMapFragment = new MapFragment(responseXml)
            mapCanvas.drawMap(mapCanvas.lastMapFragment)

            document.body.style.cursor = 'default';

        } else {
            document.body.style.cursor = 'default';
            console.log('Response received (getting mesh): ERROR. Response status: ', req.status)
        }
    }

    let body = JSON.stringify(new RouteRq(
        new MapNode(minLat, minLon),
        new MapNode(maxLat, maxLon)
    ));
    // req.setRequestHeader("Content-Type", "application/json")
    console.log('Send request (getting mesh). Url: ', url, '\nBody: ', body)
    document.body.style.cursor = 'wait';
    req.send(body)
    return req.response
}

function getRoute(topLeftPoint, botRightPoint, beginPoint, endPoint, self) {
    const url = `http://localhost:8000/route`
    const req = new XMLHttpRequest()
    req.responseType = 'json';
    req.open("POST", url, true)

    req.onload = function() {
        const status = req.status;
        if (status === 200) {
            console.log('Response received (getting route). SUCCESS.\n', req.response)
            self.drawMap(self.lastMapFragment)

            let route = req.response
            let coords = new Float32Array(route.length * 2)
            for (let i = 0; i < route.length; ++i) {
                coords[i * 2    ] = route[i].lon;
                coords[i * 2 + 1] = route[i].lat;
                self.drawPoint(route[i].lon, route[i].lat, "#ff0000")
            }
            let canvasCoords = self._convertCoordsToCanvas(coords)

            let splinePoints = calculateLineSpline(canvasCoords, canvasCoords.length, 3)
            self.mapWebGLContext.drawPolyline(splinePoints, config.routeColor)

            self.drawPoint(self.startRoute.x, self.startRoute.y, "#ff9090")
            self.drawPoint(self.endRoute.x, self.endRoute.y, "#ff0000")

        } else {
            console.log('Response received (getting route): ERROR. Response status: ', req.status)
        }
    }

    let body = JSON.stringify(new RouteRq(
        new MapNode(topLeftPoint.lat, topLeftPoint.lon),
        new MapNode(botRightPoint.lat, botRightPoint.lon),
        new MapNode(beginPoint.lat, beginPoint.lon),
        new MapNode(endPoint.lat, endPoint.lon)
    ));

    req.send(body)
    return req.response
}
