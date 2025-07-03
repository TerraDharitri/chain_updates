# Description
Python Command Line Tools for interacting with Multivers<sup>X</sup>.

## Documentation
[docs.dharitri.org](https://docs.dharitri.org/sdk-and-tools/sdk-py/)

## CLI
[CLI](CLI.md)

## Distribution
[pipx](https://docs.dharitri.org/sdk-and-tools/sdk-py/installing-drtpy/) [(PyPi)](https://pypi.org/project/dharitri-sdk-cli/#history)

## Development setup

Clone this repository and cd into it:

```
git clone https://github.com/TerraDharitri/mx-sdk-py-cli.git
cd mx-sdk-py-cli
```

### Virtual environment

Create a virtual environment and install the dependencies:

```
python3 -m venv ./venv
source ./venv/bin/activate
pip install -r ./requirements.txt --upgrade
```

Install development dependencies, as well:

```
pip install -r ./requirements-dev.txt --upgrade
```

Allow `pre-commit` to automatically run on `git commit`:
```
pre-commit install
```

Above, `requirements.txt` should mirror the **dependencies** section of `pyproject.toml`.

If using VSCode, restart it or follow these steps:
 - `Ctrl + Shift + P`
 - _Select Interpreter_
 - Choose `./venv/bin/python`.

### Using your local `drtpy`

If you want to test the modifications you locally made to `drtpy`, set `PYTHONPATH` with the path to your local repository path.

For example, if you cloned the repository at `~/mx-sdk-py-cli`, run:

```
export PYTHONPATH="~/mx-sdk-py-cli"
```

Then `drtpy` will use the code in your local repository.
