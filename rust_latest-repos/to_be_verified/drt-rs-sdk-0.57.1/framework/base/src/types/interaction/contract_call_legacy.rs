mod async_call;
mod async_call_promises;
mod contract_call_convert;
mod contract_call_exec;
mod contract_call_no_payment;
mod contract_call_trait;
mod contract_call_with_any_payment;
mod contract_call_with_rewa;
mod contract_call_with_rewa_or_single_dcdt;
mod contract_call_with_multi_dcdt;
mod contract_deploy;
mod typed_function_call;

pub use async_call::AsyncCall;
pub use async_call_promises::AsyncCallPromises;
pub use contract_call_no_payment::ContractCallNoPayment;
pub use contract_call_trait::{ContractCall, ContractCallBase};
pub use contract_call_with_any_payment::ContractCallWithAnyPayment;
pub use contract_call_with_rewa::ContractCallWithRewa;
pub use contract_call_with_rewa_or_single_dcdt::ContractCallWithRewaOrSingleDcdt;
pub use contract_call_with_multi_dcdt::ContractCallWithMultiDcdt;
pub use contract_deploy::{new_contract_deploy, ContractDeploy};
pub use typed_function_call::TypedFunctionCall;

/// Using max u64 to represent maximum possible gas,
/// so that the value zero is not reserved and can be specified explicitly.
/// Leaving the gas limit unspecified will replace it with `api.get_gas_left()`.
pub(crate) const UNSPECIFIED_GAS_LIMIT: u64 = u64::MAX;
