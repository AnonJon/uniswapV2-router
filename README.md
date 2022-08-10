# Uniswap V2 Router

Finding the rate and best trade route using the UniswapV2 factory/pool contracts.

## Server

Start RPC server

```
make
```

## RPC Methods

- Quote
- Rate

## CLI Examples

```bash
go run main.go quote --tokenIn 0xb4efd85c19999d84251304bda99e90b92300bd93 --tokenOut 0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2 --amount 1
go run main.go rate --tokenIn 0xb4efd85c19999d84251304bda99e90b92300bd93 --tokenOut 0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2 --amount 1
```

## Test

```
go test -v

```

#### Check Coverage

```
go test -coverprofile= ./...

```

## Current Strategy

Currently, the application will find a list of every pool token A/B are part of and pass an array of arrays of data type Path. Each route is then calculated on rate (getReserves function in pair contract which returns each tokens amount in wei) from token A/B and if length is greater than 1, rate is passed forwards as amount for next pair iteration. Those quotes are then saved and calculated at the end for best rate for each route.

## Next Steps

- Because the need to check for new pools is a constant task (also sub-optimal to call blockchain on every query) the use of the Uniswap subgraph to query a list of all current pools would seem to be best when currating list of possible paths to take to get quote.
- The returned quote does not include gas fees and 0.3% pool fee, so another step would be to include those into the final calculation.
- Testing needs to be done on the main algo and should be broken down into an interface to mock out client calls
