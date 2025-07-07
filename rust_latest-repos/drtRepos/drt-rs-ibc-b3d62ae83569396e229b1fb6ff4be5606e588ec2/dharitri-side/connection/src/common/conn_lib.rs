use common_types::{
    channel_types::channel::{ORDERED, UNORDERED},
    connection_types::version,
    Feature, FeatureId, FeatureVec, VersionVec,
};

dharitri_sc::imports!();

static IBC_VERSION_IDENTIFIER: &[u8] = b"1";

#[dharitri_sc::module]
pub trait ConnectionLibModule {
    /// returns the latest supported version of IBC used in connection version negotiation
    fn default_ibc_version(&self) -> version::Data<Self::Api> {
        let mut version = version::Data {
            identifier: FeatureId::from(IBC_VERSION_IDENTIFIER),
            features: FeatureVec::new(),
        };
        version.features.push(Feature::from(ORDERED));
        version.features.push(Feature::from(UNORDERED));

        version
    }

    /// sets the supported versions to a given array
    fn set_supported_versions(
        &self,
        supported_versions: VersionVec<Self::Api>,
        dest: &mut VersionVec<Self::Api>,
    ) {
        require!(dest.is_empty(), "Versions already set");

        *dest = supported_versions;
    }

    /// returns true if the proposed version has a matching version
    fn is_supported_version(
        &self,
        supported_versions: &VersionVec<Self::Api>,
        version: &version::Data<Self::Api>,
    ) -> bool {
        let opt_found = self.find_supported_version(supported_versions, version);
        match opt_found {
            Some(found_version) => self.verify_proposed_version(&found_version, version),
            None => false,
        }
    }

    fn is_supported(
        &self,
        supported_versions: &VersionVec<Self::Api>,
        feature: &Feature<Self::Api>,
    ) -> bool {
        for sup_version in supported_versions {
            if self.verify_supported_feature(&sup_version, feature) {
                return true;
            }
        }

        false
    }

    /// Verifies that the entire feature set in the proposed version is supported by this chain.
    ///
    /// If the feature set is empty it verifies that this is allowed for the specified version identifier
    fn verify_proposed_version(
        &self,
        supported_version: &version::Data<Self::Api>,
        proposed_version: &version::Data<Self::Api>,
    ) -> bool {
        if supported_version.identifier != proposed_version.identifier {
            return false;
        }
        if proposed_version.features.is_empty() {
            return false;
        }

        for feature in &proposed_version.features {
            if !self.verify_supported_feature(supported_version, &feature) {
                return false;
            }
        }

        true
    }

    /// returns the version with a matching version identifier if it exists
    fn find_supported_version(
        &self,
        supported_versions: &VersionVec<Self::Api>,
        version: &version::Data<Self::Api>,
    ) -> Option<version::Data<Self::Api>> {
        supported_versions
            .into_iter()
            .find(|supp_version| supp_version.identifier == version.identifier)
    }

    /// Iterates over the descending ordered set of compatible IBC versions
    /// and selects the first version with a version identifier that is
    /// supported by the counterparty.
    ///
    /// The returned version contains a feature set with the intersection
    /// of the features supported by the source and counterparty chains.
    ///
    /// If the feature set intersection is nil and this is not allowed for the chosen version identifier
    /// then the search for a compatible version continues.
    ///
    /// This function is called in the ConnOpenTry handshake procedure.
    fn pick_version(
        &self,
        supported_versions: &VersionVec<Self::Api>,
        counterparty_versions: &VersionVec<Self::Api>,
    ) -> version::Data<Self::Api> {
        for supp_version in supported_versions {
            let opt_counterparty_version =
                self.find_supported_version(counterparty_versions, &supp_version);
            if opt_counterparty_version.is_none() {
                continue;
            }

            let counterparty_version = unsafe { opt_counterparty_version.unwrap_unchecked() };
            let feature_set = self.get_feature_set_intersection(
                &supp_version.features,
                &counterparty_version.features,
            );
            if !feature_set.is_empty() {
                return version::Data {
                    identifier: supp_version.identifier,
                    features: feature_set,
                };
            }
        }

        sc_panic!("No matching versions found")
    }

    /// takes in a version and feature string and returns true if the feature is supported by the version
    #[inline(always)]
    fn verify_supported_feature(
        &self,
        version: &version::Data<Self::Api>,
        feature: &Feature<Self::Api>,
    ) -> bool {
        version.features.contains(feature)
    }

    fn get_feature_set_intersection(
        &self,
        source_feature_set: &FeatureVec<Self::Api>,
        counterparty_feature_set: &FeatureVec<Self::Api>,
    ) -> FeatureVec<Self::Api> {
        let mut feature_set = FeatureVec::new();
        for src_feature in source_feature_set {
            if counterparty_feature_set.contains(&src_feature) {
                feature_set.push(src_feature);
            }
        }

        feature_set
    }
}
