package vultr

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vultr/govultr"
)

const (
	appOSID      = "186"
	snapshotOSID = "164"
)

func resourceVultrBareMetalServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrBareMetalServerCreate,
		Read:   resourceVultrBareMetalServerRead,
		Update: resourceVultrBareMetalServerUpdate,
		Delete: resourceVultrBareMetalServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"os": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ram": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"main_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"default_password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"netmask_v4": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway_v4": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plan_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"v6_networks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"os_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"app_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			// options
			"startup_script_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enable_ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"ssh_key_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"userdata": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notify_activate": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVultrBareMetalServerCreate(d *schema.ResourceData, meta interface{}) error {
	appID, appOK := d.GetOk("app_id")
	osID, osOK := d.GetOk("os_id")
	snapshotID, snapshotOK := d.GetOk("snapshot_id")

	osOptions := map[string]bool{"os_id": osOK, "app_id": appOK, "snapshot_id": snapshotOK}
	osOption, err := bareMetalServerOSCheck(osOptions)
	if err != nil {
		return err
	}

	enableIPV6 := d.Get("enable_ipv6")
	ipv6 := "no"
	if enableIPV6.(bool) == true {
		ipv6 = "yes"
	}

	notifyActivate := d.Get("notify_activate")
	notify := "no"
	if notifyActivate.(bool) == true {
		notify = "yes"
	}

	keyIDs := make([]string, d.Get("ssh_key_ids.#").(int))
	for i, id := range d.Get("ssh_key_ids").([]interface{}) {
		keyIDs[i] = id.(string)
	}

	options := &govultr.BareMetalServerOptions{
		StartupScriptID: d.Get("startup_script_id").(string),
		EnableIPV6:      ipv6,
		Label:           d.Get("label").(string),
		SSHKeyIDs:       keyIDs,
		UserData:        d.Get("userdata").(string),
		NotifyActivate:  notify,
		Hostname:        d.Get("hostname").(string),
		Tag:             d.Get("tag").(string),
	}

	switch osOption {
	case "app_id":
		options.AppID = appID.(string)
		osID = appOSID
	case "snapshot_id":
		options.SnapshotID = snapshotID.(string)
		osID = snapshotOSID
	}

	client := meta.(*Client).govultrClient()

	regionID := d.Get("region_id").(string)
	planID := d.Get("plan_id").(string)

	bm, err := client.BareMetalServer.Create(context.Background(), regionID, planID, osID.(string), options)
	if err != nil {
		return fmt.Errorf("Error creating bare metal server: %v", err)
	}

	d.SetId(bm.BareMetalServerID)
	log.Printf("[INFO] Bare Metal Server ID: %s", d.Id())

	_, err = waitForBareMetalServerActiveStatus(d, meta)
	if err != nil {
		return fmt.Errorf("Error while waiting for bare metal server (%s) to be in active state: %s", d.Id(), err)
	}

	return resourceVultrBareMetalServerRead(d, meta)
}

func resourceVultrBareMetalServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	bms, err := client.BareMetalServer.GetServer(context.Background(), d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Invalid server") {
			log.Printf("[WARN] Removing bare metal server %s because it is gone", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting bare metal server: %v", err)
	}

	d.Set("os", bms.Os)
	d.Set("ram", bms.RAM)
	d.Set("disk", bms.Disk)
	d.Set("main_ip", bms.MainIP)
	d.Set("cpu_count", bms.CPUCount)
	d.Set("region_id", bms.RegionID)
	d.Set("default_password", bms.DefaultPassword)
	d.Set("date_created", bms.DateCreated)
	d.Set("status", bms.Status)
	d.Set("netmask_v4", bms.NetmaskV4)
	d.Set("gateway_v4", bms.GatewayV4)
	d.Set("plan_id", bms.BareMetalPlanID)
	d.Set("label", bms.Label)
	d.Set("tag", bms.Tag)
	d.Set("os_id", bms.OsID)
	d.Set("app_id", bms.AppID)

	var ipv6s []map[string]string
	for _, net := range bms.V6Networks {
		v6network := map[string]string{
			"v6_network":      net.Network,
			"v6_main_ip":      net.MainIP,
			"v6_network_size": net.NetworkSize,
		}
		ipv6s = append(ipv6s, v6network)
	}
	d.Set("v6_networks", ipv6s)

	return nil
}

func resourceVultrBareMetalServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	d.Partial(true)

	if d.HasChange("app_id") {
		log.Printf(`[INFO] Changing bare metal server (%s) application`, d.Id())
		_, newVal := d.GetChange("app_id")
		err := client.BareMetalServer.ChangeApp(context.Background(), d.Id(), newVal.(string))
		if err != nil {
			return fmt.Errorf("Error changing bare metal server (%s) application: %v", d.Id(), err)
		}
		_, err = waitForBareMetalServerActiveStatus(d, meta)
		if err != nil {
			return fmt.Errorf("Error while waiting for bare metal server (%s) to be in active state: %s", d.Id(), err)
		}
		d.SetPartial("app_id")
	}

	if d.HasChange("label") {
		log.Printf(`[INFO] Updating bare metal server label (%s)`, d.Id())
		_, newVal := d.GetChange("label")
		err := client.BareMetalServer.SetLabel(context.Background(), d.Id(), newVal.(string))
		if err != nil {
			return fmt.Errorf("Error updating bare metal server label (%s): %v", d.Id(), err)
		}
		d.SetPartial("label")
	}

	if d.HasChange("os_id") {
		log.Printf(`[INFO] Changing bare metal server (%s) operating system`, d.Id())
		_, newVal := d.GetChange("os_id")
		err := client.BareMetalServer.ChangeOS(context.Background(), d.Id(), newVal.(string))
		if err != nil {
			return fmt.Errorf("Error changing bare metal server (%s) operating system: %v", d.Id(), err)
		}
		_, err = waitForBareMetalServerActiveStatus(d, meta)
		if err != nil {
			return fmt.Errorf("Error while waiting for bare metal server (%s) to be in active state: %s", d.Id(), err)
		}
		d.SetPartial("os_id")
	}

	if d.HasChange("tag") {
		log.Printf(`[INFO] Updating bare metal server (%s) tag`, d.Id())
		_, newVal := d.GetChange("tag")
		err := client.BareMetalServer.SetTag(context.Background(), d.Id(), newVal.(string))
		if err != nil {
			return fmt.Errorf("Error updating bare metal server (%s) tag: %v", d.Id(), err)
		}
		d.SetPartial("tag")
	}

	if d.HasChange("userdata") {
		log.Printf(`[INFO] Updating bare metal server (%s) userdata`, d.Id())
		_, newVal := d.GetChange("userdata")
		err := client.BareMetalServer.SetUserData(context.Background(), d.Id(), newVal.(string))
		if err != nil {
			return fmt.Errorf("Error updating bare metal server (%s) userdata: %v", d.Id(), err)
		}
		d.SetPartial("userdata")
	}

	d.Partial(false)

	return resourceVultrBareMetalServerRead(d, meta)
}

func resourceVultrBareMetalServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting bare metal server: %s", d.Id())
	err := client.BareMetalServer.Destroy(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting bare metal server (%s): %v", d.Id(), err)
	}

	return nil
}

func bareMetalServerOSCheck(options map[string]bool) (string, error) {
	result := []string{}
	for k, v := range options {
		if v == true {
			result = append(result, k)
		}
	}

	if len(result) > 1 {
		return "", fmt.Errorf("Too many OS options have been selected: %v - please select one", result)
	}
	if len(result) == 0 {
		return "", errors.New("You must set one of the following: os_id, app_id, or snapshot_id")
	}

	return result[0], nil
}

func waitForBareMetalServerActiveStatus(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	log.Printf("[INFO] Waiting for bare metal server (%s) to have status of active", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"active"},
		Refresh:    newBareMetalServerStatusStateRefresh(d, meta),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,

		NotFoundChecks: 60,
	}

	return stateConf.WaitForState()
}

func newBareMetalServerStatusStateRefresh(d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*Client).govultrClient()

	return func() (interface{}, string, error) {
		bms, err := client.BareMetalServer.GetServer(context.Background(), d.Id())

		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving bare metal server %s : %s", d.Id(), err)
		}

		log.Printf("[INFO] Bare metal server (%s) status: %s", d.Id(), bms.Status)
		return bms, bms.Status, nil
	}
}
