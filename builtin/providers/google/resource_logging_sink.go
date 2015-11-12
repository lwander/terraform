package google

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"

	"google.golang.org/api/logging/v1beta3"
)


func resourceLoggingSink() *schema.Resource {
	return &schema.Resource{
		Create: resourceLoggingSinkCreate,
		Read:   resourceLoggingSinkRead,
		Update: resourceLoggingSinkUpdate,
		Delete: resourceLoggingSinkDelete,

		Schema: map[string]*schema.Schema{
			"log_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"sink_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"destination": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"filter": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"errors": &schema.Schema {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"timeNanos": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},

						"status_code": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},

						"status_message": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"status_details": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem:     schema.TypeMap,
						},
					},
				},
			},
		},
	}
}

func createLogSink(d *schema.ResourceData) *logging.LogSink {
	filter := ""
	tfilter, ok := d.GetOk("filter")
	if ok {
		filter = tfilter.(string)
	}

	return &logging.LogSink{
		Destination: d.Get("destination").(string),
		Name: d.Get("sink_name").(string),
		Filter: filter,
	}
}

func readLogSinkErrors(logSink *logging.LogSink) []interface{} {
	res := make([]interface{}, 0)
	for _, err := range(logSink.Errors) {
		ins := make(map[string]interface{})
		ins["resource"] = err.Resource
		ins["timeNanos"] = err.TimeNanos
		ins["status_code"] = err.Status.Code
		ins["status_message"] = err.Status.Message
		ins["status_details"] = make([]interface{}, 0)
		for _, detail := range(err.Status.Details) {
			ins["status_details"] = append(ins["status_details"].([]interface{}), detail)
		}
		res = append(res, ins)
	}
	return res
}

func getSinkId(sink_name, log_name string) string {
	return log_name + "-" + sink_name
}

func resourceLoggingSinkCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	logsSinksService := logging.NewProjectsLogsSinksService(config.clientLogging)
	logSink := createLogSink(d)
	res, err := logsSinksService.Create(config.Project,
		d.Get("log_name").(string), logSink).Do()
	if err != nil {
		return fmt.Errorf("Error creating sink %s: %s",
			d.Get("sink_name").(string), err)
	}

	d.Set("errors", readLogSinkErrors(res))
	d.SetId(getSinkId(d.Get("sink_name").(string), d.Get("log_name").(string)))

	return nil
}

func resourceLoggingSinkRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	logsSinksService := logging.NewProjectsLogsSinksService(config.clientLogging)
	res, err := logsSinksService.Get(config.Project,
		d.Get("log_name").(string), d.Get("sink_name").(string)).Do()

	if err != nil {
		return fmt.Errorf("Error getting sink %s: %s",
			d.Get("sink_name").(string), err)
	}

	d.Set("errors", readLogSinkErrors(res))
	d.SetId(getSinkId(d.Get("sink_name").(string), d.Get("log_name").(string)))

	return nil
}

func resourceLoggingSinkUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	old_sink := d.Get("sink_name")
	if d.HasChange("sink_name") {
		old_sink, _ = d.GetChange("sink_name")
	}

	old_sink_name := old_sink.(string)

	logsSinksService := logging.NewProjectsLogsSinksService(config.clientLogging)
	logSink := createLogSink(d)
	res, err := logsSinksService.Update(config.Project,
		d.Get("log_name").(string), old_sink_name, logSink).Do()
	if err != nil {
		return fmt.Errorf("Error updating sink %s: %s",
			d.Get("sink_name").(string), err)
	}

	d.Set("errors", readLogSinkErrors(res))
	d.SetId(getSinkId(d.Get("sink_name").(string), d.Get("log_name").(string)))

	return nil
}

func resourceLoggingSinkDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	logsSinksService := logging.NewProjectsLogsSinksService(config.clientLogging)
	_, err := logsSinksService.Delete(config.Project,
		d.Get("log_name").(string), d.Get("sink_name").(string)).Do()

	if err != nil {
		return fmt.Errorf("Error deleting sink %s: %s",
			d.Get("sink_name").(string), err)
	}

	return nil
}
