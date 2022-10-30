
class MapWay {

    vertices
    wayTags = new Map()
    color

    constructor(nodesNb) {
        this.vertices = new Float32Array(nodesNb)
    }

    addTag(key, value) {
        this.wayTags.set(key, value)
    }

    getTag(key) {
        this.wayTags.get(key)
    }
}
