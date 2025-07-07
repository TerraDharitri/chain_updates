use super::{token_properties::*, TokenPropertiesResult};

use crate::{
    api::CallTypeApi,
    types::{
        BigUint, RewaPayment, DcdtLocalRole, DcdtTokenType, FunctionCall, ManagedAddress,
        ManagedBuffer, NotPayable, OriginalResultMarker, ProxyArg, TokenIdentifier, Tx, TxEnv,
        TxFrom, TxGas, TxProxyTrait, TxTo, TxTypedCall,
    },
};

const ISSUE_FUNGIBLE_ENDPOINT_NAME: &str = "issue";
const ISSUE_NON_FUNGIBLE_ENDPOINT_NAME: &str = "issueNonFungible";
const ISSUE_SEMI_FUNGIBLE_ENDPOINT_NAME: &str = "issueSemiFungible";
const REGISTER_META_DCDT_ENDPOINT_NAME: &str = "registerMetaDCDT";
const ISSUE_AND_SET_ALL_ROLES_ENDPOINT_NAME: &str = "registerAndSetAllRoles";
const REGISTER_DYNAMIC_DCDT_ENDPOINT_NAME: &str = "registerDynamic";
const REGISTER_AND_SET_ALL_ROLES_DYNAMIC_DCDT_ENDPOINT_NAME: &str = "registerAndSetAllRolesDynamic";

/// The specific `Tx` type produces by the issue operations of the DCDTSystemSCProxy.
pub type IssueCall<Env, From, To, Gas> = Tx<
    Env,
    From,
    To,
    RewaPayment<<Env as TxEnv>::Api>,
    Gas,
    FunctionCall<<Env as TxEnv>::Api>,
    OriginalResultMarker<TokenIdentifier<<Env as TxEnv>::Api>>,
>;

/// Proxy for the DCDT system smart contract.
pub struct DCDTSystemSCProxy;

impl<Env, From, To, Gas> TxProxyTrait<Env, From, To, Gas> for DCDTSystemSCProxy
where
    Env: TxEnv,
    From: TxFrom<Env>,
    To: TxTo<Env>,
    Gas: TxGas<Env>,
{
    type TxProxyMethods = DCDTSystemSCProxyMethods<Env, From, To, Gas>;

    fn proxy_methods(self, tx: Tx<Env, From, To, (), Gas, (), ()>) -> Self::TxProxyMethods {
        DCDTSystemSCProxyMethods { wrapped_tx: tx }
    }
}

/// Method container of the DCDT system smart contract proxy.
pub struct DCDTSystemSCProxyMethods<Env, From, To, Gas>
where
    Env: TxEnv,
    From: TxFrom<Env>,
    To: TxTo<Env>,
    Gas: TxGas<Env>,
{
    wrapped_tx: Tx<Env, From, To, (), Gas, (), ()>,
}

impl<Env, From, To, Gas> DCDTSystemSCProxyMethods<Env, From, To, Gas>
where
    Env: TxEnv,
    From: TxFrom<Env>,
    To: TxTo<Env>,
    Gas: TxGas<Env>,
{
    /// Produces a contract call to the DCDT system SC,
    /// which causes it to issue a new fungible DCDT token.
    pub fn issue_fungible<
        Arg0: ProxyArg<ManagedBuffer<Env::Api>>,
        Arg1: ProxyArg<ManagedBuffer<Env::Api>>,
        Arg2: ProxyArg<BigUint<Env::Api>>,
    >(
        self,
        issue_cost: BigUint<Env::Api>,
        token_display_name: Arg0,
        token_ticker: Arg1,
        initial_supply: Arg2,
        properties: FungibleTokenProperties,
    ) -> IssueCall<Env, From, To, Gas> {
        self.issue(
            issue_cost,
            DcdtTokenType::Fungible,
            token_display_name,
            token_ticker,
            initial_supply,
            TokenProperties {
                num_decimals: properties.num_decimals,
                can_freeze: properties.can_freeze,
                can_wipe: properties.can_wipe,
                can_pause: properties.can_pause,
                can_transfer_create_role: false,
                can_mint: properties.can_mint,
                can_burn: properties.can_burn,
                can_change_owner: properties.can_change_owner,
                can_upgrade: properties.can_upgrade,
                can_add_special_roles: properties.can_add_special_roles,
            },
        )
    }

    /// Produces a contract call to the DCDT system SC,
    /// which causes it to issue a new non-fungible DCDT token.
    pub fn issue_non_fungible<
        Arg0: ProxyArg<ManagedBuffer<Env::Api>>,
        Arg1: ProxyArg<ManagedBuffer<Env::Api>>,
    >(
        self,
        issue_cost: BigUint<Env::Api>,
        token_display_name: Arg0,
        token_ticker: Arg1,
        properties: NonFungibleTokenProperties,
    ) -> IssueCall<Env, From, To, Gas> {
        let zero = &BigUint::zero();
        self.issue(
            issue_cost,
            DcdtTokenType::NonFungible,
            token_display_name,
            token_ticker,
            zero,
            TokenProperties {
                num_decimals: 0,
                can_freeze: properties.can_freeze,
                can_wipe: properties.can_wipe,
                can_pause: properties.can_pause,
                can_transfer_create_role: properties.can_transfer_create_role,
                can_mint: false,
                can_burn: false,
                can_change_owner: properties.can_change_owner,
                can_upgrade: properties.can_upgrade,
                can_add_special_roles: properties.can_add_special_roles,
            },
        )
    }

    /// Produces a contract call to the DCDT system SC,
    /// which causes it to issue a new semi-fungible DCDT token.
    pub fn issue_semi_fungible<
        Arg0: ProxyArg<ManagedBuffer<Env::Api>>,
        Arg1: ProxyArg<ManagedBuffer<Env::Api>>,
    >(
        self,
        issue_cost: BigUint<Env::Api>,
        token_display_name: Arg0,
        token_ticker: Arg1,
        properties: SemiFungibleTokenProperties,
    ) -> IssueCall<Env, From, To, Gas> {
        let zero = BigUint::zero();
        self.issue(
            issue_cost,
            DcdtTokenType::SemiFungible,
            token_display_name,
            token_ticker,
            &zero,
            TokenProperties {
                num_decimals: 0,
                can_freeze: properties.can_freeze,
                can_wipe: properties.can_wipe,
                can_pause: properties.can_pause,
                can_transfer_create_role: properties.can_transfer_create_role,
                can_mint: false,
                can_burn: false,
                can_change_owner: properties.can_change_owner,
                can_upgrade: properties.can_upgrade,
                can_add_special_roles: properties.can_add_special_roles,
            },
        )
    }

    /// Produces a contract call to the DCDT system SC,
    /// which causes it to register a new Meta DCDT token.
    pub fn register_meta_dcdt<
        Arg0: ProxyArg<ManagedBuffer<Env::Api>>,
        Arg1: ProxyArg<ManagedBuffer<Env::Api>>,
    >(
        self,
        issue_cost: BigUint<Env::Api>,
        token_display_name: Arg0,
        token_ticker: Arg1,
        properties: MetaTokenProperties,
    ) -> IssueCall<Env, From, To, Gas> {
        let zero = &BigUint::zero();
        self.issue(
            issue_cost,
            DcdtTokenType::Meta,
            token_display_name,
            token_ticker,
            zero,
            TokenProperties {
                num_decimals: properties.num_decimals,
                can_freeze: properties.can_freeze,
                can_wipe: properties.can_wipe,
                can_pause: properties.can_pause,
                can_transfer_create_role: properties.can_transfer_create_role,
                can_mint: false,
                can_burn: false,
                can_change_owner: properties.can_change_owner,
                can_upgrade: properties.can_upgrade,
                can_add_special_roles: properties.can_add_special_roles,
            },
        )
    }

    pub fn issue_and_set_all_roles<
        Arg0: ProxyArg<ManagedBuffer<Env::Api>>,
        Arg1: ProxyArg<ManagedBuffer<Env::Api>>,
    >(
        self,
        issue_cost: BigUint<Env::Api>,
        token_display_name: Arg0,
        token_ticker: Arg1,
        token_type: DcdtTokenType,
        num_decimals: usize,
    ) -> IssueCall<Env, From, To, Gas> {
        let token_type_name = match token_type {
            DcdtTokenType::Fungible => "FNG",
            DcdtTokenType::NonFungible | DcdtTokenType::DynamicNFT => "NFT",
            DcdtTokenType::SemiFungible | DcdtTokenType::DynamicSFT => "SFT",
            DcdtTokenType::Meta | DcdtTokenType::DynamicMeta => "META",
            DcdtTokenType::Invalid => "",
        };

        let endpoint = match token_type {
            DcdtTokenType::Fungible
            | DcdtTokenType::NonFungible
            | DcdtTokenType::SemiFungible
            | DcdtTokenType::Meta => ISSUE_AND_SET_ALL_ROLES_ENDPOINT_NAME,
            DcdtTokenType::DynamicNFT | DcdtTokenType::DynamicSFT | DcdtTokenType::DynamicMeta => {
                REGISTER_AND_SET_ALL_ROLES_DYNAMIC_DCDT_ENDPOINT_NAME
            },

            DcdtTokenType::Invalid => "",
        };

        let mut tx = self
            .wrapped_tx
            .raw_call(endpoint)
            .rewa(issue_cost)
            .argument(&token_display_name)
            .argument(&token_ticker)
            .argument(&token_type_name);

        if token_type != DcdtTokenType::DynamicNFT && token_type != DcdtTokenType::DynamicSFT {
            tx = tx.argument(&num_decimals);
        }

        tx.original_result()
    }

    /// Issues dynamic DCDT tokens
    pub fn issue_dynamic<
        Arg0: ProxyArg<ManagedBuffer<Env::Api>>,
        Arg1: ProxyArg<ManagedBuffer<Env::Api>>,
    >(
        self,
        issue_cost: BigUint<Env::Api>,
        token_display_name: Arg0,
        token_ticker: Arg1,
        token_type: DcdtTokenType,
        num_decimals: usize,
    ) -> IssueCall<Env, From, To, Gas> {
        let endpoint_name = match token_type {
            DcdtTokenType::DynamicNFT | DcdtTokenType::DynamicSFT | DcdtTokenType::DynamicMeta => {
                REGISTER_DYNAMIC_DCDT_ENDPOINT_NAME
            },
            _ => "",
        };

        let token_type_name = match token_type {
            DcdtTokenType::DynamicNFT => "NFT",
            DcdtTokenType::DynamicSFT => "SFT",
            DcdtTokenType::DynamicMeta => "META",
            _ => "",
        };

        let mut tx = self
            .wrapped_tx
            .raw_call(endpoint_name)
            .rewa(issue_cost)
            .argument(&token_display_name)
            .argument(&token_ticker)
            .argument(&token_type_name);

        if token_type != DcdtTokenType::DynamicNFT && token_type != DcdtTokenType::DynamicSFT {
            tx = tx.argument(&num_decimals);
        }

        tx.original_result()
    }

    /// Deduplicates code from all the possible issue functions
    fn issue<
        Arg0: ProxyArg<ManagedBuffer<Env::Api>>,
        Arg1: ProxyArg<ManagedBuffer<Env::Api>>,
        Arg2: ProxyArg<BigUint<Env::Api>>,
    >(
        self,
        issue_cost: BigUint<Env::Api>,
        token_type: DcdtTokenType,
        token_display_name: Arg0,
        token_ticker: Arg1,
        initial_supply: Arg2,
        properties: TokenProperties,
    ) -> IssueCall<Env, From, To, Gas> {
        let endpoint_name = match token_type {
            DcdtTokenType::Fungible => ISSUE_FUNGIBLE_ENDPOINT_NAME,
            DcdtTokenType::NonFungible => ISSUE_NON_FUNGIBLE_ENDPOINT_NAME,
            DcdtTokenType::SemiFungible => ISSUE_SEMI_FUNGIBLE_ENDPOINT_NAME,
            DcdtTokenType::Meta => REGISTER_META_DCDT_ENDPOINT_NAME,
            _ => "",
        };

        let mut tx = self
            .wrapped_tx
            .raw_call(endpoint_name)
            .rewa(issue_cost)
            .argument(&token_display_name)
            .argument(&token_ticker);

        if token_type == DcdtTokenType::Fungible {
            tx = tx.argument(&initial_supply);
            tx = tx.argument(&properties.num_decimals);
        } else if token_type == DcdtTokenType::Meta {
            tx = tx.argument(&properties.num_decimals);
        }

        let mut token_prop_args = TokenPropertyArguments {
            can_freeze: Some(properties.can_freeze),
            can_wipe: Some(properties.can_wipe),
            can_pause: Some(properties.can_pause),
            can_change_owner: Some(properties.can_change_owner),
            can_upgrade: Some(properties.can_upgrade),
            can_add_special_roles: Some(properties.can_add_special_roles),
            ..TokenPropertyArguments::default()
        };

        if token_type == DcdtTokenType::Fungible {
            token_prop_args.can_mint = Some(properties.can_mint);
            token_prop_args.can_burn = Some(properties.can_burn);
        } else {
            token_prop_args.can_transfer_create_role = Some(properties.can_transfer_create_role);
        }

        append_token_property_arguments(&mut tx.data, &token_prop_args);

        tx.original_result()
    }

    /// Produces a contract call to the DCDT system SC,
    /// which causes it to mint more fungible DCDT tokens.
    /// It will fail if the SC is not the owner of the token.
    pub fn mint<Arg0: ProxyArg<TokenIdentifier<Env::Api>>, Arg1: ProxyArg<BigUint<Env::Api>>>(
        self,
        token_identifier: Arg0,
        amount: Arg1,
    ) -> TxTypedCall<Env, From, To, NotPayable, Gas, ()> {
        self.wrapped_tx
            .payment(NotPayable)
            .raw_call("mint")
            .argument(&token_identifier)
            .argument(&amount)
            .original_result()
    }

    /// Produces a contract call to the DCDT system SC,
    /// which causes it to burn fungible DCDT tokens owned by the SC.
    pub fn burn<Arg0: ProxyArg<TokenIdentifier<Env::Api>>, Arg1: ProxyArg<BigUint<Env::Api>>>(
        self,
        token_identifier: Arg0,
        amount: Arg1,
    ) -> TxTypedCall<Env, From, To, NotPayable, Gas, ()> {
        self.wrapped_tx
            .payment(NotPayable)
            .raw_call("DCDTBurn")
            .argument(&token_identifier)
            .argument(&amount)
            .original_result()
    }

    /// The manager of an DCDT token may choose to suspend all transactions of the token,
    /// except minting, freezing/unfreezing and wiping.
    pub fn pause<Arg0: ProxyArg<TokenIdentifier<Env::Api>>>(
        self,
        token_identifier: Arg0,
    ) -> TxTypedCall<Env, From, To, NotPayable, Gas, ()> {
        self.wrapped_tx
            .payment(NotPayable)
            .raw_call("pause")
            .argument(&token_identifier)
            .original_result()
    }

    /// The reverse operation of `pause`.
    pub fn unpause<Arg0: ProxyArg<TokenIdentifier<Env::Api>>>(
        self,
        token_identifier: Arg0,
    ) -> TxTypedCall<Env, From, To, NotPayable, Gas, ()> {
        self.wrapped_tx
            .payment(NotPayable)
            .raw_call("unPause")
            .argument(&token_identifier)
            .original_result()
    }

    /// The manager of an DCDT token may freeze the tokens held by a specific account.
    /// As a consequence, no tokens may be transferred to or from the frozen account.
    /// Freezing and unfreezing the tokens of an account are operations designed to help token managers to comply with regulations.
    pub fn freeze<
        Arg0: ProxyArg<TokenIdentifier<Env::Api>>,
        Arg1: ProxyArg<ManagedAddress<Env::Api>>,
    >(
        self,
        token_identifier: Arg0,
        address: Arg1,
    ) -> TxTypedCall<Env, From, To, NotPayable, Gas, ()> {
        self.wrapped_tx
            .payment(NotPayable)
            .raw_call("freeze")
            .argument(&token_identifier)
            .argument(&address)
            .original_result()
    }

    /// The reverse operation of `freeze`, unfreezing, will allow further transfers to and from the account.
    pub fn unfreeze<
        Arg0: ProxyArg<TokenIdentifier<Env::Api>>,
        Arg1: ProxyArg<ManagedAddress<Env::Api>>,
    >(
        self,
        token_identifier: Arg0,
        address: Arg1,
    ) -> TxTypedCall<Env, From, To, NotPayable, Gas, ()> {
        self.wrapped_tx
            .payment(NotPayable)
            .raw_call("unFreeze")
            .argument(&token_identifier)
            .argument(&address)
            .original_result()
    }

    /// The manager of an DCDT token may wipe out all the tokens held by a frozen account.
    /// This operation is similar to burning the tokens, but the account must have been frozen beforehand,
    /// and it must be done by the token manager.
    /// Wiping the tokens of an account is an operation designed to help token managers to comply with regulations.
    pub fn wipe<
        Arg0: ProxyArg<TokenIdentifier<Env::Api>>,
        Arg1: ProxyArg<ManagedAddress<Env::Api>>,
    >(
        self,
        token_identifier: Arg0,
        address: Arg1,
    ) -> TxTypedCall<Env, From, To, NotPayable, Gas, ()> {
        self.wrapped_tx
            .payment(NotPayable)
            .raw_call("wipe")
            .argument(&token_identifier)
            .argument(&address)
            .original_result()
    }

    /// The manager of an DCDT token may freeze the NFT held by a specific Account.
    /// As a consequence, no NFT can be transferred to or from the frozen Account.
    /// Freezing and unfreezing a single NFT of an Account are operations designed to help token managers to comply with regulations.
    pub fn freeze_nft<
        Arg0: ProxyArg<TokenIdentifier<Env::Api>>,
        Arg1: ProxyArg<ManagedAddress<Env::Api>>,
    >(
        self,
        token_identifier: Arg0,
        nft_nonce: u64,
        address: Arg1,
    ) -> TxTypedCall<Env, From, To, NotPayable, Gas, ()> {
        self.wrapped_tx
            .payment(NotPayable)
            .raw_call("freezeSingleNFT")
            .argument(&token_identifier)
            .argument(&nft_nonce)
            .argument(&address)
            .original_result()
    }

    /// The reverse operation of `freeze`, unfreezing, will allow further transfers to and from the account.
    pub fn unfreeze_nft<
        Arg0: ProxyArg<TokenIdentifier<Env::Api>>,
        Arg1: ProxyArg<ManagedAddress<Env::Api>>,
    >(
        self,
        token_identifier: Arg0,
        nft_nonce: u64,
        address: Arg1,
    ) -> TxTypedCall<Env, From, To, NotPayable, Gas, ()> {
        self.wrapped_tx
            .payment(NotPayable)
            .raw_call("unFreezeSingleNFT")
            .argument(&token_identifier)
            .argument(&nft_nonce)
            .argument(&address)
            .original_result()
    }

    /// The manager of an DCDT token may wipe out a single NFT held by a frozen Account.
    /// This operation is similar to burning the quantity, but the Account must have been frozen beforehand,
    /// and it must be done by the token manager.
    /// Wiping the tokens of an Account is an operation designed to help token managers to comply with regulations.
    pub fn wipe_nft<
        Arg0: ProxyArg<TokenIdentifier<Env::Api>>,
        Arg1: ProxyArg<ManagedAddress<Env::Api>>,
    >(
        self,
        token_identifier: Arg0,
        nft_nonce: u64,
        address: Arg1,
    ) -> TxTypedCall<Env, From, To, NotPayable, Gas, ()> {
        self.wrapped_tx
            .payment(NotPayable)
            .raw_call("wipeSingleNFT")
            .argument(&token_identifier)
            .argument(&nft_nonce)
            .argument(&address)
            .original_result()
    }

    /// This function converts an SFT to a metaDCDT by adding decimals to its structure in the metachain DCDT System SC.
    /// This function as almost all in case of DCDT can be called only by the owner.
    pub fn change_sft_to_meta_dcdt<Arg0: ProxyArg<TokenIdentifier<Env::Api>>>(
        self,
        token_identifier: Arg0,
        num_decimals: usize,
    ) -> TxTypedCall<Env, From, To, NotPayable, Gas, ()> {
        self.wrapped_tx
            .payment(NotPayable)
            .raw_call("changeSFTToMetaDCDT")
            .argument(&token_identifier)
            .argument(&num_decimals)
            .original_result()
    }

    /// This function can be called only if canSetSpecialRoles was set to true.
    /// The metachain system SC will evaluate the arguments and call “DCDTSetRole@tokenId@listOfRoles” for the given address.
    /// This will be actually a cross shard call.
    /// This function as almost all in case of DCDT can be called only by the owner.
    pub fn set_special_roles<
        RoleIter: Iterator<Item = DcdtLocalRole>,
        Arg0: ProxyArg<ManagedAddress<Env::Api>>,
        Arg1: ProxyArg<TokenIdentifier<Env::Api>>,
    >(
        self,
        address: Arg0,
        token_identifier: Arg1,
        roles_iter: RoleIter,
    ) -> TxTypedCall<Env, From, To, NotPayable, Gas, ()> {
        let mut tx = self
            .wrapped_tx
            .payment(NotPayable)
            .raw_call("setSpecialRole")
            .argument(&token_identifier)
            .argument(&address);
        for role in roles_iter {
            if role != DcdtLocalRole::None {
                tx = tx.argument(&role.as_role_name());
            }
        }

        tx.original_result()
    }

    /// This function can be called to retrieve the special roles of a specific token.
    pub fn get_special_roles<Arg0: ProxyArg<TokenIdentifier<Env::Api>>>(
        self,
        token_identifier: Arg0,
    ) -> TxTypedCall<Env, From, To, NotPayable, Gas, ()> {
        let tx = self
            .wrapped_tx
            .payment(NotPayable)
            .raw_call("getSpecialRoles")
            .argument(&token_identifier);

        tx.original_result()
    }

    /// This function can be called only if canSetSpecialRoles was set to true.
    /// The metachain system SC will evaluate the arguments and call “DCDTUnsetRole@tokenId@listOfRoles” for the given address.
    /// This will be actually a cross shard call.
    /// This function as almost all in case of DCDT can be called only by the owner.
    pub fn unset_special_roles<
        RoleIter: Iterator<Item = DcdtLocalRole>,
        Arg0: ProxyArg<ManagedAddress<Env::Api>>,
        Arg1: ProxyArg<TokenIdentifier<Env::Api>>,
    >(
        self,
        address: Arg0,
        token_identifier: Arg1,
        roles_iter: RoleIter,
    ) -> TxTypedCall<Env, From, To, NotPayable, Gas, ()> {
        let mut tx = self
            .wrapped_tx
            .payment(NotPayable)
            .raw_call("unSetSpecialRole")
            .argument(&token_identifier)
            .argument(&address);
        for role in roles_iter {
            if role != DcdtLocalRole::None {
                tx = tx.argument(&role.as_role_name());
            }
        }

        tx.original_result()
    }

    pub fn transfer_ownership<
        Arg0: ProxyArg<TokenIdentifier<Env::Api>>,
        Arg1: ProxyArg<ManagedAddress<Env::Api>>,
    >(
        self,
        token_identifier: Arg0,
        new_owner: Arg1,
    ) -> TxTypedCall<Env, From, To, NotPayable, Gas, ()> {
        self.wrapped_tx
            .payment(NotPayable)
            .raw_call("transferOwnership")
            .argument(&token_identifier)
            .argument(&new_owner)
            .original_result()
    }

    pub fn transfer_nft_create_role<
        Arg0: ProxyArg<TokenIdentifier<Env::Api>>,
        Arg1: ProxyArg<ManagedAddress<Env::Api>>,
    >(
        self,
        token_identifier: Arg0,
        old_creator: Arg1,
        new_creator: Arg1,
    ) -> TxTypedCall<Env, From, To, NotPayable, Gas, ()> {
        self.wrapped_tx
            .payment(NotPayable)
            .raw_call("transferNFTCreateRole")
            .argument(&token_identifier)
            .argument(&old_creator)
            .argument(&new_creator)
            .original_result()
    }

    pub fn control_changes<Arg0: ProxyArg<TokenIdentifier<Env::Api>>>(
        self,
        token_identifier: Arg0,
        property_arguments: &TokenPropertyArguments,
    ) -> TxTypedCall<Env, From, To, NotPayable, Gas, ()> {
        let mut tx = self
            .wrapped_tx
            .payment(NotPayable)
            .raw_call("controlChanges")
            .argument(&token_identifier);
        append_token_property_arguments(&mut tx.data, property_arguments);
        tx.original_result()
    }

    /// Changes token to dynamic.
    /// Does not work for: FungibleDCDT, NonFungibleDCDT, NonFungibleDCDTv2.
    pub fn change_to_dynamic<Arg0: ProxyArg<TokenIdentifier<Env::Api>>>(
        self,
        token_id: Arg0,
    ) -> TxTypedCall<Env, From, To, NotPayable, Gas, ()> {
        self.wrapped_tx
            .payment(NotPayable)
            .raw_call("changeToDynamic")
            .argument(&token_id)
            .original_result()
    }

    /// Updates a specific token to the newest version.
    pub fn update_token<Arg0: ProxyArg<TokenIdentifier<Env::Api>>>(
        self,
        token_id: Arg0,
    ) -> TxTypedCall<Env, From, To, NotPayable, Gas, ()> {
        self.wrapped_tx
            .payment(NotPayable)
            .raw_call("updateTokenID")
            .argument(&token_id)
            .original_result()
    }

    /// Fetches token properties for a specific token.
    pub fn get_token_properties<Arg0: ProxyArg<TokenIdentifier<Env::Api>>>(
        self,
        token_id: Arg0,
    ) -> TxTypedCall<Env, From, To, NotPayable, Gas, TokenPropertiesResult> {
        self.wrapped_tx
            .payment(NotPayable)
            .raw_call("getTokenProperties")
            .argument(&token_id)
            .original_result()
    }
}

const TRUE_STR: &str = "true";
const FALSE_STR: &str = "false";

fn bool_name_bytes(b: bool) -> &'static str {
    if b {
        TRUE_STR
    } else {
        FALSE_STR
    }
}

fn set_token_property<Api>(contract_call: &mut FunctionCall<Api>, name: &str, value: bool)
where
    Api: CallTypeApi,
{
    contract_call.arg_buffer.push_multi_arg(&name);
    contract_call
        .arg_buffer
        .push_multi_arg(&bool_name_bytes(value));
}

fn append_token_property_arguments<Api>(
    contract_call: &mut FunctionCall<Api>,
    token_prop_args: &TokenPropertyArguments,
) where
    Api: CallTypeApi,
{
    if let Some(can_freeze) = token_prop_args.can_freeze {
        set_token_property(contract_call, "canFreeze", can_freeze);
    }

    if let Some(can_wipe) = token_prop_args.can_wipe {
        set_token_property(contract_call, "canWipe", can_wipe);
    }

    if let Some(can_pause) = token_prop_args.can_pause {
        set_token_property(contract_call, "canPause", can_pause);
    }

    if let Some(can_transfer_create_role) = token_prop_args.can_transfer_create_role {
        set_token_property(
            contract_call,
            "canTransferNFTCreateRole",
            can_transfer_create_role,
        );
    }

    if let Some(can_mint) = token_prop_args.can_mint {
        set_token_property(contract_call, "canMint", can_mint);
    }

    if let Some(can_burn) = token_prop_args.can_burn {
        set_token_property(contract_call, "canBurn", can_burn);
    }

    if let Some(can_change_owner) = token_prop_args.can_change_owner {
        set_token_property(contract_call, "canChangeOwner", can_change_owner);
    }

    if let Some(can_upgrade) = token_prop_args.can_upgrade {
        set_token_property(contract_call, "canUpgrade", can_upgrade);
    }

    if let Some(can_add_special_roles) = token_prop_args.can_add_special_roles {
        set_token_property(contract_call, "canAddSpecialRoles", can_add_special_roles);
    }
}
