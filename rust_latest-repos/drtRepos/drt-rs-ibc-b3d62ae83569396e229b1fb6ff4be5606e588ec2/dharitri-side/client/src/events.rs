use common_types::ClientId;

dharitri_sc::imports!();

#[dharitri_sc::module]
pub trait EventsModule {
    #[event("generatedClientIdEvent")]
    fn generated_client_id_event(&self, client_id: &ClientId<Self::Api>);
}
