use super::{ErrorApiImpl, HandleTypeInfo, ManagedTypeApiImpl};

pub trait CallValueApi: HandleTypeInfo {
    type CallValueApiImpl: CallValueApiImpl
        + HandleTypeInfo<
            ManagedBufferHandle = Self::ManagedBufferHandle,
            BigIntHandle = Self::BigIntHandle,
            BigFloatHandle = Self::BigFloatHandle,
            EllipticCurveHandle = Self::EllipticCurveHandle,
        >;

    fn call_value_api_impl() -> Self::CallValueApiImpl;
}

pub trait CallValueApiImpl: ErrorApiImpl + ManagedTypeApiImpl + Sized {
    fn check_not_payable(&self);

    /// Retrieves the REWA call value from the VM.
    /// Will return 0 in case of an DCDT transfer (cannot have both REWA and DCDT transfer simultaneously).
    fn load_rewa_value(&self, dest_handle: Self::BigIntHandle);

    /// Loads all DCDT call values into a managed vec. Overwrites destination.
    fn load_all_dcdt_transfers(&self, dest_handle: Self::ManagedBufferHandle);

    /// Gets the total number of DCDT transfers (Fungible/SFT/NFT).
    ///
    /// It is redundant, since the number can also be retrieved from `load_all_dcdt_transfers`,
    /// but it is easier and cheaper to call when the content of those transfers is of no interest.
    fn dcdt_num_transfers(&self) -> usize;
}
