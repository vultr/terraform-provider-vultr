package vultr
//
//import (
//	"context"
//	"errors"
//	"fmt"
//	"log"
//	"strings"
//	"time"
//
//	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
//	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
//	"github.com/vultr/govultr/v2"
//)
//
//func resourceVultrBareMetalServer() *schema.Resource {
//	return &schema.Resource{
//		Create: resourceVultrBareMetalServerCreate,
//		Read:   resourceVultrBareMetalServerRead,
//		Update: resourceVultrBareMetalServerUpdate,
//		Delete: resourceVultrBareMetalServerDelete,
//		Importer: &schema.ResourceImporter{
//			State: schema.ImportStatePassthrough,
//		},
//
//		Schema: map[string]*schema.Schema{
//			"region": {
//				Type:     schema.TypeString,
//				Required: true,
//				ForceNew: true,
//			},
//			"plan": {
//				Type:     schema.TypeString,
//				Required: true,
//				ForceNew: true,
//			},
//			"label": {
//				Type:     schema.TypeString,
//				Optional: true,
//				Default:  "",
//			},
//			"tag": {
//				Type:     schema.TypeString,
//				Optional: true,
//				Default:  "",
//			},
//			"script_id": {
//				Type:     schema.TypeString,
//				Optional: true,
//				ForceNew: true,
//				Default:  "",
//			},
//			"snapshot_id": {
//				Type:     schema.TypeString,
//				Optional: true,
//				ForceNew: true,
//			},
//			"enable_ipv6": {
//				Type:     schema.TypeBool,
//				Optional: true,
//				ForceNew: true,
//				Default:  false,
//			},
//			"ssh_key_ids": {
//				Type:     schema.TypeList,
//				Optional: true,
//				ForceNew: true,
//				Elem:     &schema.Schema{Type: schema.TypeString},
//				Default:  nil,
//			},
//			"user_data": {
//				Type:     schema.TypeString,
//				Optional: true,
//				Default:  "",
//			},
//			"activation_email": {
//				Type:     schema.TypeBool,
//				Optional: true,
//				ForceNew: true,
//				Default:  true,
//			},
//			"hostname": {
//				Type:     schema.TypeString,
//				Optional: true,
//				ForceNew: true,
//				Default:  "",
//			},
//			"os_id": {
//				Type:     schema.TypeInt,
//				Computed: true,
//				Optional: true,
//			},
//			"app_id": {
//				Type:     schema.TypeInt,
//				Computed: true,
//				Optional: true,
//			},
//
//			// computed
//			"os": {
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//			"ram": {
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//			"disk": {
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//			"main_ip": {
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//			"cpu_count": {
//				Type:     schema.TypeInt,
//				Computed: true,
//			},
//			"default_password": {
//				Type:      schema.TypeString,
//				Computed:  true,
//				Sensitive: true,
//			},
//			"date_created": {
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//			"status": {
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//			"netmask_v4": {
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//			"gateway_v4": {
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//			"v6_network": {
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//			"v6_main_ip": {
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//			"v6_network_size": {
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//		},
//	}
//}
//
//func resourceVultrBareMetalServerCreate(d *schema.ResourceData, meta interface{}) error {
//	appID, appOK := d.GetOk("app_id")
//	osID, osOK := d.GetOk("os_id")
//	snapshotID, snapshotOK := d.GetOk("snapshot_id")
//
//	osOptions := map[string]bool{"os_id": osOK, "app_id": appOK, "snapshot_id": snapshotOK}
//	osOption, err := bareMetalServerOSCheck(osOptions)
//	if err != nil {
//		return err
//	}
//
//	keyIDs := make([]string, d.Get("ssh_key_ids.#").(int))
//	for i, id := range d.Get("ssh_key_ids").([]interface{}) {
//		keyIDs[i] = id.(string)
//	}
//
//	req := &govultr.BareMetalReq{
//		Region:          d.Get("region").(string),
//		Plan:            d.Get("plan").(string),
//		StartupScriptID: d.Get("script_id").(string),
//		EnableIPv6:      d.Get("enable_ipv6").(bool),
//		Label:           d.Get("label").(string),
//		SSHKeyIDs:       keyIDs,
//		UserData:        d.Get("user_data").(string),
//		ActivationEmail: d.Get("activation_email").(bool),
//		Hostname:        d.Get("hostname").(string),
//		Tag:             d.Get("tag").(string),
//		//ReservedIPv4:   "",
//	}
//	switch osOption {
//	case "app_id":
//		req.AppID = appID.(int)
//	case "snapshot_id":
//		req.SnapshotID = snapshotID.(string)
//	case "os_id":
//		req.OsID = osID.(int)
//	}
//
//	client := meta.(*Client).govultrClient()
//
//	bm, err := client.BareMetalServer.Create(context.Background(), req)
//	if err != nil {
//		return fmt.Errorf("error creating bare metal server: %v", err)
//	}
//
//	d.SetId(bm.ID)
//	log.Printf("[INFO] Bare Metal Server ID: %s", d.Id())
//
//	if _, err = waitForBareMetalServerActiveStatus(d, meta); err != nil {
//		return fmt.Errorf("error while waiting for bare metal server (%s) to be in active state: %s", d.Id(), err)
//	}
//
//	return resourceVultrBareMetalServerRead(d, meta)
//}
//
//func resourceVultrBareMetalServerRead(d *schema.ResourceData, meta interface{}) error {
//	client := meta.(*Client).govultrClient()
//
//	bms, err := client.BareMetalServer.Get(context.Background(), d.Id())
//	if err != nil {
//		if strings.Contains(err.Error(), "Invalid server") {
//			log.Printf("[WARN] Removing bare metal server %s because it is gone", d.Id())
//			d.SetId("")
//			return nil
//		}
//		return fmt.Errorf("error getting bare metal server: %v", err)
//	}
//
//	d.Set("os", bms.Os)
//	d.Set("ram", bms.RAM)
//	d.Set("disk", bms.Disk)
//	d.Set("main_ip", bms.MainIP)
//	d.Set("cpu_count", bms.CPUCount)
//	d.Set("region", bms.Region)
//	d.Set("default_password", bms.DefaultPassword)
//	d.Set("date_created", bms.DateCreated)
//	d.Set("status", bms.Status)
//	d.Set("netmask_v4", bms.NetmaskV4)
//	d.Set("gateway_v4", bms.GatewayV4)
//	d.Set("plan", bms.Plan)
//	d.Set("label", bms.Label)
//	d.Set("tag", bms.Tag)
//	d.Set("os_id", bms.OsID)
//	d.Set("app_id", bms.AppID)
//	d.Set("v6_network", bms.V6Network)
//	d.Set("v6_main_ip", bms.V6MainIP)
//	d.Set("v6_network_size", bms.V6NetworkSize)
//
//	return nil
//}
//
//func resourceVultrBareMetalServerUpdate(d *schema.ResourceData, meta interface{}) error {
//	//client := meta.(*Client).govultrClient()
//	//
//	//d.Partial(true)
//	//
//	//if d.HasChange("app_id") {
//	//	log.Printf(`[INFO] Changing bare metal server (%s) application`, d.Id())
//	//	_, newVal := d.GetChange("app_id")
//	//	err := client.BareMetalServer.ChangeApp(context.Background(), d.Id(), newVal.(string))
//	//	if err != nil {
//	//		return fmt.Errorf("Error changing bare metal server (%s) application: %v", d.Id(), err)
//	//	}
//	//	_, err = waitForBareMetalServerActiveStatus(d, meta)
//	//	if err != nil {
//	//		return fmt.Errorf("Error while waiting for bare metal server (%s) to be in active state: %s", d.Id(), err)
//	//	}
//	//	d.SetPartial("app_id")
//	//}
//	//
//	//if d.HasChange("label") {
//	//	log.Printf(`[INFO] Updating bare metal server label (%s)`, d.Id())
//	//	_, newVal := d.GetChange("label")
//	//	err := client.BareMetalServer.SetLabel(context.Background(), d.Id(), newVal.(string))
//	//	if err != nil {
//	//		return fmt.Errorf("Error updating bare metal server label (%s): %v", d.Id(), err)
//	//	}
//	//	d.SetPartial("label")
//	//}
//	//
//	//if d.HasChange("os_id") {
//	//	log.Printf(`[INFO] Changing bare metal server (%s) operating system`, d.Id())
//	//	_, newVal := d.GetChange("os_id")
//	//	err := client.BareMetalServer.ChangeOS(context.Background(), d.Id(), newVal.(string))
//	//	if err != nil {
//	//		return fmt.Errorf("Error changing bare metal server (%s) operating system: %v", d.Id(), err)
//	//	}
//	//	_, err = waitForBareMetalServerActiveStatus(d, meta)
//	//	if err != nil {
//	//		return fmt.Errorf("Error while waiting for bare metal server (%s) to be in active state: %s", d.Id(), err)
//	//	}
//	//	d.SetPartial("os_id")
//	//}
//	//
//	//if d.HasChange("tag") {
//	//	log.Printf(`[INFO] Updating bare metal server (%s) tag`, d.Id())
//	//	_, newVal := d.GetChange("tag")
//	//	err := client.BareMetalServer.SetTag(context.Background(), d.Id(), newVal.(string))
//	//	if err != nil {
//	//		return fmt.Errorf("Error updating bare metal server (%s) tag: %v", d.Id(), err)
//	//	}
//	//	d.SetPartial("tag")
//	//}
//	//
//	//if d.HasChange("user_data") {
//	//	log.Printf(`[INFO] Updating bare metal server (%s) user_data`, d.Id())
//	//	_, newVal := d.GetChange("user_data")
//	//	err := client.BareMetalServer.SetUserData(context.Background(), d.Id(), newVal.(string))
//	//	if err != nil {
//	//		return fmt.Errorf("Error updating bare metal server (%s) user_data: %v", d.Id(), err)
//	//	}
//	//	d.SetPartial("user_data")
//	//}
//	//
//	//d.Partial(false)
//
//	//return resourceVultrBareMetalServerRead(d, meta)
//	return nil
//}
//
//func resourceVultrBareMetalServerDelete(d *schema.ResourceData, meta interface{}) error {
//	client := meta.(*Client).govultrClient()
//
//	log.Printf("[INFO] Deleting bare metal server: %s", d.Id())
//	if err := client.BareMetalServer.Delete(context.Background(), d.Id()); err != nil {
//		return fmt.Errorf("error deleting bare metal server (%s): %v", d.Id(), err)
//	}
//
//	return nil
//}
//
//func bareMetalServerOSCheck(options map[string]bool) (string, error) {
//	var result []string
//	for k, v := range options {
//		if v == true {
//			result = append(result, k)
//		}
//	}
//
//	if len(result) > 1 {
//		return "", fmt.Errorf("too many OS options have been selected: %v - please select one", result)
//	}
//	if len(result) == 0 {
//		return "", errors.New("you must set one of the following: os_id, app_id, or snapshot_id")
//	}
//
//	return result[0], nil
//}
//
//func waitForBareMetalServerActiveStatus(d *schema.ResourceData, meta interface{}) (interface{}, error) {
//	log.Printf("[INFO] Waiting for bare metal server (%s) to have status of active", d.Id())
//
//	stateConf := &resource.StateChangeConf{
//		Pending:    []string{"pending"},
//		Target:     []string{"active"},
//		Refresh:    newBareMetalServerStatusStateRefresh(d, meta),
//		Timeout:    60 * time.Minute,
//		Delay:      10 * time.Second,
//		MinTimeout: 3 * time.Second,
//
//		NotFoundChecks: 60,
//	}
//
//	return stateConf.WaitForState()
//}
//
//func newBareMetalServerStatusStateRefresh(d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
//	client := meta.(*Client).govultrClient()
//
//	return func() (interface{}, string, error) {
//		bms, err := client.BareMetalServer.Get(context.Background(), d.Id())
//		if err != nil {
//			return nil, "", fmt.Errorf("error retrieving bare metal server %s : %s", d.Id(), err)
//		}
//
//		log.Printf("[INFO] Bare metal server (%s) status: %s", d.Id(), bms.Status)
//		return bms, bms.Status, nil
//	}
//}
