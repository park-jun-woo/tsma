package com.example.calc;

import static org.junit.jupiter.api.Assertions.assertEquals;

import org.junit.jupiter.api.Test;

class CalculatorTest {

    @Test
    void testAdd() {
        Calculator calc = new Calculator(10);
        assertEquals(15, calc.add(2, 3));
    }

    @Test
    void testClassify() {
        Calculator calc = new Calculator(0);
        assertEquals("high", calc.classify(5, 3));
        assertEquals("low", calc.classify(1, 3));
        assertEquals("equal", calc.classify(3, 3));
    }
}
