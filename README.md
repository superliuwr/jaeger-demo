# Jaeger Demo

This is a demo application based on [one of Jaeger's official demo applications](https://github.com/uber/jaeger/tree/master/examples/hotrod) called `Hot R.O.D`.

It is an oversimplified version of Uber Eats.

In the homepage you will find a list of hardcoded customers, click on one of them to start a transaction.

The application locates the customer, finds ten nearby drivers and for each driver calculates the ETA. The route with least ETA will be returned.

A log is then shown in the homepage with the chosen driver's license plate and ETA(a request ID and the latency calculated on the frontend side are also displayed for demo purpose).

![UI](/docs/ui.png)

## Components

### frontend
It's the start point of the application. It hosts a web server for the UI and also servers the backend business logic.

The backend receives requests from the UI and sends requests to other components and returns the result to UI.

It's written in Go.

### customer
It's a Restful API application backed by Spring Boot. The application handles requests of fetching customer information. It calls `customer-delay` to get the delay value and delay the process accordingly.

It's written in Java and Spring Boot. It demonstrates how manual instrumentation works with Spring Boot.

### customer-delay
It's a Restful API application backed by Spring Boot. The API simply returns a delay value to the callers.

It's written in Java and Spring Boot. It demonstrates how automatic instrumentation with 3rd party framework works with Spring Boot.

### driver
It's a gRPC application providing driver's information. It's called by frontend and calls the mock `redis` component.

It's written in Go to demonstrated instrumentation for gRPC endpoints.

### route
It's a Restful API application backed by Express. The application handles requests of fetching route information for given two locations. It calls `route-delay` to get the delay value and delay the process accordingly.

It's written in Node.js and Express. It also demonstrates how Baggage works(the baggage item named `customer` set by `frontend`).

### route-delay
It's a Restful API application backed by Express. The API simply returns a delay value to the callers.

It's written in Node.js and Express.

## Workflow

![DAG](/docs/dag.png)

## How it is different from the official demo application

1. The microservice `customer` is rewritten in Java
2. The microservice `route` is rewritten in Node.js
3. Two new microservices `customer-delay` and `route-delay` were added to demonstrated tracing calling other services in Java and Node.js
4. The microservice `driver` is now separated into a standalone service
5. The microservices `frontend` and `driver` were simplified for beginners
6. No metrics are captured

## Features

* Discover architecture of the whole system via data-driven dependency diagram
* View request timeline & errors, understand how the app works
* Find sources of latency, lack of concurrency
* Highly contextualized logging
* Use baggage propagation to
  * Diagnose inter-request contention (queueing)
  * Attribute time spent in a service
* Use open source libraries with OpenTracing integration to get vendor-neutral instrumentation for free

## What's NOT covered

* Client side tracing
* Tracing inside DBs(built-in support or 3rd party wrappers)
* Tracing with Service Mesh(Istio for example)
* Metrics collection

## Running

1. Run `docker-compose up -d` from the root to bring up all microservices and jaeger-all-in-one.
2. Open the frontend at http://127.0.0.1:8080 and make a request by clicking on one of the customers.
3. Wait until you see the driver and ETA is returned. A message will be displayed like this: `HotROD T758836C arriving in 1min [req: 6388-1, latency: 3456ms]`.
4. Open Jaeger UI at http://localhost:16686. Use Firefox or Safari as there is currently a minor but annoying UI issue with Chrome when displaying spans.
5. Choose `frontend` from the `Service` dropbox and click `Find Traces`. You should see the trace of your last call to `HTTP Get /dispatch` appearing on the right.
6. Click on the trace and play with the spans.

![Traces](/docs/traces.png)

![Trace](/docs/trace.png)
