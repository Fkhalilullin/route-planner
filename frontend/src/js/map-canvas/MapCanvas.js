
class MapCanvas {
    
    canvas
    mapWebGLContext
    lastMapFragment
    lastElevationMesh
    points

    constructor() {
        this.canvas = document.getElementById('webgl');
        this.mapWebGLContext = new MapWebGLContext(this.canvas)
        this.points = []
    }

    drawLastMap() {
        this.drawMap(this.lastMapFragment)
    }

    drawMap(mapFragment) {
        this.lastMapFragment = mapFragment
        // this.mapWebGLContext.clear("#f2efe9")
        mapFragment.ways.forEach(way => {
            let canvasCoords = this._convertCoordsToCanvas(way.vertices)
            let color
            if (way.wayTags.get("building") != null) {
                color = "#d9d0c9"
            } else if (way.wayTags.get("natural") != null) {
                console.log(way.wayTags.get("natural"))
                const naturalValue = way.wayTags.get("natural")
                if (naturalValue === "wood")
                    color = "#aadd19"
                else if (naturalValue === "water")
                    color = config.waterColor
            } else if (way.wayTags.get("waterway") != null) {
                color = config.waterColor
            }
            if (color != null) {
                this.mapWebGLContext.drawPolygon(canvasCoords, color)
            }
        })
    }

    _componentToHex(c) {
        let hex = Math.round(c * 255).toString(16)
        return hex.length === 1 ? "0" + hex : hex
    }

    _rgbToHex(r, g, b) {
        return "#" + this._componentToHex(r) + this._componentToHex(g) + this._componentToHex(b);
    }

    _indexesFrom2dTo1d(i, j, columnCount) {
        return (i * columnCount + j)
    }

    _gradient(value, minValue, maxValue, downColor, upColor) {
        let coef = (value - minValue) / (maxValue - minValue)
        let [dr, dg, db] = this.mapWebGLContext._hexToRgb(downColor)
        let [ur, ug, ub] = this.mapWebGLContext._hexToRgb(upColor)
        let r = (ur - dr) * coef + dr
        let g = (ug - dg) * coef + dg
        let b = (ub - db) * coef + db
        return this._rgbToHex(r, g, b)
    }

    _elevationGradient(value, minValue, maxValue) {
        return this._gradient(value, minValue, maxValue, config.downColor, config.upColor)
    }

    drawLastElevations() {
        this.drawElevations(this.lastElevationMesh)
    }

    drawElevations(elevationMesh) {
        this.lastElevationMesh = elevationMesh
        for (let i = 0; i < elevationMesh.rowCount - 1; ++i) {
            for (let j = 0; j < elevationMesh.columnCount - 1; ++j) {

                let geographicCoords = new Float32Array(8)
                let point_0 = elevationMesh.points[this._indexesFrom2dTo1d(i, j, elevationMesh.columnCount)]
                let point_1 = elevationMesh.points[this._indexesFrom2dTo1d(i + 1, j, elevationMesh.columnCount)]
                let point_2 = elevationMesh.points[this._indexesFrom2dTo1d(i + 1, j + 1, elevationMesh.columnCount)]
                let point_3 = elevationMesh.points[this._indexesFrom2dTo1d(i, j + 1, elevationMesh.columnCount)]
                geographicCoords[0] = point_0.lon
                geographicCoords[1] = point_0.lat
                geographicCoords[2] = point_1.lon
                geographicCoords[3] = point_1.lat
                geographicCoords[4] = point_2.lon
                geographicCoords[5] = point_2.lat
                geographicCoords[6] = point_3.lon
                geographicCoords[7] = point_3.lat
                const minElevation = elevationMesh.minElevation
                const maxElevation = elevationMesh.maxElevation
                const hexColor_1 = this._elevationGradient(point_1.elevation, minElevation, maxElevation)
                const hexColor_2 = this._elevationGradient(point_2.elevation, minElevation, maxElevation)
                const hexColor_0 = this._elevationGradient(point_0.elevation, minElevation, maxElevation)
                const hexColor_3 = this._elevationGradient(point_3.elevation, minElevation, maxElevation)
                const hexColors = [
                    hexColor_0,
                    hexColor_1,
                    hexColor_2,
                    hexColor_3
                ]
                const alpha = 1
                const alphaArray = [
                    alpha,
                    alpha,
                    alpha,
                    alpha
                ]
                let vertices = this._convertCoordsToCanvas(geographicCoords)
                this.mapWebGLContext.drawColorPolygon(vertices, hexColors, alphaArray)
            }
        }
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
        this.mapWebGLContext.drawPolyline(vertices, config.routeColor)
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
        let mouseDown = false;
        let lastX = -1., lastY = -1.;   // Last position of the mouse
        let clickX = -1., clickY = -1.;   // Last position of the mouse
      
        let mapCanvas = this
        this.canvas.onmousedown = function(ev) {   // Mouse is pressed
            clickX = ev.clientX, clickY = ev.clientY;
            
            // Start dragging if a mouse is in <canvas>
            const rect = ev.target.getBoundingClientRect();
            if (rect.left <= clickX && clickX < rect.right && rect.top <= clickY && clickY < rect.bottom) {
                lastX = clickX; lastY = clickY;
                mouseDown = true;
            }
        };

        let self = this
        this.canvas.onmouseup = function(ev) {
            // self.drawMap(self.lastMapFragment)
            const minlon = mapCanvas.lastMapFragment.minlon
            const minlat = mapCanvas.lastMapFragment.minlat
            const maxlon = mapCanvas.lastMapFragment.maxlon
            const maxlat = mapCanvas.lastMapFragment.maxlat

            const x = ev.clientX, y = ev.clientY;
            mouseDown = false
            if (lastX - clickX === 0. && lastY - clickY === 0.) {
                console.log("click")
                const rect = ev.target.getBoundingClientRect();
                const canvasX = ((x - rect.left) - self.canvas.width/2.)/(self.canvas.width/2.);
                const canvasY = (self.canvas.height/2. - (y - rect.top))/(self.canvas.height/2.);

                let geoPoint = self._convertCanvasToGeographic(new CanvasNode(canvasX, canvasY))
                let nearestPoint = self.getNearestPoint(geoPoint.lon, geoPoint.lat)
                self.points.push(new MapNode(nearestPoint.lat, nearestPoint.lon))

                mapCanvas.drawLastElevations()
                mapCanvas.drawLastMap()
                for (const point of self.points) {
                    let canvasNearestPoint = self._convertCoordsToCanvas(
                        new Float32Array([point.lon, point.lat])
                    )
                    self.drawPoint(canvasNearestPoint[0], canvasNearestPoint[1], "#ff9090")
                }
                if (self.points.length > 1)
                    getRoute(minlon, minlat, maxlon, maxlat, self.points, self)
            } else {
                getMesh(
                    minlon,
                    minlat,
                    maxlon,
                    maxlat,
                    mapCanvas
                )
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
        
                const minLon = Number(mapCanvas.lastMapFragment.minlon) - mapDx
                const minLat = Number(mapCanvas.lastMapFragment.minlat) + mapDy
                const maxLon = Number(mapCanvas.lastMapFragment.maxlon) - mapDx
                const maxLat = Number(mapCanvas.lastMapFragment.maxlat) + mapDy

                self.points = []

                const responseXml = getMapFragment(minLon, minLat, maxLon, maxLat)
                const mapFragment = new MapFragment(responseXml)
                mapCanvas.drawMap(mapFragment)
            }
            lastX = x, lastY = y;
        };

    }

    getNearestPoint(lon, lat) {
        let nodes = this.lastElevationMesh.points
        let step = Math.abs(this.lastElevationMesh.points[1].lon - this.lastElevationMesh.points[0].lon)
        let half_step = step / 2.

        let row = Math.floor((lat - nodes[0].lat) / step)
        let column = Math.floor((lon - nodes[0].lon) / step)

        let pointIndex = this._indexesFrom2dTo1d(row, column, this.lastElevationMesh.columnCount)
        let point = nodes[pointIndex]
        let delta_lon = lon - point.lon
        let delta_lat = lat - point.lat

        if (0. <= delta_lon && delta_lon <= half_step) {
            if (0. <= delta_lat && delta_lat <= half_step)
                return nodes[this._indexesFrom2dTo1d(row, column, this.lastElevationMesh.columnCount)]
            if (half_step < delta_lat && delta_lat <= step)
                return nodes[this._indexesFrom2dTo1d(row + 1, column, this.lastElevationMesh.columnCount)]
        }
        if (half_step < delta_lon && delta_lon <= step) {
            if (0. <= delta_lat && delta_lat <= half_step)
                return nodes[this._indexesFrom2dTo1d(row, column + 1, this.lastElevationMesh.columnCount)]
            if (half_step < delta_lat && delta_lat <= step)
                return nodes[this._indexesFrom2dTo1d(row + 1, column + 1, this.lastElevationMesh.columnCount)]
        }
        return nodes[this._indexesFrom2dTo1d(row, column, this.lastElevationMesh.columnCount)]
    }

}
