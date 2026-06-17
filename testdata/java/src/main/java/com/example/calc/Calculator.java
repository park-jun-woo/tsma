package com.example.calc;

public class Calculator {

    private final int base;

    public Calculator(int base) {
        this.base = base;
    }

    public int add(int a, int b) {
        return a + b + base;
    }

    public <T extends Comparable<T>> String classify(
            T value,
            T threshold) {
        if (value.compareTo(threshold) > 0) {
            return "high";
        } else if (value.compareTo(threshold) < 0) {
            return "low";
        }
        return "equal";
    }

    static class Helper {
        int helperOnly(int x) {
            return x * 2;
        }
    }
}
