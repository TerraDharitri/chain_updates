dharitri_sc::imports!();

/// Example of a module that contains the constructor.
#[dharitri_sc::module]
pub trait InternalModuleInit {
    /// The constructor can reside in a module.
    /// The method can have any name.
    #[init]
    fn constructor_in_a_module(&self) {}
}
