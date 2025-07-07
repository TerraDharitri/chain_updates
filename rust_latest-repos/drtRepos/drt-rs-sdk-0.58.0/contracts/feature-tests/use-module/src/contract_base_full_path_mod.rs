dharitri_sc::imports!();

#[dharitri_sc::module]
pub trait ContractBaseFullPathTestModule: dharitri_sc::contract_base::ContractBase {
    #[endpoint]
    fn call_contract_base_full_path_endpoint(&self) {}
}
