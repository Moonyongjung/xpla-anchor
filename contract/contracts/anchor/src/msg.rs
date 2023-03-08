use cosmwasm_schema::{cw_serde, QueryResponses};

#[cw_serde]
pub struct InstantiateMsg {
}

#[cw_serde]
pub enum ExecuteMsg {
    Anchoring(AnchoringMsg),
}

#[cw_serde]
#[derive(QueryResponses)]
pub enum QueryMsg {
    #[returns(BlockDataResponse)]
    BlockData {
        height: String,
    },

    #[returns(LatestBlockResponse)]
    LatestBlock {},
}

// msgs
#[cw_serde]
pub struct AnchoringMsg {
    pub data: Vec<Data>,
    pub latest: String,
}

#[cw_serde]
pub struct Data {
    pub height: String,
    pub block_hash: String,
    pub data_merkle: String,
    pub timestamp: String,
}

// responses
#[cw_serde]
pub struct BlockDataResponse {
    pub height: String,
    pub block_hash: String,
    pub data_merkle: String,
    pub timestamp: String,
}

#[cw_serde]
pub struct LatestBlockResponse {
    pub latest_height: String,
}

#[cw_serde]
pub struct MigrateMsg {}

