const express = require('express');
const initTracer = require('./tracing').initTracer;
const { Tags, FORMAT_HTTP_HEADERS } = require('opentracing');

const port = 8083;
const app = express()
const tracer = initTracer('route');

app.listen(port, () => {
    console.log('Route app listening on port ' + port);
})

app.get('/route', async (req, res) => {
    const parentSpanContext = tracer.extract(FORMAT_HTTP_HEADERS, req.headers)

    const span = tracer.startSpan('http_server', {
        childOf: parentSpanContext,
        tags: {[Tags.SPAN_KIND]: Tags.SPAN_KIND_RPC_SERVER}
    });

    const pickup = req.query.pickup;
    const dropoff = req.query.dropoff;

    span.log({
        'event': 'route',
        'pickup': pickup,
        'dropoff': dropoff
    });

    // let greeting = span.getBaggageItem('greeting') || 'Hello';

    const route = {
      'Pickup': pickup,
      'Dropoff': dropoff,
      'ETA': Math.floor(Math.random() * 10) + 1,
    };

    const delay = Math.floor(Math.random() * 2000) + 1000;
    await sleep(delay);

    span.finish();

    res.json(route);
})

function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}