# cludo Change Log

Changelog for cludo The Cloud Sudo toolset

## [Unreleased]

[0.0.2-alpha] - 2022-01-28
### Added
- Support for additional SSH key types [[#20]](https://github.com/superorbital/cludo/issues/20)
- Support for passing public keys to an active SSH agent for request signing [[#39]](https://github.com/superorbital/cludo/issues/39)
- Support for reading additional SSH public keys from Github if the user has a `github_id` defined in the server-side config [[#47]](https://github.com/superorbital/cludo/issues/47)
- Server side config can now contain a user's name [[#53]](https://github.com/superorbital/cludo/issues/53)
- Github Actions for release pipeline and basic code scanning [[#60]](https://github.com/superorbital/cludo/issues/60)
- and almost certainly more small tweaks, etc...
### Changed
- Made log message format more consistent
- Converted SSH key path to a list and move to the top level of the user config [[#52]](https://github.com/superorbital/cludo/issues/52)
### Removed
- Removed the ability to pass in SSH passphrases [[#89]](https://github.com/superorbital/cludo/issues/89)

## [0.0.1-alpha] - 2022-01-04
### Added
- Initial Proof of Concept Release

[Unreleased]: https://github.com/coditory/changelog-parser-action/compare/0.0.2-alpha...HEAD
[0.0.2-alpha]: https://github.com/coditory/changelog-parser-action/releases/tag/0.0.2-alpha
[0.0.1-alpha]: https://github.com/coditory/changelog-parser-action/releases/tag/0.0.1-alpha
