## 1.4.1 (September 15, 2020)
BUG FIXES:
- resource/vultr_server: Fix bug that did not allow user-data to be passed in as a string [65](https://github.com/terraform-providers/terraform-provider-vultr/pull/65)

## 1.4.0 (August 28, 2020)
FEATURES:
- New Resource : vultr_server_ipv4 [61](https://github.com/terraform-providers/terraform-provider-vultr/pull/61)
- New Resource : vultr_reverse_ipv4 [61](https://github.com/terraform-providers/terraform-provider-vultr/pull/61)
- New Resource : vultr_reverse_ipv6 [20](https://github.com/terraform-providers/terraform-provider-vultr/pull/20)
- New Data Source : vultr_server_ipv4 [61](https://github.com/terraform-providers/terraform-provider-vultr/pull/61)
- New Data Source : vultr_reverse_ipv4 [61](https://github.com/terraform-providers/terraform-provider-vultr/pull/61)
- New Data Source : vultr_reverse_ipv6 [20](https://github.com/terraform-providers/terraform-provider-vultr/pull/20)

IMPROVEMENTS:
- resource/vultr_server: Ability to enable/disable DDOS post create [59](https://github.com/vultr/terraform-provider-vultr/pull/59)
- resource/vultr_server: Ability to detach ISO post create [60](https://github.com/vultr/terraform-provider-vultr/pull/60)

## 1.3.2 (June 17, 2020)
IMPROVEMENTS:
- resource/vultr_dns_record: New custom importer allows you to import DNS Records [51](https://github.com/terraform-providers/terraform-provider-vultr/pull/51)
- resource/vultr_firewall_rule: New custom importer allows you to import Firewall Rules [52](https://github.com/terraform-providers/terraform-provider-vultr/pull/52)

## 1.3.1 (June 11, 2020)
IMPROVEMENTS:
- resource/vultr_dns_domain: Making `server_ip` optional. If `server_ip` is not supplied terraform will now create the DNS resource with no records. [48](https://github.com/terraform-providers/terraform-provider-vultr/pull/48)

## 1.3.0 (June 03, 2020)
BUG FIXES:
- resource/vultr_dns_record: Able to create record with `priority` of `0` [45](https://github.com/terraform-providers/terraform-provider-vultr/pull/45)

FEATURES:
- New Resource : vultr_object_storage [41](https://github.com/terraform-providers/terraform-provider-vultr/pull/41)
- New Data Source : vultr_object_storage [41](https://github.com/terraform-providers/terraform-provider-vultr/pull/41)

## 1.2.0 (May 27, 2020)
BUG FIXES:
- Typo in `website/docs/index.html.markdown` [38](https://github.com/terraform-providers/terraform-provider-vultr/pull/38)

FEATURES:
- New Resource : vultr_load_balancer [37](https://github.com/terraform-providers/terraform-provider-vultr/pull/37)
- New Data Source : vultr_load_balancer [37](https://github.com/terraform-providers/terraform-provider-vultr/pull/37)

## 1.1.5 (April 07, 2020)
BUG FIXES:
- resource/vultr_server: Detach ISO prior to deletion if instance was created with ISO. [34](https://github.com/terraform-providers/terraform-provider-vultr/issues/34)

## 1.1.4 (March 30, 2020)
IMPROVEMENTS:
- resource/block_storage: Adding new optional param `live` to allow attaching/detaching of block storage to instances without restarts [31](https://github.com/terraform-providers/terraform-provider-vultr/pull/31)

## 1.1.3 (March 24, 2020)
BUG FIXES:
- resource/reserved_ip: Adding `computed: true` to `attached_id` to prevent issues when Vultr assigns this [29](https://github.com/terraform-providers/terraform-provider-vultr/pull/29)
- resource/vultr_server: Adding `ForceNew: true` to `reserved_ip` to prevent any issues where the main floating ip may get deleted and cause issues with the instance [29](https://github.com/terraform-providers/terraform-provider-vultr/pull/29/files)

## 1.1.2 (March 19, 2020)
IMPROVEMENTS:
- resource/vultr_server: New optional field `reserved_ip` which lets you assign a `reserved_ip` during server creation [#26](https://github.com/terraform-providers/terraform-provider-vultr/pull/26).
- resource/reserved_ip: During deletion any instances that are attached to the reserved_ip are detached [#27](https://github.com/terraform-providers/terraform-provider-vultr/pull/27).
- Migrated to Terraform Plugin SDK [#21](https://github.com/terraform-providers/terraform-provider-vultr/issues/21)
- docs/snapshot fixed typo in snapshot [#19](https://github.com/terraform-providers/terraform-provider-vultr/pull/19)

## 1.1.1 (December 02, 2019)
IMPROVEMENTS:
- resource/vultr_block_storage: Attaches block storage on creation. Also reattaches block storage to instances if you taint the instance.[#9](https://github.com/terraform-providers/terraform-provider-vultr/pull/9) Thanks @oogy!

## 1.1.0 (November 22, 2019)
IMPROVEMENTS:
-   provider: Retry mechanism configuration `retry_limit` was added to allow adjusting how many retries should be attempted. [#7](https://github.com/terraform-providers/terraform-provider-vultr/pull/7).

BUG FIXES:
- Fixed go module name [#4](https://github.com/terraform-providers/terraform-provider-vultr/pull/4)

## 1.0.5 (October 24, 2019)

- Initial release under the terraform-providers/ namespace

## [v1.0.4](https://github.com/vultr/terraform-provider-vultr/compare/v1.0.3..v1.0.4) (2019-08-09)
### Fixes
- Fixes issue where using a snapshot would cause drift [#96](https://github.com/vultr/terraform-provider-vultr/issues/96)
### Enhancements
- You will now not have to define the `os_id` for the following server options
    - `app_id`
    - `iso_id`
    - `snapshot_id`

## [v1.0.3](https://github.com/vultr/terraform-provider-vultr/compare/v1.0.2..v1.0.3) (2019-07-18)
### Fixes
- Fixes issue where you could not use `os_id` and `script_id` together [#92](https://github.com/vultr/terraform-provider-vultr/issues/92)
### Breaking Changes
- You will now need to provide the `os_id` on each corresponding option
    - `app_id` - uses os_id `186`
    - `iso_id` - uses os_id `159`
    - `snap_id` - uses os_id `164`
    - `script_id` - uses os_id `159` or any os specific id
    
## [v1.0.2](https://github.com/vultr/terraform-provider-vultr/compare/v1.0.1..v1.0.2) (2019-07-15)
### Dependencies
* Updated dependencies [PR #89](https://github.com/vultr/terraform-provider-vultr/pull/89)
  * Govultr `v0.1.3` -> `v0.1.4`
  
## [v1.0.1](https://github.com/vultr/terraform-provider-vultr/compare/v1.0.0..v1.0.1) (2019-07-08)
### Fixes
- Fixed bug where scriptID was not being 
properly handled in server creation [#82](https://github.com/vultr/terraform-provider-vultr/issues/82)
### Enhancements 
- Added error handler on protocol case sensitivity [#83](https://github.com/vultr/terraform-provider-vultr/issues/83)
### Docs
- Typo in doc firewall_rule doc for protocol [#83](https://github.com/vultr/terraform-provider-vultr/issues/83)

## v1.0.0 (2019-06-24)
### Features
* Initial release
* Supported Data Sources
    * Account
    * Api Key
    * Application
    * Backup
    * Bare Metal Plan
    * Bare Metal Server
    * Block Storage
    * DNS Domain
    * Firewall Group
    * Iso Private
    * Iso Public
    * Network
    * OS
    * Plan
    * Region
    * Reserved IP
    * Server
    * Snapshot
    * SSH Key
    * Startup Script 
    * User
* Supported Resources
    * Bare Metal Server
    * Block Storage
    * DNS Domain
    * DNS Record
    * Firewall Group
    * Firewall Rule
    * ISO
    * Network
    * Reserved IP
    * Server
    * Snapshot
    * SSH Key
    * Startup Scripts
    * User
