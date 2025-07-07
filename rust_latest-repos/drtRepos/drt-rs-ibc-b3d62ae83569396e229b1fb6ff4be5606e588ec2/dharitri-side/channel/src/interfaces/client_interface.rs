pub mod generic_client_proxy {
    use client_common::{GetLatestInfoResultType, VerifyMembershipArgs, VerifyNonMembershipArgs};
    use common_types::{channel_types::height, ClientId, UnixTimestamp};

    dharitri_sc::imports!();

    #[dharitri_sc::proxy]
    pub trait GenericClientProxy {
        #[view(getTimestampAtHeight)]
        fn get_timestamp_at_height(
            &self,
            client_id: &ClientId<Self::Api>,
            height: &height::Data,
        ) -> UnixTimestamp;

        #[view(getLatestInfo)]
        fn get_latest_info(&self, client_id: ClientId<Self::Api>) -> GetLatestInfoResultType;

        #[view(verifyMembership)]
        fn verify_membership(&self, args: VerifyMembershipArgs<Self::Api>) -> bool;

        #[view(verifyNonMembership)]
        fn verify_non_membership(&self, args: VerifyNonMembershipArgs<Self::Api>) -> bool;
    }
}
