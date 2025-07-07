use alloc::string::ToString;
use dharitri_chain_core::REWA_000000_TOKEN_IDENTIFIER;

use crate::{
    abi::{TypeAbi, TypeAbiFrom, TypeName},
    api::{
        const_handles, use_raw_handle, ErrorApiImpl, HandleConstraints, ManagedBufferApiImpl,
        ManagedTypeApi,
    },
    codec::*,
    err_msg,
    formatter::{FormatByteReceiver, SCDisplay, SCLowerHex},
    proxy_imports::TestTokenIdentifier,
    types::{ManagedBuffer, ManagedRef, ManagedType, TokenIdentifier},
};

/// Specialized type for handling either REWA or DCDT token identifiers.
///
/// Equivalent to a structure of the form
/// ```
/// # use dharitri_sc::{api::ManagedTypeApi, types::TokenIdentifier};
/// enum RewaOrDcdtTokenIdentifier<M: ManagedTypeApi> {
///     Rewa,
///     Dcdt(TokenIdentifier<M>),
/// }
/// ```
///
/// It is, however more optimized than that. Its implementation is based on `ManagedOption`.
///
/// REWA is indicated by a special, invalid token identifier handle.
/// This way we can fit it inside a single i32 in memory.
#[repr(transparent)]
#[derive(Clone)]
pub struct RewaOrDcdtTokenIdentifier<M: ManagedTypeApi> {
    pub(crate) buffer: ManagedBuffer<M>,
}

impl<M: ManagedTypeApi> ManagedType<M> for RewaOrDcdtTokenIdentifier<M> {
    type OwnHandle = M::ManagedBufferHandle;

    #[inline]
    unsafe fn from_handle(handle: M::ManagedBufferHandle) -> Self {
        RewaOrDcdtTokenIdentifier {
            buffer: ManagedBuffer::from_handle(handle),
        }
    }

    fn get_handle(&self) -> M::ManagedBufferHandle {
        self.buffer.get_handle()
    }

    unsafe fn forget_into_handle(self) -> Self::OwnHandle {
        self.buffer.forget_into_handle()
    }

    fn transmute_from_handle_ref(handle_ref: &M::ManagedBufferHandle) -> &Self {
        unsafe { core::mem::transmute(handle_ref) }
    }

    fn transmute_from_handle_ref_mut(handle_ref: &mut M::ManagedBufferHandle) -> &mut Self {
        unsafe { core::mem::transmute(handle_ref) }
    }
}

impl<M: ManagedTypeApi> RewaOrDcdtTokenIdentifier<M> {
    /// This special representation is interpreted as the REWA token.
    pub const REWA_REPRESENTATION: &'static [u8; 4] = b"REWA";

    /// New instance of the special REWA token representation.
    #[inline]
    pub fn rewa() -> Self {
        RewaOrDcdtTokenIdentifier {
            buffer: ManagedBuffer::from(REWA_000000_TOKEN_IDENTIFIER),
        }
    }

    /// DCDT instance, containing an DCDT token identifier.
    #[inline]
    pub fn dcdt<TI>(token_identifier: TI) -> Self
    where
        TokenIdentifier<M>: From<TI>,
    {
        let ti_obj = TokenIdentifier::from(token_identifier);
        ti_obj.data
    }

    pub fn parse(data: ManagedBuffer<M>) -> Self {
        if data == Self::REWA_REPRESENTATION {
            Self::rewa()
        } else {
            Self { buffer: data }
        }
    }

    #[inline]
    pub fn is_rewa(&self) -> bool {
        M::managed_type_impl().mb_overwrite(
            use_raw_handle(const_handles::MBUF_REWA_000000),
            REWA_000000_TOKEN_IDENTIFIER.as_bytes(),
        );
        M::managed_type_impl().mb_eq(
            use_raw_handle(const_handles::MBUF_REWA_000000),
            self.buffer.handle.clone(),
        )
    }

    #[inline]
    pub fn is_dcdt(&self) -> bool {
        !self.is_rewa()
    }

    /// Returns "REWA" or the token identifier.
    pub fn into_name(self) -> ManagedBuffer<M> {
        self.map_or_else(
            (),
            |()| ManagedBuffer::from(&Self::REWA_REPRESENTATION[..]),
            |(), token_identifier| token_identifier.into_managed_buffer(),
        )
    }

    /// Checks the DCDT token identifier for validity. REWA is considered valid, no checks needed.
    ///
    /// Will fail if it encodes an invalid DCDT token identifier.
    pub fn is_valid(&self) -> bool {
        self.map_ref_or_else(
            (),
            |()| true,
            |(), token_identifier| token_identifier.is_valid_dcdt_identifier(),
        )
    }

    #[inline]
    pub fn into_managed_buffer(self) -> ManagedBuffer<M> {
        self.buffer
    }

    #[inline]
    pub fn as_managed_buffer(&self) -> &ManagedBuffer<M> {
        &self.buffer
    }

    #[inline]
    pub fn to_boxed_bytes(&self) -> crate::types::heap::BoxedBytes {
        self.buffer.to_boxed_bytes()
    }

    pub fn map_or_else<Context, D, F, R>(self, context: Context, for_rewa: D, for_dcdt: F) -> R
    where
        D: FnOnce(Context) -> R,
        F: FnOnce(Context, TokenIdentifier<M>) -> R,
    {
        if self.is_rewa() {
            for_rewa(context)
        } else {
            unsafe { for_dcdt(context, TokenIdentifier::dcdt_unchecked(self)) }
        }
    }

    pub fn map_ref_or_else<Context, D, F, R>(&self, context: Context, for_rewa: D, for_dcdt: F) -> R
    where
        D: FnOnce(Context) -> R,
        F: FnOnce(Context, &TokenIdentifier<M>) -> R,
    {
        if self.is_rewa() {
            for_rewa(context)
        } else {
            unsafe {
                let token_identifier =
                    ManagedRef::<'_, M, TokenIdentifier<M>>::wrap_handle(self.get_handle());
                for_dcdt(context, &token_identifier)
            }
        }
    }

    pub fn unwrap_dcdt(self) -> TokenIdentifier<M> {
        self.map_or_else(
            (),
            |()| {
                M::error_api_impl().signal_error(err_msg::TOKEN_IDENTIFIER_DCDT_EXPECTED.as_bytes())
            },
            |(), token_identifier| token_identifier,
        )
    }

    /// Representation of the object as an `Option`.
    ///
    /// Because it does not consume `self` only a reference to the DCDT token identifier can be returned.
    pub fn as_dcdt_option(&self) -> Option<ManagedRef<'_, M, TokenIdentifier<M>>> {
        if self.is_rewa() {
            None
        } else {
            unsafe {
                Some(ManagedRef::<'_, M, TokenIdentifier<M>>::wrap_handle(
                    self.get_handle(),
                ))
            }
        }
    }

    /// Converts `self` into an `Option`. Consumes `self` in the process.
    pub fn into_dcdt_option(self) -> Option<TokenIdentifier<M>> {
        self.map_or_else((), |()| None, |(), token_identifier| Some(token_identifier))
    }
}

impl<M: ManagedTypeApi> From<ManagedBuffer<M>> for RewaOrDcdtTokenIdentifier<M> {
    #[inline]
    fn from(buffer: ManagedBuffer<M>) -> Self {
        RewaOrDcdtTokenIdentifier { buffer }
    }
}

impl<M: ManagedTypeApi> From<&[u8]> for RewaOrDcdtTokenIdentifier<M> {
    fn from(bytes: &[u8]) -> Self {
        RewaOrDcdtTokenIdentifier {
            buffer: ManagedBuffer::new_from_bytes(bytes),
        }
    }
}

impl<M: ManagedTypeApi> PartialEq for RewaOrDcdtTokenIdentifier<M> {
    #[inline]
    fn eq(&self, other: &Self) -> bool {
        self.buffer == other.buffer
    }
}

impl<M: ManagedTypeApi> Eq for RewaOrDcdtTokenIdentifier<M> {}

impl<M: ManagedTypeApi> PartialEq<TokenIdentifier<M>> for RewaOrDcdtTokenIdentifier<M> {
    #[inline]
    fn eq(&self, other: &TokenIdentifier<M>) -> bool {
        self.map_ref_or_else(
            (),
            |()| false,
            |(), self_dcdt_token_identifier| self_dcdt_token_identifier == other,
        )
    }
}

impl<M: ManagedTypeApi> NestedEncode for RewaOrDcdtTokenIdentifier<M> {
    #[inline]
    fn dep_encode_or_handle_err<O, H>(&self, dest: &mut O, h: H) -> Result<(), H::HandledErr>
    where
        O: NestedEncodeOutput,
        H: EncodeErrorHandler,
    {
        if self.is_rewa() {
            (&Self::REWA_REPRESENTATION[..]).dep_encode_or_handle_err(dest, h)
        } else {
            self.buffer.dep_encode_or_handle_err(dest, h)
        }
    }
}

impl<M: ManagedTypeApi> TopEncode for RewaOrDcdtTokenIdentifier<M> {
    #[inline]
    fn top_encode_or_handle_err<O, H>(&self, output: O, h: H) -> Result<(), H::HandledErr>
    where
        O: TopEncodeOutput,
        H: EncodeErrorHandler,
    {
        if self.is_rewa() {
            (&Self::REWA_REPRESENTATION[..]).top_encode_or_handle_err(output, h)
        } else {
            self.buffer.top_encode_or_handle_err(output, h)
        }
    }
}

impl<M: ManagedTypeApi> NestedDecode for RewaOrDcdtTokenIdentifier<M> {
    fn dep_decode_or_handle_err<I, H>(input: &mut I, h: H) -> Result<Self, H::HandledErr>
    where
        I: NestedDecodeInput,
        H: DecodeErrorHandler,
    {
        Ok(Self::parse(ManagedBuffer::dep_decode_or_handle_err(
            input, h,
        )?))
    }
}

impl<M: ManagedTypeApi> TopDecode for RewaOrDcdtTokenIdentifier<M> {
    fn top_decode_or_handle_err<I, H>(input: I, h: H) -> Result<Self, H::HandledErr>
    where
        I: TopDecodeInput,
        H: DecodeErrorHandler,
    {
        Ok(Self::parse(ManagedBuffer::top_decode_or_handle_err(
            input, h,
        )?))
    }
}

impl<M> TypeAbiFrom<TokenIdentifier<M>> for RewaOrDcdtTokenIdentifier<M> where M: ManagedTypeApi {}
impl<M> TypeAbiFrom<&TokenIdentifier<M>> for RewaOrDcdtTokenIdentifier<M> where M: ManagedTypeApi {}
impl<M> TypeAbiFrom<&[u8]> for RewaOrDcdtTokenIdentifier<M> where M: ManagedTypeApi {}
impl<M> TypeAbiFrom<&str> for RewaOrDcdtTokenIdentifier<M> where M: ManagedTypeApi {}

impl<M> TypeAbiFrom<TestTokenIdentifier<'_>> for RewaOrDcdtTokenIdentifier<M> where M: ManagedTypeApi
{}
impl<M> TypeAbiFrom<&TestTokenIdentifier<'_>> for RewaOrDcdtTokenIdentifier<M> where
    M: ManagedTypeApi
{
}

impl<M: ManagedTypeApi> TypeAbiFrom<Self> for RewaOrDcdtTokenIdentifier<M> {}
impl<M: ManagedTypeApi> TypeAbiFrom<&Self> for RewaOrDcdtTokenIdentifier<M> {}

impl<M: ManagedTypeApi> TypeAbi for RewaOrDcdtTokenIdentifier<M> {
    type Unmanaged = Self;

    fn type_name() -> TypeName {
        "RewaOrDcdtTokenIdentifier".into()
    }

    fn type_name_rust() -> TypeName {
        "RewaOrDcdtTokenIdentifier<$API>".into()
    }
}

impl<M: ManagedTypeApi> SCDisplay for RewaOrDcdtTokenIdentifier<M> {
    fn fmt<F: FormatByteReceiver>(&self, f: &mut F) {
        if self.is_rewa() {
            f.append_bytes(Self::REWA_REPRESENTATION);
        } else {
            let cast_handle = self.buffer.get_handle().cast_or_signal_error::<M, _>();
            let wrap_cast = unsafe { ManagedRef::wrap_handle(cast_handle) };
            f.append_managed_buffer(&wrap_cast);
        }
    }
}

const REWA_REPRESENTATION_HEX: &[u8] = b"45474C44";

impl<M: ManagedTypeApi> SCLowerHex for RewaOrDcdtTokenIdentifier<M> {
    fn fmt<F: FormatByteReceiver>(&self, f: &mut F) {
        if self.is_rewa() {
            f.append_bytes(REWA_REPRESENTATION_HEX);
        } else {
            let cast_handle = self.buffer.get_handle().cast_or_signal_error::<M, _>();
            let wrap_cast = unsafe { ManagedRef::wrap_handle(cast_handle) };
            f.append_managed_buffer_lower_hex(&wrap_cast);
        }
    }
}

impl<M> core::fmt::Debug for RewaOrDcdtTokenIdentifier<M>
where
    M: ManagedTypeApi,
{
    fn fmt(&self, f: &mut core::fmt::Formatter<'_>) -> core::fmt::Result {
        self.map_ref_or_else(
            f,
            |f| f.write_str("RewaOrDcdtTokenIdentifier::Rewa"),
            |f, token_identifier| {
                let token_id_str = token_identifier.to_string();
                f.debug_tuple("RewaOrDcdtTokenIdentifier::Dcdt")
                    .field(&token_id_str)
                    .finish()
            },
        )
    }
}
