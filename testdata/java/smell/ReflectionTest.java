package com.example.calc;

import java.lang.reflect.Field;
import java.lang.reflect.Method;

import org.junit.jupiter.api.Test;

class ReflectionTest {

    @Test
    void reachesPrivateState() throws Exception {
        Calculator calc = new Calculator(7);
        Field base = Calculator.class.getDeclaredField("base");
        base.setAccessible(true);
        base.set(calc, 99);

        Method add = Calculator.class.getDeclaredMethod("add", int.class, int.class);
        add.setAccessible(true);
        add.invoke(calc, 1, 2);
    }
}
