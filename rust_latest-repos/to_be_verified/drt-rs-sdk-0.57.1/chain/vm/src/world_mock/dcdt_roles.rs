use std::fmt::{self, Write};

#[derive(Clone, Default, Debug)]
pub struct DcdtRoles(Vec<Vec<u8>>);

impl DcdtRoles {
    pub fn new(roles: Vec<Vec<u8>>) -> Self {
        DcdtRoles(roles)
    }

    pub fn is_empty(&self) -> bool {
        self.0.is_empty()
    }

    pub fn get(&self) -> Vec<Vec<u8>> {
        self.0.clone()
    }
}

impl fmt::Display for DcdtRoles {
    fn fmt(&self, _f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let mut dcdt_buf = String::new();
        let dcdt_keys: Vec<Vec<u8>> = self.clone().0.to_vec();

        for value in &dcdt_keys {
            write!(dcdt_buf, "{}", hex::encode(value.as_slice()))?;
        }
        Ok(())
    }
}
