use crate::dns_proxy;

dharitri_sc::imports!();

/// Standard smart contract module that deals with registering usernames in a DNS contract.
///
/// Dharitri usernames/herotags need to be requested by the beneficiary.
/// For a contract, this means that they need an endpoint via which to request a username from the DNS.
///
#[dharitri_sc::module]
pub trait DnsModule {
    #[payable("REWA")]
    #[only_owner]
    #[endpoint(dnsRegister)]
    fn dns_register(&self, dns_address: ManagedAddress, name: ManagedBuffer) {
        let payment = self.call_value().rewa().clone();
        self.tx()
            .to(&dns_address)
            .typed(dns_proxy::DnsProxy)
            .register(name)
            .rewa(payment)
            .async_call_and_exit();
    }
}
