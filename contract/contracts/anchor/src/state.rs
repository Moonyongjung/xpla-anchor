use cosmwasm_std::{Addr};
use cosmwasm_schema::{cw_serde};
use cw_storage_plus::{Item, Map};

#[cw_serde]
pub struct Config {
    // the contract owner.
    pub owner: Addr,
}

impl Config {
    pub fn update(
        &mut self,
    ) -> &mut Self {
        self
    }
}

#[cw_serde]
pub struct BlockData {
    pub block_hash: String,
    pub data_merkle: String,
    pub timestamp: String,
}

pub const CONFIG: Item<Config> = Item::new("config");
pub const ANCHORING: Map<String, BlockData> = Map::new("anchoring");
pub const LATEST: Item<String> = Item::new("latest");