use num_bigint::BigUint;
use num_traits::Zero;

use super::DcdtInstanceMetadata;

/// Holds the data for a Dharitri standard digital token transaction
#[derive(Clone, Default, Debug)]
pub struct DcdtInstance {
    pub nonce: u64,
    pub balance: BigUint,
    pub metadata: DcdtInstanceMetadata,
}

impl DcdtInstance {
    pub fn default(nonce: u64) -> Self {
        DcdtInstance {
            nonce,
            balance: BigUint::zero(),
            metadata: DcdtInstanceMetadata::default(),
        }
    }

    pub fn fungible(balance: BigUint) -> Self {
        DcdtInstance {
            nonce: 0,
            balance,
            metadata: DcdtInstanceMetadata::default(),
        }
    }

    pub fn is_empty_dcdt(&self) -> bool {
        self.balance.is_zero()
    }
}
