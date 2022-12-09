
// Vertex shader program
const VSHADER_SOURCE =
    '#version 100\n' +
    'attribute vec4 a_Position;\n' +
    'attribute vec4 a_Color;\n' +
    'void main() {\n' +
    '   gl_Position = a_Position;\n' +
    '   gl_PointSize = 10.0;\n' +
    '}\n';

// Fragment shader program
const FSHADER_SOURCE =
    'precision mediump float;\n' +
    'uniform vec4 u_Color;\n' +
    'void main() {\n' +
    '   gl_FragColor = u_Color;\n' +
    '}\n';


class MapWebGLContext {

    canvas
    gl
    a_Position
    u_Color

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

        this.u_Color = gl.getUniformLocation(gl.program, 'u_Color');
        if (this.u_Color < 0)
        {
            console.log('Failed to get the storage location of u_Color');
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
        gl.uniform4f(this.u_Color, r, g, b, 1.);
        
        // Enable the assignment to a_Position variable
        gl.enableVertexAttribArray(this.a_Position);
    }

    _hexToRgb(hex) {
        var bigint = parseInt(hex.substring(1), 16)
        var r = ((bigint >> 16) & 255) / 255
        var g = ((bigint >> 8) & 255) / 255
        var b = (bigint & 255) / 255
        return [r, g, b]
    }

    drawPolygon(vertices, hexColor) {
        // Write the positions of vertices to a vertex shader
        this._initVertexBuffers(vertices, hexColor);
        // Draw three points
        this.gl.drawArrays(this.gl.TRIANGLE_FAN, 0, vertices.length / 2);
    }

    drawPoints(vertices, hexColor) {
        // Write the positions of vertices to a vertex shader
        this._initVertexBuffers(vertices, hexColor);
        // Draw three points
        this.gl.drawArrays(this.gl.POINTS, 0, vertices.length / 2);
    }

    drawPolyline(vertices, hexColor) {
        // Write the positions of vertices to a vertex shader
        this._initVertexBuffers(vertices, hexColor);
        // Draw three points
        this.gl.drawArrays(this.gl.LINE_STRIP, 0, vertices.length / 2);
    }

}
