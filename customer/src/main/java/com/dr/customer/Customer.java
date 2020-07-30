package com.dr.customer;

public class Customer {

    private final String id;
    private final String name;
    private final String location;

    public Customer(String id, String name, String location) {
        this.id = id;
        this.name = name;
        this.location = location;
    }

    public String getId() {
        return id;
    }

    public String getName() {
        return name;
    }

    public String getLocation() {
      return location;
    }
    
}