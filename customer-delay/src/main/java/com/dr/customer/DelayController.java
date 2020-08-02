package com.dr.customer;

import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class DelayController {

    @RequestMapping("/delay")
    public long get() {
        return (long) (Math.random() * 500 + 200);
    }

}