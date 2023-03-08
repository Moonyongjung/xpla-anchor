use cosmwasm_std::{MessageInfo};

use crate::{ContractError, state::{Config}};

/// check the owner of contract
pub fn check_owner(info: &MessageInfo, config: &Config) -> Result<String, ContractError> {
    if info.sender == config.owner {
        return Ok(config.owner.to_string());
    }

    Err(ContractError::Unauthorized {})
}