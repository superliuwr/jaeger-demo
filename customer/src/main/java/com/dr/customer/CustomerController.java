package com.dr.customer;

import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class CustomerController {

    private CustomerService customerService;

    public CustomerController(CustomerService customerService) {
      this.customerService = customerService;
    }

    @RequestMapping("/customer")
    public Customer get(@RequestParam(value="customer", defaultValue="") String id) {
        Customer customer = customerService.getCustomer(id);

        return customer;
    }

}