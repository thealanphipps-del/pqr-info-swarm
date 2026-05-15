# Sovereign Script: The RustScript Synthesis

The Sovereign Node integrates a high-performance Go-based implementation of **RustScript**—a synthesis of 68 years of language evolution (1958-2025).

## Design by Contract (DbC)
Inspired by Eiffel and RustScript, PQR enforces **Forensic Integrity Contracts**. Every state transition in the Ticketing Fabric is governed by mandatory pre-conditions and post-conditions.

### Contractual Ticket Commitment
```rust
// Ensure parent lineage exists before link
pre { parent.Exists() && child.IsOrphan() }

// Commit state mutation
body { 
    Fabric.Link(parent, child, "EVOLUTION") 
}

// Verify zero-divergence post-commit
post { parent.Children.Contains(child) }
```

## Effect Systems & Action Isolation
PQR utilizes **Effect Typing** to declare agent capabilities explicitly. This prevents "Black Box" agent behavior by isolating side effects like file mutation or network egress into audited streams.

### Agent Capability Declaration
```rust
effect FileMutation {
    fn write_source(path: string, content: byte[])
}

handler SelfHealingAgent performs FileMutation {
    // Agent execution is restricted to these effects
}
```

## High-Speed Registers (%q, %r)
Drawing from the MUSHcode heritage, PQR implements **Swarm Registers** for transient consensus. These allow agents to pass high-speed numeric (%q0-%q9) and string (%r0-%r9) signals across the mesh without the overhead of full ticket creation.
