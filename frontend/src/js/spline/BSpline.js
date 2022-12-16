
function findSpan(n, k, t, knot_vector)
{
    if (Math.round(t * 1000000) === Math.round(knot_vector[n + 1] * 1000000))
        return n; /* Special case */
    /* Do binary search */
    let low = k;
    let high = n + 1;
    let mid = Math.floor((low + high) / 2);
    while ((t < knot_vector[mid]) || (t >= knot_vector[mid + 1]))
    {
        if (t < knot_vector[mid])
            high = mid;
        else
            low = mid;
        mid = Math.floor((low + high) / 2);
    }
    return mid;
}

function basisFuncs(i, t, k, knot_vector, N)
{
    let left = new Array(k + 1);
    let right = new Array(k + 1);
    let saved, temp;
    N[0] = 1.0;
    for (let j = 1; j <= k; j++)
    {
        left[j] = t - knot_vector[i + 1 - j];
        right[j] = knot_vector[i + j] - t;
        saved = 0.0;
        for (let r = 0; r < j; r++)
        {
            temp = N[r] / (right[r + 1] + left[j - r]);
            N[r] = saved + right[r + 1] * temp;
            saved = left[j - r] * temp;
        }
        N[j] = saved;
    }
    return (N);
}

function calculateLineSpline(points, countSplinePoints, splineOrder, ) {
    let span, i, j;
    let pt;
    let t, dt;
    const p = splineOrder;

    let pointsCtr = []

    // if (p >= pointsCtr.length)
    //     return ;

    console.log(pointsCtr.length)
    console.log(points.length)

    for (let i = 0; i < points.length / 2; ++i) {
        pointsCtr.push(new Point(points[i * 2], points[i * 2 + 1]))
    }

    // calculating the knot vector
    let knot_vector = new Array(pointsCtr.length + p + 1);
    for (i = 0; i <= p; ++i)
        knot_vector[i] = 0;
    for (i = p + 1; i <= pointsCtr.length - 1; ++i)
        knot_vector[i] = i - p;
    for (i = pointsCtr.length; i <= pointsCtr.length + p; ++i)
        knot_vector[i] = pointsCtr.length - p;
    let t_max = pointsCtr.length - p;

    const N = countSplinePoints;
    let pointsSpline = new Array(N);

    // calculating the values of a parametric function in points
    if (pointsCtr.length > 1)
    {
        dt = (t_max - pointsCtr[0].t) / (N - 1);
        t = pointsCtr[0].t;
        for (i = 0; i < N; i++)
        {
            let x = 0, y = 0;
            let basis_func = new Array(p + 1);
            span = findSpan(pointsCtr.length - 1, p, t, knot_vector);
            basisFuncs(span, t, p, knot_vector, basis_func);
            for (let l = 0; l < p + 1; l++)
            {
                x += basis_func[l] * pointsCtr[span - p + l].x;
                y += basis_func[l] * pointsCtr[span - p + l].y;
            }
            pt = new Point(x, y);
            pointsSpline[i] = pt;
            t += dt;
        }
    }

    // filling in an array for rendering
    let verticesSpline = new Float32Array(pointsSpline.length * 2);
    for (i = 0; i < pointsSpline.length; i++) {
        verticesSpline[i * 2] = pointsSpline[i].x;
        verticesSpline[i * 2 + 1] = pointsSpline[i].y;
    }
    return verticesSpline
}