# CHAST

*CHAST* short for *Change Stuff!*

[//]: # ([![CI]&#40;https://github.com/tj-actions/coverage-badge-go/workflows/CI/badge.svg&#41;]&#40;https://github.com/chast-io/chast-core/actions&#41;)
![Coverage](https://img.shields.io/badge/Coverage-40.7%25-yellow)
![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)

[![Ubuntu](https://img.shields.io/badge/Ubuntu%20(Tested)-E95420?logo=ubuntu\&logoColor=white)](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#jobsjob_idruns-on)
[![Ubuntu](https://img.shields.io/badge/Other%20Linux%20(Untested)-white?logo=linux\&logoColor=black)](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#jobsjob_idruns-on)
[![Mac OS](https://img.shields.io/badge/macOS%20(Planned)-000000?logo=apple\&logoColor=F0F0F0)](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#jobsjob_idruns-on)
[![Windows](https://img.shields.io/badge/Windows%20(Planned)-0078D6?logo=windows\&logoColor=white)](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#jobsjob_idruns-on)

This is the core of chast.
Run refactorings and other commands through a unified system no matter which operating system, installer or programming
language.

> This project is in its early stages of development so no further information is currently available.

## Required tools

- **General**
  - [unionfs-fuse](https://github.com/rpodgorny/unionfs-fuse) (Linux only, for Apple see MacOS support section in their README)
  - user namespace support required
  - (For OverlayFs-MergerFs-Isolation-Strategy: OverlayFs, Fuse, MergerFs required)
- **For development:**
  - [Go](https://golang.org/doc/install) (1.19.2 or higher)
  - [GolangCI-Lint](https://golangci-lint.run/usage/install/)
