# Changes from V1 to V2

The v2 Vultr Terraform plugin is backed by our new [V2 API](https://www.vultr.com/api/) and introduces a several key changes. 

#### 1.  All resource IDs are now uuids and references to the previous ID format are not utilized. This affects the following resources:
- Instances
- Bare Metal
- Block Storage
- Object Storage
- Reserved IP
- Load Balancer
- Snapshots & Backups
- DNS Records
- ISOs - Private & Public
- Private Networks
- Startup Scripts
- SSH Keys
- Firewall Groups

#### 2. Newly created resources will return a uuid and you must use uuids to interact with the v2 plugin
- Data Sources, when filtered by ID or any other resource ID, must use the uuid of the resource
- All resources that accept other resources as arguments must provide the uuid (e.g, snapshot_id when creating a new instance)

#### 3. Plan IDs have a new, more descriptive format and must be used, e.g  "vc2-1c-1gb", "vbm-4c-32gb". See [V2 Plans](https://api.vultr.com/v2/plans) for a complete list.

#### 4. Region IDs have a new, more descriptive format and must be used. E.g, "ewr", "sea", "atl". See [V2 Regions](https://api.vultr.com/v2/regions) for a complete list.

#### 5. Arguments added, removed and name changes throughout to accommodate new features within the V2 API. [Review the Terraform Docs](https://registry.terraform.io/providers/vultr/vultr/latest/docs) for more details.
