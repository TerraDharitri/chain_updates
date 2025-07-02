import os
import sys
from typing import Any, Dict

from dharitri_sdk_cli.localnet.config_root import ConfigRoot
from dharitri_sdk_cli.localnet.config_software import SoftwareResolution

sys.path = [os.getcwd() + "/.."] + sys.path


def test_override_config() -> None:
    config = ConfigRoot()

    # Check a few default values
    assert config.general.rounds_per_epoch == 100
    assert config.general.round_duration_milliseconds == 6000
    assert config.metashard.consensus_size == 1
    assert config.networking.port_proxy == 7950
    assert config.software.drt_go_chain.resolution == SoftwareResolution.Remote
    assert (
        config.software.drt_go_chain.archive_url
        == "https://github.com/TerraDharitri/drt-go-chain/archive/refs/heads/master.zip"
    )

    # Now partly override the config
    config_patch: Dict[str, Any] = dict()
    config_patch["general"] = {
        "rounds_per_epoch": 200,
        "round_duration_milliseconds": 4000,
    }
    config_patch["metashard"] = {
        "consensus_size": 2,
    }
    config_patch["networking"] = {
        "port_proxy": 7951,
    }
    config_patch["software"] = {
        "drt_go_chain": {"archive_url": "https://github.com/TerraDharitri/drt-go-chain/archive/refs/tags/v1.5.1.zip"}
    }

    config.override(config_patch)

    # Check the overridden values
    assert config.general.rounds_per_epoch == 200
    assert config.general.round_duration_milliseconds == 4000
    assert config.metashard.consensus_size == 2
    assert config.networking.port_proxy == 7951
    assert config.software.drt_go_chain.resolution == SoftwareResolution.Remote
    assert (
        config.software.drt_go_chain.archive_url
        == "https://github.com/TerraDharitri/drt-go-chain/archive/refs/tags/v1.5.1.zip"
    )
