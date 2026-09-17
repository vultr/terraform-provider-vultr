package vultr

import (
	"context"
	"encoding/base64"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"go.yaml.in/yaml/v4"
)

type kubeConfig struct {
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

func dataSourceVultrKubernetesKubeConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrKubernetesKubeConfigRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"kube_config": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
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

func dataSourceVultrKubernetesKubeConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	config, _, err := client.Kubernetes.GetKubeConfig(ctx, d.Get("cluster_id").(string))
	if err != nil {
		return diag.Errorf("error getting kubernetes kubeconfig : %v", err)
	}

	ca, cert, key, err := getCertsFromKubeConfig(config.KubeConfig)
	if err != nil {
		return diag.Errorf("error reading kubernetes kubeconfig : %v", err)
	}

	if err := d.Set("kube_config", config.KubeConfig); err != nil {
		return diag.Errorf("unable to set kubernetes kubeconfig `kube_config` read value : %v", err)
	}

	if err := d.Set("cluster_ca_certificate", ca); err != nil {
		return diag.Errorf("unable to set kubernetes kubeconfig `cluster_ca_certificate` read value : %v", err)
	}

	if err := d.Set("client_certificate", cert); err != nil {
		return diag.Errorf("unable to set kubernetes kubeconfig `client_certificate` read value : %v", err)
	}

	if err := d.Set("client_key", key); err != nil {
		return diag.Errorf("unable to set kubernetes `client_key` read value : %v", err)
	}

	d.SetId(d.Get("cluster_id").(string))

	return nil
}

func getCertsFromKubeConfig(kubeconfig string) (ca string, cert string, key string, err error) {
	decodedKC, err := base64.StdEncoding.DecodeString(kubeconfig)
	if err != nil {
		return "", "", "", err
	}

	config := kubeConfig{}
	err = yaml.Unmarshal(decodedKC, &config)
	if err != nil {
		return "", "", "", err
	}

	return config.Clusters[0].Cluster.CaCert, config.Users[0].User.ClientCert, config.Users[0].User.ClientKey, nil
}
