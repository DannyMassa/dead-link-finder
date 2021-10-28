<h1 align="center">
  <a href="https://github.com/DannyMassa/dead-link-linter">
    <img src="docs/images/logo.png" alt="Logo" width="320" height="320">
  </a>
</h1>

<div align="center">
  Dead Link Linter
  <br />
  <a href="https://github.com/DannyMassa/dead-link-linter/issues/new?assignees=&labels=bug&template=01_BUG_REPORT.md&title=bug%3A+">Report a Bug</a>
  ·
  <a href="https://github.com/DannyMassa/dead-link-linter/issues/new?assignees=&labels=enhancement&template=02_FEATURE_REQUEST.md&title=feat%3A+">Request a Feature</a>
</div>
<div align="center">
<br />

[![license](https://img.shields.io/github/license/DannyMassa/dead-link-linter.svg?style=flat-square)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/DannyMassa/dead-link-linter)](https://goreportcard.com/report/github.com/DannyMassa/dead-link-linter)
[![build](https://img.shields.io/github/workflow/status/DannyMassa/dead-link-linter/Continuous%20Integration%20Gate)](https://img.shields.io/github/workflow/status/DannyMassa/dead-link-linter/Continuous%20Integration%20Gate)
[![release](https://img.shields.io/github/v/release/DannyMassa/dead-link-linter)](https://img.shields.io/github/v/release/DannyMassa/dead-link-linter)
</div>

<details open="open">
<summary>Table of Contents</summary>

- [Synopsis](#synopsis)
- [TL;DR](#TL;DR)
- [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
      - [Single Use Script (Requires cURL)](#single-use-script-requires-curl)
      - [Portable Executable](#portable-executable)
- [Usage](#usage)
  - [Command Line Interface](#command-line-interface)
    - [Examples](#examples)
  - [.deadlink](#deadlink)
    - [Example .deadlink File](#example-deadlink-file)
  - [Configuration](#configuration)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [Authors & contributors](#authors--contributors)
- [Security](#security)
- [License](#license)

</details>

---

## Synopsis

The Dead Link Linter is a small executable that will crawl through your directory(ies) to ensure that all external links embedded in your documentation work as intended. This can be used to ensure embedded media works as intended, or can ensure that embedded hyperlinks land where they are supposed to.

#### Why?

Software projects typically have additional software projects that they are built upon, or utilize in some way. This usually leads to referencing documentation from other projects in your projects. Since we often aren't maintainers of the projects referenced, the Dead Link Linter can be added as a Continous Integration gate in your CI/CD system to ensure these external references to not go stale, and if they are stale, they can be identified and updated quickly.

## TL;DR

The following will download, run, report, and remove the Dead Link Linter, leaving you with the results

```shell
# For (Debian) Linux
curl https://raw.githubusercontent.com/DannyMassa/dead-link-linter/master/scripts/linux/dead-link-linter-single-run.sh | bash

# For Windows
curl https://raw.githubusercontent.com/DannyMassa/dead-link-linter/master/scripts/windows/dead-link-linter-single-run.sh | bash
```

## Getting Started

### Prerequisites

**None!** The Dead Link Linter runs as a standalone executable, but cURL or a similar alternative is required if running the application using the command provided in the [TL;DR](#TL;DR)

### Installation

There are two simple ways to install/use the Dead Link Linter: As a single use script or as a portable executable located in the directory you are executing from.

#### Single Use Script (Requires cURL)

Briefly referenced in the [TL;DR](#TL;DR), this is the simpliest way to try the Dead Link Linter. This script was created to allow users to quickly download the executable, run the service, and remove all related files. The limitation being that since the linter is not executed directly by the user, CLI overrides will not be available. ```.deadlink``` files are utilized if they reside in the directory the script is called from.

TODO INSERT GIF

#### Portable Executable

From the [Releases](https://github.com/DannyMassa/dead-link-linter/releases) page, download the artifact that matches your machine. It is always recommended that you download the latest release. Unpack the resulting compressed artifact, and move the executable from the extracted directory to the directory you intend on executing it from. From your shell, navigate to that directory and utilize the [CLI commands](#Command-Line-Interface) to execute 

TODO INSERT GIF

## Usage

The Dead Link Linter offers a usable set of default values but is highly configurable. Configuration using CLI flags takes priority over a ```.deadlink``` file, and a ```.deadlink``` file will take priority over default values.

### Command Line Interface 

The Dead Link Linter does not require any parameters to be supplied via CLI in order to run, but supplying command line arguements will overwrite any configuration settings provided by a ```.deadlink``` file, and will also override any defaults. CLI command syntax follows POSIX/GNU-style flags.

#### Examples

```bash
# Run with defaults
$ ./dead-link-linter

# Allow up to 10 links to be dead, and only scan files with an '.md' extension
$ ./dead-link-linter --max_failures=10 --file_extensions=[".md"]

# Use an alternative golden URL if google.com doesn't suit your needs
$ ./dead-link-linter --golden_url="https://bing.com"

# Ignore search engine links
$ ./dead-link-linter --ignored=["https://bing.com", "https://google.com"]

# Only print failure logs
$ ./dead-link-linter --log_verbosity=0
```

### .deadlink

Automatically, the Dead Link Linter looks for a file named ```.deadlink``` in the directory the linter is being executed from. ```.deadlink``` should be a valid YAML file. If the file is found, the configuration will be loaded and parsed. If the YAML is not successfully parsed, the application will exit with an unhealthy status code before any links are checked.

#### Example .deadlink File

```yaml
directories:
  - "./"
fileExtensions:
  - ".md"
goldenURL: "https://google.com"
ignored:
  - "https://google.com"
maxFailures: 0
logVerbosity: 1

# In this case, individualTimeout is omitted, therefore the timeout will remain the default value of 10 seconds, unless the CLI flag (--individual_timeout) is set to another value
```

### Configuration

YAML Key               | CLI                  |Default             | [Domain] Description 
------            | ------               |------              |----------
directories       | --directories        | ["./"]             | A list of (string) directories to run the linter against relative and absolute paths. Overlapping files will be checked in each directory section (multiple times) and all subdirectories will be recursively included. 
fileExtensions    | --file_extensions    | [".markdown", ".mdown", ".mkdn", ".md", ".mkd", ".mdwn", ".mdtxt", ".mdtext", ".text", ".txt", ".rst"]     | A list of (string) file extensions that indicates what file types will be scanned during the linter's run. All strings must begin with a ```.``` 
goldenURL         | --golden_url         | https://google.com | A url that will be tested first, and upon failure will cause an unhealthy exit from the program (exit code ```5```). This is a 'fail fast' mechanism. 
ignored           | --ignored            | []                 |A list of (string) URLs that will be ignored if they are found by the URL parser. This is the suggested way to manage false positives. i.e. example URLs & example IP Addresses from your documentation.
individualTimeout | --individual_timeout | 10                 | **[1, 2147483647)** The amount of time, in seconds, allowed for a web request 
maxFailures       | --max_failures       | 0                  | **[0, 2147483647)** The number of URLs allowed to have a ```FAILURE``` status before the linter exits with an unhealthy code. (exit code ```1```) 
logVerbosity      | --log_verbosity      | 1                  | **[0, 2]** a numerical value describing if detailed logs are printed. ```0 ``` indicates failure only logs, ```1``` indicates skipped urls will be included, ```2``` indicates success urls will be included

## Roadmap

See the [open issues](https://github.com/DannyMassa/dead-link-linter/issues) for a list of proposed features (and known issues).

- [Top Feature Requests](https://github.com/DannyMassa/dead-link-linter/issues?q=label%3Aenhancement+is%3Aopen+sort%3Areactions-%2B1-desc) (Add your votes using the 👍 reaction)
- [Top Bugs](https://github.com/DannyMassa/dead-link-linter/issues?q=is%3Aissue+is%3Aopen+label%3Abug+sort%3Areactions-%2B1-desc) (Add your votes using the 👍 reaction)

## Contributing

First off, thanks for taking the time to contribute! Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make will benefit everybody else and are **greatly appreciated**.

We have set up a separate document containing our [contribution guidelines](docs/CONTRIBUTING.md).

Thank you for being involved!

## Authors & contributors

The original setup of this repository is by [Danny Massa](https://github.com/DannyMassa).

For a full list of all authors and contributors, check [the contributor's page](https://github.com/DannyMassa/dead-link-linter/contributors).

## Security

Dead Link Linter follows good practices of security, but 100% security can't be granted in software.
Dead Link Linter is provided **"as is"** without any **warranty**. Use at your own risk.

## License

This project is licensed under the **MIT license**.

See [LICENSE](LICENSE) for more information.
