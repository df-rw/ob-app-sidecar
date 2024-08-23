---
theme: dashboard
---

```js
const db = DuckDBClient.of({ data: FileAttachment("./data/rand-xy.csv") });
```

# Observable app

This application demonstrates an Observable Framework application that houses a
Go backend and communicates to it using HTMX.

## Simple API calls

<!-- a couple of calls to the backend api -->
<div class="grid grid-cols-2">
    <div class="card">
        <h2>server time</h2>
        <button
            hx-get="/api/now"
            hx-target="#now"
            hx-swap="innerHTML">hit me</button>
        <span id="now"></span>
    </div>
    <div class="card">
        <h2>server time + 10minutes</h2>
        <button
            hx-get="/api/then"
            hx-target="#then"
            hx-swap="innerHTML">hit me</button>
        <span id="then"></span>
    </div>
</div>

<!-- making sure deployment works with file attachments -->

```js
const data = db.sql`SELECT * FROM data`;
const div2 = FileAttachment("./data/rand-xy-div2.csv").csv({ typed: true });
```

<!-- and making sure deployment works with plot -->

```js
const plot = resize((width) =>
  Plot.plot({
    width,
    height: 200,
    x: {
      domain: [0, 100],
    },
    y: {
      domain: [0, 100],
    },
    marks: [
      Plot.axisX({
        ticks: d3.ticks(0, 100, 10),
      }),
      Plot.axisY({
        ticks: d3.ticks(0, 100, 10),
      }),
      Plot.dot(data, {
        x: "x",
        y: "y",
        stroke: "green",
      }),
      Plot.dot(div2, {
        x: "x",
        y: "y",
        stroke: "blue",
      }),
    ],
  })
);
```

<div class="card">
    ${plot}
</div>
