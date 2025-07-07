use super::DcdtLocalRoleFlags;
use crate::codec::{
    self,
    derive::{NestedDecode, NestedEncode, TopDecode, TopEncode},
};

const DCDT_ROLE_NONE: &str = "";
const DCDT_ROLE_LOCAL_MINT: &str = "DCDTRoleLocalMint";
const DCDT_ROLE_LOCAL_BURN: &str = "DCDTRoleLocalBurn";
const DCDT_ROLE_NFT_CREATE: &str = "DCDTRoleNFTCreate";
const DCDT_ROLE_NFT_ADD_QUANTITY: &str = "DCDTRoleNFTAddQuantity";
const DCDT_ROLE_NFT_BURN: &str = "DCDTRoleNFTBurn";
const DCDT_ROLE_NFT_ADD_URI: &str = "DCDTRoleNFTAddURI";
const DCDT_ROLE_NFT_UPDATE_ATTRIBUTES: &str = "DCDTRoleNFTUpdateAttributes";
const DCDT_ROLE_SET_NEW_URI: &str = "DCDTRoleSetNewURI";
const DCDT_ROLE_MODIFY_ROYALTIES: &str = "DCDTRoleModifyRoyalties";
const DCDT_ROLE_MODIFY_CREATOR: &str = "DCDTRoleModifyCreator";
const DCDT_ROLE_NFT_RECREATE: &str = "DCDTRoleNFTRecreate";
const DCDT_ROLE_TRANSFER: &str = "DCDTTransferRole";

#[derive(TopDecode, TopEncode, NestedDecode, NestedEncode, Clone, PartialEq, Eq, Debug, Copy)]
pub enum DcdtLocalRole {
    None,
    Mint,
    Burn,
    NftCreate,
    NftAddQuantity,
    NftBurn,
    NftUpdateAttributes,
    NftAddUri,
    NftRecreate,
    ModifyRoyalties,
    ModifyCreator,
    SetNewUri,
    Transfer,
}

impl DcdtLocalRole {
    pub fn as_u16(&self) -> u16 {
        match self {
            Self::None => 0,
            Self::Mint => 1,
            Self::Burn => 2,
            Self::NftCreate => 3,
            Self::NftAddQuantity => 4,
            Self::NftBurn => 5,
            Self::NftUpdateAttributes => 6,
            Self::NftAddUri => 7,
            Self::NftRecreate => 8,
            Self::ModifyCreator => 9,
            Self::ModifyRoyalties => 10,
            Self::SetNewUri => 11,
            Self::Transfer => 12,
        }
    }

    pub fn as_role_name(&self) -> &'static [u8] {
        self.name().as_bytes()
    }

    pub fn name(&self) -> &'static str {
        match self {
            Self::None => DCDT_ROLE_NONE,
            Self::Mint => DCDT_ROLE_LOCAL_MINT,
            Self::Burn => DCDT_ROLE_LOCAL_BURN,
            Self::NftCreate => DCDT_ROLE_NFT_CREATE,
            Self::NftAddQuantity => DCDT_ROLE_NFT_ADD_QUANTITY,
            Self::NftBurn => DCDT_ROLE_NFT_BURN,
            Self::NftUpdateAttributes => DCDT_ROLE_NFT_UPDATE_ATTRIBUTES,
            Self::NftAddUri => DCDT_ROLE_NFT_ADD_URI,
            Self::NftRecreate => DCDT_ROLE_NFT_RECREATE,
            Self::ModifyRoyalties => DCDT_ROLE_MODIFY_ROYALTIES,
            Self::ModifyCreator => DCDT_ROLE_MODIFY_CREATOR,
            Self::SetNewUri => DCDT_ROLE_SET_NEW_URI,
            Self::Transfer => DCDT_ROLE_TRANSFER,
        }
    }

    pub fn to_flag(&self) -> DcdtLocalRoleFlags {
        match self {
            Self::None => DcdtLocalRoleFlags::NONE,
            Self::Mint => DcdtLocalRoleFlags::MINT,
            Self::Burn => DcdtLocalRoleFlags::BURN,
            Self::NftCreate => DcdtLocalRoleFlags::NFT_CREATE,
            Self::NftAddQuantity => DcdtLocalRoleFlags::NFT_ADD_QUANTITY,
            Self::NftBurn => DcdtLocalRoleFlags::NFT_BURN,
            Self::NftUpdateAttributes => DcdtLocalRoleFlags::NFT_UPDATE_ATTRIBUTES,
            Self::NftAddUri => DcdtLocalRoleFlags::NFT_ADD_URI,
            Self::NftRecreate => DcdtLocalRoleFlags::NFT_RECREATE,
            Self::ModifyRoyalties => DcdtLocalRoleFlags::MODIFY_ROYALTIES,
            Self::ModifyCreator => DcdtLocalRoleFlags::MODIFY_CREATOR,
            Self::SetNewUri => DcdtLocalRoleFlags::SET_NEW_URI,
            Self::Transfer => DcdtLocalRoleFlags::TRANSFER,
        }
    }
}

// TODO: can be done with macros, but I didn't find a public library that does it and is no_std
// we can implement it, it's easy
const ALL_ROLES: [DcdtLocalRole; 12] = [
    DcdtLocalRole::Mint,
    DcdtLocalRole::Burn,
    DcdtLocalRole::NftCreate,
    DcdtLocalRole::NftAddQuantity,
    DcdtLocalRole::NftBurn,
    DcdtLocalRole::NftUpdateAttributes,
    DcdtLocalRole::NftAddUri,
    DcdtLocalRole::NftRecreate,
    DcdtLocalRole::ModifyRoyalties,
    DcdtLocalRole::ModifyCreator,
    DcdtLocalRole::SetNewUri,
    DcdtLocalRole::Transfer,
];

impl DcdtLocalRole {
    pub fn iter_all() -> core::slice::Iter<'static, DcdtLocalRole> {
        ALL_ROLES.iter()
    }
}

impl From<u16> for DcdtLocalRole {
    #[inline]
    fn from(value: u16) -> Self {
        match value {
            1 => Self::Mint,
            2 => Self::Burn,
            3 => Self::NftCreate,
            4 => Self::NftAddQuantity,
            5 => Self::NftBurn,
            6 => Self::NftUpdateAttributes,
            7 => Self::NftAddUri,
            8 => Self::NftRecreate,
            9 => Self::ModifyRoyalties,
            10 => Self::ModifyCreator,
            11 => Self::SetNewUri,
            12 => Self::Transfer,
            _ => Self::None,
        }
    }
}

impl<'a> From<&'a [u8]> for DcdtLocalRole {
    #[inline]
    fn from(byte_slice: &'a [u8]) -> Self {
        if byte_slice == DCDT_ROLE_LOCAL_MINT.as_bytes() {
            Self::Mint
        } else if byte_slice == DCDT_ROLE_LOCAL_BURN.as_bytes() {
            Self::Burn
        } else if byte_slice == DCDT_ROLE_NFT_CREATE.as_bytes() {
            Self::NftCreate
        } else if byte_slice == DCDT_ROLE_NFT_ADD_QUANTITY.as_bytes() {
            Self::NftAddQuantity
        } else if byte_slice == DCDT_ROLE_NFT_BURN.as_bytes() {
            Self::NftBurn
        } else if byte_slice == DCDT_ROLE_NFT_UPDATE_ATTRIBUTES.as_bytes() {
            Self::NftUpdateAttributes
        } else if byte_slice == DCDT_ROLE_NFT_ADD_URI.as_bytes() {
            Self::NftAddUri
        } else if byte_slice == DCDT_ROLE_NFT_RECREATE.as_bytes() {
            Self::NftRecreate
        } else if byte_slice == DCDT_ROLE_MODIFY_ROYALTIES.as_bytes() {
            Self::ModifyRoyalties
        } else if byte_slice == DCDT_ROLE_MODIFY_CREATOR.as_bytes() {
            Self::ModifyCreator
        } else if byte_slice == DCDT_ROLE_SET_NEW_URI.as_bytes() {
            Self::SetNewUri
        } else if byte_slice == DCDT_ROLE_TRANSFER.as_bytes() {
            Self::Transfer
        } else {
            Self::None
        }
    }
}
