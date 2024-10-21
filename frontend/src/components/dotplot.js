import {resize} from "npm:@observablehq/stdlib";
import * as d3 from "npm:d3";
import * as Plot from "npm:@observablehq/plot";

export const dotplot = (input, x, y) => {
  const opacity = d3.scaleLinear([0, 9], [0, 255])
  return resize((width) =>
    Plot.plot({
      width,
      height: 200,
      marks: [
        Plot.dot(input, {
          x: x,
          y: y,
          r: (d) => d.r,
          stroke: (d) => d3.schemeTableau10[Number(d.r)],
          strokeOpacity: (d) => opacity(d.r),
          strokeWidth: 1,
        }),
      ],
    })
  );
};
