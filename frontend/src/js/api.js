
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
            let canvasCoords = self._convertCoordsToCanvas(coords)

            console.log(canvasCoords)
            let splinePoints = calculateLineSpline(canvasCoords, canvasCoords.length, 3)
            self.mapWebGLContext.drawPolyline(splinePoints, "#ff0000")

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
    let canvasCoords = [0.013128580525517464, 0.3997208774089813, 0.013128580525517464, 0.4226090610027313, -0.00009570312249707058, 0.4226090610027313, -0.013319986872375011, 0.4226090610027313, -0.013319986872375011, 0.44464951753616333, -0.013319986872375011, 0.4666900038719177, -0.013319986872375011, 0.48873046040534973, -0.013319986872375011, 0.5107709169387817, -0.013319986872375011, 0.5336591005325317, -0.013319986872375011, 0.5556995868682861, -0.026544271036982536, 0.5556995868682861, -0.03976855427026749, 0.5556995868682861, -0.03976855427026749, 0.5777400732040405, -0.03976855427026749, 0.5997805595397949, -0.052992839366197586, 0.5997805595397949, -0.052992839366197586, 0.6218210458755493, -0.052992839366197586, 0.6447092294692993, -0.06621712446212769, 0.6447092294692993, -0.08045865595340729, 0.6447092294692993, -0.09368294477462769, 0.6447092294692993, -0.10690722614526749, 0.6447092294692993, -0.12013150751590729, 0.6447092294692993, -0.13335579633712769, 0.6447092294692993, -0.13335579633712769, 0.6667496562004089, -0.14658008515834808, 0.6667496562004089, -0.1598043590784073, 0.6667496562004089, -0.1598043590784073, 0.6887901425361633, -0.17302864789962769, 0.6887901425361633, -0.18625293672084808, 0.6887901425361633, -0.20049446821212769, 0.6667496562004089, -0.21371875703334808, 0.6667496562004089, -0.2269430309534073, 0.6447092294692993, -0.24016731977462769, 0.6447092294692993, -0.2533915936946869, 0.6447092294692993, -0.2666158974170685, 0.6218210458755493, -0.2798401713371277, 0.5997805595397949, -0.2930644452571869, 0.5777400732040405, -0.3062887489795685, 0.5556995868682861, -0.3195130228996277, 0.5336591005325317, -0.3337545692920685, 0.5336591005325317]
    let splinePoints = calculateLineSpline(canvasCoords, canvasCoords.length * 2, 3)
    self.mapWebGLContext.drawPolyline(splinePoints, "#ff0000")
    return req.response
}
