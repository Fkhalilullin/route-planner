
class MapFragment {

    minlat = 0.
    minlon = 0.
    maxlat = 0.
    maxlon = 0.
    nodes = new Map()
    ways = new Map()

    constructor(osmXml) {
        this.minlat = osmXml.querySelector("bounds").attributes.minlat.value
        this.minlon = osmXml.querySelector("bounds").attributes.minlon.value
        this.maxlat = osmXml.querySelector("bounds").attributes.maxlat.value
        this.maxlon = osmXml.querySelector("bounds").attributes.maxlon.value
        osmXml
            .querySelectorAll("node")
            .forEach(nodeXmlTag => this.addNode(nodeXmlTag))
        osmXml
            .querySelectorAll("way")
            .forEach(wayXmlTag => this.addWay(wayXmlTag))
    }

    addNode(nodeXmlTag) {
        const lat = nodeXmlTag.attributes.lat.value
        const lon = nodeXmlTag.attributes.lon.value
        this.nodes.set(nodeXmlTag.id, new MapNode(lat, lon))
    }

    addWay(wayXmlTag) {
        const nodesNb = this._countNodes(wayXmlTag)
        const wayId = wayXmlTag.id
        const way = new MapWay(nodesNb * 2)
        this.ways.set(wayId, way)

        let i = 0
        for (let elem of wayXmlTag.children) {
            if (elem.localName === "nd") {
                const nodeId = this.nodes.get(elem.attributes.ref.value)
                way.vertices[i    ] = nodeId.lon
                way.vertices[i + 1] = nodeId.lat
                i += 2
            } else if (elem.localName === "tag") {
                const tagValue = elem.attributes.v.value
                const tagKey = elem.attributes.k.value
                way.addTag(tagKey, tagValue)
            }
        }
    }

    _countNodes(wayXmlTag) {
        let count = 0;
        for (let elem of wayXmlTag.children) {
            if (elem.localName === "nd")
                count++
        }
        return count
    }

}