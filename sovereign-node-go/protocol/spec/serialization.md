# Canonical Binary Serialization (v1)

## Endianness
All numeric variables are encoded in BIG-ENDIAN format.

## Floats
- All floating-point fields must be encoded as standard IEEE-754 double-precision float64 values.
- Conversions must occur directly on bit representations (`math.Float64bits`) to preserve absolute cross-platform consistency.

## Ordering Rules
- **Protein States:** Sorted lexicographically by `ID` (ascending).
- **Failures:** Sorted lexicographically (ascending).
- **Tissues:** Sorted by key names (ascending).
- **Conditions:** Sorted primarily by Match Score descending, and secondarily by condition Name lexicographically (ascending) in case of equal scores.

## Binary Layout Structure
### STATE
`[ID_LEN:2][ID_BYTES][ACTIVATION:8][STABILITY:8][BINDING:8]`

### FAILURES
`[COUNT:2][STRINGS...]`

### TISSUES
`[COUNT:2][KEY][VALUE]`

### CONDITIONS
`[COUNT:2][NAME][SCORE]`
