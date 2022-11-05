
class ColorUtils {

    static hexToRgb(hex) {
        const bigint = parseInt(hex.substring(1), 16)
        const r = ((bigint >> 16) & 255) / 255
        const g = ((bigint >> 8) & 255) / 255
        const b = (bigint & 255) / 255
        return [r, g, b]
    }

    static getGradientColor(minColor, maxColor, minValue, maxValue, value) {
        const deltaValue = maxValue - minValue
        const [minR, minG, minB] = this.hexToRgb(minColor)
        const [maxR, maxG, maxB] = this.hexToRgb(maxColor)
        const deltaR = (maxR - minR) / deltaValue
        const deltaG = (maxG - minG) / deltaValue
        const deltaB = (maxB - minB) / deltaValue
        const valueFromMin = value - minValue
        const r = valueFromMin * deltaR + minR
        const g = valueFromMin * deltaG + minG
        const b = valueFromMin * deltaB + minB
        return [r, g, b]
    }
}
