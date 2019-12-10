### Multicall

Wrapper for MakerDao's [Multicall](https://github.com/makerdao/multicall) which batches calls to contract
view functions into a single call and reads all the state in one EVM round-trip.