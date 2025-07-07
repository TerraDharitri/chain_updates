use crate::{types::RawHandle, vm_err_msg, vm_hooks::VMHooksHandlerSource};
use num_traits::Zero;

use super::VMHooksManagedTypes;

pub trait VMHooksCallValue: VMHooksHandlerSource + VMHooksManagedTypes {
    fn check_not_payable(&self) {
        if self.input_ref().rewa_value > num_bigint::BigUint::zero() {
            self.vm_error(vm_err_msg::NON_PAYABLE_FUNC_REWA);
        }
        if self.dcdt_num_transfers() > 0 {
            self.vm_error(vm_err_msg::NON_PAYABLE_FUNC_DCDT);
        }
    }

    fn load_rewa_value(&self, dest: RawHandle) {
        let value = self.input_ref().received_rewa().clone();
        self.m_types_lock().bi_overwrite(dest, value.into());
    }

    fn load_all_dcdt_transfers(&self, dest_handle: RawHandle) {
        let transfers = self.input_ref().received_dcdt();
        self.m_types_lock()
            .mb_set_vec_of_dcdt_payments(dest_handle, transfers);
    }

    fn dcdt_num_transfers(&self) -> usize {
        self.input_ref().received_dcdt().len()
    }
}
