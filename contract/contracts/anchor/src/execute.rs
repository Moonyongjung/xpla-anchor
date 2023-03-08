#[cfg(not(feature = "library"))]
use cosmwasm_std::entry_point;
use cosmwasm_std::{DepsMut, Env, MessageInfo, Response, StdError};
use cw2::{set_contract_version, get_contract_version};

use crate::error::{ContractError};
use crate::handler::{check_owner};
use crate::msg::{ExecuteMsg, InstantiateMsg, MigrateMsg, AnchoringMsg};
use crate::state::{Config, CONFIG, ANCHORING, BlockData, LATEST};

// version info for migration info
const CONTRACT_NAME: &str = env!("CARGO_PKG_NAME");
const CONTRACT_VERSION: &str = env!("CARGO_PKG_VERSION");

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn instantiate(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    _msg: InstantiateMsg,
) -> Result<Response, ContractError> {
    let config = Config {
        owner: info.sender,
    };

    set_contract_version(deps.storage, CONTRACT_NAME, CONTRACT_VERSION)?;

    CONFIG.save(deps.storage, &config)?;
    LATEST.save(deps.storage, &"0".to_string())?;

    Ok(Response::new()
        .add_attribute("method", "instantiate")
        .add_attribute("contract_owner", config.owner)
    )        
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn execute(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response, ContractError> {
    let config = CONFIG.load(deps.storage)?;
    
    match msg {
        ExecuteMsg::Anchoring(msg) => anchoring(deps, info, env, config, msg),
    }
}

// execute anchoring
pub fn anchoring(
    deps: DepsMut,
    info: MessageInfo,
    _env: Env,
    config: Config,
    msg: AnchoringMsg,
) -> Result<Response, ContractError> {
    check_owner(&info, &config)?;

    let anchoring = msg.data;
    
    let _: Vec<_> = anchoring
        .iter()
        .map(|x| {
            let recorded = BlockData {
                block_hash: x.block_hash.to_string(),
                data_merkle: x.data_merkle.to_string(),
                timestamp: x.timestamp.to_string(),
            };
            ANCHORING.save(deps.storage, x.height.to_string(), &recorded)
        })
        .collect();
    
    LATEST.save(deps.storage, &msg.latest)?;

    Ok(Response::new()
        .add_attribute("method", "anchoring")
    )
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn migrate(
    deps: DepsMut, 
    _env: Env, 
    _msg: MigrateMsg
) -> Result<Response, ContractError> {
    let ver = get_contract_version(deps.storage)?;
    if ver.contract != CONTRACT_NAME {
        return Err(StdError::generic_err("Can only upgrade from same type").into());
    }

    #[allow(clippy::cmp_owned)]
    if ver.version >= CONTRACT_VERSION.to_string() {
        return Err(StdError::generic_err("Cannot upgrade from a newer version").into());
    }

    // set the new version
    cw2::set_contract_version(deps.storage, CONTRACT_NAME, CONTRACT_VERSION)?;

    Ok(Response::default())
}