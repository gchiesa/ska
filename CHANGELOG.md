# Changelog

## [0.1.0](https://github.com/gchiesa/ska/compare/v0.0.24...v0.1.0) (2025-08-26)


### ⚠ BREAKING CHANGES

* bumping to v1.x release model

### Features

* bumping to v1.x release model ([439cd60](https://github.com/gchiesa/ska/commit/439cd60e3c16828f265ff6c9e1c85d79c3c9973a))


### Docs

* update pkg docs ([2240bfb](https://github.com/gchiesa/ska/commit/2240bfba091ce169b2e7d80ec0895f44b892d9a8))
* update README.md and add use-cases ([9c9f8c2](https://github.com/gchiesa/ska/commit/9c9f8c2f3e44f2bd896d8c312d13bdf1938fb0a4))

## [0.0.24](https://github.com/gchiesa/ska/compare/v0.0.23...v0.0.24) (2025-07-06)


### Bug Fixes

* logging issue and missing path default for config ([d83da86](https://github.com/gchiesa/ska/commit/d83da862287ad0310c0cd6fdc7af1b49272ef7a0))

## [0.0.23](https://github.com/gchiesa/ska/compare/v0.0.22...v0.0.23) (2025-07-06)


### Bug Fixes

* adding ignores ([785621e](https://github.com/gchiesa/ska/commit/785621e017705b3a29d5d71dc85e5715c4d12b12))
* removing clutter ([9c5a15b](https://github.com/gchiesa/ska/commit/9c5a15b3de1f9a728f7377da7780d2e4a310beb1))
* update test to let goproxy able to zip the pkg ([f21709f](https://github.com/gchiesa/ska/commit/f21709f2ed314f6dc7ef01b2e1a3adb57b538614))

## [0.0.22](https://github.com/gchiesa/ska/compare/v0.0.21...v0.0.22) (2025-07-05)


### Bug Fixes

* add goignore for skipping invalid files while consuming pkg ([f834981](https://github.com/gchiesa/ska/commit/f834981d26e2b7b9443571a2f49047946781881e))
* wrong path ([297a822](https://github.com/gchiesa/ska/commit/297a822b2985a15e9732b05d2822ea6eee42c0a0))

## [0.0.21](https://github.com/gchiesa/ska/compare/v0.0.20...v0.0.21) (2025-07-05)


### Bug Fixes

* downgrade bubbles pkg for fixing UI issue ([9b19f1f](https://github.com/gchiesa/ska/commit/9b19f1fcc691c028c8957a9f85c1d24d0a8cd53a))
* linter issue ([68efac5](https://github.com/gchiesa/ska/commit/68efac5ac2093097033f9114302dd5c4f809c602))

## [0.0.20](https://github.com/gchiesa/ska/compare/v0.0.19...v0.0.20) (2025-07-05)


### Features

* banner optional when used as lib and improved ui ([9d2bf39](https://github.com/gchiesa/ska/commit/9d2bf397fb7eae3e7a17ce0d0a3e62031cd8b7aa))
* update gitlab client and add some test for gitignore ([6388313](https://github.com/gchiesa/ska/commit/6388313bf93622f519ed15c51bfbd91fde79bf4c))


### Bug Fixes

* **security:** update goviper ([3ee82b3](https://github.com/gchiesa/ska/commit/3ee82b38f12f17e5272bbf35c0c619196893024c))

## [0.0.19](https://github.com/gchiesa/ska/compare/v0.0.18...v0.0.19) (2025-05-11)


### Features

* fix linting ([daaebf3](https://github.com/gchiesa/ska/commit/daaebf3c5dc79345a0e1a4756b6de5588d205165))
* new sprig function optional + additional tests ([ad7c163](https://github.com/gchiesa/ska/commit/ad7c163c9e292f7d5d9df0174b29fb83cb7d51b7))


### Bug Fixes

* security vuln update ([1fce751](https://github.com/gchiesa/ska/commit/1fce751054878b654223b800cff060652f83e66e))

## [0.0.18](https://github.com/gchiesa/ska/compare/v0.0.17...v0.0.18) (2024-11-29)


### Features

* add support for gitlab public/private blueprints ([d3b364f](https://github.com/gchiesa/ska/commit/d3b364ff50815b02a30945abfda6372690ff704c))
* add support for jinja2 like templates ([e82f9f7](https://github.com/gchiesa/ska/commit/e82f9f7757422d7f1807bab9914bc7dc11383a8a))
* extracted to public pkg templateprovider ([d54ba43](https://github.com/gchiesa/ska/commit/d54ba43228450e9ea2d6985db573236f00423efd))
* implement config rename ([0e67d51](https://github.com/gchiesa/ska/commit/0e67d5136dff66e851f4e1118b60adde3ad86ac1))
* implement delete command ([f1cd067](https://github.com/gchiesa/ska/commit/f1cd067a560dc2ccd555c46350683de8bf8d2cda))
* implement ignorepaths ([e00d711](https://github.com/gchiesa/ska/commit/e00d7117411743b80c0e54bd9ae706dc81451375))
* implement json output, non interactive mode and better arguments for CLI ([e2b64eb](https://github.com/gchiesa/ska/commit/e2b64eb2fdadc9dd720c5a9216d38f39d1204a1c))
* implement minLength support for accepting empty variables ([4f74c95](https://github.com/gchiesa/ska/commit/4f74c95f97f2ab90bb302a7b61738570e5f12a91))
* implement support for automatically add ignorePaths in generated ska-config ([500768c](https://github.com/gchiesa/ska/commit/500768c995813ad8428a3c2cab1d28f2675e9f92))
* implement support for path inside remote repository ([f4f3edf](https://github.com/gchiesa/ska/commit/f4f3edff764c47b032420f72a10dbbf019a97d16))
* implement update template capability ([8523155](https://github.com/gchiesa/ska/commit/8523155bacce729aa56c0d5fc60314a8f49367a0))
* implemented create command ([5911e29](https://github.com/gchiesa/ska/commit/5911e29215fec9ac64411048e10ea1a67fd8f5c8))
* introduce writeOnce fields ([b7cb8e0](https://github.com/gchiesa/ska/commit/b7cb8e001f15fa88e6ce34a12dc9d9f0d47f1635))
* lint issues ([845bd3c](https://github.com/gchiesa/ska/commit/845bd3c60697609b4cfc29155cce75ab3d9892aa))
* minor updates and extended readme with demo ([807c318](https://github.com/gchiesa/ska/commit/807c318a6e1d6730c6d539afbe0b48a712a17004))
* refactoring to implement config subcommands - implemented config list ([819446b](https://github.com/gchiesa/ska/commit/819446b48aa3c3e08bafde1a48e31d10312cf3f6))


### Bug Fixes

* add missing secret ([1235ff2](https://github.com/gchiesa/ska/commit/1235ff296936534285f89f1a98790e01e739fb15))
* configure homebrew integration ([7db05b1](https://github.com/gchiesa/ska/commit/7db05b1e35ecf3d799b6fa05fbd115f25aa0aa40))
* goreleaser ([d7da3d9](https://github.com/gchiesa/ska/commit/d7da3d98a2c80134ae86a53ef1b8ed8fbae9b020))
* goreleaser action ([69ece7f](https://github.com/gchiesa/ska/commit/69ece7feddb1def0e6fc27cb3c8ed0db7aabe3cd))
* implement support for multiple local ska-config.yaml ([af43a23](https://github.com/gchiesa/ska/commit/af43a234ffcb6213446da1f0297e0d6456fa2e2a))
* lint issues ([10f6cad](https://github.com/gchiesa/ska/commit/10f6cad9c0d6c15bf81294df8dc44dbe8a139c52))
* lint issues ([0f94080](https://github.com/gchiesa/ska/commit/0f94080b951b5095f5b14b45797d96429bfc8955))
* lint issues ([d939d5c](https://github.com/gchiesa/ska/commit/d939d5c4d36f86ddf25b4fbd82ea1645e172129d))
* lint issues ([266fd9a](https://github.com/gchiesa/ska/commit/266fd9af7cf40986d7eb5025fe03638eaf6f6e45))
* lint issues ([a68751a](https://github.com/gchiesa/ska/commit/a68751a2996a710df5850d6bbe76f6afb00f5a6c))
* lint passing ([1b958b6](https://github.com/gchiesa/ska/commit/1b958b6e24f682f22a656f092c0c86d0c6df095d))
* missing release please manifest ([06a9aa7](https://github.com/gchiesa/ska/commit/06a9aa7d617dba30099e59ea49534df3934233dc))
* removed unused variable ([8aff9e1](https://github.com/gchiesa/ska/commit/8aff9e118d8f7378a72c3527983c749a5bb27472))
* template error issues and better reporting. ([51eb60c](https://github.com/gchiesa/ska/commit/51eb60c95a0f4cfbd601d398ac94b17f36d134a2))
* update url for gitlab test repository ([d136e00](https://github.com/gchiesa/ska/commit/d136e00768de7b1ce2f2a3e1b77d41dece6d77a7))


### Docs

* smaller terminal demo ([921ef7f](https://github.com/gchiesa/ska/commit/921ef7fbcca152c0712c71bc982b3a1a7c14761f))
* update demo ([4606065](https://github.com/gchiesa/ska/commit/4606065a35ecd3462c2f7989dd566552e5d325d3))
* update README ([69af4bd](https://github.com/gchiesa/ska/commit/69af4bd700af311e26808875dbe93f2f79a639e4))
* update README ([84126a7](https://github.com/gchiesa/ska/commit/84126a7e5fa7b87227bdeb774ec89a2d924ec72e))
* update README ([72cd21b](https://github.com/gchiesa/ska/commit/72cd21b646c10238be0e79521a20f0e7eda8decb))
* update README ([25bc2e1](https://github.com/gchiesa/ska/commit/25bc2e10b5541300c7046e55ebbf44a66594ba90))


### Other

* fix goland version ([9619d92](https://github.com/gchiesa/ska/commit/9619d921ad045e0278e53d72bb271fec3b30b0d4))
* fix golang-ci ([0d60139](https://github.com/gchiesa/ska/commit/0d601394f418c2285cd5acabb7742fd04d730dab))
* fix pipelines ([fca2931](https://github.com/gchiesa/ska/commit/fca29314f2a9f6fb28f42716d254c085fa7f99a3))
* fix workflows and config ([8dff677](https://github.com/gchiesa/ska/commit/8dff6770235a69f6fcda5a0a9811cedaaf0473ac))
* **main:** release 0.0.10 ([ae2e336](https://github.com/gchiesa/ska/commit/ae2e3366252dd3c86850d6ca88590ed7a4b2880b))
* **main:** release 0.0.11 ([bbaab65](https://github.com/gchiesa/ska/commit/bbaab65264539df326e325b56b29dcd58008f5c3))
* **main:** release 0.0.12 ([860b8d7](https://github.com/gchiesa/ska/commit/860b8d7ccc47fa00c113a84b6475f148c729dab4))
* **main:** release 0.0.13 ([d17d2d5](https://github.com/gchiesa/ska/commit/d17d2d551edcb61e71855f9fe628d1ca809a7b70))
* **main:** release 0.0.14 ([61bf00c](https://github.com/gchiesa/ska/commit/61bf00c88594390aa96bafca54be1454b81cac95))
* **main:** release 0.0.15 ([d8e48b5](https://github.com/gchiesa/ska/commit/d8e48b5cef3e44b2c8ea31ede2258b6e3e20abf8))
* **main:** release 0.0.16 ([c745f4e](https://github.com/gchiesa/ska/commit/c745f4e3db39acf6ffbc7abbd9ce1d0b4a394150))
* **main:** release 0.0.17 ([c438045](https://github.com/gchiesa/ska/commit/c438045d97982911e74fb02987a98fe4f3eeb8d8))
* **main:** release 0.0.2 ([1e59766](https://github.com/gchiesa/ska/commit/1e5976691ad5109dc91ac090dd872db49e8403ec))
* **main:** release 0.0.3 ([bda8652](https://github.com/gchiesa/ska/commit/bda8652c4733ef65ec986d7af0b4c34367e9a505))
* **main:** release 0.0.4 ([971b575](https://github.com/gchiesa/ska/commit/971b575c0b38f47e8dc988c49653d609c10d3418))
* **main:** release 0.0.5 ([3d5a571](https://github.com/gchiesa/ska/commit/3d5a57156a42c267dd7cb92d1e375501066aa0be))
* **main:** release 0.0.6 ([8df70f5](https://github.com/gchiesa/ska/commit/8df70f53297edfff03d3da924396139896bc7735))
* **main:** release 0.0.7 ([1c82c76](https://github.com/gchiesa/ska/commit/1c82c7696d0e7833c79bc907b2a9d94ac582941d))
* **main:** release 0.0.8 ([3134e48](https://github.com/gchiesa/ska/commit/3134e48dee82c4de1d6b1381afa27f31f67a1e47))
* **main:** release 0.0.9 ([e65e53c](https://github.com/gchiesa/ska/commit/e65e53cdccce00138b19c52d8ab2201cadea65f7))
* remove idea folder ([5e80c81](https://github.com/gchiesa/ska/commit/5e80c8147733995a7292f23cd9f49100cc145166))
* update configuration ([9bd68b7](https://github.com/gchiesa/ska/commit/9bd68b728003bfc20d68542803ae9b753307875e))
* update dependencies ([d9f2cb6](https://github.com/gchiesa/ska/commit/d9f2cb6d47fb28c601884bc26608f982ad871d17))
* update dependencies ([f7164cc](https://github.com/gchiesa/ska/commit/f7164ccc074372d883d16292958f10f85d602ec2))
* update linter config ([2556554](https://github.com/gchiesa/ska/commit/255655456bf355560a4b177569d03b631528bab9))
* update linter gh action ([610fc25](https://github.com/gchiesa/ska/commit/610fc25fff7a34d29af209cefb967939a4edfd69))

## [0.0.17](https://github.com/gchiesa/ska/compare/v0.0.16...v0.0.17) (2024-11-29)


### Features

* extracted to public pkg templateprovider ([d54ba43](https://github.com/gchiesa/ska/commit/d54ba43228450e9ea2d6985db573236f00423efd))
* introduce writeOnce fields ([b7cb8e0](https://github.com/gchiesa/ska/commit/b7cb8e001f15fa88e6ce34a12dc9d9f0d47f1635))


### Other

* update dependencies ([d9f2cb6](https://github.com/gchiesa/ska/commit/d9f2cb6d47fb28c601884bc26608f982ad871d17))

## [0.0.16](https://github.com/gchiesa/ska/compare/v0.0.15...v0.0.16) (2024-10-22)


### Features

* implement config rename ([0e67d51](https://github.com/gchiesa/ska/commit/0e67d5136dff66e851f4e1118b60adde3ad86ac1))
* implement delete command ([f1cd067](https://github.com/gchiesa/ska/commit/f1cd067a560dc2ccd555c46350683de8bf8d2cda))
* refactoring to implement config subcommands - implemented config list ([819446b](https://github.com/gchiesa/ska/commit/819446b48aa3c3e08bafde1a48e31d10312cf3f6))


### Bug Fixes

* lint issues ([10f6cad](https://github.com/gchiesa/ska/commit/10f6cad9c0d6c15bf81294df8dc44dbe8a139c52))


### Other

* update configuration ([9bd68b7](https://github.com/gchiesa/ska/commit/9bd68b728003bfc20d68542803ae9b753307875e))
* update linter config ([2556554](https://github.com/gchiesa/ska/commit/255655456bf355560a4b177569d03b631528bab9))
* update linter gh action ([610fc25](https://github.com/gchiesa/ska/commit/610fc25fff7a34d29af209cefb967939a4edfd69))

## [0.0.15](https://github.com/gchiesa/ska/compare/v0.0.14...v0.0.15) (2024-09-28)


### Features

* implement support for automatically add ignorePaths in generated ska-config ([500768c](https://github.com/gchiesa/ska/commit/500768c995813ad8428a3c2cab1d28f2675e9f92))


### Bug Fixes

* lint issues ([0f94080](https://github.com/gchiesa/ska/commit/0f94080b951b5095f5b14b45797d96429bfc8955))
* lint issues ([d939d5c](https://github.com/gchiesa/ska/commit/d939d5c4d36f86ddf25b4fbd82ea1645e172129d))

## [0.0.14](https://github.com/gchiesa/ska/compare/v0.0.13...v0.0.14) (2024-09-25)


### Bug Fixes

* implement support for multiple local ska-config.yaml ([af43a23](https://github.com/gchiesa/ska/commit/af43a234ffcb6213446da1f0297e0d6456fa2e2a))
* removed unused variable ([8aff9e1](https://github.com/gchiesa/ska/commit/8aff9e118d8f7378a72c3527983c749a5bb27472))

## [0.0.13](https://github.com/gchiesa/ska/compare/v0.0.12...v0.0.13) (2024-09-23)


### Features

* implement minLength support for accepting empty variables ([4f74c95](https://github.com/gchiesa/ska/commit/4f74c95f97f2ab90bb302a7b61738570e5f12a91))

## [0.0.12](https://github.com/gchiesa/ska/compare/v0.0.11...v0.0.12) (2024-08-30)


### Features

* implement support for path inside remote repository ([f4f3edf](https://github.com/gchiesa/ska/commit/f4f3edff764c47b032420f72a10dbbf019a97d16))


### Docs

* smaller terminal demo ([921ef7f](https://github.com/gchiesa/ska/commit/921ef7fbcca152c0712c71bc982b3a1a7c14761f))
* update demo ([4606065](https://github.com/gchiesa/ska/commit/4606065a35ecd3462c2f7989dd566552e5d325d3))


### Other

* update dependencies ([f7164cc](https://github.com/gchiesa/ska/commit/f7164ccc074372d883d16292958f10f85d602ec2))

## [0.0.11](https://github.com/gchiesa/ska/compare/v0.0.10...v0.0.11) (2024-08-16)


### Features

* add support for gitlab public/private blueprints ([d3b364f](https://github.com/gchiesa/ska/commit/d3b364ff50815b02a30945abfda6372690ff704c))


### Bug Fixes

* update url for gitlab test repository ([d136e00](https://github.com/gchiesa/ska/commit/d136e00768de7b1ce2f2a3e1b77d41dece6d77a7))


### Docs

* update README ([69af4bd](https://github.com/gchiesa/ska/commit/69af4bd700af311e26808875dbe93f2f79a639e4))

## [0.0.10](https://github.com/gchiesa/ska/compare/v0.0.9...v0.0.10) (2024-08-09)


### Features

* minor updates and extended readme with demo ([807c318](https://github.com/gchiesa/ska/commit/807c318a6e1d6730c6d539afbe0b48a712a17004))

## [0.0.9](https://github.com/gchiesa/ska/compare/v0.0.8...v0.0.9) (2024-08-09)


### Features

* add support for jinja2 like templates ([e82f9f7](https://github.com/gchiesa/ska/commit/e82f9f7757422d7f1807bab9914bc7dc11383a8a))
* lint issues ([845bd3c](https://github.com/gchiesa/ska/commit/845bd3c60697609b4cfc29155cce75ab3d9892aa))

## [0.0.8](https://github.com/gchiesa/ska/compare/v0.0.7...v0.0.8) (2024-08-03)


### Bug Fixes

* template error issues and better reporting. ([51eb60c](https://github.com/gchiesa/ska/commit/51eb60c95a0f4cfbd601d398ac94b17f36d134a2))

## [0.0.7](https://github.com/gchiesa/ska/compare/v0.0.6...v0.0.7) (2024-08-02)


### Features

* implement ignorepaths ([e00d711](https://github.com/gchiesa/ska/commit/e00d7117411743b80c0e54bd9ae706dc81451375))
* implement json output, non interactive mode and better arguments for CLI ([e2b64eb](https://github.com/gchiesa/ska/commit/e2b64eb2fdadc9dd720c5a9216d38f39d1204a1c))


### Bug Fixes

* lint issues ([266fd9a](https://github.com/gchiesa/ska/commit/266fd9af7cf40986d7eb5025fe03638eaf6f6e45))


### Docs

* update README ([84126a7](https://github.com/gchiesa/ska/commit/84126a7e5fa7b87227bdeb774ec89a2d924ec72e))
* update README ([72cd21b](https://github.com/gchiesa/ska/commit/72cd21b646c10238be0e79521a20f0e7eda8decb))
* update README ([25bc2e1](https://github.com/gchiesa/ska/commit/25bc2e10b5541300c7046e55ebbf44a66594ba90))

## [0.0.6](https://github.com/gchiesa/ska/compare/v0.0.5...v0.0.6) (2024-07-26)


### Bug Fixes

* add missing secret ([1235ff2](https://github.com/gchiesa/ska/commit/1235ff296936534285f89f1a98790e01e739fb15))

## [0.0.5](https://github.com/gchiesa/ska/compare/v0.0.4...v0.0.5) (2024-07-26)


### Bug Fixes

* configure homebrew integration ([7db05b1](https://github.com/gchiesa/ska/commit/7db05b1e35ecf3d799b6fa05fbd115f25aa0aa40))

## [0.0.4](https://github.com/gchiesa/ska/compare/v0.0.3...v0.0.4) (2024-07-26)


### Bug Fixes

* goreleaser action ([69ece7f](https://github.com/gchiesa/ska/commit/69ece7feddb1def0e6fc27cb3c8ed0db7aabe3cd))

## [0.0.3](https://github.com/gchiesa/ska/compare/v0.0.2...v0.0.3) (2024-07-26)


### Bug Fixes

* goreleaser ([d7da3d9](https://github.com/gchiesa/ska/commit/d7da3d98a2c80134ae86a53ef1b8ed8fbae9b020))

## [0.0.2](https://github.com/gchiesa/ska/compare/v0.0.1...v0.0.2) (2024-07-26)


### Bug Fixes

* lint issues ([a68751a](https://github.com/gchiesa/ska/commit/a68751a2996a710df5850d6bbe76f6afb00f5a6c))
* missing release please manifest ([06a9aa7](https://github.com/gchiesa/ska/commit/06a9aa7d617dba30099e59ea49534df3934233dc))


### Other

* fix goland version ([9619d92](https://github.com/gchiesa/ska/commit/9619d921ad045e0278e53d72bb271fec3b30b0d4))
* fix golang-ci ([0d60139](https://github.com/gchiesa/ska/commit/0d601394f418c2285cd5acabb7742fd04d730dab))
* fix pipelines ([fca2931](https://github.com/gchiesa/ska/commit/fca29314f2a9f6fb28f42716d254c085fa7f99a3))
* fix workflows and config ([8dff677](https://github.com/gchiesa/ska/commit/8dff6770235a69f6fcda5a0a9811cedaaf0473ac))
