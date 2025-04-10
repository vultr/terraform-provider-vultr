package vultr

import (
	"encoding/base64"

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

func nodePoolSchema(isNodePool bool) map[string]*schema.Schema {
	s := map[string]*schema.Schema{
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
		//computed fields
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
	}

	if isNodePool {
		s["cluster_id"] = &schema.Schema{
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
			ForceNew:     true,
		}
		s["tag"] = &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		}
	} else {
		// Make tags unmodifiable for the vultr_kubernetes resource
		// This lets us know which node pool was part of the vultr_kubernetes resource
		s["tag"] = &schema.Schema{
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
