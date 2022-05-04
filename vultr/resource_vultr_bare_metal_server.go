package vultr

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
)

func resourceVultrBareMetalServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrBareMetalServerCreate,
		ReadContext:   resourceVultrBareMetalServerRead,
		UpdateContext: resourceVultrBareMetalServerUpdate,
		DeleteContext: resourceVultrBareMetalServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"plan": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"tag": {
				Type:       schema.TypeString,
				Optional:   true,
				Default:    "",
				Deprecated: "tag has been deprecated and should no longer be used. Instead, use tags",
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Default:  nil,
			},
			"script_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "",
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
				Default:  nil,
			},
			"user_data": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"activation_email": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "",
			},
			"os_id": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"app_id": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
				Optional: true,
			},
			"reserved_ipv4": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			// computed
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
			"v6_network": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"v6_main_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"v6_network_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"mac_address": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
		},
	}
}

func resourceVultrBareMetalServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appID, appOK := d.GetOk("app_id")
	osID, osOK := d.GetOk("os_id")
	imageID, imageOK := d.GetOk("image_id")
	snapshotID, snapshotOK := d.GetOk("snapshot_id")

	osOptions := map[string]bool{"os_id": osOK, "app_id": appOK, "snapshot_id": snapshotOK, "image_id": imageOK}
	osOption, err := bareMetalServerOSCheck(osOptions)
	if err != nil {
		return diag.FromErr(err)
	}

	keyIDs := make([]string, d.Get("ssh_key_ids.#").(int))
	for i, id := range d.Get("ssh_key_ids").([]interface{}) {
		keyIDs[i] = id.(string)
	}

	req := &govultr.BareMetalCreate{
		Region:          d.Get("region").(string),
		Plan:            d.Get("plan").(string),
		StartupScriptID: d.Get("script_id").(string),
		EnableIPv6:      govultr.BoolToBoolPtr(d.Get("enable_ipv6").(bool)),
		Label:           d.Get("label").(string),
		SSHKeyIDs:       keyIDs,
		UserData:        base64.StdEncoding.EncodeToString([]byte(d.Get("user_data").(string))),
		ActivationEmail: govultr.BoolToBoolPtr(d.Get("activation_email").(bool)),
		Hostname:        d.Get("hostname").(string),
		Tag:             d.Get("tag").(string),
		ReservedIPv4:    d.Get("reserved_ipv4").(string),
	}

	switch osOption {
	case "app_id":
		req.AppID = appID.(int)
	case "snapshot_id":
		req.SnapshotID = snapshotID.(string)
	case "os_id":
		req.OsID = osID.(int)
	case "image_id":
		req.ImageID = imageID.(string)
	}

	if tagsIDs, tagsOK := d.GetOk("tags"); tagsOK {
		for _, v := range tagsIDs.(*schema.Set).List() {
			req.Tags = append(req.Tags, v.(string))
		}
	}

	client := meta.(*Client).govultrClient()

	bm, err := client.BareMetalServer.Create(ctx, req)
	if err != nil {
		return diag.Errorf("error creating bare metal server: %v", err)
	}

	d.SetId(bm.ID)
	log.Printf("[INFO] Bare Metal Server ID: %s", d.Id())

	if _, err = waitForBareMetalServerActiveStatus(ctx, d, meta); err != nil {
		return diag.Errorf("error while waiting for bare metal server (%s) to be in active state: %s", d.Id(), err)
	}

	return resourceVultrBareMetalServerRead(ctx, d, meta)
}

func resourceVultrBareMetalServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	bms, err := client.BareMetalServer.Get(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "Invalid server") {
			log.Printf("[WARN] Removing bare metal server %s because it is gone", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting bare metal server: %v", err)
	}

	d.SetId(bms.ID)
	d.Set("os", bms.Os)
	d.Set("ram", bms.RAM)
	d.Set("disk", bms.Disk)
	d.Set("main_ip", bms.MainIP)
	d.Set("cpu_count", bms.CPUCount)
	d.Set("region", bms.Region)
	d.Set("default_password", bms.DefaultPassword)
	d.Set("date_created", bms.DateCreated)
	d.Set("status", bms.Status)
	d.Set("netmask_v4", bms.NetmaskV4)
	d.Set("gateway_v4", bms.GatewayV4)
	d.Set("plan", bms.Plan)
	d.Set("label", bms.Label)
	d.Set("tag", bms.Tag)
	d.Set("tags", bms.Tags)
	d.Set("mac_address", bms.MacAddress)
	d.Set("os_id", bms.OsID)
	d.Set("app_id", bms.AppID)
	d.Set("image_id", bms.ImageID)
	d.Set("v6_network", bms.V6Network)
	d.Set("v6_main_ip", bms.V6MainIP)
	d.Set("v6_network_size", bms.V6NetworkSize)

	return nil
}

func resourceVultrBareMetalServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	req := &govultr.BareMetalUpdate{
		Label:      d.Get("label").(string),
		Tag:        d.Get("tag").(string),
		Tags:       []string{},
		EnableIPv6: govultr.BoolToBoolPtr(d.Get("enable_ipv6").(bool)),
	}

	if d.HasChange("app_id") {
		log.Printf(`[INFO] Changing bare metal server (%s) application`, d.Id())
		_, newVal := d.GetChange("app_id")

		appID := newVal.(int)
		req.AppID = appID
	}

	if d.HasChange("os_id") {
		log.Printf(`[INFO] Changing bare metal server (%s) operating system`, d.Id())
		_, newVal := d.GetChange("os_id")

		osID := newVal.(int)
		req.OsID = osID
	}

	if d.HasChange("tags") {
		oldTags, newTags := tfChangeToSlices("tags", d)
		for _, v := range diffSlice(oldTags, newTags) {
			req.Tags = append(req.Tags, v)
		}
	}

	if _, err := client.BareMetalServer.Update(ctx, d.Id(), req); err != nil {
		return diag.Errorf("error updating bare metal %s : %s", d.Id(), err.Error())
	}

	return resourceVultrBareMetalServerRead(ctx, d, meta)
}

func resourceVultrBareMetalServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting bare metal server: %s", d.Id())
	if err := client.BareMetalServer.Delete(ctx, d.Id()); err != nil {
		return diag.Errorf("error deleting bare metal server (%s): %v", d.Id(), err)
	}

	return nil
}

func bareMetalServerOSCheck(options map[string]bool) (string, error) {
	var result []string
	for k, v := range options {
		if v == true {
			result = append(result, k)
		}
	}

	if len(result) > 1 {
		return "", fmt.Errorf("too many OS options have been selected: %v - please select one", result)
	}
	if len(result) == 0 {
		return "", errors.New("you must set one of the following: os_id, app_id, or snapshot_id")
	}

	return result[0], nil
}

func waitForBareMetalServerActiveStatus(ctx context.Context, d *schema.ResourceData, meta interface{}) (interface{}, error) {
	log.Printf("[INFO] Waiting for bare metal server (%s) to have status of active", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"active"},
		Refresh:    newBareMetalServerStatusStateRefresh(ctx, d, meta),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,

		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func newBareMetalServerStatusStateRefresh(ctx context.Context, d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*Client).govultrClient()

	return func() (interface{}, string, error) {
		bms, err := client.BareMetalServer.Get(ctx, d.Id())
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving bare metal server %s : %s", d.Id(), err)
		}

		log.Printf("[INFO] Bare metal server (%s) status: %s", d.Id(), bms.Status)
		return bms, bms.Status, nil
	}
}
