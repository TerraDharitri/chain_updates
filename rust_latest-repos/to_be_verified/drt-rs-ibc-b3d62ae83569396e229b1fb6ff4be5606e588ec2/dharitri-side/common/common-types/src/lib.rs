#![no_std]

use connection_types::version;

dharitri_sc::imports!();

pub mod channel_types;
pub mod connection_types;

pub const HASH_LENGTH: usize = 32;
pub type Hash<M> = ManagedByteArray<M, HASH_LENGTH>;

pub type UnixTimestamp = u64;
pub type Sequence = u64;

pub type ChainId<M> = ManagedByteArray<M, 32>;
pub type ClientId<M> = ManagedBuffer<M>;
pub type ClientType<M> = ManagedBuffer<M>;
pub type ConnectionId<M> = ManagedBuffer<M>;
pub type ChannelId<M> = ManagedBuffer<M>;
pub type Feature<M> = ManagedBuffer<M>;
pub type FeatureId<M> = ManagedBuffer<M>;
pub type PortId<M> = ManagedBuffer<M>;
pub type Path<M> = ManagedBuffer<M>;
pub type Version<M> = ManagedBuffer<M>;

pub type VersionVec<M> = ManagedVec<M, version::Data<M>>;
pub type FeatureVec<M> = ManagedVec<M, Feature<M>>;
pub type ConnectionHops<M> = ManagedVec<M, ConnectionId<M>>;
