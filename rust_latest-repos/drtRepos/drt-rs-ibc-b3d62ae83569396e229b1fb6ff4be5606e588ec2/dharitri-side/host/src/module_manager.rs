use common_types::{ChannelId, PortId};

dharitri_sc::imports!();

#[dharitri_sc::module]
pub trait ModuleManagerModule: crate::storage::StorageModule {
    fn claim_port_capability(&self, port_id: &PortId<Self::Api>, address: &ManagedAddress) {
        let mapper = self.port_capabilities(port_id);
        require!(mapper.is_empty(), "Port already claimed");

        mapper.set(address);
    }

    fn claim_channel_capability(
        &self,
        port_id: &PortId<Self::Api>,
        channel_id: &ChannelId<Self::Api>,
        address: &ManagedAddress,
    ) {
        let mapper = self.channel_capabilities(port_id, channel_id);
        require!(mapper.is_empty(), "Channel already claimed");

        mapper.set(address);
    }

    fn authenticate_channel_capability(
        &self,
        port_id: &PortId<Self::Api>,
        channel_id: &ChannelId<Self::Api>,
        user: &ManagedAddress,
    ) {
        let mapper = self.channel_capabilities(port_id, channel_id);
        require!(!mapper.is_empty(), "Channel not claimed");

        let stored_addr = mapper.get();
        require!(user == &stored_addr, "Not allowed to use this port");
    }

    fn lookup_module_by_port(&self, port_id: &PortId<Self::Api>) -> ManagedAddress {
        let mapper = self.port_capabilities(port_id);
        require!(mapper.is_empty(), "Port not found");

        mapper.get()
    }

    fn lookup_module_by_channel(
        &self,
        port_id: &PortId<Self::Api>,
        channel_id: &ChannelId<Self::Api>,
    ) -> ManagedAddress {
        let mapper = self.channel_capabilities(port_id, channel_id);
        require!(!mapper.is_empty(), "Channel not found");

        mapper.get()
    }

    // TODO: Do we even need something like this?

    /*
     function canTransitionToFlushComplete(
        Channel.Order ordering,
        string calldata portId,
        string calldata channelId,
        uint64 upgradeSequence
    ) internal view virtual returns (bool) {
        if (ordering == Channel.Order.ORDER_ORDERED) {
            ChannelStorage storage channelStorage = getChannelStorage()[portId][channelId];
            if (channelStorage.nextSequenceSend == channelStorage.nextSequenceAck) {
                return true;
            }
        }
        return lookupUpgradableModuleByPortUnchecked(portId).canTransitionToFlushComplete(
            portId, channelId, upgradeSequence, _msgSender()
        );
    }
     */
}
