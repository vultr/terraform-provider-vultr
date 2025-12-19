package vultr

import (
	"context"
	"encoding/base64"
	"fmt"
	"maps"
	"slices"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vultr/govultr/v3"
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

func updateNodePoolOptions(ctx context.Context, client *govultr.Client, clusterID, nodePoolID, optionKind string, oldData, newData []interface{}) error { //nolint:lll
	type optionData struct {
		Create   bool
		Delete   bool
		OptionID string
		Key      string
		Value    string
		Effect   string
	}

	oldOptionData := []optionData{}
	optionRequests := []optionData{}
	for i := range oldData {
		oldOption := oldData[i].(map[string]interface{})

		oldData := optionData{
			OptionID: oldOption["id"].(string),
			Key:      oldOption["key"].(string),
		}

		oldOptionData = append(oldOptionData, oldData)
	}

	for i := range newData {
		newOption := newData[i].(map[string]interface{})

		newRequest := optionData{
			Key:    newOption["key"].(string),
			Value:  newOption["value"].(string),
			Create: true,
			Delete: false,
		}

		if optionKind == "taints" {
			newRequest.Effect = newOption["effect"].(string)
		}

		oldIndex := slices.IndexFunc(oldOptionData, func(o optionData) bool { return o.Key == newRequest.Key })

		if oldIndex >= 0 {
			// delete the old option in the process
			newRequest.Delete = true
			newRequest.OptionID = oldOptionData[oldIndex].OptionID
		}

		optionRequests = append(optionRequests, newRequest)
	}

	// mark delete options not in changed data
	for i := range oldOptionData {
		if !slices.ContainsFunc(optionRequests, func(o optionData) bool { return o.Key == oldOptionData[i].Key }) {
			optionRequests = append(optionRequests, optionData{
				OptionID: oldOptionData[i].OptionID,
				Delete:   true,
				Create:   false,
			})
		}
	}

	for i := range optionRequests {
		switch optionKind {
		case "labels":
			if optionRequests[i].Delete {
				err := client.Kubernetes.DeleteNodePoolLabel(
					ctx,
					clusterID,
					nodePoolID,
					optionRequests[i].OptionID,
				)
				if err != nil {
					return fmt.Errorf(
						"error deleting label %q from vke %q node pool %q during option update : %w",
						optionRequests[i].Key,
						clusterID,
						nodePoolID,
						err,
					)
				}
			}

			if optionRequests[i].Create {
				_, _, err := client.Kubernetes.CreateNodePoolLabel(
					ctx,
					clusterID,
					nodePoolID,
					&govultr.NodePoolLabelReq{
						Key:   optionRequests[i].Key,
						Value: optionRequests[i].Value,
					},
				)
				if err != nil {
					return fmt.Errorf(
						"error creating label %q on vke %q node pool %q during option update : %w",
						optionRequests[i].Key,
						clusterID,
						nodePoolID,
						err,
					)
				}
			}
		case "taints":
			if optionRequests[i].Delete {
				err := client.Kubernetes.DeleteNodePoolTaint(
					ctx,
					clusterID,
					nodePoolID,
					optionRequests[i].OptionID,
				)
				if err != nil {
					return fmt.Errorf(
						"error deleting taint %q from vke %q node pool %q during option update : %w",
						optionRequests[i].Key,
						clusterID,
						nodePoolID,
						err,
					)
				}
			}

			if optionRequests[i].Create {
				_, _, err := client.Kubernetes.CreateNodePoolTaint(
					ctx,
					clusterID,
					nodePoolID,
					&govultr.NodePoolTaintReq{
						Key:    optionRequests[i].Key,
						Value:  optionRequests[i].Value,
						Effect: optionRequests[i].Effect,
					},
				)
				if err != nil {
					return fmt.Errorf(
						"error creating label %q on vke %q node pool %q during update : %w",
						optionRequests[i].Key,
						clusterID,
						nodePoolID,
						err,
					)
				}
			}
		}
	}

	return nil
}
