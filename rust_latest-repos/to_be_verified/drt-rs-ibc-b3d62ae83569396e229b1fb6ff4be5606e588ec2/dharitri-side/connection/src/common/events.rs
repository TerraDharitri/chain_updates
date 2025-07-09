use common_types::ConnectionId;

dharitri_sc::imports!();

#[dharitri_sc::module]
pub trait EventsModule {
    #[event("generatedConnectionIdEvent")]
    fn generated_connection_id_event(&self, connection_id: &ConnectionId<Self::Api>);
}
