import logging
import shutil
import urllib.request
from pathlib import Path

from dharitri_sdk_cli import dependencies
from dharitri_sdk_cli.localnet.config_root import ConfigRoot
from dharitri_sdk_cli.localnet.config_software import (
    SoftwareComponent,
    SoftwareResolution,
)

logger = logging.getLogger("localnet")


def fetch_prerequisites(configfile: Path):
    logger.info("fetch_prerequisites()")

    config = ConfigRoot.from_file(configfile)

    dependencies.install_module("testwallets", overwrite=True)

    if config.software.drt_go_chain.resolution == SoftwareResolution.Remote:
        download_software_component(config.software.drt_go_chain)

    if config.software.drt_go_chain_proxy.resolution == SoftwareResolution.Remote:
        download_software_component(config.software.drt_go_chain_proxy)

    config.software.drt_go_chain.node_config_must_exist()
    config.software.drt_go_chain.seednode_config_must_exist()
    config.software.drt_go_chain_proxy.proxy_config_must_exist()

    is_node_built = config.software.drt_go_chain.is_node_built()
    is_seednode_built = config.software.drt_go_chain.is_seednode_built()
    is_proxy_built = config.software.drt_go_chain_proxy.is_proxy_built()

    is_golang_needed = not (is_node_built and is_seednode_built and is_proxy_built)
    if is_golang_needed:
        dependencies.install_module("golang")


def download_software_component(component: SoftwareComponent):
    download_folder = component.get_archive_download_folder()
    extraction_folder = component.get_archive_extraction_folder()
    url = component.archive_url

    shutil.rmtree(str(download_folder), ignore_errors=True)
    shutil.rmtree(str(extraction_folder), ignore_errors=True)

    download_folder.mkdir(parents=True, exist_ok=True)
    extraction_folder.mkdir(parents=True, exist_ok=True)
    archive_extension = url.split(".")[-1]
    download_path = download_folder / f"archive.{archive_extension}"

    logger.info(f"Downloading archive {url} to {download_path}")
    urllib.request.urlretrieve(url, download_path)

    logger.info(f"Unpacking archive {download_path} to {extraction_folder}")
    shutil.unpack_archive(download_path, extraction_folder, format="zip")
