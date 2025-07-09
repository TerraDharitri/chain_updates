#![no_std]

dharitri_sc::imports!();

#[dharitri_sc::contract]
pub trait SendTxRepeat {
    #[init]
    fn init(&self) {}

    #[payable("REWA")]
    #[endpoint]
    fn repeat(&self, to: ManagedAddress, amount: BigUint, times: usize) {
        for _ in 0..times {
            self.tx().to(&to).rewa(&amount).transfer();
        }
    }
}
