//! Arithmetic helpers exercising free fns, impl methods, generics, and a nested
//! module — plus an in-file #[cfg(test)] module that must NOT be indexed.

/// Adds two integers, saturating on the sign of the operands (branchy so D3 has
/// real branches to cover).
pub fn add(a: i32, b: i32) -> i32 {
    if a > 0 && b > 0 {
        a + b
    } else {
        a.wrapping_add(b)
    }
}

/// Subtracts b from a.
pub fn sub(a: i32, b: i32) -> i32 {
    a - b
}

/// A private helper that is still indexed (non-pub functions are first-class
/// targets — the reason in-file mod injection is required for D5).
fn private_helper(x: i32) -> i32 {
    if x < 0 {
        -x
    } else {
        x
    }
}

/// A generic free function with a where-style bound, multi-line signature.
pub fn max_of<T>(a: T, b: T) -> T
where
    T: PartialOrd,
{
    if a >= b {
        a
    } else {
        b
    }
}

/// A stateful accumulator.
pub struct Calculator {
    total: i32,
}

impl Calculator {
    /// Constructs a zeroed calculator.
    pub fn new() -> Self {
        Calculator { total: 0 }
    }

    /// Folds a slice into the running total, returning the new total. The
    /// multi-line signature is exactly what the line-based path mis-bounds and
    /// the tree-sitter path captures precisely.
    pub fn compute(
        &mut self,
        values: &[i32],
    ) -> i32 {
        for v in values {
            self.total = add(self.total, private_helper(*v));
        }
        self.total
    }
}

/// A nested module whose function is indexed with a `nested::` qualifier.
pub mod nested {
    /// Doubles its argument.
    pub fn double(n: i32) -> i32 {
        n * 2
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn add_branches() {
        assert_eq!(add(2, 3), 5);
        assert_eq!(add(-1, -2), -3);
    }

    #[test]
    fn sub_works() {
        assert_eq!(sub(5, 2), 3);
    }

    #[test]
    fn private_helper_is_reachable() {
        assert_eq!(private_helper(-4), 4);
        assert_eq!(private_helper(4), 4);
    }

    #[test]
    fn compute_folds() {
        let mut c = Calculator::new();
        assert_eq!(c.compute(&[1, 2, 3]), 6);
    }

    #[test]
    fn max_and_double() {
        assert_eq!(max_of(1, 9), 9);
        assert_eq!(nested::double(4), 8);
    }
}
