use cosmwasm_schema::write_api;

use xpla_anchor_contract::{msg::ExecuteMsg, msg::InstantiateMsg, msg::QueryMsg};

fn main() {
    write_api! {
        instantiate: InstantiateMsg,
        execute: ExecuteMsg,
        query: QueryMsg,
    }
}
