package com.dr.customer;

import java.util.HashMap;
import java.util.Map;

import org.springframework.stereotype.Service;

@Service
public class CustomerService {
  private static final Map<String, Customer> demoCustomers = new HashMap<String, Customer>();

  static {
    demoCustomers.put("123", new Customer("123", "Rachel's Floral Designs", "115,277"));
    demoCustomers.put("567", new Customer("567", "Amazing Coffee Roasters", "211,653"));
    demoCustomers.put("392", new Customer("392", "Trom Chocolatier", "577,322"));
    demoCustomers.put("731", new Customer("731", "Japanese Desserts", "728,326"));
  }

  public Customer getCustomer(String id) {
    Customer customer = demoCustomers.get(id);

    if (customer == null) {
      customer = demoCustomers.get("123");
    }

    long delay = (long) (Math.random() * 5000 + 1000);

    try {
      Thread.sleep(delay);
    } catch (InterruptedException e) {
      e.printStackTrace();
    }

    return customer;
  }
}