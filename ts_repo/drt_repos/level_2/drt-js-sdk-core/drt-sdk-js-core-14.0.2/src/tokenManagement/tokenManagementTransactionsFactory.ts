import { AddressValue, ArgSerializer, BigUIntValue, BytesValue, StringValue } from "../abi";
import { Address } from "../core/address";
import { DCDT_CONTRACT_ADDRESS_HEX } from "../core/constants";
import { ErrBadUsage } from "../core/errors";
import { Logger } from "../core/logger";
import { Transaction } from "../core/transaction";
import { TransactionBuilder } from "../core/transactionBuilder";
import * as resources from "./resources";

interface IConfig {
    chainID: string;
    addressHrp: string;
    minGasLimit: bigint;
    gasLimitPerByte: bigint;
    gasLimitIssue: bigint;
    gasLimitToggleBurnRoleGlobally: bigint;
    gasLimitDcdtLocalMint: bigint;
    gasLimitDcdtLocalBurn: bigint;
    gasLimitSetSpecialRole: bigint;
    gasLimitPausing: bigint;
    gasLimitFreezing: bigint;
    gasLimitWiping: bigint;
    gasLimitDcdtNftCreate: bigint;
    gasLimitDcdtNftUpdateAttributes: bigint;
    gasLimitDcdtNftAddQuantity: bigint;
    gasLimitDcdtNftBurn: bigint;
    gasLimitStorePerByte: bigint;
    gasLimitDcdtModifyRoyalties: bigint;
    gasLimitDcdtModifyCreator: bigint;
    gasLimitDcdtMetadataUpdate: bigint;
    gasLimitSetNewUris: bigint;
    gasLimitNftMetadataRecreate: bigint;
    gasLimitNftChangeToDynamic: bigint;
    gasLimitUpdateTokenId: bigint;
    gasLimitRegisterDynamic: bigint;
    issueCost: bigint;
}

/**
 * Use this class to create token management transactions like issuing DCDTs, creating NFTs, setting roles, etc.
 */
export class TokenManagementTransactionsFactory {
    private readonly config: IConfig;
    private readonly argSerializer: ArgSerializer;
    private readonly trueAsString: string;
    private readonly falseAsString: string;
    private readonly dcdtContractAddress: Address;

    constructor(options: { config: IConfig }) {
        this.config = options.config;
        this.argSerializer = new ArgSerializer();
        this.trueAsString = "true";
        this.falseAsString = "false";
        this.dcdtContractAddress = Address.newFromHex(DCDT_CONTRACT_ADDRESS_HEX, this.config.addressHrp);
    }

    createTransactionForIssuingFungible(sender: Address, options: resources.IssueFungibleInput): Transaction {
        this.notifyAboutUnsettingBurnRoleGlobally();

        const args = [
            new StringValue(options.tokenName),
            new StringValue(options.tokenTicker),
            new BigUIntValue(options.initialSupply),
            new BigUIntValue(options.numDecimals),
            new StringValue("canFreeze"),
            new StringValue(this.boolToString(options.canFreeze)),
            new StringValue("canWipe"),
            new StringValue(this.boolToString(options.canWipe)),
            new StringValue("canPause"),
            new StringValue(this.boolToString(options.canPause)),
            new StringValue("canChangeOwner"),
            new StringValue(this.boolToString(options.canChangeOwner)),
            new StringValue("canUpgrade"),
            new StringValue(this.boolToString(options.canUpgrade)),
            new StringValue("canAddSpecialRoles"),
            new StringValue(this.boolToString(options.canAddSpecialRoles)),
        ];

        const dataParts = ["issue", ...this.argSerializer.valuesToStrings(args)];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: this.dcdtContractAddress,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitIssue,
            addDataMovementGas: true,
            amount: this.config.issueCost,
        }).build();
    }

    createTransactionForIssuingSemiFungible(sender: Address, options: resources.IssueSemiFungibleInput): Transaction {
        this.notifyAboutUnsettingBurnRoleGlobally();

        const args = [
            new StringValue(options.tokenName),
            new StringValue(options.tokenTicker),
            new StringValue("canFreeze"),
            new StringValue(this.boolToString(options.canFreeze)),
            new StringValue("canWipe"),
            new StringValue(this.boolToString(options.canWipe)),
            new StringValue("canPause"),
            new StringValue(this.boolToString(options.canPause)),
            new StringValue("canTransferNFTCreateRole"),
            new StringValue(this.boolToString(options.canTransferNFTCreateRole)),
            new StringValue("canChangeOwner"),
            new StringValue(this.boolToString(options.canChangeOwner)),
            new StringValue("canUpgrade"),
            new StringValue(this.boolToString(options.canUpgrade)),
            new StringValue("canAddSpecialRoles"),
            new StringValue(this.boolToString(options.canAddSpecialRoles)),
        ];

        const dataParts = ["issueSemiFungible", ...this.argSerializer.valuesToStrings(args)];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: this.dcdtContractAddress,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitIssue,
            addDataMovementGas: true,
            amount: this.config.issueCost,
        }).build();
    }

    createTransactionForIssuingNonFungible(sender: Address, options: resources.IssueNonFungibleInput): Transaction {
        this.notifyAboutUnsettingBurnRoleGlobally();

        const args = [
            new StringValue(options.tokenName),
            new StringValue(options.tokenTicker),
            new StringValue("canFreeze"),
            new StringValue(this.boolToString(options.canFreeze)),
            new StringValue("canWipe"),
            new StringValue(this.boolToString(options.canWipe)),
            new StringValue("canPause"),
            new StringValue(this.boolToString(options.canPause)),
            new StringValue("canTransferNFTCreateRole"),
            new StringValue(this.boolToString(options.canTransferNFTCreateRole)),
            new StringValue("canChangeOwner"),
            new StringValue(this.boolToString(options.canChangeOwner)),
            new StringValue("canUpgrade"),
            new StringValue(this.boolToString(options.canUpgrade)),
            new StringValue("canAddSpecialRoles"),
            new StringValue(this.boolToString(options.canAddSpecialRoles)),
        ];

        const dataParts = ["issueNonFungible", ...this.argSerializer.valuesToStrings(args)];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: this.dcdtContractAddress,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitIssue,
            addDataMovementGas: true,
            amount: this.config.issueCost,
        }).build();
    }

    createTransactionForRegisteringMetaDCDT(sender: Address, options: resources.RegisterMetaDCDTInput): Transaction {
        this.notifyAboutUnsettingBurnRoleGlobally();

        const args = [
            new StringValue(options.tokenName),
            new StringValue(options.tokenTicker),
            new BigUIntValue(options.numDecimals),
            new StringValue("canFreeze"),
            new StringValue(this.boolToString(options.canFreeze)),
            new StringValue("canWipe"),
            new StringValue(this.boolToString(options.canWipe)),
            new StringValue("canPause"),
            new StringValue(this.boolToString(options.canPause)),
            new StringValue("canTransferNFTCreateRole"),
            new StringValue(this.boolToString(options.canTransferNFTCreateRole)),
            new StringValue("canChangeOwner"),
            new StringValue(this.boolToString(options.canChangeOwner)),
            new StringValue("canUpgrade"),
            new StringValue(this.boolToString(options.canUpgrade)),
            new StringValue("canAddSpecialRoles"),
            new StringValue(this.boolToString(options.canAddSpecialRoles)),
        ];

        const dataParts = ["registerMetaDCDT", ...this.argSerializer.valuesToStrings(args)];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: this.dcdtContractAddress,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitIssue,
            addDataMovementGas: true,
            amount: this.config.issueCost,
        }).build();
    }

    createTransactionForRegisteringAndSettingRoles(
        sender: Address,
        options: resources.RegisterRolesInput,
    ): Transaction {
        this.notifyAboutUnsettingBurnRoleGlobally();

        const dataParts = [
            "registerAndSetAllRoles",
            ...this.argSerializer.valuesToStrings([
                new StringValue(options.tokenName),
                new StringValue(options.tokenTicker),
                new StringValue(options.tokenType),
                new BigUIntValue(options.numDecimals),
            ]),
        ];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: this.dcdtContractAddress,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitIssue,
            addDataMovementGas: true,
            amount: this.config.issueCost,
        }).build();
    }

    createTransactionForSettingBurnRoleGlobally(
        sender: Address,
        options: resources.BurnRoleGloballyInput,
    ): Transaction {
        const dataParts = [
            "setBurnRoleGlobally",
            ...this.argSerializer.valuesToStrings([new StringValue(options.tokenIdentifier)]),
        ];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: this.dcdtContractAddress,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitToggleBurnRoleGlobally,
            addDataMovementGas: true,
        }).build();
    }

    createTransactionForUnsettingBurnRoleGlobally(
        sender: Address,
        options: resources.BurnRoleGloballyInput,
    ): Transaction {
        const dataParts = [
            "unsetBurnRoleGlobally",
            ...this.argSerializer.valuesToStrings([new StringValue(options.tokenIdentifier)]),
        ];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: this.dcdtContractAddress,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitToggleBurnRoleGlobally,
            addDataMovementGas: true,
        }).build();
    }

    createTransactionForSettingSpecialRoleOnFungibleToken(
        sender: Address,
        options: resources.FungibleSpecialRoleInput,
    ): Transaction {
        const args = [new StringValue(options.tokenIdentifier), new AddressValue(options.user)];

        options.addRoleLocalMint ? args.push(new StringValue("DCDTRoleLocalMint")) : 0;
        options.addRoleLocalBurn ? args.push(new StringValue("DCDTRoleLocalBurn")) : 0;
        options.addRoleDCDTTransferRole ? args.push(new StringValue("DCDTTransferRole")) : 0;

        const dataParts = ["setSpecialRole", ...this.argSerializer.valuesToStrings(args)];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: this.dcdtContractAddress,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitSetSpecialRole,
            addDataMovementGas: true,
        }).build();
    }

    createTransactionForUnsettingSpecialRoleOnFungibleToken(
        sender: Address,
        options: resources.UnsetFungibleSpecialRoleInput,
    ): Transaction {
        const args = [new StringValue(options.tokenIdentifier), new AddressValue(options.user)];

        options.removeRoleLocalMint ? args.push(new StringValue("DCDTRoleLocalMint")) : 0;
        options.removeRoleDCDTTransferRole ? args.push(new StringValue("DCDTRoleLocalBurn")) : 0;
        options.removeRoleDCDTTransferRole ? args.push(new StringValue("DCDTTransferRole")) : 0;

        const dataParts = ["unSetSpecialRole", ...this.argSerializer.valuesToStrings(args)];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: this.dcdtContractAddress,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitSetSpecialRole,
            addDataMovementGas: true,
        }).build();
    }

    createTransactionForSettingSpecialRoleOnSemiFungibleToken(
        sender: Address,
        options: resources.SemiFungibleSpecialRoleInput,
    ): Transaction {
        const args = [new StringValue(options.tokenIdentifier), new AddressValue(options.user)];

        options.addRoleNFTCreate ? args.push(new StringValue("DCDTRoleNFTCreate")) : 0;
        options.addRoleNFTBurn ? args.push(new StringValue("DCDTRoleNFTBurn")) : 0;
        options.addRoleNFTAddQuantity ? args.push(new StringValue("DCDTRoleNFTAddQuantity")) : 0;
        options.addRoleDCDTTransferRole ? args.push(new StringValue("DCDTTransferRole")) : 0;

        const dataParts = ["setSpecialRole", ...this.argSerializer.valuesToStrings(args)];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: this.dcdtContractAddress,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitSetSpecialRole,
            addDataMovementGas: true,
        }).build();
    }

    createTransactionForUnsettingSpecialRoleOnSemiFungibleToken(
        sender: Address,
        options: resources.UnsetSemiFungibleSpecialRoleInput,
    ): Transaction {
        const args = [new StringValue(options.tokenIdentifier), new AddressValue(options.user)];

        options.removeRoleNFTBurn ? args.push(new StringValue("DCDTRoleNFTBurn")) : 0;
        options.removeRoleNFTAddQuantity ? args.push(new StringValue("DCDTRoleNFTAddQuantity")) : 0;
        options.removeRoleDCDTTransferRole ? args.push(new StringValue("DCDTTransferRole")) : 0;

        const dataParts = ["unSetSpecialRole", ...this.argSerializer.valuesToStrings(args)];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: this.dcdtContractAddress,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitSetSpecialRole,
            addDataMovementGas: true,
        }).build();
    }

    createTransactionForSettingSpecialRoleOnMetaDCDT(
        sender: Address,
        options: resources.SemiFungibleSpecialRoleInput,
    ): Transaction {
        return this.createTransactionForSettingSpecialRoleOnSemiFungibleToken(sender, options);
    }

    createTransactionForUnsettingSpecialRoleOnMetaDCDT(
        sender: Address,
        options: resources.UnsetSemiFungibleSpecialRoleInput,
    ): Transaction {
        return this.createTransactionForUnsettingSpecialRoleOnSemiFungibleToken(sender, options);
    }

    createTransactionForSettingSpecialRoleOnNonFungibleToken(
        sender: Address,
        options: resources.SpecialRoleInput,
    ): Transaction {
        const args = [new StringValue(options.tokenIdentifier), new AddressValue(options.user)];

        options.addRoleNFTCreate ? args.push(new StringValue("DCDTRoleNFTCreate")) : 0;
        options.addRoleNFTBurn ? args.push(new StringValue("DCDTRoleNFTBurn")) : 0;
        options.addRoleNFTUpdateAttributes ? args.push(new StringValue("DCDTRoleNFTUpdateAttributes")) : 0;
        options.addRoleNFTAddURI ? args.push(new StringValue("DCDTRoleNFTAddURI")) : 0;
        options.addRoleDCDTTransferRole ? args.push(new StringValue("DCDTTransferRole")) : 0;
        options.addRoleDCDTModifyCreator ? args.push(new StringValue("DCDTRoleModifyCreator")) : 0;
        options.addRoleNFTRecreate ? args.push(new StringValue("DCDTRoleNFTRecreate")) : 0;
        options.addRoleDCDTSetNewURI ? args.push(new StringValue("DCDTRoleSetNewURI")) : 0;
        options.addRoleDCDTModifyRoyalties ? args.push(new StringValue("DCDTRoleModifyRoyalties")) : 0;

        const dataParts = ["setSpecialRole", ...this.argSerializer.valuesToStrings(args)];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: this.dcdtContractAddress,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitSetSpecialRole,
            addDataMovementGas: true,
        }).build();
    }

    createTransactionForUnsettingSpecialRoleOnNonFungibleToken(
        sender: Address,
        options: resources.UnsetSpecialRoleInput,
    ): Transaction {
        const args = [new StringValue(options.tokenIdentifier), new AddressValue(options.user)];

        options.removeRoleNFTBurn ? args.push(new StringValue("DCDTRoleNFTBurn")) : 0;
        options.removeRoleNFTUpdateAttributes ? args.push(new StringValue("DCDTRoleNFTUpdateAttributes")) : 0;
        options.removeRoleNFTAddURI ? args.push(new StringValue("DCDTRoleNFTAddURI")) : 0;
        options.removeRoleDCDTTransferRole ? args.push(new StringValue("DCDTTransferRole")) : 0;
        options.removeRoleDCDTModifyCreator ? args.push(new StringValue("DCDTRoleModifyCreator")) : 0;
        options.removeRoleNFTRecreate ? args.push(new StringValue("DCDTRoleNFTRecreate")) : 0;
        options.removeRoleDCDTSetNewURI ? args.push(new StringValue("DCDTRoleSetNewURI")) : 0;
        options.removeRoleDCDTModifyRoyalties ? args.push(new StringValue("DCDTRoleModifyRoyalties")) : 0;

        const dataParts = ["unSetSpecialRole", ...this.argSerializer.valuesToStrings(args)];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: this.dcdtContractAddress,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitSetSpecialRole,
            addDataMovementGas: true,
        }).build();
    }

    createTransactionForCreatingNFT(sender: Address, options: resources.MintInput): Transaction {
        const dataParts = [
            "DCDTNFTCreate",
            ...this.argSerializer.valuesToStrings([
                new StringValue(options.tokenIdentifier),
                new BigUIntValue(options.initialQuantity ?? 1n),
                new StringValue(options.name),
                new BigUIntValue(options.royalties),
                new StringValue(options.hash),
                new BytesValue(Buffer.from(options.attributes)),
                ...options.uris.map((uri) => new StringValue(uri)),
            ]),
        ];

        // Note that the following is an approximation (a reasonable one):
        const nftData = options.name + options.hash + options.attributes + options.uris.join("");
        const storageGasLimit = this.config.gasLimitStorePerByte + BigInt(nftData.length);

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: sender,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitDcdtNftCreate + storageGasLimit,
            addDataMovementGas: true,
        }).build();
    }

    createTransactionForPausing(sender: Address, options: resources.PausingInput): Transaction {
        const dataParts = ["pause", ...this.argSerializer.valuesToStrings([new StringValue(options.tokenIdentifier)])];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: sender,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitPausing,
            addDataMovementGas: true,
        }).build();
    }

    createTransactionForUnpausing(sender: Address, options: resources.PausingInput): Transaction {
        const dataParts = [
            "unPause",
            ...this.argSerializer.valuesToStrings([new StringValue(options.tokenIdentifier)]),
        ];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: sender,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitPausing,
            addDataMovementGas: true,
        }).build();
    }

    createTransactionForFreezing(sender: Address, options: resources.ManagementInput): Transaction {
        const dataParts = [
            "freeze",
            ...this.argSerializer.valuesToStrings([
                new StringValue(options.tokenIdentifier),
                new AddressValue(options.user),
            ]),
        ];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: sender,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitFreezing,
            addDataMovementGas: true,
        }).build();
    }

    createTransactionForUnfreezing(sender: Address, options: resources.ManagementInput): Transaction {
        const dataParts = [
            "UnFreeze",
            ...this.argSerializer.valuesToStrings([
                new StringValue(options.tokenIdentifier),
                new AddressValue(options.user),
            ]),
        ];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: sender,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitFreezing,
            addDataMovementGas: true,
        }).build();
    }

    createTransactionForWiping(sender: Address, options: resources.ManagementInput): Transaction {
        const dataParts = [
            "wipe",
            ...this.argSerializer.valuesToStrings([
                new StringValue(options.tokenIdentifier),
                new AddressValue(options.user),
            ]),
        ];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: sender,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitWiping,
            addDataMovementGas: true,
        }).build();
    }

    createTransactionForLocalMint(sender: Address, options: resources.LocalMintInput): Transaction {
        const dataParts = [
            "DCDTLocalMint",
            ...this.argSerializer.valuesToStrings([
                new StringValue(options.tokenIdentifier),
                new BigUIntValue(options.supplyToMint),
            ]),
        ];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: sender,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitDcdtLocalMint,
            addDataMovementGas: true,
        }).build();
    }

    createTransactionForLocalBurning(sender: Address, options: resources.LocalBurnInput): Transaction {
        const dataParts = [
            "DCDTLocalBurn",
            ...this.argSerializer.valuesToStrings([
                new StringValue(options.tokenIdentifier),
                new BigUIntValue(options.supplyToBurn),
            ]),
        ];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: sender,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitDcdtLocalBurn,
            addDataMovementGas: true,
        }).build();
    }

    createTransactionForUpdatingAttributes(sender: Address, options: resources.UpdateAttributesInput): Transaction {
        const dataParts = [
            "DCDTNFTUpdateAttributes",
            ...this.argSerializer.valuesToStrings([
                new StringValue(options.tokenIdentifier),
                new BigUIntValue(options.tokenNonce),
                new BytesValue(Buffer.from(options.attributes)),
            ]),
        ];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: sender,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitDcdtNftUpdateAttributes,
            addDataMovementGas: true,
        }).build();
    }

    createTransactionForAddingQuantity(sender: Address, options: resources.UpdateQuantityInput): Transaction {
        const dataParts = [
            "DCDTNFTAddQuantity",
            ...this.argSerializer.valuesToStrings([
                new StringValue(options.tokenIdentifier),
                new BigUIntValue(options.tokenNonce),
                new BigUIntValue(options.quantity),
            ]),
        ];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: sender,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitDcdtNftAddQuantity,
            addDataMovementGas: true,
        }).build();
    }

    createTransactionForBurningQuantity(sender: Address, options: resources.UpdateQuantityInput): Transaction {
        const dataParts = [
            "DCDTNFTBurn",
            ...this.argSerializer.valuesToStrings([
                new StringValue(options.tokenIdentifier),
                new BigUIntValue(options.tokenNonce),
                new BigUIntValue(options.quantity),
            ]),
        ];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: sender,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitDcdtNftBurn,
            addDataMovementGas: true,
        }).build();
    }

    createTransactionForModifyingRoyalties(sender: Address, options: resources.ModifyRoyaltiesInput): Transaction {
        const dataParts = [
            "DCDTModifyRoyalties",
            ...this.argSerializer.valuesToStrings([
                new StringValue(options.tokenIdentifier),
                new BigUIntValue(options.tokenNonce),
                new BigUIntValue(options.newRoyalties),
            ]),
        ];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: sender,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitDcdtModifyRoyalties,
            addDataMovementGas: true,
        }).build();
    }

    createTransactionForSettingNewUris(sender: Address, options: resources.SetNewUriInput): Transaction {
        if (!options.newUris.length) {
            throw new ErrBadUsage("No URIs provided");
        }

        const dataParts = [
            "DCDTSetNewURIs",
            ...this.argSerializer.valuesToStrings([
                new StringValue(options.tokenIdentifier),
                new BigUIntValue(options.tokenNonce),
                ...options.newUris.map((uri) => new StringValue(uri)),
            ]),
        ];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: sender,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitSetNewUris,
            addDataMovementGas: true,
        }).build();
    }

    createTransactionForModifyingCreator(sender: Address, options: resources.ModifyCreatorInput): Transaction {
        const dataParts = [
            "DCDTModifyCreator",
            ...this.argSerializer.valuesToStrings([
                new StringValue(options.tokenIdentifier),
                new BigUIntValue(options.tokenNonce),
            ]),
        ];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: sender,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitDcdtModifyCreator,
            addDataMovementGas: true,
        }).build();
    }

    createTransactionForUpdatingMetadata(sender: Address, options: resources.ManageMetadataInput): Transaction {
        const dataParts = [
            "DCDTMetaDataUpdate",
            ...this.argSerializer.valuesToStrings([
                new StringValue(options.tokenIdentifier),
                new BigUIntValue(options.tokenNonce),
                ...(options.newTokenName ? [new StringValue(options.newTokenName)] : []),
                ...(options.newRoyalties ? [new BigUIntValue(options.newRoyalties)] : []),
                ...(options.newHash ? [new StringValue(options.newHash)] : []),
                ...(options.newAttributes ? [new BytesValue(Buffer.from(options.newAttributes))] : []),
                ...(options.newUris ? options.newUris.map((uri) => new StringValue(uri)) : []),
            ]),
        ];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: sender,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitDcdtMetadataUpdate,
            addDataMovementGas: true,
        }).build();
    }

    createTransactionForMetadataRecreate(sender: Address, options: resources.ManageMetadataInput): Transaction {
        const dataParts = [
            "DCDTMetaDataRecreate",
            ...this.argSerializer.valuesToStrings([
                new StringValue(options.tokenIdentifier),
                new BigUIntValue(options.tokenNonce),
                ...(options.newTokenName ? [new StringValue(options.newTokenName)] : []),
                ...(options.newRoyalties ? [new BigUIntValue(options.newRoyalties)] : []),
                ...(options.newHash ? [new StringValue(options.newHash)] : []),
                ...(options.newAttributes ? [new BytesValue(Buffer.from(options.newAttributes))] : []),
                ...(options.newUris ? options.newUris.map((uri) => new StringValue(uri)) : []),
            ]),
        ];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: sender,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitNftMetadataRecreate,
            addDataMovementGas: true,
        }).build();
    }

    createTransactionForChangingTokenToDynamic(
        sender: Address,
        options: resources.ChangeTokenToDynamicInput,
    ): Transaction {
        const dataParts = [
            "changeToDynamic",
            ...this.argSerializer.valuesToStrings([new StringValue(options.tokenIdentifier)]),
        ];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: this.dcdtContractAddress,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitNftChangeToDynamic,
            addDataMovementGas: true,
        }).build();
    }

    createTransactionForUpdatingTokenId(sender: Address, options: resources.UpdateTokenIDInput): Transaction {
        const dataParts = [
            "updateTokenID",
            ...this.argSerializer.valuesToStrings([new StringValue(options.tokenIdentifier)]),
        ];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: this.dcdtContractAddress,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitUpdateTokenId,
            addDataMovementGas: true,
        }).build();
    }

    createTransactionForRegisteringDynamicToken(
        sender: Address,
        options: resources.RegisteringDynamicTokenInput,
    ): Transaction {
        const dataParts = [
            "registerDynamic",
            ...this.argSerializer.valuesToStrings([
                new StringValue(options.tokenName),
                new StringValue(options.tokenTicker),
                new StringValue(options.tokenType),
            ]),
        ];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: this.dcdtContractAddress,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitRegisterDynamic,
            addDataMovementGas: true,
            amount: this.config.issueCost,
        }).build();
    }

    createTransactionForRegisteringDynamicAndSettingRoles(
        sender: Address,
        options: resources.RegisteringDynamicTokenInput,
    ): Transaction {
        const dataParts = [
            "registerAndSetAllRolesDynamic",
            ...this.argSerializer.valuesToStrings([
                new StringValue(options.tokenName),
                new StringValue(options.tokenTicker),
                new StringValue(options.tokenType),
            ]),
        ];

        return new TransactionBuilder({
            config: this.config,
            sender: sender,
            receiver: this.dcdtContractAddress,
            dataParts: dataParts,
            gasLimit: this.config.gasLimitRegisterDynamic,
            addDataMovementGas: true,
            amount: this.config.issueCost,
        }).build();
    }

    private notifyAboutUnsettingBurnRoleGlobally() {
        Logger.info(`
==========
IMPORTANT!
==========
You are about to issue (register) a new token. This will set the role "DCDTRoleBurnForAll" (globally).
Once the token is registered, you can unset this role by calling "unsetBurnRoleGlobally" (in a separate transaction).`);
    }

    private boolToString(value: boolean): string {
        if (value) {
            return this.trueAsString;
        }

        return this.falseAsString;
    }
}
