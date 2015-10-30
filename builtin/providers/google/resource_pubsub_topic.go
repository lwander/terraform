package google

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/pubsub/v1"
)

func resourcePubSubTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourcePubSubTopicCreate,
		Read:   resourcePubSubTopicRead,
		Delete: resourcePubSubTopicDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func formatTopicName(project, name string) string {
	return fmt.Sprintf("projects/%s/topics/%s", project, name);
}

func resourcePubSubTopicCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	name := formatTopicName(config.Project, d.Get("name").(string))

	topic := &pubsub.Topic{ Name: name }

	_, err := config.clientPubSub.Projects.Topics.Create(name, topic).Do()
	if err != nil {
		return fmt.Errorf("Error creating topic %s: %s", name, err)
	}

	return resourcePubSubTopicRead(d, meta)
}

func resourcePubSubTopicRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	name := formatTopicName(config.Project, d.Get("name").(string))

	topic, err := config.clientPubSub.Projects.Topics.Get(name).Do()
	if err != nil {
		return fmt.Errorf("Error getting topic %s: %s", name, err)
	}

	d.Set("name", topic.Name);
	d.SetId(topic.Name);

	return nil
}

func resourcePubSubTopicDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	name := formatTopicName(config.Project, d.Get("name").(string))

	_, err := config.clientPubSub.Projects.Topics.Delete(name).Do()
	if err != nil {
		return fmt.Errorf("Error deleting topic %s: %s", name, err)
	}

	d.SetId("")
	return nil
}
