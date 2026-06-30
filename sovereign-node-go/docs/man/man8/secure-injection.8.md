# SECURE-INJECTION(8) | Sovereign Node Manual

## NAME
secure-injection - Low-Level Code Injection and Epoch Management

## DESCRIPTION
The **Secure Injection Layer** (managed by `secure_injection.py`) is a system-level component that handles the atomic injection of code modifications into the running Sovereign Node. It is the primary mechanism for transitioning between system "Epochs".

## FORENSIC GUARD
All injections are monitored by the **Forensic Guard**. This layer ensures that meta-programming cycles do not violate system integrity or security constraints.

## LOGIC
- **Meta-Hook**: Initializes the connection between the MetaGo engine and the running process.
- **Epoch Transition**: Manages the shift from one system state to the next (e.g., Epoch 3 to Epoch 4).

## FILES
- **secure_injection.py**: The primary script for injection logic.
- **~/Jovian_Archives**: Storage for forensic logs related to injection events.

## SEE ALSO
metago(1), gsh(1)
