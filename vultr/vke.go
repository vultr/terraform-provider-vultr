package vultr

import (
	"encoding/base64"
	"maps"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"gopkg.in/yaml.v2"
)

type KubeConfig struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Clusters   []struct {
		Name    string `yaml:"name"`
		Cluster struct {
			CaCert string `yaml:"certificate-authority-data"`
			Server string `yaml:"server"`
		} `yaml:"cluster"`
	} `yaml:"clusters"`
	Users []struct {
		Name string `yaml:"name"`
		User struct {
			ClientCert string `yaml:"client-certificate-data"`
			ClientKey  string `yaml:"client-key-data"`
		} `yaml:"user"`
	} `yaml:"users"`
}

func resourceVultrKubernetesV1() map[string]*schema.Schema {
	schemaV0 := resourceVultrKubernetesV0().Schema
	schemaV1 := map[string]*schema.Schema{}
	maps.Copy(schemaV1, schemaV0)

	schemaV1["node_pools"].Elem.(*schema.Resource).Schema = resourceVultrKubernetesNodePoolsV1(false)

	return schemaV1
}

func resourceVultrKubernetesNodePoolsV1(isNodePool bool) map[string]*schema.Schema {
	schemaV0 := resourceVultrKubernetesNodePoolsV0(isNodePool).Schema
	schemaV1 := map[string]*schema.Schema{}

	schemaLabels := map[string]*schema.Schema{
		"labels": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Required: true,
					},
					"value": {
						Type:     schema.TypeString,
						Required: true,
					},
					"id": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}

	schemaTaints := map[string]*schema.Schema{
		"taints": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Required: true,
					},
					"value": {
						Type:     schema.TypeString,
						Required: true,
					},
					"effect": {
						Type:     schema.TypeString,
						Required: true,
					},
					"id": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}

	maps.Copy(schemaV0, schemaLabels)
	maps.Copy(schemaV0, schemaTaints)
	maps.Copy(schemaV1, schemaV0)

	return schemaV1
}

func resourceVultrKubernetesV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				DiffSuppressFunc: IgnoreCase,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ha_controlplanes": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"enable_firewall": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"node_pools": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     resourceVultrKubernetesNodePoolsV0(false),
			},
			// cluster computed fields
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"firewall_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kube_config": {
				Description: "Base64 encoded KubeConfig",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
			"cluster_ca_certificate": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"client_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"client_certificate": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceVultrKubernetesNodePoolsV0(isNodePool bool) *schema.Resource {
	s := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"plan": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_quantity": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntAtLeast(1),
				Required:     true,
			},
			"auto_scaler": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"min_nodes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"max_nodes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"taints": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"effect": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// computed fields
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"date_created": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}

	if isNodePool {
		s.Schema["cluster_id"] = &schema.Schema{
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
			ForceNew:     true,
		}
		s.Schema["tag"] = &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		}
	} else {
		// Make tags unmodifiable for the vultr_kubernetes resource
		// This lets us know which node pool was part of the vultr_kubernetes resource
		s.Schema["tag"] = &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		}
	}

	return s
}

func getCertsFromKubeConfig(kubeconfig string) (ca string, cert string, key string, err error) {
	decodedKC, err := base64.StdEncoding.DecodeString(kubeconfig)
	if err != nil {
		return "", "", "", err
	}

	var kc KubeConfig

	err = yaml.Unmarshal(decodedKC, &kc)
	if err != nil {
		return "", "", "", err
	}

	return kc.Clusters[0].Cluster.CaCert, kc.Users[0].User.ClientCert, kc.Users[0].User.ClientKey, nil
}
