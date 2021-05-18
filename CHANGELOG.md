<a name="unreleased"></a>
## [Unreleased]


<a name="v1.0.0"></a>
## v1.0.0 - 2021-05-18

### Features
- **cli:** add short version option
- **cli:** allow config file without extension forcing yaml type
- **cli:** display version and keep build info data
- **cli:** add cli
- **merge:** add merge functionality

### Code Refactoring
- **cli:** check errors when binding viper with cobra
- **cli:** create error with recommended function

### Build changes
- **changelog:** generate changelog from git log
- **lint:** lint source code
- **release:** generate rpm, deb and apk packages for linux
- **release:** use goreleaser to build release
- **release-notes:** generate release-notes for current version from git log
- **version:** generate next semantic version based on git log

### Documentation changes
- add deb, rpm and apk installation info
- add development info
- add usage info
- **help:** capitalize all short command messages


[Unreleased]: https://github.com/amplia-iiot/yutil/compare/v1.0.0...HEAD
