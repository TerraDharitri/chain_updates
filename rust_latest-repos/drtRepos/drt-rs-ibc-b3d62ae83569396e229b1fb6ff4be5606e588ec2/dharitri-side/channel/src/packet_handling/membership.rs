use client_common::{VerifyMembershipArgs, VerifyNonMembershipArgs};
use common_types::{
    channel_types::{channel, channel_counterparty},
    connection_types::connection_end,
    ConnectionHops,
};

use crate::{
    channel_libs::packet_types::{MsgTimeoutOnClose, TimeoutArgs},
    interfaces::client_interface,
};

use super::errors::UNKNOW_CHANNEL_ORDER_ERR_MSG;

dharitri_sc::imports!();

#[dharitri_sc::module]
pub trait MembershipModule:
    host::module_manager::ModuleManagerModule
    + host::storage::StorageModule
    + host::commitment::CommitmentModule
    + common_modules::utils::UtilsModule
{
    fn check_channel_membership(
        &self,
        ordering: channel::Order,
        client_impl: ManagedAddress,
        connection_info: &connection_end::Data<Self::Api>,
        timeout_args: &dyn TimeoutArgs<Self::Api>,
    ) {
        match ordering {
            channel::Order::Ordered => {
                self.check_channel_ordered_membership(client_impl, connection_info, timeout_args)
            }
            channel::Order::Unordered => {
                self.check_channel_unordered_membership(client_impl, connection_info, timeout_args)
            }
            channel::Order::NoneUnspecified => sc_panic!(UNKNOW_CHANNEL_ORDER_ERR_MSG),
        };
    }

    fn check_channel_ordered_membership(
        &self,
        client_impl: ManagedAddress,
        connection_info: &connection_end::Data<Self::Api>,
        timeout_args: &dyn TimeoutArgs<Self::Api>,
    ) {
        let packet = timeout_args.get_packet();
        require!(
            packet.seq >= timeout_args.get_next_seq_recv(),
            "Packet may already be received"
        );

        let encoded_value = self.encode_to_buffer(&timeout_args.get_next_seq_recv());
        let membership_args = VerifyMembershipArgs {
            client_id: connection_info.client_id.clone(),
            height: timeout_args.get_proof_height(),
            delay_time_period: connection_info.delay_period,
            delay_block_period: self.calculate_block_delay(connection_info.delay_period),
            proof: timeout_args.get_proof().clone(),
            prefix: connection_info.counterparty.prefix.key_prefix.clone(),
            path: self.get_next_seq_recv_commitment_path(&packet.dest_port, &packet.dest_channel),
            value: encoded_value,
        };
        let membership_result: bool = self
            .generic_client_proxy_impl_membership(client_impl)
            .verify_membership(membership_args)
            .execute_on_dest_context();
        require!(membership_result, "Failed to verify next seq receive");

        self.channel_info(&packet.src_port, &packet.src_channel)
            .update(|channel_info| channel_info.channel.state = channel::State::Closed);
    }

    fn check_channel_unordered_membership(
        &self,
        client_impl: ManagedAddress,
        connection_info: &connection_end::Data<Self::Api>,
        timeout_args: &dyn TimeoutArgs<Self::Api>,
    ) {
        let packet = timeout_args.get_packet();
        let path = self.get_packet_receipt_commitment_path(
            &packet.dest_port,
            &packet.dest_channel,
            packet.seq,
        );
        let non_membership_args = VerifyNonMembershipArgs {
            client_id: connection_info.client_id.clone(),
            height: timeout_args.get_proof_height(),
            delay_time_period: connection_info.delay_period,
            delay_block_period: self.calculate_block_delay(connection_info.delay_period),
            proof: timeout_args.get_proof().clone(),
            prefix: connection_info.counterparty.prefix.key_prefix.clone(),
            path,
        };
        let non_membership_result: bool = self
            .generic_client_proxy_impl_membership(client_impl)
            .verify_non_membership(non_membership_args)
            .execute_on_dest_context();
        require!(
            non_membership_result,
            "Failed to verify packet receipt absence"
        );
    }

    fn check_expected_channel_membership(
        &self,
        client_impl: ManagedAddress,
        channel: &channel::Data<Self::Api>,
        connection_info: &connection_end::Data<Self::Api>,
        args: &MsgTimeoutOnClose<Self::Api>,
    ) {
        let expected_channel = channel::Data {
            state: channel::State::Closed,
            ordering: channel.ordering,
            counterparty: channel_counterparty::Data {
                port_id: args.packet.src_port.clone(),
                channel_id: args.packet.src_channel.clone(),
            },
            connection_hops: ConnectionHops::from_single_item(
                connection_info.counterparty.connection_id.clone(),
            ),
            version: channel.version.clone(),
            upgrade_sequence: args.counterparty_upgrade_seq,
        };

        let encoded_value = self.encode_to_buffer(&expected_channel);
        let membership_args = VerifyMembershipArgs {
            client_id: connection_info.client_id.clone(),
            height: args.proof_height,
            delay_time_period: connection_info.delay_period,
            delay_block_period: self.calculate_block_delay(connection_info.delay_period),
            proof: args.proof_close.clone(),
            prefix: connection_info.counterparty.prefix.key_prefix.clone(),
            path: self.get_next_seq_recv_commitment_path(
                &args.packet.dest_port,
                &args.packet.dest_channel,
            ),
            value: encoded_value,
        };
        let membership_result: bool = self
            .generic_client_proxy_impl_membership(client_impl)
            .verify_membership(membership_args)
            .execute_on_dest_context();
        require!(membership_result, "Failed to verify channel state");
    }

    #[proxy]
    fn generic_client_proxy_impl_membership(
        &self,
        sc_address: ManagedAddress,
    ) -> client_interface::generic_client_proxy::GenericClientProxy<Self::Api>;
}
