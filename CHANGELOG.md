<a name="unreleased"></a>
## [Unreleased]


<a name="v2.0.0"></a>
## [v2.0.0] - 2021-07-21

### Bug Fixes
- **cli:** don't export unwanted methods and types
- **merge:** fix typo in error message
- **merge:** return ignored errors

### Features
- **format:** continue formatting files on error
- **format:** add format functionality

### Code Refactoring
- convert internal functions to variables to allow mocking
- **merge:** wrap options to avoid clashes with other cmds

### Build changes
- add option to include coverage in check and watch tasks
- **docs:** add godoc to display documentation
- **release:** add dispatch workflow option

### Documentation changes
- update info for release workflow
- remove license header from docs
- add yutil code documentation (description)
- **cli:** add code documentation
- **format:** add usage information
- **format:** add code documentation
- **merge:** add code documentation
- **testing:** document internal testing package

### BREAKING CHANGE

Some parts of cmd package are no longer exported to be used, as it was never
intended to leave them open.


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


[Unreleased]: https://github.com/amplia-iiot/yutil/compare/v2.0.0...HEAD
[v2.0.0]: https://github.com/amplia-iiot/yutil/compare/v1.0.0...v2.0.0
