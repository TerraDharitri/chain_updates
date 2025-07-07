use dharitri_sc::api::{CallValueApi, CallValueApiImpl, HandleConstraints};

use crate::api::{VMHooksApi, VMHooksApiBackend};

impl<VHB: VMHooksApiBackend> CallValueApi for VMHooksApi<VHB> {
    type CallValueApiImpl = Self;

    fn call_value_api_impl() -> Self::CallValueApiImpl {
        Self::api_impl()
    }
}

impl<VHB: VMHooksApiBackend> CallValueApiImpl for VMHooksApi<VHB> {
    fn check_not_payable(&self) {
        self.with_vm_hooks(|vh| vh.check_no_payment())
    }

    fn load_rewa_value(&self, dest: Self::BigIntHandle) {
        self.assert_live_handle(&dest);
        self.with_vm_hooks(|vh| vh.big_int_get_call_value(dest.get_raw_handle_unchecked()));
    }

    fn load_all_dcdt_transfers(&self, dest_handle: Self::ManagedBufferHandle) {
        self.assert_live_handle(&dest_handle);
        self.with_vm_hooks(|vh| {
            vh.managed_get_multi_dcdt_call_value(dest_handle.get_raw_handle_unchecked())
        });
    }

    fn dcdt_num_transfers(&self) -> usize {
        self.with_vm_hooks(|vh| vh.get_num_dcdt_transfers()) as usize
    }
}
