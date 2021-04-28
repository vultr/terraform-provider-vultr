package vultr

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/vultr/govultr/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVultrIsoPrivate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrIsoCreate,
		ReadContext:   resourceVultrIsoRead,
		DeleteContext: resourceVultrIsoDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filename": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"md5sum": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sha512sum": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVultrIsoCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Creating new ISO")

	isoReq := &govultr.ISOReq{URL: d.Get("url").(string)}
	iso, err := client.ISO.Create(ctx, isoReq)
	if err != nil {
		return diag.Errorf("Error creating ISO : %v", err)
	}

	d.SetId(iso.ID)

	_, err = waitForIsoAvailable(ctx, d, "complete", []string{"pending"}, "status", meta)
	if err != nil {
		return diag.Errorf(
			"error while waiting for ISO %s detach to be completed: %s", d.Id(), err)
	}

	return resourceVultrIsoRead(ctx, d, meta)
}

func resourceVultrIsoRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	iso, err := client.ISO.Get(ctx, d.Id())
	if err != nil {
		if strings.Contains("Invalid iso", err.Error()) {
			log.Printf("[WARN] Removing ISO (%s) because it is gone", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error getting ISO %s : %v", d.Id(), err)
	}

	d.Set("date_created", iso.DateCreated)
	d.Set("filename", iso.FileName)
	d.Set("size", iso.Size)
	d.Set("md5sum", iso.MD5Sum)
	d.Set("sha512sum", iso.SHA512Sum)
	d.Set("status", iso.Status)

	return nil
}

func resourceVultrIsoDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting iso : %s", d.Id())

	if err := client.ISO.Delete(ctx, d.Id()); err != nil {
		// decode the error
		var attachedErr struct {
			Error  string `json:"error"`
			Status int    `json:"status"`
		}

		if unmarshalError := json.Unmarshal([]byte(err.Error()), &attachedErr); unmarshalError != nil {
			return diag.Errorf("error deleting ISO %s: parsing error %v in deleting ISO %s : %v", d.Id(), err.Error(), unmarshalError)
		}

		if !strings.Contains(attachedErr.Error, "is still attached to") {
			return diag.Errorf("error deleting ISO %s: delete ISO error not related to attachment: delete error %+v", d.Id(), attachedErr)
		}

		parts := strings.Split(attachedErr.Error, " ")
		ip := parts[len(parts)-1]
		if parsedIP := net.ParseIP(ip); parsedIP == nil {
			return diag.Errorf("error deleting ISO %s : failed to parse IP to which ISO is attached: %s", d.Id(), ip)
		}

		// default is 100 instances
		var options govultr.ListOptions
		for {
			instances, responseMeta, err := client.Instance.List(ctx, &options)
			if err != nil {
				return diag.Errorf("error deleting ISO %s : failed to list instances for detaching ISO", d.Id(), err)
			}

			// check for the instance with this IP, return on failure or discovery
			for _, instance := range instances {
				if instance.MainIP == ip {
					if err := client.Instance.DetachISO(ctx, instance.ID); err != nil {
						return diag.Errorf("error deleting ISO %s : failed to detach from instances %s : %s", d.Id(), instance.ID, err)
					}
					_, err := waitForIsoDetached(ctx, instance.ID, "ready", []string{"isomounted"}, "status", meta)
					if err != nil {
						return diag.Errorf("error deleting ISO %s: failed to wait for ISO to detach from instance %s: %s", d.Id(), instance.ID, err)
					}
					if err = client.ISO.Delete(ctx, d.Id()); err != nil {
						return diag.Errorf("error deleting ISO %s: failed to delete ISO: %s", d.Id(), err)
					}
					return nil
				}
			}

			// no more instances to check
			if responseMeta.Links.Next == "" {
				break
			}
			options.Cursor = responseMeta.Links.Next
		}
		return diag.Errorf("failed to identify instance associated with IP %s for deleting ISO %s", ip, d.Id())
	}

	return nil
}

func waitForIsoAvailable(ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for ISO (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &resource.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    newIsoStateRefresh(ctx, d, meta),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,

		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func newIsoStateRefresh(ctx context.Context,
	d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*Client).govultrClient()

	return func() (interface{}, string, error) {

		log.Printf("[INFO] Creating Private ISO")
		iso, err := client.ISO.Get(ctx, d.Id())
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving ISO %s : %s", d.Id(), err)
		}

		log.Printf("[INFO] The ISO Status is %s", iso.Status)
		return &iso, iso.Status, nil
	}
}

func waitForIsoDetached(ctx context.Context, instanceID string, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for ISO to detach from %s",
		instanceID)

	stateConf := &resource.StateChangeConf{
		Pending:        pending,
		Target:         []string{target},
		Refresh:        isoDetachStateRefresh(ctx, instanceID, meta, attribute),
		Timeout:        60 * time.Minute,
		Delay:          10 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isoDetachStateRefresh(ctx context.Context, instanceID string, meta interface{}, attr string) resource.StateRefreshFunc {
	client := meta.(*Client).govultrClient()
	return func() (interface{}, string, error) {

		log.Printf("[INFO] Detaching ISO")
		iso, err := client.Instance.ISOStatus(ctx, instanceID)
		if err != nil {
			return nil, "", fmt.Errorf("error getting ISO status for instance %s : %s", instanceID, err)
		}
		return &iso, iso.State, nil
	}
}
