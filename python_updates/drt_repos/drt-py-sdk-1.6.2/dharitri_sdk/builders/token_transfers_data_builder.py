from dharitri_sdk.abi import BigUIntValue, Serializer
from dharitri_sdk.abi.string_value import StringValue
from dharitri_sdk.core import Address, TokenComputer, TokenTransfer


class TokenTransfersDataBuilder:
    """
    **FOR INTERNAL USE ONLY.**
    Used for the transactions factories.
    """

    def __init__(self, token_computer: TokenComputer) -> None:
        self.token_computer = token_computer
        self.serializer = Serializer()

    def build_args_for_dcdt_transfer(self, transfer: TokenTransfer) -> list[str]:
        args = ["DCDTTransfer"]

        serialized_args = self.serializer.serialize_to_parts(
            [StringValue(transfer.token.identifier), BigUIntValue(transfer.amount)]
        )
        args.extend([arg.hex() for arg in serialized_args])

        return args

    def build_args_for_single_dcdt_nft_transfer(self, transfer: TokenTransfer, receiver: Address) -> list[str]:
        args = ["DCDTNFTTransfer"]
        token = transfer.token
        identifier = self.token_computer.extract_identifier_from_extended_identifier(token.identifier)

        serialized_args = self.serializer.serialize_to_parts(
            [
                StringValue(identifier),
                BigUIntValue(token.nonce),
                BigUIntValue(transfer.amount),
            ]
        )
        args.extend([arg.hex() for arg in serialized_args])
        args.append(receiver.to_hex())

        return args

    def build_args_for_multi_dcdt_nft_transfer(self, receiver: Address, transfers: list[TokenTransfer]) -> list[str]:
        serialized_num_of_transfers = self.serializer.serialize([BigUIntValue(len(transfers))])
        args = ["MultiDCDTNFTTransfer", receiver.to_hex(), serialized_num_of_transfers]

        for transfer in transfers:
            identifier = self.token_computer.extract_identifier_from_extended_identifier(transfer.token.identifier)
            serialized_args = self.serializer.serialize_to_parts(
                [
                    StringValue(identifier),
                    BigUIntValue(transfer.token.nonce),
                    BigUIntValue(transfer.amount),
                ]
            )
            args.extend([arg.hex() for arg in serialized_args])

        return args
