const express = require('express')
const bent = require('bent')
const { initTracerFromEnv } = require("jaeger-client")
const opentracing = require('opentracing')

const port = process.env.PORT || 8083
const serviceName = process.env.SERVICE_NAME || 'route'

const tracer = initTracer(serviceName)
opentracing.initGlobalTracer(tracer)

// ----- Express handlers -----
async function getRoute (req, res) {
  const tracer = opentracing.globalTracer()
  const span = tracer.startSpan('getRoute', { childOf: req.span })

  const pickup = req.query.pickup
  const dropoff = req.query.dropoff

  const customerInBaggage = span.getBaggageItem('customer')

  span.log({
      'event': 'request_params_parsed',
      'pickup': pickup,
      'dropoff': dropoff,
      'customer': customerInBaggage
  })

  const delay = await fetchDelay(span)
  await sleep(delay)

  const response = {
    'Pickup': pickup,
    'Dropoff': dropoff,
    'ETA': (Math.floor(Math.random() * 10) + 1) * (1000000 * 1000 * 60),
  }

  span.setTag('delay', delay)
  span.setTag('response', response)

  span.finish()

  res.json(response)
}

// ----- Calling another API -----
async function fetchDelay(parentSpan) {
  const tracer = opentracing.globalTracer()
  const span = tracer.startSpan('fetchDelay', { childOf: parentSpan })
  span.log({ event: 'fetch_delay', message: 'about to fetch delay for route service' })

  const service = process.env.DELAY_SERVICE_HOST || 'route-delay'
  const servicePort = process.env.DELAY_SERVICE_PORT || '8084'

  const url = `http://${service}:${servicePort}/delay`

  const headers = {}
  tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers)

  const request = bent('string', headers)

  let response = { delay: 500 }
  try {
    response = await request(url)
  } catch (e) {
    span.log({
      event: 'error',
      err: e.message
    })
  }
  
  span.finish()
  return response.delay || 500
}

// ----- Tracing initialization -----
function initTracer(serviceName) {
  const config = {
    serviceName: serviceName,
    // Sample every request
    sampler: {
      type: "const",
      param: 1
    },
    reporter: {
      logSpans: true,
    }
  }

  const options = {
    // Tracer level tags
    tags: {
      'app.name': serviceName,
      'app.version': 'v1.0.0'
    },
    logger: {
      info(msg) {
        console.log("INFO ", msg)
      },
      error(msg) {
        console.log("ERROR", msg)
      }
    }
  }
  
  return initTracerFromEnv(config, options)
}

// ----- Tracing Middleware -----
function tracingMiddleWare(req, res, next) {
  const tracer = opentracing.globalTracer()
  // Extracting the tracing headers from the incoming http request
  const wireCtx = tracer.extract(opentracing.FORMAT_HTTP_HEADERS, req.headers)
  // Creating our span with context from incoming request
  const span = tracer.startSpan(req.path, { childOf: wireCtx })
  // Use the log api to capture a log
  span.log({ event: 'request_received' })

  // Use the setTag api to capture standard span tags for http traces
  span.setTag(opentracing.Tags.HTTP_METHOD, req.method)
  span.setTag(opentracing.Tags.SPAN_KIND, opentracing.Tags.SPAN_KIND_RPC_SERVER)
  span.setTag(opentracing.Tags.HTTP_URL, req.path)

  // include trace ID in headers so that we can debug slow requests we see in
  // the browser by looking up the trace ID found in response headers
  const responseHeaders = {}
  tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, responseHeaders)
  res.set(responseHeaders)

  // add the span to the request object for any other handler to use the span
  Object.assign(req, { span })

  // finalize the span when the response is completed
  const finishSpan = () => {
    if (res.statusCode >= 500) {
      // Force the span to be collected for http errors
      span.setTag(opentracing.Tags.SAMPLING_PRIORITY, 1)
      // If error then set the span to error
      span.setTag(opentracing.Tags.ERROR, true)

      // Response should have meaning info to further troubleshooting
      span.log({ event: 'error', message: res.statusMessage })
    }
    // Capture the status code
    span.setTag(opentracing.Tags.HTTP_STATUS_CODE, res.statusCode)
    span.log({ event: 'request_end' })
    span.finish()
  }

  res.on('finish', finishSpan)

  next()
}

// ------ Utils -----
function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms))
}

// ----- App -----
const app = express()
app.use(tracingMiddleWare)
app.get('/route', getRoute)
app.disable('etag')
app.listen(port, () => {
  console.log('Route app listening on port ' + port)
})