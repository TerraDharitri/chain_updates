use pausable::State;

dharitri_sc::imports!();

#[dharitri_sc::module]
pub trait CommonMethodsModule {
    #[inline]
    fn is_state_active(&self, state: State) -> bool {
        state == State::Active || state == State::PartialActive
    }

    #[inline]
    fn can_swap(&self, state: State) -> bool {
        state == State::Active
    }
}
