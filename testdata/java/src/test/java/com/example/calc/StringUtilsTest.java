package com.example.calc;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

import org.junit.jupiter.api.Test;

class StringUtilsTest {

    @Test
    void testIsBlank() {
        assertTrue(StringUtils.isBlank("   "));
    }

    @Test
    void testRepeat() {
        assertEquals("ababab", StringUtils.repeat("ab", 3));
    }
}
