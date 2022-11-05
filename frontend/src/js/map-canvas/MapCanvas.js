
class MapCanvas {
    
    canvas
    mapWebGLContext
    lastMapFragment

    constructor() {
        this.canvas = document.getElementById('webgl');
        this.mapWebGLContext = new MapWebGLContext(this.canvas)
    }

    drawMap(mapFragment) {
        this.lastMapFragment = mapFragment
        this.mapWebGLContext.clear("#f2efe9")
        mapFragment.ways.forEach(way => {
            const canvasCoords = this._convertCoordsToCanvas(way.vertices, mapFragment)
            let color
            if (way.wayTags.get("building") != null) {
                color = "#d9d0c9"
            } else if (way.wayTags.get("natural") != null) {
                console.log(way.wayTags.get("natural"))
                const naturalValue = way.wayTags.get("natural")
                if (naturalValue == "wood")
                    color = "#aadd19e"
                else if (naturalValue == "water")
                    color = "#aad3df"
            }
            if (color != null)
                this.mapWebGLContext.drawPolygon(canvasCoords, color)
        })
    }

    drawReliefMap(elevationRs) {
        let verticesInfoArr = new Float32Array(elevationRs.elevationGrid.length * 7)
        for (let i = 0; i < elevationRs.elevationGrid.length; ++i) {
            let [r, g, b] = ColorUtils.getGradientColor(
                '#ffffff', 
                '#ff0000',
                elevationRs.minElevation,
                elevationRs.maxElevation,
                elevationRs.elevationGrid[i].elevation)
            verticesInfoArr[i * 7    ] = elevationRs.elevationGrid[i].location.lat
            verticesInfoArr[i * 7 + 1] = elevationRs.elevationGrid[i].location.lon
            verticesInfoArr[i * 7 + 2] = 50.
            verticesInfoArr[i * 7 + 3] = r
            verticesInfoArr[i * 7 + 4] = g
            verticesInfoArr[i * 7 + 5] = b
            verticesInfoArr[i * 7 + 6] = 1.
        }
        console.log(verticesInfoArr);
        this.mapWebGLContext.drawPoints(verticesInfoArr)
    }

    _convertCoordsToCanvas(geographicCoords, mapFragment) {
        const deltaLat = mapFragment.maxlat - mapFragment.minlat
        const deltaLon = mapFragment.maxlon - mapFragment.minlon
        const canvasCoords = new Float32Array(geographicCoords.length)
        for (let i = 0; i < geographicCoords.length / 2; ++i) {
            const factor = (Math.sin(Math.abs((geographicCoords[i * 2 + 1] * Math.PI / 180.))) + 1)
            canvasCoords[i * 2    ] = (geographicCoords[i * 2    ] - mapFragment.minlon) * 2 / deltaLon - 1
            canvasCoords[i * 2 + 1] = ((geographicCoords[i * 2 + 1] - mapFragment.minlat) * 2 / deltaLat - 1) * factor
        }
        return canvasCoords
    }

    drawRoute(vertices) {
        this.mapWebGLContext.drawPolyline(vertices, "#ff0000")
    }

    getWidth() {
        return this.canvas.width
    }

    getHeight() {
        return this.canvas.height
    }

    enableMouseEvents() {
        var mouseDown = false;
        var lastX = -1., lastY = -1.;   // Last position of the mouse
        var clickX = -1., clickY = -1.;   // Last position of the mouse
      
        let mapCanvas = this
        this.canvas.onmousedown = function(ev) {   // Mouse is pressed
            clickX = ev.clientX, clickY = ev.clientY;
            
            // Start dragging if a mouse is in <canvas>
            var rect = ev.target.getBoundingClientRect();
            if (rect.left <= clickX && clickX < rect.right && rect.top <= clickY && clickY < rect.bottom) {
                lastX = clickX; lastY = clickY;
                mouseDown = true;
            }
        };
    
        this.canvas.onmouseup = function(ev) {
            const x = ev.clientX, y = ev.clientY;
            mouseDown = false
            if (lastX - clickX == 0. && lastY - clickY == 0.) {
                console.log("click")
                console.log(document.activeElement)
            }
            document.body.style.cursor = 'default';
        }; // Mouse is released
    
        this.canvas.onmousemove = function(ev) { // Mouse is moved
            const x = ev.clientX, y = ev.clientY;
            if (mouseDown) {
                document.body.style.cursor = 'grabbing';
                const dx = (x - lastX);
                const dy = (y - lastY);
    
                const mapDy = (mapCanvas.lastMapFragment.maxlat - mapCanvas.lastMapFragment.minlat) / mapCanvas.canvas.height * dy
                const mapDx = (mapCanvas.lastMapFragment.maxlon - mapCanvas.lastMapFragment.minlon) / mapCanvas.canvas.width * dx
        
                const minlon = Number(mapCanvas.lastMapFragment.minlon) - mapDx
                const minlat = Number(mapCanvas.lastMapFragment.minlat) + mapDy
                const maxlon = Number(mapCanvas.lastMapFragment.maxlon) - mapDx
                const maxlat = Number(mapCanvas.lastMapFragment.maxlat) + mapDy

                const responseXml = getMapFragment(minlon, minlat, maxlon, maxlat)
                const mapFragment = new MapFragment(responseXml)
                mapCanvas.drawMap(mapFragment)
            }
            lastX = x, lastY = y;
        };

    }
    
}
