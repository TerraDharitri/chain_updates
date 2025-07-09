use crate::{
    api::{use_raw_handle, BlockchainApiImpl, CallTypeApi, StaticVarApiImpl, StorageWriteApi},
    codec::TopDecodeMulti,
    contract_base::SendRawWrapper,
    formatter::SCLowerHex,
    types::{
        decode_result, AsyncCall, AsyncCallPromises, BigUint, DcdtTokenPayment, ManagedBuffer,
        ManagedBufferBuilder, ManagedType, ManagedVec, Tx, TRANSFER_EXECUTE_DEFAULT_LEFTOVER,
    },
};

use super::{ContractCallNoPayment, ContractCallWithRewa, UNSPECIFIED_GAS_LIMIT};
use crate::api::managed_types::handles::HandleConstraints;

impl<SA, OriginalResult> ContractCallWithRewa<SA, OriginalResult>
where
    SA: CallTypeApi + 'static,
{
    pub fn resolve_gas_limit(&self) -> u64 {
        if self.basic.explicit_gas_limit == UNSPECIFIED_GAS_LIMIT {
            SA::blockchain_api_impl().get_gas_left()
        } else {
            self.basic.explicit_gas_limit
        }
    }

    #[inline]
    pub fn get_back_transfers(&self) -> (BigUint<SA>, ManagedVec<SA, DcdtTokenPayment<SA>>) {
        let dcdt_transfer_value_handle: SA::BigIntHandle =
            use_raw_handle(SA::static_var_api_impl().next_handle());
        let call_value_handle: SA::BigIntHandle =
            use_raw_handle(SA::static_var_api_impl().next_handle());

        SA::blockchain_api_impl().managed_get_back_transfers(
            dcdt_transfer_value_handle.get_raw_handle(),
            call_value_handle.get_raw_handle(),
        );

        unsafe {
            (
                BigUint::from_raw_handle(call_value_handle.get_raw_handle()),
                ManagedVec::from_raw_handle(dcdt_transfer_value_handle.get_raw_handle()),
            )
        }
    }

    pub fn to_call_data_string(&self) -> ManagedBuffer<SA> {
        let mut result = ManagedBufferBuilder::default();
        result.append_managed_buffer(&self.basic.function_call.function_name);
        for arg in self.basic.function_call.arg_buffer.raw_arg_iter() {
            result.append_bytes(b"@");
            SCLowerHex::fmt(&*arg, &mut result);
        }
        result.into_managed_buffer()
    }
}

impl<SA, OriginalResult> ContractCallWithRewa<SA, OriginalResult>
where
    SA: CallTypeApi + StorageWriteApi + 'static,
{
    pub(super) fn build_async_call(self) -> AsyncCall<SA> {
        Tx::new_tx_from_sc()
            .to(self.basic.to)
            .rewa(self.rewa_payment)
            .raw_data(self.basic.function_call)
            .callback(None)
    }
}

impl<SA, OriginalResult> ContractCallWithRewa<SA, OriginalResult>
where
    SA: CallTypeApi + 'static,
{
    pub(super) fn build_async_call_promise(self) -> AsyncCallPromises<SA> {
        AsyncCallPromises {
            to: self.basic.to,
            rewa_payment: self.rewa_payment,
            function_call: self.basic.function_call,
            explicit_gas_limit: self.basic.explicit_gas_limit,
            extra_gas_for_callback: 0,
            callback_call: None,
        }
    }

    /// Executes immediately, synchronously, and returns contract call result.
    /// Only works if the target contract is in the same shard.
    pub(super) fn execute_on_dest_context<RequestedResult>(self) -> RequestedResult
    where
        RequestedResult: TopDecodeMulti,
    {
        let raw_result = SendRawWrapper::<SA>::new().execute_on_dest_context_raw(
            self.resolve_gas_limit(),
            &self.basic.to,
            &self.rewa_payment,
            &self.basic.function_call.function_name,
            &self.basic.function_call.arg_buffer,
        );

        SendRawWrapper::<SA>::new().clean_return_data();

        decode_result(raw_result)
    }

    pub(super) fn execute_on_dest_context_readonly<RequestedResult>(self) -> RequestedResult
    where
        RequestedResult: TopDecodeMulti,
    {
        let raw_result = SendRawWrapper::<SA>::new().execute_on_dest_context_readonly_raw(
            self.resolve_gas_limit(),
            &self.basic.to,
            &self.basic.function_call.function_name,
            &self.basic.function_call.arg_buffer,
        );

        SendRawWrapper::<SA>::new().clean_return_data();

        decode_result(raw_result)
    }

    pub(super) fn execute_on_same_context<RequestedResult>(self) -> RequestedResult
    where
        RequestedResult: TopDecodeMulti,
    {
        let raw_result = SendRawWrapper::<SA>::new().execute_on_same_context_raw(
            self.resolve_gas_limit(),
            &self.basic.to,
            &self.rewa_payment,
            &self.basic.function_call.function_name,
            &self.basic.function_call.arg_buffer,
        );

        SendRawWrapper::<SA>::new().clean_return_data();

        decode_result(raw_result)
    }
}

impl<SA, OriginalResult> ContractCallNoPayment<SA, OriginalResult>
where
    SA: CallTypeApi + 'static,
{
    pub(super) fn resolve_gas_limit_with_leftover(&self) -> u64 {
        if self.explicit_gas_limit == UNSPECIFIED_GAS_LIMIT {
            let mut gas_left = SA::blockchain_api_impl().get_gas_left();
            if gas_left > TRANSFER_EXECUTE_DEFAULT_LEFTOVER {
                gas_left -= TRANSFER_EXECUTE_DEFAULT_LEFTOVER;
            }
            gas_left
        } else {
            self.explicit_gas_limit
        }
    }

    pub(super) fn transfer_execute_rewa(self, rewa_payment: BigUint<SA>) {
        let gas_limit = self.resolve_gas_limit_with_leftover();

        let _ = SendRawWrapper::<SA>::new().direct_rewa_execute(
            &self.to,
            &rewa_payment,
            gas_limit,
            &self.function_call.function_name,
            &self.function_call.arg_buffer,
        );
    }

    pub(super) fn transfer_execute_single_dcdt(self, payment: DcdtTokenPayment<SA>) {
        let gas_limit = self.resolve_gas_limit_with_leftover();

        if payment.token_nonce == 0 {
            // fungible DCDT
            let _ = SendRawWrapper::<SA>::new().transfer_dcdt_execute(
                &self.to,
                &payment.token_identifier,
                &payment.amount,
                gas_limit,
                &self.function_call.function_name,
                &self.function_call.arg_buffer,
            );
        } else {
            // non-fungible/semi-fungible DCDT
            let _ = SendRawWrapper::<SA>::new().transfer_dcdt_nft_execute(
                &self.to,
                &payment.token_identifier,
                payment.token_nonce,
                &payment.amount,
                gas_limit,
                &self.function_call.function_name,
                &self.function_call.arg_buffer,
            );
        }
    }

    pub(super) fn transfer_execute_multi_dcdt(
        self,
        payments: ManagedVec<SA, DcdtTokenPayment<SA>>,
    ) {
        let gas_limit = self.resolve_gas_limit_with_leftover();
        let _ = SendRawWrapper::<SA>::new().multi_dcdt_transfer_execute(
            &self.to,
            &payments,
            gas_limit,
            &self.function_call.function_name,
            &self.function_call.arg_buffer,
        );
    }

    pub(super) fn transfer_execute_dcdt(self, payments: ManagedVec<SA, DcdtTokenPayment<SA>>) {
        match payments.len() {
            0 => self.transfer_execute_rewa(BigUint::zero()),
            1 => self.transfer_execute_single_dcdt(payments.get(0).clone()),
            _ => self.transfer_execute_multi_dcdt(payments),
        }
    }
}
