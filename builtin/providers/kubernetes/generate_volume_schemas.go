package kubernetes

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func genPersistentVolumeSource() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"host_path": genHostPathVolumeSource(),

				"gce_persistent_disk": genGcePersistentDiskVolumeSource(),

				"aws_elastic_block_store": genAwsElasticBlockStoreVolumeSource(),

				"nfs": genNfsVolumeSource(),

				"iscsi": genIscsiVolumeSource(),

				"glusterfs": genGlusterfsVolumeSource(),

				"rbd": genRbdVolumeSource(),

				"cinder": genCinderVolumeSource(),

				"ceph_fs": genCephVolumeSource(),

				"flocker": genFlockerVolumeSource(),

				"fc": genFcVolumeSource(),
			},
		},
	}
}

func genObjectReference() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"kind": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},

				"namespace": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},

				"name": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},

				"uid": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},

				"api_version": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},

				"resource_version": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},

				"field_path": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}
