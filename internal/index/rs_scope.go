//ff:type feature=index type=model
//ff:what Holds an enclosing Rust scope (impl receiver or module) tied to a brace depth
package index

// rsScope records an enclosing impl/mod context opened at a given brace depth.
type rsScope struct {
	depth    int    // brace depth at which this scope's body was opened
	receiver string // impl type name, or "" for a plain module scope
	module   string // module name, or "" for an impl scope
	cfgTest  bool   // true if guarded by #[cfg(test)] (its fns are tests)
}
