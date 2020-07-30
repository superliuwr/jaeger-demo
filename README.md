# Hot R.O.D. - Rides on Demand

This is a demo application based on [Jaeger's official demo application](https://github.com/uber/jaeger/tree/master/examples/hotrod) with the microservice `customer` rewritten in Java and `route` rewritten in Node.js.

## About Jaeger's official demo application

It is a demo application that consists of several microservices and illustrates the use of the OpenTracing API. It can be run standalone, but requires Jaeger backend to view the traces.

## Features

* Discover architecture of the whole system via data-driven dependency diagram
* View request timeline & errors, understand how the app works
* Find sources of latency, lack of concurrency
* Highly contextualized logging
* Use baggage propagation to
  * Diagnose inter-request contention (queueing)
  * Attribute time spent in a service
* Use open source libraries with OpenTracing integration to get vendor-neutral instrumentation for free

## Running

Run `docker-compose up -d` from the root to bring up all microservices and jaeger-all-in-one.

Jaeger UI can be accessed at http://localhost:16686.

Then open the frontend at http://127.0.0.1:8080

## Microservices

- frontend: Go
- customer: Java/Spring Boot Application
- driver: Go
- route: Node.js