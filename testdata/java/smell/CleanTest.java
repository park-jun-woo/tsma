package com.example.calc;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;

import org.junit.jupiter.api.Test;

class CleanTest {

    @Test
    void exercisesPublicApi() {
        Calculator calc = new Calculator(1);
        assertEquals(4, calc.add(1, 2));
        // The string "getDeclaredMethod" here is an argument, not a reflective call.
        assertFalse(StringUtils.isBlank("getDeclaredMethod"));
    }

    @Test
    void restoresGuardIsNotASmell() throws Exception {
        java.lang.reflect.Field f = Calculator.class.getClass().getFields()[0];
        // setAccessible(false) restores the guard; it must never fire.
        f.setAccessible(false);
    }
}
