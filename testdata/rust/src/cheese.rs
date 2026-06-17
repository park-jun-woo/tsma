//! A module whose production code is clean, but whose in-file #[cfg(test)]
//! module reaches into memory via unsafe / transmute / std::ptr — the D4
//! escape-hatch positives. Production code here is deliberately safe so the
//! detectors prove they only fire inside test scopes (false-positive zero).

/// Returns the raw byte view of a u32 as four bytes — a *safe*, legitimate API
/// (no unsafe), so production code never trips the smell detectors.
pub fn to_bytes(n: u32) -> [u8; 4] {
    n.to_le_bytes()
}

#[cfg(test)]
mod tests {
    use super::*;

    // TS-REFL-RS-001: an unsafe block forcing a memory reinterpretation.
    #[test]
    fn unsafe_block_cheese() {
        let n: u32 = 0x0102_0304;
        let b = unsafe {
            let p = &n as *const u32 as *const u8;
            *p
        };
        assert_eq!(b, to_bytes(n)[0]);
    }

    // TS-REFL-RS-002: std::mem::transmute reinterpret.
    #[test]
    fn transmute_cheese() {
        let bits: u32 = unsafe { std::mem::transmute(1.0f32) };
        assert!(bits > 0);
    }

    // TS-REFL-RS-003: std::ptr / core::ptr forced access + as_ptr.
    #[test]
    fn ptr_cheese() {
        let v = vec![1u8, 2, 3];
        let head = v.as_ptr();
        let first = unsafe { std::ptr::read(head) };
        assert_eq!(first, 1);
    }
}
