use crate::dharitri_sc::{
    codec::{top_encode_to_vec_u8_or_panic, TopEncode},
    types::heap::Address,
};
use dharitri_chain_vm::tx_mock::TxTokenTransfer;
use num_traits::Zero;

pub struct ScCallDenali {
    pub(crate) from: Address,
    pub(crate) to: Address,
    pub(crate) rewa_value: num_bigint::BigUint,
    pub(crate) dcdt: Vec<TxTokenTransfer>,
    pub(crate) function: String,
    pub(crate) arguments: Vec<Vec<u8>>,
    pub(crate) gas_limit: u64,
    pub(crate) gas_price: u64,
}

impl ScCallDenali {
    pub fn new(from: &Address, to: &Address, function: &str) -> Self {
        ScCallDenali {
            from: from.clone(),
            to: to.clone(),
            rewa_value: num_bigint::BigUint::zero(),
            dcdt: Vec::new(),
            function: function.to_owned(),
            arguments: Vec::new(),
            gas_limit: u64::MAX,
            gas_price: 0,
        }
    }

    pub fn add_rewa_value(&mut self, rewa_value: &num_bigint::BigUint) {
        self.rewa_value.clone_from(rewa_value);
    }

    pub fn add_dcdt_transfer(
        &mut self,
        token_id: &[u8],
        nonce: u64,
        dcdt_value: &num_bigint::BigUint,
    ) {
        self.dcdt.push(TxTokenTransfer {
            token_identifier: token_id.to_vec(),
            nonce,
            value: dcdt_value.clone(),
        });
    }

    pub fn add_argument<T: TopEncode>(&mut self, arg: &T) {
        self.arguments.push(top_encode_to_vec_u8_or_panic(arg));
    }

    pub fn set_gas_limit(&mut self, gas_limit: u64) {
        self.gas_limit = gas_limit;
    }

    pub fn set_gas_price(&mut self, gas_price: u64) {
        self.gas_price = gas_price;
    }
}

pub struct ScQueryDenali {
    pub(crate) to: Address,
    pub(crate) function: String,
    pub(crate) arguments: Vec<Vec<u8>>,
}

impl ScQueryDenali {
    pub fn new(to: &Address, function: &str) -> Self {
        ScQueryDenali {
            to: to.clone(),
            function: function.to_owned(),
            arguments: Vec::new(),
        }
    }

    pub fn add_argument<T: TopEncode>(&mut self, arg: &T) {
        self.arguments.push(top_encode_to_vec_u8_or_panic(arg));
    }
}

pub struct TxExpectDenali {
    pub(crate) out: Vec<Vec<u8>>,
    pub(crate) status: u64,
    pub(crate) message: String,
    // TODO: Add logs?
}

impl TxExpectDenali {
    pub fn new(status: u64) -> Self {
        TxExpectDenali {
            out: Vec::new(),
            status,
            message: String::new(),
        }
    }

    pub fn add_out_value<T: TopEncode>(&mut self, out_val: &T) {
        self.out.push(top_encode_to_vec_u8_or_panic(out_val));
    }

    pub fn set_message(&mut self, msg: &str) {
        self.message = msg.to_string();
    }
}
