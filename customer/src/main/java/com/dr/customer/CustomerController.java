package com.dr.customer;

import java.net.URI;
import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.Map;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.client.RestTemplate;
import org.springframework.web.util.UriComponentsBuilder;

import io.opentracing.Scope;
import io.opentracing.Span;
import io.opentracing.Tracer;

@RestController
public class CustomerController {
    private static final Map<String, Customer> demoCustomers = new LinkedHashMap<String, Customer>();

    static {
        demoCustomers.put("123", new Customer("123", "Rachel's Floral Designs", "115,277"));
        demoCustomers.put("567", new Customer("567", "Amazing Coffee Roasters", "211,653"));
        demoCustomers.put("392", new Customer("392", "Trom Chocolatier", "577,322"));
        demoCustomers.put("731", new Customer("731", "Japanese Desserts", "728,326"));
    }

    @Autowired
    private RestTemplate restTemplate;

    @Autowired
    private Tracer tracer;

    @GetMapping("/customer")
    public Customer get(@RequestParam(value="customer", defaultValue="") String id) {
        try (Scope scope = tracer.buildSpan("get-customer-handler").startActive(true)) {
          Span span = scope.span();
          Map<String, String> fields = new LinkedHashMap<>();
          fields.put("event", "request_params_parsed");
          fields.put("customer_id", id);
          span.log(fields);

          Customer customer = demoCustomers.get(id);

          if (customer == null) {
            customer = demoCustomers.get("123");
          }
      
          long delay = fetchDelay();
      
          try {
            Thread.sleep(delay);
          } catch (InterruptedException e) {
            e.printStackTrace();
          }
      
          span.setTag("response", customer.toString());
          
          return customer;
      }
    }

    private long fetchDelay() {
        try (Scope scope = tracer.buildSpan("fetch-delay").startActive(true)) {
            Span span = scope.span();
            span.log("fetching delay for customer service");

            String serviceName = System.getenv("DELAY_SERVICE_HOST");
            String servicePort = System.getenv("DELAY_SERVICE_PORT");
            if (serviceName == null) {
                serviceName = "customer-delay";
            }
            if (servicePort == null) {
                servicePort = "8085";
            }

            String urlPath = "http://" + serviceName + ":" + servicePort + "/delay";

            URI uri = UriComponentsBuilder
                    .fromHttpUrl(urlPath)
                    .build(Collections.emptyMap());

            ResponseEntity<Long> response = restTemplate.getForEntity(uri, Long.class);

            return response.getBody();
        }
    }
}