# XPLA Anchor

The anchor provides secure to the private chain by sending transactions which include parameters such as block height, block hash and etc in their blocks to main chain. The private chain based on Tendermint core is anchoring to public chain as XPLA.

## Prerequisites

### Generate the config file (`./config.yaml`)
```yaml
Anchor:
    CollectBlockCount: 10
    RequestPeriod: 50

PublicChain:
    ChainID: "dimension_37-1"
    LCD: https://dimension-lcd.xpla.dev
    GasAdj: 1.75
    GasLimit:
    BroadcastMode: block

PrivateChain:
    ChainID: privatechain-1
    LCD: http://localhost:1317
```
- `CollectBlockCount`: The number of the private chain's blocks to be included in one anchoring transaction.
- `RequestPeriod`: The duration time to request the block info of the private chain. If a query is requested to the private chain too often, it may be blocked, so it is recommended to set the period with acceptable time.  
- `PublicChain`: The main chain as XPLA.
- `PrivateChain`: The private chain would be anchoring to the main chain.

The parameters that `ChainId` and `LCD` are mandatory, but `GasAdj`, `GasLimit` and `BroadcastMode` are optional.

### Generate the account of the main chain.
The owner of the anchor should generate the account with `axpla` balance. The anchor uses this account for sending transactions.

### Set the DB (optional)
If the owner of the anchor need to record logs by using database, prepare DB as `MySQL`. It is able to apply to the anchor by using flag (`--log db`) when start the gateway.

## Start

### Automatic
```sh
$ ./init.sh
```

### Manual
1. Installation
```sh 
$ cd xpla-anchor
$ make install
```

2. Initialization
```sh
$ anc init
$ anc key recover
```

3. Set the anchor contract
```sh
$ anc execute contract store
$ anc execute contract instantiate
```

4. Start the anchor gateway
```sh
$ anc execute start
```

## Interaction
### Query
The anchor can interact to the anchor contract by querying.
```sh
# Query the latest block is recorded in the anchor contract.
$ anc query contract latest

# Query the block info.
$ anc query contract block [block_height]
```

Also, can check base information of the account which used to the anchor.
```sh
# Can check balances.
$ anc query account balance

# Can check account info.
$ anc query account info
```

### Verify
The anchor can verify the consistency between block info that is recorded in the anchor contract and query response from the private chain. 

```sh
# Check the block height.
$ anc query verify [block_height]
```