# Changelog
## [v2.28.0](https://github.com/vultr/terraform-provider-vultr/compare/v2.27.1...v2.28.0) (2025-12-18)
### Enhancements
* data source/bare_metal_server: Add snapshot_id field [PR 616](https://github.com/vultr/terraform-provider-vultr/pull/616)
* resource/bare_metal_server: Add snapshot_id field [PR 616](https://github.com/vultr/terraform-provider-vultr/pull/616)

* data source/instance: Add snapshot_id field [PR 616](https://github.com/vultr/terraform-provider-vultr/pull/616)
* data source/instances: Add snapshot_id field [PR 616](https://github.com/vultr/terraform-provider-vultr/pull/616)
* resource/instance: Add snapshot_id & image_id fields [PR 616](https://github.com/vultr/terraform-provider-vultr/pull/616)

* data source/database: Add ca_certificate field [PR 616](https://github.com/vultr/terraform-provider-vultr/pull/616)
* resource/database: Add ca_certificate field [PR 616](https://github.com/vultr/terraform-provider-vultr/pull/616)

* resource/kubernetes: Migrate node pool labels and taints [PR 649](https://github.com/vultr/terraform-provider-vultr/pull/649)
* data source/kubernetes: Migrate node pool labels and taints [PR 649](https://github.com/vultr/terraform-provider-vultr/pull/649)
* resource/kubernetes_node_pool: Migrate node pool labels and taints [PR 649](https://github.com/vultr/terraform-provider-vultr/pull/649)
* data source/kubernetes_node_pool: Migrate node pool labels and taints [PR 649](https://github.com/vultr/terraform-provider-vultr/pull/649)

* resource/user: Add service_user field [PR 651](https://github.com/vultr/terraform-provider-vultr/pull/651)
* data source/user: Add service_user field [PR 651](https://github.com/vultr/terraform-provider-vultr/pull/651)

### Dependencies
* Bump github.com/vultr/govultr/v3 from 3.22.1 to 3.23.0 [PR 615](https://github.com/vultr/terraform-provider-vultr/pull/615)
* Bump golang.org/x/oauth2 from 0.30.0 to 0.32.0 [PR 621](https://github.com/vultr/terraform-provider-vultr/pull/621)
* Bump github.com/vultr/govultr/v3 from 3.23.0 to 3.24.0 [PR 619](https://github.com/vultr/terraform-provider-vultr/pull/619)
* Bump golang.org/x/oauth2 from 0.32.0 to 0.33.0 [PR 629](https://github.com/vultr/terraform-provider-vultr/pull/629)
* Bump golang.org/x/crypto from 0.36.0 to 0.45.0 [PR 634](https://github.com/vultr/terraform-provider-vultr/pull/634)
* Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.36.1 to 2.38.1 [PR 620](https://github.com/vultr/terraform-provider-vultr/pull/620)
* Bump github.com/hashicorp/terraform-plugin-log from 0.9.0 to 0.10.0 [PR 633](https://github.com/vultr/terraform-provider-vultr/pull/633)
* Update govultr from v3.25.0 to v3.26.0 [PR 648](https://github.com/vultr/terraform-provider-vultr/pull/648)
* Bump golang.org/x/oauth2 from 0.33.0 to 0.34.0 [PR 643](https://github.com/vultr/terraform-provider-vultr/pull/643)
* Update govultr from v3.26.0 to v3.26.1 [PR 650](https://github.com/vultr/terraform-provider-vultr/pull/650)

### Automation
* Ignore static check deprecation error on VPC2 [PR 636](https://github.com/vultr/terraform-provider-vultr/pull/636)

## [v2.27.1](https://github.com/vultr/terraform-provider-vultr/compare/v2.27.0...v2.27.1) (2025-08-14)
### Bug Fix
* resource/kubernetes: Add user data on new default node pool creation [PR 613](https://github.com/vultr/terraform-provider-vultr/pull/613)
* resource/kubernetes_nodepool: Handle user data changes in update context [PR 613](https://github.com/vultr/terraform-provider-vultr/pull/613)

## [v2.27.0](https://github.com/vultr/terraform-provider-vultr/compare/v2.26.0...v2.27.0) (2025-08-11)
### Enhancements
* data source/database: Add database backup schedule fields [PR 590](https://github.com/vultr/terraform-provider-vultr/pull/590)
* resource/database: Add database backup schedule fields [PR 590](https://github.com/vultr/terraform-provider-vultr/pull/590)
* resource/database_replica: Add database backup schedule fields [PR 590](https://github.com/vultr/terraform-provider-vultr/pull/590)
* data source/database: Add support for additional Kafka features [PR 599](https://github.com/vultr/terraform-provider-vultr/pull/599)
* resource/database: Add support for additional Kafka features [PR 599](https://github.com/vultr/terraform-provider-vultr/pull/599)
* resource/database_connector: Add support for additional Kafka features [PR 599](https://github.com/vultr/terraform-provider-vultr/pull/599)
* resource/database_quota: Add support for additional Kafka features [PR 599](https://github.com/vultr/terraform-provider-vultr/pull/599)
* data source/kubernetes: Add node pool user_data field [PR 607](https://github.com/vultr/terraform-provider-vultr/pull/607)
* resource/kubernetes: Add node pool user_data field [PR 607](https://github.com/vultr/terraform-provider-vultr/pull/607)
* resource/kubernetes_nodepool: Add user_data field [PR 607](https://github.com/vultr/terraform-provider-vultr/pull/607)
* resource/snapshot_from_url: Add use_uefi field [PR 609](https://github.com/vultr/terraform-provider-vultr/pull/609)

### Documentation
* data source/object_storage_tier: Add link to object storage tier API documentation page [PR 592](https://github.com/vultr/terraform-provider-vultr/pull/592)
* resource/database_connection_pool: Fix connection pool size type in example usage [PR 594](https://github.com/vultr/terraform-provider-vultr/pull/594)
* data source/kubernetes: Document node pool user_data field and usage [PR 611](https://github.com/vultr/terraform-provider-vultr/pull/611)
* resource/kubernetes: Document node pool user_data attribute and usage [PR 611](https://github.com/vultr/terraform-provider-vultr/pull/611)
* resource/snapshot_from_url: Document use_uefi attribute and usage [PR 611](https://github.com/vultr/terraform-provider-vultr/pull/611)

### Dependencies
* Bump github.com/vultr/govultr/v3 from 3.19.1 to 3.20.0 [PR 588](https://github.com/vultr/terraform-provider-vultr/pull/588)
* Bump golang.org/x/oauth2 from 0.29.0 to 0.30.0 [PR 587](https://github.com/vultr/terraform-provider-vultr/pull/587)
* Bump golang.org/x/net from 0.36.0 to 0.38.0 [PR 586](https://github.com/vultr/terraform-provider-vultr/pull/586)
* Bump github.com/vultr/govultr/v3 from 3.20.0 to 3.21.0 [PR 598](https://github.com/vultr/terraform-provider-vultr/pull/598)
* Bump github.com/vultr/govultr/v3 from 3.21.0 to 3.21.1 [PR 602](https://github.com/vultr/terraform-provider-vultr/pull/602)
* Update govultr from v3.21.1 to v3.22.0 [PR 605](https://github.com/vultr/terraform-provider-vultr/pull/605)
* Update govultr from v3.22.0 to v3.22.1 [PR 608](https://github.com/vultr/terraform-provider-vultr/pull/608)
* Bump github.com/cloudflare/circl from 1.3.7 to 1.6.1 [PR 595](https://github.com/vultr/terraform-provider-vultr/pull/595)

### Automation
* Migrate golangci-lint configuration to v2 [PR 606](https://github.com/vultr/terraform-provider-vultr/pull/606)

## [v2.26.0](https://github.com/vultr/terraform-provider-vultr/compare/v2.25.0...v2.26.0) (2025-04-09)
### Enhancements
* resource/instance: Add support for iPXE URL create param [PR 573](https://github.com/vultr/terraform-provider-vultr/pull/573)
* resource/bare_metal_server: Add support for VPC [PR 578](https://github.com/vultr/terraform-provider-vultr/pull/578)
* resource/kubernetes: Add support for VPC [PR 579](https://github.com/vultr/terraform-provider-vultr/pull/579)
* resource/kubernetes: Add labels to kubernetes node pools [PR 581](https://github.com/vultr/terraform-provider-vultr/pull/581)
* data source/kubernetes: Add labels to kubernetes node pools [PR 581](https://github.com/vultr/terraform-provider-vultr/pull/581)
* resource/kubernetes: Add taints to kubernetes node pools [PR 584](https://github.com/vultr/terraform-provider-vultr/pull/584)
* data source/kubernetes: Add taints to kubernetes node pools [PR 584](https://github.com/vultr/terraform-provider-vultr/pull/584)

### Bug Fixes
* resource/kubernetes: Fix default node pool labels [PR 582](https://github.com/vultr/terraform-provider-vultr/pull/582)
* data source/kubernetes: Fix default node pool labels [PR 582](https://github.com/vultr/terraform-provider-vultr/pull/582)

### Documentation
* data source/kubernetes: Add missing labels docs [PR 583](https://github.com/vultr/terraform-provider-vultr/pull/583)
* resource/kubernetes: Add missing labels docs [PR 583](https://github.com/vultr/terraform-provider-vultr/pull/583)

### Dependencies
* Bump github.com/vultr/govultr/v3 from 3.17.0 to 3.18.0 [PR 574](https://github.com/vultr/terraform-provider-vultr/pull/574)
* Update govultr from v3.18.0 to v3.19.0 [PR 577](https://github.com/vultr/terraform-provider-vultr/pull/577)
* Update govultr from v3.19.0 to v3.19.1 [PR 580](https://github.com/vultr/terraform-provider-vultr/pull/580)
* Bump golang.org/x/oauth2 from 0.28.0 to 0.29.0 [PR 575](https://github.com/vultr/terraform-provider-vultr/pull/575)

## [v2.25.0](https://github.com/vultr/terraform-provider-vultr/compare/v2.24.0...v2.25.0) (2025-03-14)
### Enhancements
* Add support for virtual file system storages [PR 571](https://github.com/vultr/terraform-provider-vultr/pull/571)

### Dependencies
* Update govultr from v3.16.1 to v3.17.0 [PR 569](https://github.com/vultr/terraform-provider-vultr/pull/569)
* Bump golang.org/x/net from 0.34.0 to 0.36.0 [PR 568](https://github.com/vultr/terraform-provider-vultr/pull/568)

### Documentation
* resource/vpc2: Add deprecation notice for all VPC2 elements [PR 570](https://github.com/vultr/terraform-provider-vultr/pull/570)
* resource/instance: Add deprecation notice for all VPC2 fields [PR 570](https://github.com/vultr/terraform-provider-vultr/pull/570)
* resource/bare_metal_server: Add deprecation notice for all VPC2 fields [PR 570](https://github.com/vultr/terraform-provider-vultr/pull/570)
* data_source/vpc2: Add deprecation notice for all VPC2 elements [PR 570](https://github.com/vultr/terraform-provider-vultr/pull/570)
* data_source/instance: Add deprecation notice for all VPC2 fields [PR 570](https://github.com/vultr/terraform-provider-vultr/pull/570)
* data_source/bare_metal_server: Add deprecation notice for all VPC2 fields [PR 570](https://github.com/vultr/terraform-provider-vultr/pull/570)

## [v2.24.0](https://github.com/vultr/terraform-provider-vultr/compare/v2.23.1...v2.24.0) (2025-03-10)
### Enhancements/Bug Fixes
* data source/object_storage: Add obeject storage tier data source [PR 565](https://github.com/vultr/terraform-provider-vultr/pull/565)
* resource/object_storage: Add tier param [PR 565](https://github.com/vultr/terraform-provider-vultr/pull/565)

### Documentation
* data source/object_storage: Add docs for the tier changes [PR 566](https://github.com/vultr/terraform-provider-vultr/pull/566)
* resource/object_storage: Add docs for the tier changes [PR 566](https://github.com/vultr/terraform-provider-vultr/pull/566)

### Clean Up
* resource/database: Remove managed Redis references [PR 549](https://github.com/vultr/terraform-provider-vultr/pull/549)
* Add terraform config to gitignore [PR 563](https://github.com/vultr/terraform-provider-vultr/pull/563)

### Dependencies
* Update govultr from v3.12.0 to v3.14.1 [PR 546](https://github.com/vultr/terraform-provider-vultr/pull/546)
* Update govultr from v3.14.1 to v3.16.0 [PR 561](https://github.com/vultr/terraform-provider-vultr/pull/561)
* Update govultr from v3.16.0 to v3.16.1 [PR 564](https://github.com/vultr/terraform-provider-vultr/pull/564)
* Bump golang.org/x/oauth2 from 0.24.0 to 0.25.0 [PR 542](https://github.com/vultr/terraform-provider-vultr/pull/542)
* Bump golang.org/x/net from 0.28.0 to 0.33.0 [PR 547](https://github.com/vultr/terraform-provider-vultr/pull/547)
* Bump golang.org/x/oauth2 from 0.25.0 to 0.28.0 [PR 558](https://github.com/vultr/terraform-provider-vultr/pull/558)
* Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.35.0 to 2.36.1 [PR 554](https://github.com/vultr/terraform-provider-vultr/pull/554)
* Update Go from v1.23 to v1.24 [PR 559](https://github.com/vultr/terraform-provider-vultr/pull/559)

## [v2.23.1](https://github.com/vultr/terraform-provider-vultr/compare/v2.23.0...v2.23.1) (2024-12-06)
### Automation
* Update github workflows from go 1.22 to 1.23 [PR 538](https://github.com/vultr/terraform-provider-vultr/pull/538)
* Update goreleaser action from v5 to v6 [PR 539](https://github.com/vultr/terraform-provider-vultr/pull/539)

## [v2.23.0](https://github.com/vultr/terraform-provider-vultr/compare/v2.22.1...v2.23.0) (2024-12-06)
### Enhancements
* resource/snapshot: Add default timeout on create [PR 536](https://github.com/vultr/terraform-provider-vultr/pull/536)

### Bug Fixes
* resource/kubernetes: Remove the set on enable_firewall [PR 531](https://github.com/vultr/terraform-provider-vultr/pull/531)
* resource/kubernetes: Fix node pool default tag lookup [PR 528](https://github.com/vultr/terraform-provider-vultr/pull/528)

### Deprecations
* resource/database: Deprecate Redis-named fields [PR 533](https://github.com/vultr/terraform-provider-vultr/pull/533)
* data_source/database: Deprecate Redis-named fields [PR 533](https://github.com/vultr/terraform-provider-vultr/pull/533)

### Dependencies
* Update govultr from v3.11.2 to v3.12.0 [PR 532](https://github.com/vultr/terraform-provider-vultr/pull/532)
* Bump golang.org/x/oauth2 from 0.23.0 to 0.24.0 [PR 527](https://github.com/vultr/terraform-provider-vultr/pull/527)

### Documentation
* Add provider installation instructions to README [PR 535](https://github.com/vultr/terraform-provider-vultr/pull/535)

### Automation
* Add goreleaser config version for v2 [PR 534](https://github.com/vultr/terraform-provider-vultr/pull/534)

### New Contributors
* @timurbabs made their first contribution in [PR 528](https://github.com/vultr/terraform-provider-vultr/pull/528)

## [v2.22.1](https://github.com/vultr/terraform-provider-vultr/compare/v2.22.0...v2.22.1) (2024-11-07)
### Bug Fixes
* resource/bare_metal_server: Set default value for user_scheme [PR 525](https://github.com/vultr/terraform-provider-vultr/pull/525)
* resource/instance: Set default value for user_scheme [PR 525](https://github.com/vultr/terraform-provider-vultr/pull/525)

## [v2.22.0](https://github.com/vultr/terraform-provider-vultr/compare/v2.21.0...v2.22.0) (2024-11-06)
### Enhancements
* resource/bare_metal_server: Add `user_scheme` options [PR 514](https://github.com/vultr/terraform-provider-vultr/pull/514)
* data_source/bare_metal_server: Add `user_scheme` options [PR 514](https://github.com/vultr/terraform-provider-vultr/pull/514)
* resource/instance: Add `user_scheme` options [PR 514](https://github.com/vultr/terraform-provider-vultr/pull/514)
* data_source/instance: Add `user_scheme` options [PR 514](https://github.com/vultr/terraform-provider-vultr/pull/514)
* resource/kubernetes: Add/improve existing resource import [PR 503](https://github.com/vultr/terraform-provider-vultr/pull/503)
* resource/database: Add support for Kafka [PR 522](https://github.com/vultr/terraform-provider-vultr/pull/522)
* data_source/database: Add support for Kafka [PR 522](https://github.com/vultr/terraform-provider-vultr/pull/522)

### Dependencies
* Update govultr from v3.8.1 to v3.9.0 [PR 504](https://github.com/vultr/terraform-provider-vultr/pull/504)
* Bump github.com/vultr/govultr/v3 from 3.9.0 to 3.10.0 [PR 511](https://github.com/vultr/terraform-provider-vultr/pull/511)
* Bump golang.org/x/oauth2 from 0.21.0 to 0.23.0 [PR 510](https://github.com/vultr/terraform-provider-vultr/pull/510)
* Update go from 1.21 to 1.23 [PR 513](https://github.com/vultr/terraform-provider-vultr/pull/513)
* Bump github.com/vultr/govultr/v3 from 3.10.0 to 3.11.0 [PR 516](https://github.com/vultr/terraform-provider-vultr/pull/516)
* Update govultr from v3.11.0 to v3.11.1 [PR 518](https://github.com/vultr/terraform-provider-vultr/pull/518)
* Bump github.com/vultr/govultr/v3 from 3.11.1 to 3.11.2 [PR 519](https://github.com/vultr/terraform-provider-vultr/pull/519)
* Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.34.0 to 2.35.0 [PR 521](https://github.com/vultr/terraform-provider-vultr/pull/521)

### Bug Fixes
* Ensure format strings for TF diag errors [PR 512](https://github.com/vultr/terraform-provider-vultr/pull/512)

### Clean up
* Remove unused travis CI config [PR 515](https://github.com/vultr/terraform-provider-vultr/pull/515)
* Remove vendored code [PR 520](https://github.com/vultr/terraform-provider-vultr/pull/520)

### Automation
* Add github CODEOWNERS file [PR 517](https://github.com/vultr/terraform-provider-vultr/pull/517)

## [v2.21.0](https://github.com/vultr/terraform-provider-vultr/compare/v2.20.1...v2.21.0) (2024-06-10)
### Enhancements
* resource/container_registry: add resource support [PR 445](https://github.com/vultr/terraform-provider-vultr/pull/445)
* resource/container_registry: add registry name validation [PR 493](https://github.com/vultr/terraform-provider-vultr/pull/493)
* data_source/container_registry: add data source support [PR 493](https://github.com/vultr/terraform-provider-vultr/pull/493)
* resource/user: change ACL schema to set to state drift [PR 495](https://github.com/vultr/terraform-provider-vultr/pull/495)
* resource/inference: add resource [PR 501](https://github.com/vultr/terraform-provider-vultr/pull/501)
* data_source/inference: add data source [PR 501](https://github.com/vultr/terraform-provider-vultr/pull/501)

### Deprecations
* resource/private_network: removed from provider [PR 496](https://github.com/vultr/terraform-provider-vultr/pull/496)
* data_source/private_network: removed from provider [PR 496](https://github.com/vultr/terraform-provider-vultr/pull/496)

### Dependencies
* Bump golang.org/x/oauth2 from 0.20.0 to 0.21.0 [PR 497](https://github.com/vultr/terraform-provider-vultr/pull/497)
* Update govultr from v2.7.0 to v2.8.1 [PR 500](https://github.com/vultr/terraform-provider-vultr/pull/500)

### Automation
* Lint fixes; add an updated lint configuration [PR 498](https://github.com/vultr/terraform-provider-vultr/pull/498)

### New Contributors
* @im6h made their first contribution in [PR 445](https://github.com/vultr/terraform-provider-vultr/pull/445)
* @F21 made their first contribution in [PR 495](https://github.com/vultr/terraform-provider-vultr/pull/495)

## [v2.20.1](https://github.com/vultr/terraform-provider-vultr/compare/v2.20.0...v2.20.1) (2024-05-29)
### Automation
* Update GPG import github action [PR 491](https://github.com/vultr/terraform-provider-vultr/pull/491)

## [v2.20.0](https://github.com/vultr/terraform-provider-vultr/compare/v2.19.0...v2.20.0) (2024-05-29)

### Enhancements
* resource/bare_metal_server: add support for mdisk_mode option [PR 489](https://github.com/vultr/terraform-provider-vultr/pull/489)

### Bug Fixes
* Stop using deprecated terraform helper resource for retries [PR 456](https://github.com/vultr/terraform-provider-vultr/pull/456)

### Dependencies
* Bump github.com/cloudflare/circl from 1.3.3 to 1.3.7 [PR 457](https://github.com/vultr/terraform-provider-vultr/pull/457)
* Bump github.com/vultr/govultr/v3 from 3.6.0 to 3.6.2 [PR 465](https://github.com/vultr/terraform-provider-vultr/pull/465)
* Bump golang.org/x/oauth2 from 0.15.0 to 0.17.0 [PR 464](https://github.com/vultr/terraform-provider-vultr/pull/464)
* Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.31.0 to 2.32.0 [PR 461](https://github.com/vultr/terraform-provider-vultr/pull/461)
* Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.32.0 to 2.33.0 [PR 466](https://github.com/vultr/terraform-provider-vultr/pull/466)
* Bump github.com/vultr/govultr/v3 from 3.6.2 to 3.6.3 [PR 467](https://github.com/vultr/terraform-provider-vultr/pull/467)
* Bump github.com/vultr/govultr/v3 from 3.6.3 to 3.6.4 [PR 471](https://github.com/vultr/terraform-provider-vultr/pull/471)
* Bump golang.org/x/net from 0.21.0 to 0.23.0 [PR 478](https://github.com/vultr/terraform-provider-vultr/pull/478)
* Bump golang.org/x/oauth2 from 0.17.0 to 0.19.0 [PR 475](https://github.com/vultr/terraform-provider-vultr/pull/475)
* Bump google.golang.org/protobuf from 1.32.0 to 1.33.0 [PR 473](https://github.com/vultr/terraform-provider-vultr/pull/473)
* Bump golang.org/x/oauth2 from 0.19.0 to 0.20.0 [PR 483](https://github.com/vultr/terraform-provider-vultr/pull/483)
* Bump github.com/vultr/govultr/v3 from 3.6.4 to 3.7.0 [PR 488](https://github.com/vultr/terraform-provider-vultr/pull/488)
* Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.33.0 to 2.34.0 [PR 485](https://github.com/vultr/terraform-provider-vultr/pull/485)
* Update go from v1.20 to v1.21 [PR 486](https://github.com/vultr/terraform-provider-vultr/pull/486)

### Automation
* Update notify-pr.yml [PR 481](https://github.com/vultr/terraform-provider-vultr/pull/481)
* Fix missing step on go-checks action [PR 480](https://github.com/vultr/terraform-provider-vultr/pull/480)
* Fix mattermost notifications [PR 484](https://github.com/vultr/terraform-provider-vultr/pull/484)
* CI & automation actions updates [PR 479](https://github.com/vultr/terraform-provider-vultr/pull/479)

### New Contributors
* @fjoenichols made their first contribution in [PR 489](https://github.com/vultr/terraform-provider-vultr/pull/489)

##[v2.19.0](https://github.com/vultr/terraform-provider-vultr/compare/v2.18.0...v2.19.0) (2024-01-03)
### Enhancements
* resource/instances: Allow creation without public IP [PR 450](https://github.com/vultr/terraform-provider-vultr/pull/450)
* resource/instances: Add marketplace app variables support [PR 448](https://github.com/vultr/terraform-provider-vultr/pull/448)
* resource/bare_metal_server: Add marketplace app variables support [PR 448](https://github.com/vultr/terraform-provider-vultr/pull/448)
* resource/load_balancers: Add retry to delete [PR 451](https://github.com/vultr/terraform-provider-vultr/pull/451)

### Bug Fixes
* resource/bare_metal_server: Fix nil interface panic on creation [PR 452](https://github.com/vultr/terraform-provider-vultr/pull/452)

### Documentation
* resource/instances: Add disable_public_ipv4 field to webdocs [PR 453](https://github.com/vultr/terraform-provider-vultr/pull/453)

### Dependencies
* Update govultr from v3.5.0 to v3.6.0 [PR 444](https://github.com/vultr/terraform-provider-vultr/pull/444)
* Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.30.0 to 2.31.0 [PR 442](https://github.com/vultr/terraform-provider-vultr/pull/442)
* Bump golang.org/x/crypto from 0.16.0 to 0.17.0 [PR 447](https://github.com/vultr/terraform-provider-vultr/pull/447)

### Automation
* Use GITHUB_OUTPUT envvar instead of set-output command [PR 449](https://github.com/vultr/terraform-provider-vultr/pull/449)

### New Contributors
* @OpenGLShaders made their first contribution in [PR 450](https://github.com/vultr/terraform-provider-vultr/pull/450)
* @arunsathiya made their first contribution in [PR 449](https://github.com/vultr/terraform-provider-vultr/pull/449)

##[v2.18.0](https://github.com/vultr/terraform-provider-vultr/compare/v2.17.1...v2.18.0) (2023-12-11)
### Enhancements
* resource/bare_metal_server: Add Persistent PXE field [PR 368](https://github.com/vultr/terraform-provider-vultr/pull/368)
* data_source/instances: Add instances data source [PR 296](https://github.com/vultr/terraform-provider-vultr/pull/296)
* data_source/ssh_key: Export the key ID [PR 338](https://github.com/vultr/terraform-provider-vultr/pull/338)
* resource/kubernetes: Add firewall field [PR 434](https://github.com/vultr/terraform-provider-vultr/pull/434)
* data_source/kubernetes: Add firewall field [PR 434](https://github.com/vultr/terraform-provider-vultr/pull/434)
* resource/database: Add redis user access control [PR 439](https://github.com/vultr/terraform-provider-vultr/pull/439)

### Bug Fixes
* Remove deprecated SDK meta version function usage [PR 432](https://github.com/vultr/terraform-provider-vultr/pull/432)
* data_source/database: Fix bug with flattening non-FerretDB replicas [PR 427](https://github.com/vultr/terraform-provider-vultr/pull/427)

### Documentation
* Add documentation for the instances data source [PR 431](https://github.com/vultr/terraform-provider-vultr/pull/431)

### Dependencies
* Bump github.com/vultr/govultr/v3 from 3.3.4 to 3.4.0 [PR 430](https://github.com/vultr/terraform-provider-vultr/pull/430)
* Bump golang.org/x/oauth2 from 0.13.0 to 0.14.0 [PR 429](https://github.com/vultr/terraform-provider-vultr/pull/429)
* Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.29.0 to 2.30.0 [PR 428](https://github.com/vultr/terraform-provider-vultr/pull/428)
* Update govultr from v3.4.0 to v3.5.0 [PR 438](https://github.com/vultr/terraform-provider-vultr/pull/438)
* Bump golang.org/x/oauth2 from 0.14.0 to 0.15.0 [PR 435](https://github.com/vultr/terraform-provider-vultr/pull/435)

### New Contributors
* @neilmock made their first contribution in [PR 368](https://github.com/vultr/terraform-provider-vultr/pull/368)
* @aarani made their first contribution in [PR 296](https://github.com/vultr/terraform-provider-vultr/pull/296)
* @Byteflux made their first contribution in [PR 434](https://github.com/vultr/terraform-provider-vultr/pull/434)

## [v2.17.1](https://github.com/vultr/terraform-provider-vultr/compare/v2.17.0...v2.17.1) (2023-10-31)
### Enhancements
* resource/database: Add FerretDB Support [PR 422](https://github.com/vultr/terraform-provider-vultr/pull/422)
* data_source/database: Add FerretDB Support [PR 422](https://github.com/vultr/terraform-provider-vultr/pull/422)
* resource/kubernetes: Add support for the VKE HA control plane option [PR 423](https://github.com/vultr/terraform-provider-vultr/pull/423)
* data_source/kubernetes: Add support for the VKE HA control plane option [PR 423](https://github.com/vultr/terraform-provider-vultr/pull/423)

### Bug Fixes
* resource/vpc2: Fix ForceNew when optional fields not set [PR 424](https://github.com/vultr/terraform-provider-vultr/pull/424)

### Dependencies
* Update govultr to v3.3.4 [PR 421](https://github.com/vultr/terraform-provider-vultr/pull/421)
* Bump google.golang.org/grpc from 1.57.0 to 1.57.1 [PR 419](https://github.com/vultr/terraform-provider-vultr/pull/419)

## [v2.17.0](https://github.com/vultr/terraform-provider-vultr/compare/v2.16.4...v2.17.0) (2023-10-25)
### Enhancement
* Database: Add support for public/private hostnames [PR 416](https://github.com/vultr/terraform-provider-vultr/pull/416)

### Documentation
* Update some invalid plans in documentation [PR 414](https://github.com/vultr/terraform-provider-vultr/pull/414)

### Dependencies
* Update govultr to v3.3.2 [PR 415](https://github.com/vultr/terraform-provider-vultr/pull/415)
* Update govultr to v3.3.3 [PR 417](https://github.com/vultr/terraform-provider-vultr/pull/417)

## [v2.16.4](https://github.com/vultr/terraform-provider-vultr/compare/v2.16.3...v2.16.4) (2023-10-16)
### Documentation 
* kubernetes: fix typo in VKE plan example [PR 412](https://github.com/vultr/terraform-provider-vultr/pull/412)

## [v2.16.3](https://github.com/vultr/terraform-provider-vultr/compare/v2.16.2...v2.16.3) (2023-10-13)
### Documentation 
* Update a few of the resource documentation pages [PR 409](https://github.com/vultr/terraform-provider-vultr/pull/409)

### Dependencies
* Bump golang.org/x/oauth2 from 0.12.0 to 0.13.0 [PR 407](https://github.com/vultr/terraform-provider-vultr/pull/407)
* Bump golang.org/x/net from 0.15.0 to 0.17.0 [PR 408](https://github.com/vultr/terraform-provider-vultr/pull/408)

### Automation
* Extend the nightly acceptance test timeout by one hour [PR 405](https://github.com/vultr/terraform-provider-vultr/pull/405)

## [v2.16.2](https://github.com/vultr/terraform-provider-vultr/compare/v2.16.1...v2.16.2) (2023-09-25)
# Enhancement
* data_source/instance: Add a per-page param for instance data source [PR 384](https://github.com/vultr/terraform-provider-vultr/pull/384)

# Bug Fix
* resource/database: Add missing vpc_id property to managed database read replicas [PR 403](https://github.com/vultr/terraform-provider-vultr/pull/403)
* data_source/database: Add missing vpc_id property to managed database read replicas [PR 403](https://github.com/vultr/terraform-provider-vultr/pull/403)

## [v2.16.1](https://github.com/vultr/terraform-provider-vultr/compare/v2.16.0...v2.16.1) (2023-09-22)
### Bug Fix
* resource/vpc2: Fix delete retries and detach errors with VPC2s [PR 399](https://github.com/vultr/terraform-provider-vultr/pull/399)
* resource/bare_metal_server: Revert BM update delay for detach VPC2 [PR 400](https://github.com/vultr/terraform-provider-vultr/pull/400)

### Dependencies
* Bump github.com/vultr/govultr/v3 from 3.2.0 to 3.3.1 [PR 391](https://github.com/vultr/terraform-provider-vultr/pull/391)
* Bump golang.org/x/oauth2 from 0.10.0 to 0.12.0 [PR 393](https://github.com/vultr/terraform-provider-vultr/pull/393)
* Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.27.0 to 2.29.0 [PR 394](https://github.com/vultr/terraform-provider-vultr/pull/394)

## [v2.16.0](https://github.com/vultr/terraform-provider-vultr/compare/v2.15.1...v2.16.0) (2023-09-21)
### Enhancements
* resource/instance: Add/update retry on create/delete actions [PR 362](https://github.com/vultr/terraform-provider-vultr/pull/362)
* resource/vpc: Add/update retry on create/delete actions [PR 362](https://github.com/vultr/terraform-provider-vultr/pull/362)
* resource/database: Add support for DBaaS VPC networks [PR 385](https://github.com/vultr/terraform-provider-vultr/pull/385)
* resource/vpc2: Add VPC 2.0 [PR 389](https://github.com/vultr/terraform-provider-vultr/pull/389)
* data_source/vpc2: Add VPC 2.0 [PR 389](https://github.com/vultr/terraform-provider-vultr/pull/389)
* resource/bare_metal_server: Wait for VPC 2.0 detachments on BM [PR 396](https://github.com/vultr/terraform-provider-vultr/pull/396)

### Documentation
* load balancer: Set non-required fields with default values [PR 365](https://github.com/vultr/terraform-provider-vultr/pull/365)
* kubernetes: Update resource docs [PR 369](https://github.com/vultr/terraform-provider-vultr/pull/369)

### Dependencies
* Update govultr to v3.1.0 [PR 380](https://github.com/vultr/terraform-provider-vultr/pull/380)
* Bump github.com/hashicorp/terraform-plugin-log from 0.8.0 to 0.9.0 [PR 364](https://github.com/vultr/terraform-provider-vultr/pull/364)
* Bump github.com/vultr/govultr/v3 from 3.0.2 to 3.0.3 [PR 366](https://github.com/vultr/terraform-provider-vultr/pull/366)
* Bump github.com/vultr/govultr/v3 from 3.1.0 to 3.2.0 [PR 382](https://github.com/vultr/terraform-provider-vultr/pull/382)
* Bump golang.org/x/oauth2 from 0.8.0 to 0.10.0 [PR 376](https://github.com/vultr/terraform-provider-vultr/pull/376)
* Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.26.1 to 2.27.0 [PR 370](https://github.com/vultr/terraform-provider-vultr/pull/370)

### Automation
* Update goreleaser-action to v4 and fix `release` args [PR 361](https://github.com/vultr/terraform-provider-vultr/pull/361)

### New Contributors
* @tavsec made their first contribution in [PR 365](https://github.com/vultr/terraform-provider-vultr/pull/365)
* @ogawa0071 made their first contribution in [PR 389](https://github.com/vultr/terraform-provider-vultr/pull/389)

## [v2.15.1](https://github.com/vultr/terraform-provider-vultr/compare/v2.15.0...v2.15.1) (2023-05-10)
### Enhancement
* Add documentation for Vultr managed database data sources and resources [PR 356](https://github.com/vultr/terraform-provider-vultr/pull/356)
* Add VPC delete retries [PR 358](https://github.com/vultr/terraform-provider-vultr/pull/358)

### Dependencies
* Bump golang.org/x/oauth2 from 0.7.0 to 0.8.0 [PR 357](https://github.com/vultr/terraform-provider-vultr/pull/357)

## [v2.15.0](https://github.com/vultr/terraform-provider-vultr/compare/v2.14.1...v2.15.0) (2023-05-04)
### Enhancement
* resource/database: Add Support for Vultr Managed Databases [PR 352](https://github.com/vultr/terraform-provider-vultr/pull/352)
* data source/database: Add Support for Vultr Managed Databases [PR 352](https://github.com/vultr/terraform-provider-vultr/pull/352)

### Automation
* Update acceptance test configurations [PR 353](https://github.com/vultr/terraform-provider-vultr/pull/353)

### New Contributors
* @christhemorse made their first contribution in [PR 352](https://github.com/vultr/terraform-provider-vultr/pull/352)

## [v2.14.1](https://github.com/vultr/terraform-provider-vultr/compare/v2.14.0...v2.14.1) (2023-04-28)
### Enhancement
* resource/kuberneters: added vke certs as exported atrributes [PR 349](https://github.com/vultr/terraform-provider-vultr/pull/349)
* data source/kuberneters: added vke certs as exported atrributes [PR 349](https://github.com/vultr/terraform-provider-vultr/pull/349)

### New Contributors
* @happytreees made their first contribution in [PR 349](https://github.com/vultr/terraform-provider-vultr/pull/349)

## [v2.14.0](https://github.com/vultr/terraform-provider-vultr/compare/v2.13.0...v2.14.0) (2023-04-13)
### Enhancement
* resource/kubernetes: Add VKE k8s version upgrade functionality [PR 344](https://github.com/vultr/terraform-provider-vultr/pull/344)
* resource/kubernetes: Mark the kube_config schema value as sensitive [PR 346](https://github.com/vultr/terraform-provider-vultr/pull/346)

### Dependencies
* Bump golang.org/x/oauth2 from 0.6.0 to 0.7.0 [PR 343](https://github.com/vultr/terraform-provider-vultr/pull/343)

## [v2.13.0](https://github.com/vultr/terraform-provider-vultr/compare/v2.12.1...v2.13.0) (2023-04-03)
### Enhancement
* resource/reserved ip: Add missing resource warning for reserved IP [PR 327](https://github.com/vultr/terraform-provider-vultr/pull/327)
* resource/dns domain: Add missing resource warnings [PR 323](https://github.com/vultr/terraform-provider-vultr/pull/323)
* resource/block storage: Add missing resource warnings [PR 323](https://github.com/vultr/terraform-provider-vultr/pull/323)
* resource/user: Add missing resource warnings [PR 323](https://github.com/vultr/terraform-provider-vultr/pull/323)
* resource/startup script: Add missing resource warnings [PR 323](https://github.com/vultr/terraform-provider-vultr/pull/323)
* resource/ssh key: Add missing resource warnings [PR 323](https://github.com/vultr/terraform-provider-vultr/pull/323)
* resource/firewall rule: Add missing resource warnings [PR 323](https://github.com/vultr/terraform-provider-vultr/pull/323)

### Dependencies
* Bump golang.org/x/text from 0.3.7 to 0.3.8 [PR 324](https://github.com/vultr/terraform-provider-vultr/pull/324)
* Update govultr to v3.0.1 [PR 336](https://github.com/vultr/terraform-provider-vultr/pull/336)
* Bump github.com/vultr/govultr/v3 from 3.0.1 to 3.0.2 [PR 339](https://github.com/vultr/terraform-provider-vultr/pull/339)

### Automation
* Fix broken workflows resulting from go version 1.20 [PR 340](https://github.com/vultr/terraform-provider-vultr/pull/340)
* bump setup-go in github workflow [PR 337](https://github.com/vultr/terraform-provider-vultr/pull/337)

### New Contributors
* @mondragonfx made their first contribution in [PR 336](https://github.com/vultr/terraform-provider-vultr/pull/336)

## [v2.12.1](https://github.com/vultr/terraform-provider-vultr/compare/v2.12.0...v2.12.1) (2023-02-10)
### Enhancement
* resource/instance: Add check for & detach of ISO on instance delete [312](https://github.com/vultr/terraform-provider-vultr/pull/312)
* All resources that use "region":
    - Add DiffSuppressFunc to ignore case [318](https://github.com/vultr/terraform-provider-vultr/pull/318)

## [v2.12.0](https://github.com/vultr/terraform-provider-vultr/compare/v2.11.4...v2.12.0) (2022-12-08)
### Enhancement
* resource/instance: remove deprecated tag fields [297](https://github.com/vultr/terraform-provider-vultr/pull/297)
* resource/bare_metal_server: remove deprecated tag fields [297](https://github.com/vultr/terraform-provider-vultr/pull/297)
* data source/instance: remove deprecated tag fields [297](https://github.com/vultr/terraform-provider-vultr/pull/297)
* data source/bare_metal_server: remove deprecated tag fields [297](https://github.com/vultr/terraform-provider-vultr/pull/297)

### Bug Fix
* everything: golangci-lint fixes [302](https://github.com/vultr/terraform-provider-vultr/pull/302)

### Documentation
* Fixed typo [279](https://github.com/vultr/terraform-provider-vultr/pull/279)
* Update rate-limit documentation [283](https://github.com/vultr/terraform-provider-vultr/pull/283)
* resource/instance_ipv4 fix type error on reboot [292](https://github.com/vultr/terraform-provider-vultr/pull/292)
* resource/bare_metal_server: update floating IP description [293](https://github.com/vultr/terraform-provider-vultr/pull/293)
* resource/instance: remove the tag field from the docs [297](https://github.com/vultr/terraform-provider-vultr/pull/297)
* resource/bare_metal_server: remove the tag field from the docs [297](https://github.com/vultr/terraform-provider-vultr/pull/297)
* data source/instance: remove the tag field from the docs [297](https://github.com/vultr/terraform-provider-vultr/pull/297)
* data source/bare_metal_server: remove the tag field from the docs [297](https://github.com/vultr/terraform-provider-vultr/pull/297)

### Dependency
* update terraform-sdk from 2.19.0 to 2.21.0 [280](https://github.com/vultr/terraform-provider-vultr/pull/280)
* update terraform-sdk from 2.21.0 to 2.24.0 [294](https://github.com/vultr/terraform-provider-vultr/pull/294)
* update terraform-sdk from 2.24.0 to 2.24.1 [298](https://github.com/vultr/terraform-provider-vultr/pull/298)
* update go to v1.19 [303](https://github.com/vultr/terraform-provider-vultr/pull/303)
* update goreleaser to v1.19 [305](https://github.com/vultr/terraform-provider-vultr/pull/305)

### New Contributors
* @nschlemm made their first contribution in [279](https://github.com/vultr/terraform-provider-vultr/pull/279)
* @jesseorr made their first contribution in [292](https://github.com/vultr/terraform-provider-vultr/pull/292)
* @jasites made their first contribution in [293](https://github.com/vultr/terraform-provider-vultr/pull/293)

## [v2.11.4](https://github.com/vultr/terraform-provider-vultr/compare/v2.11.3...v2.11.4) (2022-07-25)
### Enhancement
* data source/object storage cluster: add datasource for object storage cluster [275](https://github.com/vultr/terraform-provider-vultr/pull/275)

### Documentatio
* data source/object storage cluster: add docs for object storage cluster [275](https://github.com/vultr/terraform-provider-vultr/pull/275)

### Dependency
* update terraform-sdk to v2.18.0 [273](https://github.com/vultr/terraform-provider-vultr/pull/273)
* update terraform-plugin-sdk from 2.18.0 to 2.19.0 [274](https://github.com/vultr/terraform-provider-vultr/pull/274)

## [v2.11.3](https://github.com/vultr/terraform-provider-vultr/compare/v2.11.2...v2.11.3) (2022-06-14)
### Enchancement
* resource/reserved_ip: Add support for reserved IP label updates [268](https://github.com/vultr/terraform-provider-vultr/pull/268)

### Documentation
* resource/instance: Fix typo [268](https://github.com/vultr/terraform-provider-vultr/pull/268)
* resource/reverse_ip: Fix type [268](https://github.com/vultr/terraform-provider-vultr/pull/268)

## [v2.11.2](https://github.com/vultr/terraform-provider-vultr/compare/v2.11.1...v2.11.2) (2022-06-03)
### Enhancement
* data source/plan: Add GPU fields [264](https://github.com/vultr/terraform-provider-vultr/pull/264)

### Bug Fix
* Fix acceptance tests [260](https://github.com/vultr/terraform-provider-vultr/pull/260)

### Dependency
* update govultr to v2.17.1 [262](https://github.com/vultr/terraform-provider-vultr/pull/262)
* update github.com/hashicorp/terraform-plugin-sdk/v2 from 2.16.0 to 2.17.0 [261](https://github.com/vultr/terraform-provider-vultr/pull/261)

## [v2.11.1](https://github.com/vultr/terraform-provider-vultr/compare/v2.11.0...v2.11.1) (2022-05-18)
### Documentation
* resource/instance: fix incorrect import example [251](https://github.com/vultr/terraform-provider-vultr/pull/251)
* resource/instance_ipv4: fix vultr_instance_ipv4 resource doc and argument reference [253](https://github.com/vultr/terraform-provider-vultr/pull/253)

### Dependency
* updated govultr from v1.16.0 to v1.17.0 [255](https://github.com/vultr/terraform-provider-vultr/pull/255)

### Bug Fix
* resource/kubernetes_nodepool: fix `tag` so that it can be deleted [255](https://github.com/vultr/terraform-provider-vultr/pull/255)
* resource/instance: fix `tag` so that it can be deleted [255](https://github.com/vultr/terraform-provider-vultr/pull/255)
* resource/bare_metal_server: fix `tag` so that it can be deleted [255](https://github.com/vultr/terraform-provider-vultr/pull/255)

### New Contributors
* @NicolasCARPi made their first contribution in [251](https://github.com/vultr/terraform-provider-vultr/pull/251)

## [v2.11.0](https://github.com/vultr/terraform-provider-vultr/compare/v2.10.1..v2.11.0) (2022-05-11)
### Documentation
* resource/instance: add additional examples for backups [246](https://github.com/vultr/terraform-provider-vultr/pull/246)
* resource/kubernetes: update examples for default optional node pool [249](https://github.com/vultr/terraform-provider-vultr/pull/249)
* readme: add link to quickstart guide [244](https://github.com/vultr/terraform-provider-vultr/pull/244)

### Dependency
* updated terraform-plugin-sdk from 2.15.0 to 2.16.0 [245](https://github.com/vultr/terraform-provider-vultr/pull/245)
* updated terraform-plugin-sdk from 2.12.0 to 2.15.0 [242](https://github.com/vultr/terraform-provider-vultr/pull/242)
* updated Go v1.16 -> v1.17  [221](https://github.com/vultr/terraform-provider-vultr/pull/221)
* updated govultr from 2.14.2 to 2.15.1 [233](https://github.com/vultr/terraform-provider-vultr/pull/233)
* updated govultr from 2.15.1 to 2.16.0 [241](https://github.com/vultr/terraform-provider-vultr/pull/241)

### Enhancement
* resource/kubernetes: allow removal of default node pool after resource creation [248](https://github.com/vultr/terraform-provider-vultr/pull/248)
* resource/kubernetes: add support for auto scaler options on node pools [247](https://github.com/vultr/terraform-provider-vultr/pull/247)
* resource/kubernetes node pools: add support for auto scaler options on node pools [247](https://github.com/vultr/terraform-provider-vultr/pull/247)
* data source/kubernetes: add auto scaler fields[247](https://github.com/vultr/terraform-provider-vultr/pull/247)
* data source/kubernetes node pools: add auto scaler fields [247](https://github.com/vultr/terraform-provider-vultr/pull/247)
* resource/block storage: add block type [238](https://github.com/vultr/terraform-provider-vultr/pull/238)
* data source/block storage: add block type field [238](https://github.com/vultr/terraform-provider-vultr/pull/238)
* resource/instance: add VPC support [237](https://github.com/vultr/terraform-provider-vultr/pull/237)
* resource/load balancer: add VPC support [237](https://github.com/vultr/terraform-provider-vultr/pull/237)
* data source/instance: add VPC fields[237](https://github.com/vultr/terraform-provider-vultr/pull/237)
* data source/load balancer: add VPC support [237](https://github.com/vultr/terraform-provider-vultr/pull/237)
* resource/kubernetes: add better error handling to reads [236](https://github.com/vultr/terraform-provider-vultr/pull/236)
* resource/kubernetes node pools: add better error handling to reads [236](https://github.com/vultr/terraform-provider-vultr/pull/236)
* resource/instance: add support for tags [240](https://github.com/vultr/terraform-provider-vultr/pull/240)
* resource/bare metal: add support for tags [240](https://github.com/vultr/terraform-provider-vultr/pull/240)
* data source/instance: add support for tags [240](https://github.com/vultr/terraform-provider-vultr/pull/240)
* data source/bare metal: add support for tags [240](https://github.com/vultr/terraform-provider-vultr/pull/240)

### Bug Fix
* resource/kubernetes: fix labeling on cluster updates [239](https://github.com/vultr/terraform-provider-vultr/pull/239)
* resource/firewall rule: read from correct govultr data [243](https://github.com/vultr/terraform-provider-vultr/pull/243)

### New Contributors
* @optik-aper made their first contribution in [238](https://github.com/vultr/terraform-provider-vultr/pull/238)
* @dfinr made their first contribution in [244](https://github.com/vultr/terraform-provider-vultr/pull/244)
* @travispaul made their first contribution in [246](https://github.com/vultr/terraform-provider-vultr/pull/246)

## 2.10.1 (March 30, 2022)
Enhancement:
* vultr_resource_instance : Updating hostname will result in a forcenew change [226](https://github.com/vultr/terraform-provider-vultr/pull/226)

## 2.10.0 (March 25, 2022)
### Dependency
* updated Go v1.16 -> v1.17  [221](https://github.com/vultr/terraform-provider-vultr/pull/221)
* updated terraform-plugin-sdk from 2.10.1 to 2.12.0 [218](https://github.com/vultr/terraform-provider-vultr/pull/218)
* updated govultr from 2.14.1 to 2.14.2 [219](https://github.com/vultr/terraform-provider-vultr/pull/219)

### Enhancement
* vultr_resource_block : add waits for active status [222](https://github.com/vultr/terraform-provider-vultr/pull/222)

## 2.9.1 (February 2, 2022)
### Dependency
* updated govultr to v2.14.0 -> v2.14.1  [210](https://github.com/vultr/terraform-provider-vultr/pull/210)

## 2.9.0 (January 21, 2022)
### Enhancement
* datasource/kubernetes: New datasource for VKE [198](https://github.com/vultr/terraform-provider-vultr/pull/198)
* Updated all datasources deprecations read -> readContext [204](https://github.com/vultr/terraform-provider-vultr/pull/204)

### Bug Fix
* datasource/backups : fix scheme mismatch [201](https://github.com/vultr/terraform-provider-vultr/pull/201)

### Dependency
* updated govultr to v2.12.0 -> v2.14.0  [206](https://github.com/vultr/terraform-provider-vultr/pull/206)

## 2.8.1 (December 20, 2021)
### Bug Fix
* resource/instance: fix importer [192](https://github.com/vultr/terraform-provider-vultr/pull/192) Thanks @vincentbernat

## 2.8.0 (December 08, 2021)
### Enhancement
* Implement datasource filtering on string lists [188](https://github.com/vultr/terraform-provider-vultr/pull/188) Thanks @kaorihinata

## 2.7.0 (December 6, 2021)
### Dependencies
* Bump govultr to v2.12.0, adjust monthly charges to float [182](https://github.com/vultr/terraform-provider-vultr/pull/182)

## 2.6.0 (November 19, 2021)
### Enhancement
* resource/bare_metal: Add timeout options [175](https://github.com/vultr/terraform-provider-vultr/pull/175)

### Bug Fix
* datasource/account : Fix type mismatch for billing fields [174](https://github.com/vultr/terraform-provider-vultr/pull/174)
* resource/instance : Fix invalid error message change [178](https://github.com/vultr/terraform-provider-vultr/pull/178)
* resource/instance : Fix issue where changing hostname didn't trigger hostname change [180](https://github.com/vultr/terraform-provider-vultr/pull/180)

### Documentatio
* resource/snapshots : fix typo [167](https://github.com/vultr/terraform-provider-vultr/pull/167)
* resources/vultr_kubernetes : Add description that kubeconfigs are base64 encoded [169](https://github.com/vultr/terraform-provider-vultr/pull/169)

### Dependency
* updated govultr to v2.9.2 -> v2.10.0  [179](https://github.com/vultr/terraform-provider-vultr/pull/179)

## 2.5.0 (October 22, 2021)
### Enhancement
* resource/vultr_kubernetes: New resource that allows for deployment of VKE clusters [165](https://github.com/vultr/terraform-provider-vultr/pull/165)
* resource/vultr_kubernetes_node_pools: New resource that allows for deployment of node pools to existing VKE Cluster[165](https://github.com/vultr/terraform-provider-vultr/pull/165)


## 2.4.2 (September 15, 2021)
### Bug Fix
* resource/load_balancer: added missing `region` and `ssl_redirect` values from being set [163](https://github.com/vultr/terraform-provider-vultr/pull/163)

## 2.4.1 (August 13, 2021)
### Enhancement
* resource/instance: increased default timeout for create/update from 20 to 60 minutes [160](https://github.com/vultr/terraform-provider-vultr/pull/160)

## 2.4.0 (August 02, 2021)
### Enhancement
* resource/instance: add marketplace support with `image_id` [150](https://github.com/vultr/terraform-provider-vultr/pull/150)
* resource/bare_metal: add marketplace support with `image_id` [150](https://github.com/vultr/terraform-provider-vultr/pull/150)
* datasource/applications: adds marketplace support [150](https://github.com/vultr/terraform-provider-vultr/pull/150)
* Add openBSD to builds [155](https://github.com/vultr/terraform-provider-vultr/pull/155)

### Bug Fix
* resource/bare_metal: fix importer [157](https://github.com/vultr/terraform-provider-vultr/pull/157)
* Doc updates [152](https://github.com/vultr/terraform-provider-vultr/pull/152) [146](https://github.com/vultr/terraform-provider-vultr/pull/146) [147](https://github.com/vultr/terraform-provider-vultr/pull/147)

### Dependency
* updated terraform-plugin-sdk to v2.6.0 -> v2.7.0  [149](https://github.com/vultr/terraform-provider-vultr/pull/149)
* updated govultr to v2.5.1 -> v2.7.1  [150](https://github.com/vultr/terraform-provider-vultr/pull/150)


## 2.3.3 (June 25, 2021)
### Enhancement
* resource/instance: adding wait if a plan is being upgrade [144](https://github.com/vultr/terraform-provider-vultr/pull/144)

## 2.3.2 (June 10, 2021)
### Enhancement
* resource/instance: allow plan changes to do in-place upgrades [142](https://github.com/vultr/terraform-provider-vultr/pull/142)

## 2.3.1 (June 2, 2021)
### Bug Fix
* resource/bare_metal: fix type issue on `v6_network_size` [140](https://github.com/vultr/terraform-provider-vultr/pull/140)
* resource/bare_metal: fix missing `mac_address` definition in scheme [140](https://github.com/vultr/terraform-provider-vultr/pull/140)

## 2.3.0 (May 11, 2021)
### Enchancements
* resource/vultr_instances: allow the configuration of `backups_schedule` [134](https://github.com/vultr/terraform-provider-vultr/pull/134) [136](https://github.com/vultr/terraform-provider-vultr/pull/136)
* resource/vultr_load_balancers: add support for new LB features `private_network` and `firewall_rules` [137](https://github.com/vultr/terraform-provider-vultr/pull/137)
* resource/vultr_iso: support detaching during deletion  [131](https://github.com/vultr/terraform-provider-vultr/pull/131) Thanks @johnrichardrinehart
* resource/vultr_instances: `private_network_ids` are now tracked in statefile  [135](https://github.com/vultr/terraform-provider-vultr/pull/135)
* resource/vultr_block_storage: new field added `mount_id`  [135](https://github.com/vultr/terraform-provider-vultr/pull/135)
* resource/vultr_plans: new field added `disk_count`  [135](https://github.com/vultr/terraform-provider-vultr/pull/135)

### Dependency
* updated terraform-plugin-sdk to v2.4.0 -> v2.6.0  [134](https://github.com/vultr/terraform-provider-vultr/pull/134)
* updated govultr to v2.3.1 -> v2.5.1  [135](https://github.com/vultr/terraform-provider-vultr/pull/135)

## 2.2.0 (March 30, 2021)
### Feature
* Updated to Go 1.16 to support `darwin_arm64` [125](https://github.com/vultr/terraform-provider-vultr/pull/125)

## 2.1.4 (March 23, 2021)
### Bug Fix
* Fix issue with vultr_instance.reserved_ip_id and vultr_reserved_ip.attached_id conflicting [122](https://github.com/vultr/terraform-provider-vultr/pull/122)

## 2.1.3 (January 29, 2021)
### Dependency
* updated terraform-plugin-sdk to v1.8.0 -> v2.4.0  [111](https://github.com/vultr/terraform-provider-vultr/pull/111)

## 2.1.2 (January 05, 2021)
### Dependency
* updated GoVultr to v2.3.1 (fixes #102 #105) [106](https://github.com/vultr/terraform-provider-vultr/pull/106)

### Enhancements
* datasource/vultr_instance_ip4 & reverse_ipv4 improved filter and cleaned up docs [107](https://github.com/vultr/terraform-provider-vultr/pull/107)

## 2.1.1 (December 09, 2020)
### Enhancements
* resource/vultr_instances: Private-networks will be detached prior to deletion [93](https://github.com/vultr/terraform-provider-vultr/pull/93)
* resource/vultr_instances: Removal of `forcenew` on `activiation_email` [84](https://github.com/vultr/terraform-provider-vultr/pull/84)

## 2.1.0 (December 03, 2020)
### BUG FIXES
* resource/vultr_instances: In v2 the ID of the Reserved IP, not the IP itself, is required for creation. [79](https://github.com/vultr/terraform-provider-vultr/pull/79)

### Breaking Change
* resource/vultr_instances: Changing `reservered_ip` to `reservered_ip_id` to make it clear that the ID should be passed [79](https://github.com/vultr/terraform-provider-vultr/pull/79)

## 2.0.0 (December 01, 2020)

### Breaking Changes
* The TF Vultr provider v2.0.0 is a large change that uses the new Vultr API v2. This change resolves quite a few limitations and improves overall performance of tooling. These changes include field and resource/datasource name updates to match the API and offer a consistent experience.

### Dependency
* updated GoVultr to v2.1.0

## 1.5.0 (November 09, 2020)
### Breaking Change
* resource/vultr_server: Changing `user_data` will now trigger a `force_new` [70](https://github.com/vultr/terraform-provider-vultr/pull/70)

### Dependency
* updated GoVultr to v1.1.1 [70](https://github.com/vultr/terraform-provider-vultr/pull/70)

## 1.4.1 (September 15, 2020)
### BUG FIXES
* resource/vultr_server: Fix bug that did not allow user-data to be passed in as a string [65](https://github.com/vultr/terraform-provider-vultr/pull/65)

## 1.4.0 (August 28, 2020)
### FEATURES
* New Resource : vultr_server_ipv4 [61](https://github.com/vultr/terraform-provider-vultr/pull/61)
* New Resource : vultr_reverse_ipv4 [61](https://github.com/vultr/terraform-provider-vultr/pull/61)
* New Resource : vultr_reverse_ipv6 [20](https://github.com/vultr/terraform-provider-vultr/pull/20)
* New Data Source : vultr_server_ipv4 [61](https://github.com/vultr/terraform-provider-vultr/pull/61)
* New Data Source : vultr_reverse_ipv4 [61](https://github.com/vultr/terraform-provider-vultr/pull/61)
* New Data Source : vultr_reverse_ipv6 [20](https://github.com/vultr/terraform-provider-vultr/pull/20)

### IMPROVEMENTS
* resource/vultr_server: Ability to enable/disable DDOS post create [59](https://github.com/vultr/terraform-provider-vultr/pull/59)
* resource/vultr_server: Ability to detach ISO post create [60](https://github.com/vultr/terraform-provider-vultr/pull/60)

## 1.3.2 (June 17, 2020)
### IMPROVEMENTS
* resource/vultr_dns_record: New custom importer allows you to import DNS Records [51](https://github.com/vultr/terraform-provider-vultr/pull/51)
* resource/vultr_firewall_rule: New custom importer allows you to import Firewall Rules [52](https://github.com/vultr/terraform-provider-vultr/pull/52)

## 1.3.1 (June 11, 2020)
### IMPROVEMENTS
* resource/vultr_dns_domain: Making `server_ip` optional. If `server_ip` is not supplied terraform will now create the DNS resource with no records. [48](https://github.com/vultr/terraform-provider-vultr/pull/48)

## 1.3.0 (June 03, 2020)
### BUG FIXES
* resource/vultr_dns_record: Able to create record with `priority` of `0` [45](https://github.com/vultr/terraform-provider-vultr/pull/45)

### FEATURES
* New Resource : vultr_object_storage [41](https://github.com/vultr/terraform-provider-vultr/pull/41)
* New Data Source : vultr_object_storage [41](https://github.com/vultr/terraform-provider-vultr/pull/41)

## 1.2.0 (May 27, 2020)
### BUG FIXES
* Typo in `website/docs/index.html.markdown` [38](https://github.com/vultr/terraform-provider-vultr/pull/38)

### FEATURES
* New Resource : vultr_load_balancer [37](https://github.com/vultr/terraform-provider-vultr/pull/37)
* New Data Source : vultr_load_balancer [37](https://github.com/vultr/terraform-provider-vultr/pull/37)

## 1.1.5 (April 07, 2020)
### BUG FIXES
* resource/vultr_server: Detach ISO prior to deletion if instance was created with ISO. [34](https://github.com/vultr/terraform-provider-vultr/issues/34)

## 1.1.4 (March 30, 2020)
### IMPROVEMENTS
* resource/block_storage: Adding new optional param `live` to allow attaching/detaching of block storage to instances without restarts [31](https://github.com/vultr/terraform-provider-vultr/pull/31)

## 1.1.3 (March 24, 2020)
### BUG FIXES
* resource/reserved_ip: Adding `computed: true` to `attached_id` to prevent issues when Vultr assigns this [29](https://github.com/vultr/terraform-provider-vultr/pull/29)
* resource/vultr_server: Adding `ForceNew: true` to `reserved_ip` to prevent any issues where the main floating ip may get deleted and cause issues with the instance [29](https://github.com/vultr/terraform-provider-vultr/pull/29/files)

## 1.1.2 (March 19, 2020)
### IMPROVEMENTS
* resource/vultr_server: New optional field `reserved_ip` which lets you assign a `reserved_ip` during server creation [#26](https://github.com/vultr/terraform-provider-vultr/pull/26).
* resource/reserved_ip: During deletion any instances that are attached to the reserved_ip are detached [#27](https://github.com/vultr/terraform-provider-vultr/pull/27).
* Migrated to Terraform Plugin SDK [#21](https://github.com/vultr/terraform-provider-vultr/issues/21)
* docs/snapshot fixed typo in snapshot [#19](https://github.com/vultr/terraform-provider-vultr/pull/19)

## 1.1.1 (December 02, 2019)
### IMPROVEMENTS
* resource/vultr_block_storage: Attaches block storage on creation. Also reattaches block storage to instances if you taint the instance.[#9](https://github.com/vultr/terraform-provider-vultr/pull/9) Thanks @oogy!

## 1.1.0 (November 22, 2019)
### IMPROVEMENTS
*   provider: Retry mechanism configuration `retry_limit` was added to allow adjusting how many retries should be attempted. [#7](https://github.com/vultr/terraform-provider-vultr/pull/7).

### BUG FIXES
* Fixed go module name [#4](https://github.com/vultr/terraform-provider-vultr/pull/4)

## 1.0.5 (October 24, 2019)

* Initial release under the terraform-providers/ namespace

## [v1.0.4](https://github.com/vultr/terraform-provider-vultr/compare/v1.0.3..v1.0.4) (2019-08-09)
### Fixes
* Fixes issue where using a snapshot would cause drift [#96](https://github.com/vultr/terraform-provider-vultr/issues/96)
### Enhancements
* You will now not have to define the `os_id` for the following server options
    - `app_id`
    - `iso_id`
    - `snapshot_id`

## [v1.0.3](https://github.com/vultr/terraform-provider-vultr/compare/v1.0.2..v1.0.3) (2019-07-18)
### Fixes
* Fixes issue where you could not use `os_id` and `script_id` together [#92](https://github.com/vultr/terraform-provider-vultr/issues/92)
### Breaking Changes
* You will now need to provide the `os_id` on each corresponding option
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
* Fixed bug where scriptID was not being
properly handled in server creation [#82](https://github.com/vultr/terraform-provider-vultr/issues/82)
### Enhancements
* Added error handler on protocol case sensitivity [#83](https://github.com/vultr/terraform-provider-vultr/issues/83)
### Docs
* Typo in doc firewall_rule doc for protocol [#83](https://github.com/vultr/terraform-provider-vultr/issues/83)

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
