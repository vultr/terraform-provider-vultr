package vultr

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
	iso, err := client.ISO.CreateFromURL(context.Background(), d.Get("url").(string))

	if err != nil {
		return fmt.Errorf("Error creating ISO : %v", err)
	}

	d.SetId(strconv.Itoa(iso.ISOID))

	_, err = waitForIsoAvailable(d, "complete", []string{"pending"}, "status", meta)
	if err != nil {
		return fmt.Errorf(
			"Error while waiting for ISO %s to be completed: %s", d.Id(), err)
	}

	return resourceVultrIsoRead(d, meta)
}

func resourceVultrIsoRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	isoList, err := client.ISO.GetList(context.Background())

	if err != nil {
		fmt.Errorf("Error getting ISO %s : %v", d.Id(), err)
	}

	exists := false
	counter := 0
	for _, v := range isoList {
		if strconv.Itoa(v.ISOID) == d.Id() {
			exists = true
			break
		}
		counter++
	}

	if !exists {
		log.Printf("[WARN] Removing ISO (%s) because it is gone", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("date_created", isoList[counter].DateCreated)
	d.Set("filename", isoList[counter].FileName)
	d.Set("size", isoList[counter].Size)
	d.Set("md5sum", isoList[counter].MD5Sum)
	d.Set("sha512sum", isoList[counter].SHA512Sum)
	d.Set("status", isoList[counter].Status)

	return nil
}

func resourceVultrIsoDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Destroying iso : %s", d.Id())

	id, err := strconv.Atoi(d.Id())

	if err != nil {
		return fmt.Errorf("Error occuring while retreiving ISO id : %v", err)
	}

	if err := client.ISO.Delete(context.Background(), id); err != nil {
		return fmt.Errorf("Error destroying ISO %d : %v", id, err)
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
		isoList, err := client.ISO.GetList(context.Background())

		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving ISO %s : %s", d.Id(), err)
		}

		counter := 0
		for _, v := range isoList {
			if strconv.Itoa(v.ISOID) == d.Id() {
				break
			}
			counter++
		}

		log.Printf("[INFO] The ISO Status is %s", isoList[counter].Status)
		return &isoList[counter], isoList[counter].Status, nil
	}
}
