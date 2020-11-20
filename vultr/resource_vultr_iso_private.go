package vultr

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/vultr/govultr/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceVultrIsoPrivate() *schema.Resource {
	return &schema.Resource{
		Create: resourceVultrIsoCreate,
		Read:   resourceVultrIsoRead,
		Delete: resourceVultrIsoDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceVultrIsoCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Creating new ISO")

	isoReq := &govultr.ISOReq{URL: d.Get("url").(string)}
	iso, err := client.ISO.Create(context.Background(), isoReq)
	if err != nil {
		return fmt.Errorf("Error creating ISO : %v", err)
	}

	d.SetId(iso.ID)

	_, err = waitForIsoAvailable(d, "complete", []string{"pending"}, "status", meta)
	if err != nil {
		return fmt.Errorf(
			"error while waiting for ISO %s to be completed: %s", d.Id(), err)
	}

	return resourceVultrIsoRead(d, meta)
}

func resourceVultrIsoRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	iso, err := client.ISO.Get(context.Background(), d.Id())
	if err != nil {
		if strings.Contains("Invalid iso", err.Error()) {
			log.Printf("[WARN] Removing ISO (%s) because it is gone", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting ISO %s : %v", d.Id(), err)
	}

	d.Set("date_created", iso.DateCreated)
	d.Set("filename", iso.FileName)
	d.Set("size", iso.Size)
	d.Set("md5sum", iso.MD5Sum)
	d.Set("sha512sum", iso.SHA512Sum)
	d.Set("status", iso.Status)

	return nil
}

func resourceVultrIsoDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting iso : %s", d.Id())

	if err := client.ISO.Delete(context.Background(), d.Id()); err != nil {
		return fmt.Errorf("error destroying ISO %s : %v", d.Id(), err)
	}

	return nil
}

func waitForIsoAvailable(d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for ISO (%s) to have %s of %s",
		d.Id(), attribute, target)

	stateConf := &resource.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    newIsoStateRefresh(d, meta),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,

		NotFoundChecks: 60,
	}

	return stateConf.WaitForState()
}

func newIsoStateRefresh(
	d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*Client).govultrClient()

	return func() (interface{}, string, error) {

		log.Printf("[INFO] Creating Private ISO")
		iso, err := client.ISO.Get(context.Background(), d.Id())
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving ISO %s : %s", d.Id(), err)
		}

		log.Printf("[INFO] The ISO Status is %s", iso.Status)
		return &iso, iso.Status, nil
	}
}
