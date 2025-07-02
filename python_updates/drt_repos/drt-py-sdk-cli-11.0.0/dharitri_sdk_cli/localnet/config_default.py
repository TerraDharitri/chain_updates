from pathlib import Path

from dharitri_sdk_cli.localnet.config_general import General
from dharitri_sdk_cli.localnet.config_networking import Networking
from dharitri_sdk_cli.localnet.config_sharding import Metashard, RegularShards
from dharitri_sdk_cli.localnet.config_software import (
    Software,
    SoftwareChainGo,
    SoftwareChainProxyGo,
    SoftwareResolution,
)

general = General(
    log_level="*:DEBUG",
    genesis_delay_seconds=10,
    rounds_per_epoch=100,
    round_duration_milliseconds=6000,
)

software = Software(
    drt_go_chain=SoftwareChainGo(
        resolution=SoftwareResolution.Remote,
        archive_url="https://github.com/TerraDharitri/drt-go-chain/archive/refs/heads/master.zip",
        archive_download_folder=Path("~/dharitri-sdk") / "localnet_software_remote" / "downloaded" / "drt-go-chain",
        archive_extraction_folder=Path("~/dharitri-sdk") / "localnet_software_remote" / "extracted" / "drt-go-chain",
        local_path=Path("~/dharitri-sdk") / "localnet_software_local" / "drt-go-chain",
    ),
    drt_go_chain_proxy=SoftwareChainProxyGo(
        resolution=SoftwareResolution.Remote,
        archive_url="https://github.com/TerraDharitri/drt-go-chain-proxy/archive/refs/heads/master.zip",
        archive_download_folder=Path("~/dharitri-sdk")
        / "localnet_software_remote"
        / "downloaded"
        / "drt-go-chain-proxy",
        archive_extraction_folder=Path("~/dharitri-sdk")
        / "localnet_software_remote"
        / "extracted"
        / "drt-go-chain-proxy",
        local_path=Path("~/dharitri-sdk") / "localnet_software_local" / "drt-go-chain-proxy",
    ),
)

metashard = Metashard(
    consensus_size=1,
    num_observers=0,
    num_validators=1,
)

shards = RegularShards(
    num_shards=2,
    consensus_size=1,
    num_observers_per_shard=0,
    num_validators_per_shard=1,
)

networking = Networking(
    host="127.0.0.1",
    port_seednode=9999,
    port_seednode_rest_api=10000,
    p2p_id_seednode="16Uiu2HAkx4QqgXXDdHdUWbLu5kxhd3Uo2hqB2FfCxmxH5Sd7bZFk",
    port_proxy=7950,
    port_first_observer=21100,
    port_first_observer_rest_api=10100,
    port_first_validator=21500,
    port_first_validator_rest_api=10200,
)
