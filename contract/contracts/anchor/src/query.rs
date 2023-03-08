#[cfg(not(feature = "library"))]
use cosmwasm_std::entry_point;
use cosmwasm_std::{to_binary, Binary, Env, StdResult, Deps, StdError};
use crate::msg::{QueryMsg, BlockDataResponse, LatestBlockResponse};
use crate::state::{ANCHORING, LATEST};

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn query(
    deps: Deps,
    _env: Env,
    msg: QueryMsg
) -> StdResult<Binary> {
    match msg {
        QueryMsg::BlockData { height } => to_binary(&block_data(deps, height)?),
        QueryMsg::LatestBlock {} => to_binary(&latest_block(deps)?),
    }
}

// query configuration.
fn block_data(deps: Deps, height: String) -> StdResult<BlockDataResponse> {
    let anchoring = ANCHORING.may_load(deps.storage, height.to_string())?;
    if anchoring.is_none() {
        return Err(StdError::GenericErr { msg: "invalid block height".to_string() })
    }

    let anchoring = anchoring.unwrap();

    Ok(BlockDataResponse { 
        height, 
        block_hash: anchoring.block_hash, 
        data_merkle: anchoring.data_merkle, 
        timestamp: anchoring.timestamp 
    })
}

// query saved latest block
fn latest_block(deps: Deps) -> StdResult<LatestBlockResponse> {
    let latest_height = LATEST.load(deps.storage)?;
    Ok(LatestBlockResponse {
        latest_height
    })

}