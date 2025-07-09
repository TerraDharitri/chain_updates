use common_types::PortId;

dharitri_sc::imports!();

const MIN_PORT_LEN: usize = 2;
const MAX_PORT_LEN: usize = 128;
const SPECIAL_CHARS: &[u8] = b"._+-#[]<>";

#[dharitri_sc::module]
pub trait HostLibModule {
    /// check if the string consist of characters in one of the following categories only:
    ///
    /// - Alphanumeric
    ///
    /// - `.`, `_`, `+`, `-`, `#`
    ///
    /// - `[`, `]`, `<`, `>`
    fn is_valid_port_id(&self, port_id: &PortId<Self::Api>) -> bool {
        let port_len = port_id.len();

        // Clippy suggestion makes code unreadable
        #[allow(clippy::manual_range_contains)]
        if port_len < MIN_PORT_LEN || port_len > MAX_PORT_LEN {
            return false;
        }

        let mut as_array = [0; MAX_PORT_LEN];
        let slice = port_id.load_to_byte_array(&mut as_array);
        for character in slice {
            if !character.is_ascii_alphanumeric() && !SPECIAL_CHARS.contains(character) {
                return false;
            }
        }

        true
    }
}
