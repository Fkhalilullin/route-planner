
// Vertex shader program
const VSHADER_SOURCE =
    '#version 100\n' +
    'attribute vec4 a_Position;\n' +
    'attribute vec4 a_Color;\n' +
    'varying vec4 v_Color;\n' +
    'void main() {\n' +
    '   gl_Position = a_Position;\n' +
    '   gl_PointSize = 10.0;\n' +
    '   v_Color = a_Color;\n' +
    '}\n';

// Fragment shader program
const FSHADER_SOURCE =
    'precision mediump float;\n' +
    'varying vec4 v_Color;\n' +
    'void main() {\n' +
    '   gl_FragColor = v_Color;\n' +
    '}\n';


class MapWebGLContext {

    canvas
    gl
    a_Position
    a_Color

    constructor(canvas) {
        // Retrieve <canvas> element
        this.canvas = canvas
    
        // Get the rendering context for WebGL
        this.gl = getWebGLContext(this.canvas);
        const gl = this.gl
        if (!gl)
        {
            console.log('Failed to get the rendering context for WebGL');
            return;
        }
    
        // Initialize shaders
        if (!initShaders(gl, VSHADER_SOURCE, FSHADER_SOURCE))
        {
            console.log('Failed to intialize shaders.');
            return;
        }

        // Clear canvas
        this.clear("#000000")

        gl.enable(gl.BLEND)

        // Create a buffer object
        const vertexBuffer = gl.createBuffer();
        if (!vertexBuffer)
        {
            console.log('Failed to create the buffer object');
            return -1;
        }
        // Bind the buffer object to target
        gl.bindBuffer(gl.ARRAY_BUFFER, vertexBuffer);

        this.a_Position = gl.getAttribLocation(gl.program, 'a_Position');
        if (this.a_Position < 0)
        {
            console.log('Failed to get the storage location of a_Position');
            return -1;
        }

        this.a_Color = gl.getAttribLocation(gl.program, 'a_Color');
        if(this.a_Color < 0) {
            console.log('Failed to get the storage location of a_Color');
            return -1;
        }

    }

    clear(hexColor) {
        // Specify the color for clearing <canvas>
        let [r, g, b] = this._hexToRgb(hexColor)
        this.gl.clearColor(r, g, b, 1);

        // Clear <canvas>
        this.gl.clear(this.gl.COLOR_BUFFER_BIT);
    }

    _initVertexBuffers(vertices, hexColor)
    {
        const gl = this.gl
        // Write date into the buffer object
        gl.bufferData(gl.ARRAY_BUFFER, vertices, gl.STATIC_DRAW);

        const FSIZE = vertices.BYTES_PER_ELEMENT;
        // Assign the buffer object to a_Position variable
        gl.vertexAttribPointer(this.a_Position,  2, gl.FLOAT, false, FSIZE * 2, 0);
        let [r, g, b] = this._hexToRgb(hexColor)
        gl.vertexAttrib4f(this.a_Color, r, g, b, 1.);
        
        // Enable the assignment to a_Position variable
        gl.enableVertexAttribArray(this.a_Position);
    }

    _initColorVertexBuffers(vertices, hexColors, alpha)
    {
        const gl = this.gl
        // Write date into the buffer object
        let verticesInfo = new Float32Array(hexColors.length * 6)
        for (let i = 0; i < hexColors.length; ++i) {
            let [r, g, b] = this._hexToRgb(hexColors[i])
            verticesInfo[i * 6    ] = vertices[i * 2    ]
            verticesInfo[i * 6 + 1] = vertices[i * 2 + 1]
            verticesInfo[i * 6 + 2] = r
            verticesInfo[i * 6 + 3] = g
            verticesInfo[i * 6 + 4] = b
            verticesInfo[i * 6 + 5] = alpha[i]
        }
        gl.bufferData(gl.ARRAY_BUFFER, verticesInfo, gl.STATIC_DRAW);

        const FSIZE = verticesInfo.BYTES_PER_ELEMENT;
        // Assign the buffer object to a_Position variable
        gl.vertexAttribPointer(this.a_Position, 2, gl.FLOAT, false, FSIZE * 6, 0);
        gl.vertexAttribPointer(this.a_Color, 4, gl.FLOAT, false, FSIZE * 6, FSIZE * 2);

        gl.enableVertexAttribArray(this.a_Position);
        gl.enableVertexAttribArray(this.a_Color);
    }

    _hexToRgb(hex) {
        const bigint = parseInt(hex.substring(1), 16);
        const r = ((bigint >> 16) & 255) / 255;
        const g = ((bigint >> 8) & 255) / 255;
        const b = (bigint & 255) / 255;
        return [r, g, b]
    }

    drawPolygon(vertices, hexColor) {
        let colorArray = this._createArray(vertices.length / 2, hexColor)
        let alphaArray = this._createArray(vertices.length / 2, 1.)
        // Write the positions of vertices to a vertex shader
        this._initColorVertexBuffers(vertices, colorArray, alphaArray)
        // Draw three points
        this.gl.drawArrays(this.gl.TRIANGLE_FAN, 0, vertices.length / 2);
    }

    _createArray(length, value) {
        let array = Array(length)
        for (let i = 0; i < length; ++i) {
            array[i] = value
        }
        return array
    }

    drawColorPolygon(vertices, hexColors, alpha) {
        // Write the positions of vertices to a vertex shader
        this._initColorVertexBuffers(vertices, hexColors, alpha)
        // Draw three points
        this.gl.drawArrays(this.gl.TRIANGLE_FAN, 0, vertices.length / 2)
    }

    drawPoints(vertices, hexColor) {
        // Write the positions of vertices to a vertex shader
        this._initVertexBuffers(vertices, hexColor);
        // Draw three points
        this.gl.drawArrays(this.gl.POINTS, 0, vertices.length / 2);
    }

    drawPolyline(vertices, hexColor) {
        let colorArray = this._createArray(vertices.length / 2, hexColor)
        let alphaArray = this._createArray(vertices.length / 2, 1.)
        // Write the positions of vertices to a vertex shader
        this._initColorVertexBuffers(vertices, colorArray, alphaArray)
        // Draw three points
        this.gl.drawArrays(this.gl.LINE_STRIP, 0, vertices.length / 2);
    }

}
