
class MapCanvas {
    
    canvas
    mapWebGLContext
    lastMapFragment
    startRoute
    endRoute

    constructor() {
        this.canvas = document.getElementById('webgl');
        this.mapWebGLContext = new MapWebGLContext(this.canvas)
    }

    drawMap(mapFragment) {
        this.lastMapFragment = mapFragment
        this.mapWebGLContext.clear("#f2efe9")
        mapFragment.ways.forEach(way => {
            let canvasCoords = this._convertCoordsToCanvas(way.vertices)
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
            } else if (way.wayTags.get("waterway") != null) {
                color = "#aad3df"
            }
            if (color != null) {
                this.mapWebGLContext.drawPolygon(canvasCoords, color)
            }
        })
    }

    _convertCoordsToCanvas(geographicCoords) {
        const deltaLat = this.lastMapFragment.maxlat - this.lastMapFragment.minlat
        const deltaLon = this.lastMapFragment.maxlon - this.lastMapFragment.minlon
        const canvasCoords = new Float32Array(geographicCoords.length)
        for (let i = 0; i < geographicCoords.length / 2; ++i) {
            const factor = 1//(Math.sin(Math.abs((geographicCoords[i * 2 + 1] * Math.PI / 180.))) + 1)
            canvasCoords[i * 2    ] = (geographicCoords[i * 2    ] - this.lastMapFragment.minlon) * 2 / deltaLon - 1
            canvasCoords[i * 2 + 1] = ((geographicCoords[i * 2 + 1] - this.lastMapFragment.minlat) * 2 / deltaLat - 1) * factor
        }
        return canvasCoords
    }

    _convertCanvasToGeographic(canvasNode) {
        const deltaLat = this.lastMapFragment.maxlat - this.lastMapFragment.minlat
        const deltaLon = this.lastMapFragment.maxlon - this.lastMapFragment.minlon
        const factor = 1//(Math.sin(Math.abs((canvasNode.y * Math.PI / 180.))) + 1)
        const lon = (canvasNode.x + 1) * deltaLon / 2 + this.lastMapFragment.minlon;
        const lat = (canvasNode.y / factor + 1) * deltaLat / 2 + this.lastMapFragment.minlat;
        return new MapNode(lat, lon);
    }

    drawRoute(vertices) {
        this.mapWebGLContext.drawPolyline(vertices, "#ff0000")
    }

    drawPoint(canvasX, canvasY, color) {
        this.mapWebGLContext.drawPoints(new Float32Array([canvasX, canvasY]), color)
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

        let self = this
        this.canvas.onmouseup = function(ev) {
            self.drawMap(self.lastMapFragment)
            const x = ev.clientX, y = ev.clientY;
            mouseDown = false
            if (lastX - clickX == 0. && lastY - clickY == 0.) {
                console.log("click")
                const rect = ev.target.getBoundingClientRect();
                const canvasX = ((x - rect.left) - self.canvas.width/2.)/(self.canvas.width/2.);
                const canvasY = (self.canvas.height/2. - (y - rect.top))/(self.canvas.height/2.);

                const newNode = new CanvasNode(canvasX, canvasY)
                if (self.startRoute == null)
                    self.startRoute = newNode
                else if (self.endRoute == null)
                    self.endRoute = newNode
                else {
                    self.startRoute = newNode
                    self.endRoute = null
                }

                self.drawPoint(self.startRoute.x, self.startRoute.y, "#ff9090")
                if (self.endRoute != null) {
                    self.drawPoint(self.endRoute.x, self.endRoute.y, "#ff0000")

                    let topLeftPoint = new MapNode(self.lastMapFragment.minlat, self.lastMapFragment.minlon)
                    let botRightPoint = new MapNode(self.lastMapFragment.maxlat, self.lastMapFragment.maxlon)
                    let beginPoint = self._convertCanvasToGeographic(self.startRoute)
                    let endPoint = self._convertCanvasToGeographic(self.endRoute)

                    let route = getRoute(topLeftPoint, botRightPoint, beginPoint, endPoint, self)

                }

                // self.drawPoint(canvasX, canvasY, "#ff0000")
            }
            document.body.style.cursor = 'default';
        }; // Mouse is released
    
        this.canvas.onmousemove = function(ev) { // Mouse is moved
            const x = ev.clientX, y = ev.clientY;
            if (mouseDown) {
                console.log("onmousemove")
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
