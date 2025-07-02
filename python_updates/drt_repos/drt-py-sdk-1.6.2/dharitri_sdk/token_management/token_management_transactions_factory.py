import logging
from enum import Enum
from typing import Optional

from dharitri_sdk.abi import Serializer
from dharitri_sdk.abi.biguint_value import BigUIntValue
from dharitri_sdk.abi.bytes_value import BytesValue
from dharitri_sdk.abi.string_value import StringValue
from dharitri_sdk.builders.transaction_builder import TransactionBuilder
from dharitri_sdk.core import Address, Transaction
from dharitri_sdk.core.errors import BadUsageError
from dharitri_sdk.core.transactions_factory_config import TransactionsFactoryConfig

logger = logging.getLogger(__name__)


class TokenType(Enum):
    NFT = "NFT"
    SFT = "SFT"
    META = "META"
    FNG = "FNG"


class TokenManagementTransactionsFactory:
    def __init__(self, config: TransactionsFactoryConfig):
        self._config = config
        self.serializer = Serializer()

    def create_transaction_for_issuing_fungible(
        self,
        sender: Address,
        token_name: str,
        token_ticker: str,
        initial_supply: int,
        num_decimals: int,
        can_freeze: bool,
        can_wipe: bool,
        can_pause: bool,
        can_change_owner: bool,
        can_upgrade: bool,
        can_add_special_roles: bool,
    ) -> Transaction:
        self._notify_about_unsetting_burn_role_globally()

        parts = [
            "issue",
        ]

        serialized_parts = self.serializer.serialize_to_parts(
            [
                StringValue(token_name),
                StringValue(token_ticker),
                BigUIntValue(initial_supply),
                BigUIntValue(num_decimals),
                *[StringValue("canFreeze"), self._bool_to_typed_string(can_freeze)],
                *[StringValue("canWipe"), self._bool_to_typed_string(can_wipe)],
                *[StringValue("canPause"), self._bool_to_typed_string(can_pause)],
                *[
                    StringValue("canChangeOwner"),
                    self._bool_to_typed_string(can_change_owner),
                ],
                *[StringValue("canUpgrade"), self._bool_to_typed_string(can_upgrade)],
                *[
                    StringValue("canAddSpecialRoles"),
                    self._bool_to_typed_string(can_add_special_roles),
                ],
            ]
        )

        parts.extend([part.hex() for part in serialized_parts])

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=self._config.issue_cost,
            gas_limit=self._config.gas_limit_issue,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def _notify_about_unsetting_burn_role_globally(self) -> None:
        logger.info(
            """
==========
IMPORTANT!
==========
You are about to issue (register) a new token. This will set the role "DCDTRoleBurnForAll" (globally).
Once the token is registered, you can unset this role by calling "unsetBurnRoleGlobally" (in a separate transaction)."""
        )

    def create_transaction_for_issuing_semi_fungible(
        self,
        sender: Address,
        token_name: str,
        token_ticker: str,
        can_freeze: bool,
        can_wipe: bool,
        can_pause: bool,
        can_transfer_nft_create_role: bool,
        can_change_owner: bool,
        can_upgrade: bool,
        can_add_special_roles: bool,
    ) -> Transaction:
        self._notify_about_unsetting_burn_role_globally()

        parts = ["issueSemiFungible"]

        serialized_parts = self.serializer.serialize_to_parts(
            [
                StringValue(token_name),
                StringValue(token_ticker),
                *[StringValue("canFreeze"), self._bool_to_typed_string(can_freeze)],
                *[StringValue("canWipe"), self._bool_to_typed_string(can_wipe)],
                *[StringValue("canPause"), self._bool_to_typed_string(can_pause)],
                *[
                    StringValue("canTransferNFTCreateRole"),
                    self._bool_to_typed_string(can_transfer_nft_create_role),
                ],
                *[
                    StringValue("canChangeOwner"),
                    self._bool_to_typed_string(can_change_owner),
                ],
                *[StringValue("canUpgrade"), self._bool_to_typed_string(can_upgrade)],
                *[
                    StringValue("canAddSpecialRoles"),
                    self._bool_to_typed_string(can_add_special_roles),
                ],
            ]
        )

        parts.extend([part.hex() for part in serialized_parts])

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=self._config.issue_cost,
            gas_limit=self._config.gas_limit_issue,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_issuing_non_fungible(
        self,
        sender: Address,
        token_name: str,
        token_ticker: str,
        can_freeze: bool,
        can_wipe: bool,
        can_pause: bool,
        can_transfer_nft_create_role: bool,
        can_change_owner: bool,
        can_upgrade: bool,
        can_add_special_roles: bool,
    ) -> Transaction:
        self._notify_about_unsetting_burn_role_globally()

        parts = ["issueNonFungible"]

        serialized_parts = self.serializer.serialize_to_parts(
            [
                StringValue(token_name),
                StringValue(token_ticker),
                *[StringValue("canFreeze"), self._bool_to_typed_string(can_freeze)],
                *[StringValue("canWipe"), self._bool_to_typed_string(can_wipe)],
                *[StringValue("canPause"), self._bool_to_typed_string(can_pause)],
                *[
                    StringValue("canTransferNFTCreateRole"),
                    self._bool_to_typed_string(can_transfer_nft_create_role),
                ],
                *[
                    StringValue("canChangeOwner"),
                    self._bool_to_typed_string(can_change_owner),
                ],
                *[StringValue("canUpgrade"), self._bool_to_typed_string(can_upgrade)],
                *[
                    StringValue("canAddSpecialRoles"),
                    self._bool_to_typed_string(can_add_special_roles),
                ],
            ]
        )

        parts.extend([part.hex() for part in serialized_parts])

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=self._config.issue_cost,
            gas_limit=self._config.gas_limit_issue,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_registering_meta_dcdt(
        self,
        sender: Address,
        token_name: str,
        token_ticker: str,
        num_decimals: int,
        can_freeze: bool,
        can_wipe: bool,
        can_pause: bool,
        can_transfer_nft_create_role: bool,
        can_change_owner: bool,
        can_upgrade: bool,
        can_add_special_roles: bool,
    ) -> Transaction:
        self._notify_about_unsetting_burn_role_globally()

        parts = ["registerMetaDCDT"]

        serialized_parts = self.serializer.serialize_to_parts(
            [
                StringValue(token_name),
                StringValue(token_ticker),
                BigUIntValue(num_decimals),
                *[StringValue("canFreeze"), self._bool_to_typed_string(can_freeze)],
                *[StringValue("canWipe"), self._bool_to_typed_string(can_wipe)],
                *[StringValue("canPause"), self._bool_to_typed_string(can_pause)],
                *[
                    StringValue("canTransferNFTCreateRole"),
                    self._bool_to_typed_string(can_transfer_nft_create_role),
                ],
                *[
                    StringValue("canChangeOwner"),
                    self._bool_to_typed_string(can_change_owner),
                ],
                *[StringValue("canUpgrade"), self._bool_to_typed_string(can_upgrade)],
                *[
                    StringValue("canAddSpecialRoles"),
                    self._bool_to_typed_string(can_add_special_roles),
                ],
            ]
        )

        parts.extend([part.hex() for part in serialized_parts])

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=self._config.issue_cost,
            gas_limit=self._config.gas_limit_issue,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_registering_and_setting_roles(
        self,
        sender: Address,
        token_name: str,
        token_ticker: str,
        token_type: TokenType,
        num_decimals: int,
    ) -> Transaction:
        self._notify_about_unsetting_burn_role_globally()

        parts = ["registerAndSetAllRoles"]

        serialized_parts = self.serializer.serialize_to_parts(
            [
                StringValue(token_name),
                StringValue(token_ticker),
                StringValue(token_type.value),
                BigUIntValue(num_decimals),
            ]
        )

        parts.extend([part.hex() for part in serialized_parts])

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=self._config.issue_cost,
            gas_limit=self._config.gas_limit_issue,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_setting_burn_role_globally(self, sender: Address, token_identifier: str) -> Transaction:
        parts = [
            "setBurnRoleGlobally",
            self.serializer.serialize([StringValue(token_identifier)]),
        ]

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=None,
            gas_limit=self._config.gas_limit_toggle_burn_role_globally,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_unsetting_burn_role_globally(
        self, sender: Address, token_identifier: str
    ) -> Transaction:
        parts = [
            "unsetBurnRoleGlobally",
            self.serializer.serialize([StringValue(token_identifier)]),
        ]

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=None,
            gas_limit=self._config.gas_limit_toggle_burn_role_globally,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_setting_special_role_on_fungible_token(
        self,
        sender: Address,
        user: Address,
        token_identifier: str,
        add_role_local_mint: bool = False,
        add_role_local_burn: bool = False,
        add_role_dcdt_transfer_role: bool = False,
    ) -> Transaction:
        parts = [
            "setSpecialRole",
            self.serializer.serialize([StringValue(token_identifier)]),
            user.to_hex(),
        ]

        serialized_parts = self.serializer.serialize_to_parts(
            [
                *([StringValue("DCDTRoleLocalMint")] if add_role_local_mint else []),
                *([StringValue("DCDTRoleLocalBurn")] if add_role_local_burn else []),
                *([StringValue("DCDTTransferRole")] if add_role_dcdt_transfer_role else []),
            ]
        )

        parts.extend([part.hex() for part in serialized_parts])

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=None,
            gas_limit=self._config.gas_limit_set_special_role,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_unsetting_special_role_on_fungible_token(
        self,
        sender: Address,
        user: Address,
        token_identifier: str,
        remove_role_local_mint: bool = False,
        remove_role_local_burn: bool = False,
        remove_role_dcdt_transfer_role: bool = False,
    ) -> Transaction:
        parts = [
            "unSetSpecialRole",
            self.serializer.serialize([StringValue(token_identifier)]),
            user.to_hex(),
        ]

        serialized_parts = self.serializer.serialize_to_parts(
            [
                *([StringValue("DCDTRoleLocalMint")] if remove_role_local_mint else []),
                *([StringValue("DCDTRoleLocalBurn")] if remove_role_local_burn else []),
                *([StringValue("DCDTTransferRole")] if remove_role_dcdt_transfer_role else []),
            ]
        )

        parts.extend([part.hex() for part in serialized_parts])

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=None,
            gas_limit=self._config.gas_limit_set_special_role,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_setting_special_role_on_semi_fungible_token(
        self,
        sender: Address,
        user: Address,
        token_identifier: str,
        add_role_nft_create: bool = False,
        add_role_nft_burn: bool = False,
        add_role_nft_add_quantity: bool = False,
        add_role_dcdt_transfer_role: bool = False,
        add_role_nft_update: bool = False,
        add_role_dcdt_modify_royalties: bool = False,
        add_role_dcdt_set_new_uri: bool = False,
        add_role_dcdt_modify_creator: bool = False,
        add_role_nft_recreate: bool = False,
    ) -> Transaction:
        parts = [
            "setSpecialRole",
            self.serializer.serialize([StringValue(token_identifier)]),
            user.to_hex(),
        ]

        serialized_parts = self.serializer.serialize_to_parts(
            [
                *([StringValue("DCDTRoleNFTCreate")] if add_role_nft_create else []),
                *([StringValue("DCDTRoleNFTBurn")] if add_role_nft_burn else []),
                *([StringValue("DCDTRoleNFTAddQuantity")] if add_role_nft_add_quantity else []),
                *([StringValue("DCDTTransferRole")] if add_role_dcdt_transfer_role else []),
                *([StringValue("DCDTRoleNFTUpdate")] if add_role_nft_update else []),
                *([StringValue("DCDTRoleModifyRoyalties")] if add_role_dcdt_modify_royalties else []),
                *([StringValue("DCDTRoleSetNewURI")] if add_role_dcdt_set_new_uri else []),
                *([StringValue("DCDTRoleModifyCreator")] if add_role_dcdt_modify_creator else []),
                *([StringValue("DCDTRoleNFTRecreate")] if add_role_nft_recreate else []),
            ]
        )

        parts.extend([part.hex() for part in serialized_parts])

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=None,
            gas_limit=self._config.gas_limit_set_special_role,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_unsetting_special_role_on_semi_fungible_token(
        self,
        sender: Address,
        user: Address,
        token_identifier: str,
        remove_role_nft_burn: bool = False,
        remove_role_nft_add_quantity: bool = False,
        remove_role_dcdt_transfer_role: bool = False,
        remove_role_nft_update: bool = False,
        remove_role_dcdt_modify_royalties: bool = False,
        remove_role_dcdt_set_new_uri: bool = False,
        remove_role_dcdt_modify_creator: bool = False,
        remove_role_nft_recreate: bool = False,
    ) -> Transaction:
        parts = [
            "unSetSpecialRole",
            self.serializer.serialize([StringValue(token_identifier)]),
            user.to_hex(),
        ]

        serialized_parts = self.serializer.serialize_to_parts(
            [
                *([StringValue("DCDTRoleNFTBurn")] if remove_role_nft_burn else []),
                *([StringValue("DCDTRoleNFTAddQuantity")] if remove_role_nft_add_quantity else []),
                *([StringValue("DCDTTransferRole")] if remove_role_dcdt_transfer_role else []),
                *([StringValue("DCDTRoleNFTUpdate")] if remove_role_nft_update else []),
                *([StringValue("DCDTRoleModifyRoyalties")] if remove_role_dcdt_modify_royalties else []),
                *([StringValue("DCDTRoleSetNewURI")] if remove_role_dcdt_set_new_uri else []),
                *([StringValue("DCDTRoleModifyCreator")] if remove_role_dcdt_modify_creator else []),
                *([StringValue("DCDTRoleNFTRecreate")] if remove_role_nft_recreate else []),
            ]
        )

        parts.extend([part.hex() for part in serialized_parts])

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=None,
            gas_limit=self._config.gas_limit_set_special_role,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_setting_special_role_on_meta_dcdt(
        self,
        sender: Address,
        user: Address,
        token_identifier: str,
        add_role_nft_create: bool = False,
        add_role_nft_burn: bool = False,
        add_role_nft_add_quantity: bool = False,
        add_role_dcdt_transfer_role: bool = False,
    ) -> Transaction:
        return self.create_transaction_for_setting_special_role_on_semi_fungible_token(
            sender,
            user,
            token_identifier,
            add_role_nft_create,
            add_role_nft_burn,
            add_role_nft_add_quantity,
            add_role_dcdt_transfer_role,
        )

    def create_transaction_for_unsetting_special_role_on_meta_dcdt(
        self,
        sender: Address,
        user: Address,
        token_identifier: str,
        remove_role_nft_burn: bool = False,
        remove_role_nft_add_quantity: bool = False,
        remove_role_dcdt_transfer_role: bool = False,
    ) -> Transaction:
        return self.create_transaction_for_unsetting_special_role_on_semi_fungible_token(
            sender,
            user,
            token_identifier,
            remove_role_nft_burn,
            remove_role_nft_add_quantity,
            remove_role_dcdt_transfer_role,
        )

    def create_transaction_for_setting_special_role_on_non_fungible_token(
        self,
        sender: Address,
        user: Address,
        token_identifier: str,
        add_role_nft_create: bool = False,
        add_role_nft_burn: bool = False,
        add_role_nft_update_attributes: bool = False,
        add_role_nft_add_uri: bool = False,
        add_role_dcdt_transfer_role: bool = False,
        add_role_nft_update: bool = False,
        add_role_dcdt_modify_royalties: bool = False,
        add_role_dcdt_set_new_uri: bool = False,
        add_role_dcdt_modify_creator: bool = False,
        add_role_nft_recreate: bool = False,
    ) -> Transaction:
        parts = [
            "setSpecialRole",
            self.serializer.serialize([StringValue(token_identifier)]),
            user.to_hex(),
        ]

        serialized_parts = self.serializer.serialize_to_parts(
            [
                *([StringValue("DCDTRoleNFTCreate")] if add_role_nft_create else []),
                *([StringValue("DCDTRoleNFTBurn")] if add_role_nft_burn else []),
                *([StringValue("DCDTRoleNFTUpdateAttributes")] if add_role_nft_update_attributes else []),
                *([StringValue("DCDTRoleNFTAddURI")] if add_role_nft_add_uri else []),
                *([StringValue("DCDTTransferRole")] if add_role_dcdt_transfer_role else []),
                *([StringValue("DCDTRoleNFTUpdate")] if add_role_nft_update else []),
                *([StringValue("DCDTRoleModifyRoyalties")] if add_role_dcdt_modify_royalties else []),
                *([StringValue("DCDTRoleSetNewURI")] if add_role_dcdt_set_new_uri else []),
                *([StringValue("DCDTRoleModifyCreator")] if add_role_dcdt_modify_creator else []),
                *([StringValue("DCDTRoleNFTRecreate")] if add_role_nft_recreate else []),
            ]
        )

        parts.extend([part.hex() for part in serialized_parts])

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=None,
            gas_limit=self._config.gas_limit_set_special_role,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_unsetting_special_role_on_non_fungible_token(
        self,
        sender: Address,
        user: Address,
        token_identifier: str,
        remove_role_nft_burn: bool = False,
        remove_role_nft_update_attributes: bool = False,
        remove_role_nft_remove_uri: bool = False,
        remove_role_dcdt_transfer_role: bool = False,
        remove_role_nft_update: bool = False,
        remove_role_dcdt_modify_royalties: bool = False,
        remove_role_dcdt_set_new_uri: bool = False,
        remove_role_dcdt_modify_creator: bool = False,
        remove_role_nft_recreate: bool = False,
    ) -> Transaction:
        parts = [
            "unSetSpecialRole",
            self.serializer.serialize([StringValue(token_identifier)]),
            user.to_hex(),
        ]

        serialized_parts = self.serializer.serialize_to_parts(
            [
                *([StringValue("DCDTRoleNFTBurn")] if remove_role_nft_burn else []),
                *([StringValue("DCDTRoleNFTUpdateAttributes")] if remove_role_nft_update_attributes else []),
                *([StringValue("DCDTRoleNFTAddURI")] if remove_role_nft_remove_uri else []),
                *([StringValue("DCDTTransferRole")] if remove_role_dcdt_transfer_role else []),
                *([StringValue("DCDTRoleNFTUpdate")] if remove_role_nft_update else []),
                *([StringValue("DCDTRoleModifyRoyalties")] if remove_role_dcdt_modify_royalties else []),
                *([StringValue("DCDTRoleSetNewURI")] if remove_role_dcdt_set_new_uri else []),
                *([StringValue("DCDTRoleModifyCreator")] if remove_role_dcdt_modify_creator else []),
                *([StringValue("DCDTRoleNFTRecreate")] if remove_role_nft_recreate else []),
            ]
        )

        parts.extend([part.hex() for part in serialized_parts])

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=None,
            gas_limit=self._config.gas_limit_set_special_role,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_creating_nft(
        self,
        sender: Address,
        token_identifier: str,
        initial_quantity: int,
        name: str,
        royalties: int,
        hash: str,
        attributes: bytes,
        uris: list[str],
    ) -> Transaction:
        if not uris:
            raise BadUsageError("No URIs provided")

        parts = ["DCDTNFTCreate"]

        serialized_parts = self.serializer.serialize_to_parts(
            [
                StringValue(token_identifier),
                BigUIntValue(initial_quantity),
                StringValue(name),
                BigUIntValue(royalties),
                StringValue(hash),
                BytesValue(attributes),
                *map(StringValue, uris),
            ]
        )

        parts.extend([part.hex() for part in serialized_parts])

        # Note that the following is an approximation (a reasonable one):
        nft_data = name + hash + attributes.hex() + "".join(uris)
        storage_gas_limit = len(nft_data) * self._config.gas_limit_store_per_byte

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=sender,
            amount=None,
            gas_limit=self._config.gas_limit_dcdt_nft_create + storage_gas_limit,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_pausing(self, sender: Address, token_identifier: str) -> Transaction:
        parts = ["pause", self.serializer.serialize([StringValue(token_identifier)])]

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=None,
            gas_limit=self._config.gas_limit_pausing,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_unpausing(self, sender: Address, token_identifier: str) -> Transaction:
        parts = ["unPause", self.serializer.serialize([StringValue(token_identifier)])]

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=None,
            gas_limit=self._config.gas_limit_pausing,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_freezing(self, sender: Address, user: Address, token_identifier: str) -> Transaction:
        """Can be used for FungibleDCDT"""
        parts = [
            "freeze",
            self.serializer.serialize([StringValue(token_identifier)]),
            user.to_hex(),
        ]

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=None,
            gas_limit=self._config.gas_limit_freezing,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_unfreezing(self, sender: Address, user: Address, token_identifier: str) -> Transaction:
        """Can be used for FungibleDCDT"""
        parts = [
            "unFreeze",
            self.serializer.serialize([StringValue(token_identifier)]),
            user.to_hex(),
        ]

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=None,
            gas_limit=self._config.gas_limit_freezing,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_wiping(self, sender: Address, user: Address, token_identifier: str) -> Transaction:
        parts = [
            "wipe",
            self.serializer.serialize([StringValue(token_identifier)]),
            user.to_hex(),
        ]

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=None,
            gas_limit=self._config.gas_limit_wiping,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_local_minting(
        self, sender: Address, token_identifier: str, supply_to_mint: int
    ) -> Transaction:
        parts = [
            "DCDTLocalMint",
            self.serializer.serialize([StringValue(token_identifier)]),
            self.serializer.serialize([BigUIntValue(supply_to_mint)]),
        ]

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=sender,
            amount=None,
            gas_limit=self._config.gas_limit_dcdt_local_mint,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_local_burning(
        self, sender: Address, token_identifier: str, supply_to_burn: int
    ) -> Transaction:
        parts = [
            "DCDTLocalBurn",
            self.serializer.serialize([StringValue(token_identifier)]),
            self.serializer.serialize([BigUIntValue(supply_to_burn)]),
        ]

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=sender,
            amount=None,
            gas_limit=self._config.gas_limit_dcdt_local_burn,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_updating_attributes(
        self,
        sender: Address,
        token_identifier: str,
        token_nonce: int,
        attributes: bytes,
    ) -> Transaction:
        parts = ["DCDTNFTUpdateAttributes"]

        serialized_parts = self.serializer.serialize_to_parts(
            [
                StringValue(token_identifier),
                BigUIntValue(token_nonce),
                BytesValue(attributes),
            ]
        )

        parts.extend([part.hex() for part in serialized_parts])

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=sender,
            amount=None,
            gas_limit=self._config.gas_limit_dcdt_nft_update_attributes,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_adding_quantity(
        self,
        sender: Address,
        token_identifier: str,
        token_nonce: int,
        quantity_to_add: int,
    ) -> Transaction:
        parts = ["DCDTNFTAddQuantity"]

        serialized_parts = self.serializer.serialize_to_parts(
            [
                StringValue(token_identifier),
                BigUIntValue(token_nonce),
                BigUIntValue(quantity_to_add),
            ]
        )

        parts.extend([part.hex() for part in serialized_parts])

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=sender,
            amount=None,
            gas_limit=self._config.gas_limit_dcdt_nft_add_quantity,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_burning_quantity(
        self,
        sender: Address,
        token_identifier: str,
        token_nonce: int,
        quantity_to_burn: int,
    ) -> Transaction:
        parts = ["DCDTNFTBurn"]

        serialized_parts = self.serializer.serialize_to_parts(
            [
                StringValue(token_identifier),
                BigUIntValue(token_nonce),
                BigUIntValue(quantity_to_burn),
            ]
        )

        parts.extend([part.hex() for part in serialized_parts])

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=sender,
            amount=None,
            gas_limit=self._config.gas_limit_dcdt_nft_burn,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_modifying_royalties(
        self,
        sender: Address,
        token_identifier: str,
        token_nonce: int,
        new_royalties: int,
    ) -> Transaction:
        parts = ["DCDTModifyRoyalties"]

        serialized_parts = self.serializer.serialize_to_parts(
            [
                StringValue(token_identifier),
                BigUIntValue(token_nonce),
                BigUIntValue(new_royalties),
            ]
        )

        parts.extend([part.hex() for part in serialized_parts])

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=sender,
            amount=None,
            gas_limit=self._config.gas_limit_dcdt_modify_royalties,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_setting_new_uris(
        self,
        sender: Address,
        token_identifier: str,
        token_nonce: int,
        new_uris: list[str],
    ) -> Transaction:
        if not new_uris:
            raise BadUsageError("No URIs provided")

        parts = ["DCDTSetNewURIs"]

        serialized_parts = self.serializer.serialize_to_parts(
            [
                StringValue(token_identifier),
                BigUIntValue(token_nonce),
                *map(StringValue, new_uris),
            ]
        )

        parts.extend([part.hex() for part in serialized_parts])

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=sender,
            amount=None,
            gas_limit=self._config.gas_limit_set_new_uris,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_modifying_creator(
        self, sender: Address, token_identifier: str, token_nonce: int
    ) -> Transaction:
        parts = [
            "DCDTModifyCreator",
            self.serializer.serialize([StringValue(token_identifier)]),
            self.serializer.serialize([BigUIntValue(token_nonce)]),
        ]

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=sender,
            amount=None,
            gas_limit=self._config.gas_limit_dcdt_modify_creator,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_updating_metadata(
        self,
        sender: Address,
        token_identifier: str,
        token_nonce: int,
        new_token_name: str,
        new_royalties: int,
        new_hash: str,
        new_attributes: bytes,
        new_uris: list[str],
    ) -> Transaction:
        parts = ["DCDTMetaDataUpdate"]

        serialized_parts = self.serializer.serialize_to_parts(
            [
                StringValue(token_identifier),
                BigUIntValue(token_nonce),
                StringValue(new_token_name),
                BigUIntValue(new_royalties),
                StringValue(new_hash),
                BytesValue(new_attributes),
                *map(StringValue, new_uris),
            ]
        )

        parts.extend([part.hex() for part in serialized_parts])

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=sender,
            amount=None,
            gas_limit=self._config.gas_limit_dcdt_metadata_update,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_nft_metadata_recreate(
        self,
        sender: Address,
        token_identifier: str,
        token_nonce: int,
        new_token_name: str,
        new_royalties: int,
        new_hash: str,
        new_attributes: bytes,
        new_uris: list[str],
    ) -> Transaction:
        parts = ["DCDTMetaDataRecreate"]

        serialized_parts = self.serializer.serialize_to_parts(
            [
                StringValue(token_identifier),
                BigUIntValue(token_nonce),
                StringValue(new_token_name),
                BigUIntValue(new_royalties),
                StringValue(new_hash),
                BytesValue(new_attributes),
                *map(StringValue, new_uris),
            ]
        )

        parts.extend([part.hex() for part in serialized_parts])

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=sender,
            amount=None,
            gas_limit=self._config.gas_limit_nft_metadata_recreate,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_changing_token_to_dynamic(self, sender: Address, token_identifier: str) -> Transaction:
        """The following token types cannot be changed to dynamic: FungibleDCDT, NonFungibleDCDT, NonFungibleDCDTv2"""
        parts = [
            "changeToDynamic",
            self.serializer.serialize([StringValue(token_identifier)]),
        ]

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=None,
            gas_limit=self._config.gas_limit_nft_change_to_dynamic,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_updating_token_id(self, sender: Address, token_identifier: str) -> Transaction:
        parts = [
            "updateTokenID",
            self.serializer.serialize([StringValue(token_identifier)]),
        ]

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=None,
            gas_limit=self._config.gas_limit_update_token_id,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_registering_dynamic_token(
        self,
        sender: Address,
        token_name: str,
        token_ticker: str,
        token_type: TokenType,
        denominator: Optional[int] = None,
    ) -> Transaction:
        if token_type == TokenType.FNG:
            raise Exception("Cannot register fungible token as dynamic")

        parts = ["registerDynamic"]

        serialized_parts = self.serializer.serialize_to_parts(
            [
                StringValue(token_name),
                StringValue(token_ticker),
                StringValue(token_type.value),
            ]
        )

        parts.extend([part.hex() for part in serialized_parts])

        if token_type == TokenType.META and denominator is not None:
            parts.append(self.serializer.serialize([BigUIntValue(denominator)]))

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=self._config.issue_cost,
            gas_limit=self._config.gas_limit_register_dynamic,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_registering_dynamic_and_setting_roles(
        self,
        sender: Address,
        token_name: str,
        token_ticker: str,
        token_type: TokenType,
        denominator: Optional[int] = None,
    ) -> Transaction:
        if token_type == TokenType.FNG:
            raise Exception("Cannot register fungible token as dynamic")

        parts = ["registerAndSetAllRolesDynamic"]

        serialized_parts = self.serializer.serialize_to_parts(
            [
                StringValue(token_name),
                StringValue(token_ticker),
                StringValue(token_type.value),
            ]
        )

        parts.extend([part.hex() for part in serialized_parts])

        if token_type == TokenType.META and denominator is not None:
            parts.append(self.serializer.serialize([BigUIntValue(denominator)]))

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=self._config.issue_cost,
            gas_limit=self._config.gas_limit_register_dynamic,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_transferring_ownership(
        self, sender: Address, token_identifier: str, new_owner: Address
    ) -> Transaction:
        parts = [
            "transferOwnership",
            self.serializer.serialize([StringValue(token_identifier)]),
            new_owner.to_hex(),
        ]

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=0,
            gas_limit=self._config.gas_limit_transfer_ownership,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_freezing_single_nft(
        self, sender: Address, token_identifier: str, token_nonce: int, user: Address
    ) -> Transaction:
        parts = [
            "freezeSingleNFT",
            self.serializer.serialize([StringValue(token_identifier)]),
            self.serializer.serialize([BigUIntValue(token_nonce)]),
            user.to_hex(),
        ]

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=0,
            gas_limit=self._config.gas_limit_freeze_single_nft,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_unfreezing_single_nft(
        self, sender: Address, token_identifier: str, token_nonce: int, user: Address
    ) -> Transaction:
        parts = [
            "unFreezeSingleNFT",
            self.serializer.serialize([StringValue(token_identifier)]),
            self.serializer.serialize([BigUIntValue(token_nonce)]),
            user.to_hex(),
        ]

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=0,
            gas_limit=self._config.gas_limit_unfreeze_single_nft,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_changing_sft_to_meta_dcdt(
        self, sender: Address, collection: str, num_decimals: int
    ) -> Transaction:
        parts = [
            "changeSFTToMetaDCDT",
            self.serializer.serialize([StringValue(collection)]),
            self.serializer.serialize([BigUIntValue(num_decimals)]),
        ]

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=0,
            gas_limit=self._config.gas_limit_change_sft_to_meta_dcdt,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_transferring_nft_create_role(
        self, sender: Address, token_identifier: str, user: Address
    ) -> Transaction:
        """This role can be transferred only if the `canTransferNFTCreateRole` property of the token is set to `true`."""
        parts = [
            "transferNFTCreateRole",
            self.serializer.serialize([StringValue(token_identifier)]),
            sender.to_hex(),
            user.to_hex(),
        ]

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=0,
            gas_limit=self._config.gas_limit_transfer_nft_create_role,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_stopping_nft_creation(self, sender: Address, token_identifier: str) -> Transaction:
        parts = [
            "stopNFTCreate",
            self.serializer.serialize([StringValue(token_identifier)]),
        ]

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=0,
            gas_limit=self._config.gas_limit_stop_nft_create,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transaction_for_wiping_single_nft(
        self, sender: Address, token_identifier: str, token_nonce: int, user: Address
    ) -> Transaction:
        parts = [
            "wipeSingleNFT",
            self.serializer.serialize([StringValue(token_identifier)]),
            self.serializer.serialize([BigUIntValue(token_nonce)]),
            user.to_hex(),
        ]

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=self._config.dcdt_contract_address,
            amount=0,
            gas_limit=self._config.gas_limit_wipe_single_nft,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def create_transction_for_adding_uris(
        self,
        sender: Address,
        token_identifier: str,
        token_nonce: int,
        uris: list[str],
    ) -> Transaction:
        parts = ["DCDTNFTAddURI"]

        serialized_parts = self.serializer.serialize_to_parts(
            [StringValue(token_identifier), BigUIntValue(token_nonce), *map(StringValue, uris)],
        )

        parts.extend([part.hex() for part in serialized_parts])

        return TransactionBuilder(
            config=self._config,
            sender=sender,
            receiver=sender,
            amount=0,
            gas_limit=self._config.gas_limit_dcdt_nft_add_uri,
            add_data_movement_gas=True,
            data_parts=parts,
        ).build()

    def _bool_to_typed_string(self, value: bool) -> StringValue:
        if value:
            return StringValue("true")
        return StringValue("false")
