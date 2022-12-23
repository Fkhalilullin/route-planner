
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
            mapCanvas.drawLastMap()

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

function getRoute(minlon, minlat, maxlon, maxlat, paths, mapCanvas) {
    const url = `http://localhost:8000/route`
    const req = new XMLHttpRequest()
    req.responseType = 'json';
    req.open("POST", url, true)

    req.onload = function() {
        const status = req.status;
        if (status === 200) {
            console.log('Response received (getting route). SUCCESS.\n', req.response)

            mapCanvas.drawLastElevations()
            mapCanvas.drawLastMap()

            let routes = req.response
            for (let i = 0; i < routes.length; ++i) {
                let route = routes[i]
                let coords = new Float32Array(route.length * 2)
                for (let j = 0; j < route.length; ++j) {
                    coords[j * 2] = route[j].lon;
                    coords[j * 2 + 1] = route[j].lat;
                    mapCanvas.drawPoint(route[j].lon, route[j].lat, "#ff0000")
                }
                let canvasCoords = mapCanvas._convertCoordsToCanvas(coords)

                let splinePoints = calculateLineSpline(canvasCoords, canvasCoords.length, 3)
                mapCanvas.mapWebGLContext.drawPolyline(splinePoints, config.routeColor)

            }
            for (const point of mapCanvas.points) {
                let canvasPoint = mapCanvas._convertCoordsToCanvas(new Float32Array([point.lon, point.lat]))
                mapCanvas.drawPoint(canvasPoint[0], canvasPoint[1], config.routePointColor)
            }

        } else {
            console.log('Response received (getting route): ERROR. Response status: ', req.status)
        }
    }

    let body = JSON.stringify(new RouteRq(
        new MapNode(minlat, minlon),
        new MapNode(maxlat, maxlon),
        paths
    ));

    req.send(body)
    return req.response
}
