# Sszgen

`sszgen` is a code generator for the [SimpleSerialize (SSZ)](https://github.com/ethereum/eth2.0-specs/blob/dev/specs/simple-serialize.md)
protocol. Since SSZ requires the codecs to be aware of the data structures that are being encoded/decoded, 
it makes sense to generate the code in advance instead of doing it during runtime. 

## Saszy

The generated codecs will require the base library [`saszy`](https://github.com/holiman/saszy).

## Status

This is work in progress.