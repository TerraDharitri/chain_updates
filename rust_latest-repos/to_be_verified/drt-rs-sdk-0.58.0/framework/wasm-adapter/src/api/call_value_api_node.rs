use super::VmApiImpl;
use dharitri_sc::api::{CallValueApi, CallValueApiImpl};

unsafe extern "C" {
    fn checkNoPayment();

    fn bigIntGetCallValue(dest: i32);

    fn managedGetMultiDCDTCallValue(resultHandle: i32);

    fn getNumDCDTTransfers() -> i32;
}

impl CallValueApi for VmApiImpl {
    type CallValueApiImpl = VmApiImpl;

    #[inline]
    fn call_value_api_impl() -> Self::CallValueApiImpl {
        VmApiImpl {}
    }
}

impl CallValueApiImpl for VmApiImpl {
    #[inline]
    fn check_not_payable(&self) {
        unsafe {
            checkNoPayment();
        }
    }

    fn load_rewa_value(&self, dest: Self::BigIntHandle) {
        unsafe {
            bigIntGetCallValue(dest);
        }
    }

    fn load_all_dcdt_transfers(&self, dest_handle: Self::ManagedBufferHandle) {
        unsafe {
            managedGetMultiDCDTCallValue(dest_handle);
        }
    }

    fn dcdt_num_transfers(&self) -> usize {
        unsafe { getNumDCDTTransfers() as usize }
    }
}
