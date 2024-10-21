---
theme: dashboard
title: points
---

```js
import { dotplot } from "/components/dotplot.js";
const sql = DuckDBClient.sql({
    points: FileAttachment(`/data/points-xy-${observable.params.num}.csv`)
})
```

```sql id=points
SELECT * FROM points
```

You asked for ${observable.params.num} points!

<div class="card">
    ${points.numRows.toLocaleString()} points
    ${x1y1}
</div>

```js
const x1y1 = dotplot(points, "x", "y");
```
