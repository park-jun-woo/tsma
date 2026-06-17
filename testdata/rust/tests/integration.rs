//! Integration tests live in a separate crate (the `tests/` directory) and can
//! only reach `pub` items — the D2 `tests/` fallback path. A clean test (no
//! escape hatches) so the smell detectors stay silent here too.

use tsma_rust_fixture::calc;

#[test]
fn add_via_integration() {
    assert_eq!(calc::add(40, 2), 42);
}

#[test]
fn double_via_integration() {
    assert_eq!(calc::nested::double(21), 42);
}
