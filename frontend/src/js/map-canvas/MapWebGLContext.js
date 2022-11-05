
// Vertex shader program
const VSHADER_SOURCE =
    '#version 100\n' +
    'attribute vec4  a_Position;\n' +
    'attribute float a_Size;\n' +
    'attribute vec4  a_Color;\n' +
    'varying   vec4  v_Color;\n' +
    '\n' +
    'void main() {\n' +
    '   gl_Position = a_Position;\n' +
    '   gl_PointSize = a_Size;\n' +
    '   v_Color = a_Color;\n' +
    '}\n';

// Fragment shader program
const FSHADER_SOURCE =
    '#ifdef GL_ES\n' +
    'precision mediump float;\n' +
    '#endif\n' +
    'varying vec4 v_Color;\n' +
    '\n' +
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
        this.clear("#ffffff")

        gl.enable(gl.BLEND)
        gl.blendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

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

        this.a_Size = gl.getAttribLocation(gl.program, 'a_Size');
        if (this.a_Size < 0)
        {
            console.log('Failed to get the storage location of a_Size');
            return -1;
        }

        this.a_Color = gl.getAttribLocation(gl.program, 'a_Color');
        if (this.a_Color < 0)
        {
            console.log('Failed to get the storage location of a_Color');
            return -1;
        }

    }

    clear(hexColor) {
        // Specify the color for clearing <canvas>
        let [r, g, b] = ColorUtils.hexToRgb(hexColor)
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
        let [r, g, b] = ColorUtils.hexToRgb(hexColor)
        gl.vertexAttrib4f(this.a_Color, r, g, b, 1.);
        
        // Enable the assignment to a_Position variable
        gl.enableVertexAttribArray(this.a_Position);
    }

    drawPolygon(vertices, hexColor) {
        // Write the positions of vertices to a vertex shader
        this._initVertexBuffers(vertices, hexColor);
        // Draw three points
        this.gl.drawArrays(this.gl.TRIANGLE_FAN, 0, vertices.length / 2);
    }

    drawPolyline(vertices, hexColor) {
        // Write the positions of vertices to a vertex shader
        this._initVertexBuffers(vertices, hexColor);
        // Draw three points
        this.gl.drawArrays(this.gl.LINE_STRIP, 0, vertices.length / 2);
        this.gl.disableVertexAttribArray(this.a_Position);
    }

    drawPoints(verticesInfoArr) {
        // Write the positions of vertices to a vertex shader
        this._initVertexBuffersForPoints(verticesInfoArr);
        // Draw three points
        this.gl.drawArrays(this.gl.POINTS, 0, verticesInfoArr.length / 7);

        this.gl.disableVertexAttribArray(this.a_Position);
        this.gl.disableVertexAttribArray(this.a_Size);
        this.gl.disableVertexAttribArray(this.a_Color);
    }

    _initVertexBuffersForPoints(verticesInfoArr)
    {
        const gl = this.gl
        const FSIZE = verticesInfoArr.BYTES_PER_ELEMENT;

        // Write date into the buffer object
        gl.bufferData(gl.ARRAY_BUFFER, verticesInfoArr, gl.STATIC_DRAW);

        // Assign the buffer object to a_Position variable
        gl.vertexAttribPointer(this.a_Position, 2, gl.FLOAT, false, FSIZE * 7, 0);
        gl.vertexAttribPointer(this.a_Size,     1, gl.FLOAT, false, FSIZE * 7, FSIZE * 2);
        gl.vertexAttribPointer(this.a_Color,    4, gl.FLOAT, false, FSIZE * 7, FSIZE * 3);
    
        // Enable the assignment to a_Position variable
        gl.enableVertexAttribArray(this.a_Position);
        gl.enableVertexAttribArray(this.a_Size);
        gl.enableVertexAttribArray(this.a_Color);
    }

}
