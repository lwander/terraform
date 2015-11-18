package kubernetes

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"k8s.io/kubernetes/pkg/api"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/util/yaml"
)

func resourceKubernetesPersistentVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesPersistentVolumeCreate,
		Read:   resourceKubernetesPersistentVolumeRead,
		Update: resourceKubernetesPersistentVolumeUpdate,
		Delete: resourceKubernetesPersistentVolumeDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"labels": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},

			"capacity": genResourceList(),

			"persistent_volume_source": genPersistentVolumeSource(),

			"access_modes": genStringList(),

			"claim_ref": genObjectReference(),

			"persistent_volume_reclaim_policy": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceKubernetesPersistentVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client.Client)
	log.Printf("[DEBUG] preparing to create persistent volume");

	_name := d.Get("name").(string)

	spec := api.PersistentVolumeSpec{}

	if v, ok := d.GetOk("capacity"); ok {
		spec.Capacity = createResourceList(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("persistent_volume_source"); ok {
		source := createPersistentVolumeSource(v.([]interface{}))

		// Persistent volume source is an inline attribute
		spec.HostPath = source.HostPath
		spec.GCEPersistentDisk = source.GCEPersistentDisk
		spec.AWSElasticBlockStore = source.AWSElasticBlockStore
		spec.NFS = source.NFS
		spec.ISCSI = source.ISCSI
		spec.Glusterfs = source.Glusterfs
		spec.Cinder = source.Cinder
		spec.CephFS = source.CephFS
		spec.Flocker = source.Flocker
		spec.FC = source.FC
	}

	if v, ok := d.GetOk("access_modes"); ok {
		spec.AccessModes = createStringList(v.([]interface{}))
	}

	if v, ok := d.GetOk("claim_ref"); ok {
		spec.ClaimRef = createObjectReference(v.([]interface{}))
	}

	if v, ok := d.GetOk("persistent_volume_reclaim_policy"); ok {
		spec.PersitentVolumeReclaimPolicy = v.(string)
	}

	_labels := d.Get("labels").(map[string]interface{})
	labels := make(map[string]string, len(_labels))
	for k, v := range _labels {
		labels[k] = v.(string)
	}

	req := api.PersitentVolume{
		ObjectMeta: api.ObjectMeta{
			Name:   _name,
			Labels: labels,
		},
		Spec: spec,
	}

	_namespace := d.Get("namespace").(string)

	_, err := c.PersistentVolumes().Create(&req)
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to create pod %s: %s", _name, err)
	}

	return resourceKubernetesPersistentVolumeRead(d, meta)
}

func resourceKubernetesPersistentVolumeRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client.Client)
	vol, err := c.PersistentVolumes().Get(d.Get("name").(string))
	if err != nil {
		return err
	}

	d.Set("capacity", readResourceList(vol.Capacity))
	d.Set("persistent_volume_source", readPersistentVolumeSource(vol))
	d.Set("access_modes", readAccessMode(vol.AccessModes))
	d.Set("claim_ref", readObjectReference(vol.ClaimRef))
	d.Set("persistent_volume_reclaim_policy", vol.PersistentVolumeReclaimPolicy)

	labels := vol.ObjectMeta.Labels
	_labels := make(map[string]interface{}, len(labels))
	for k, v := range labels {
		_labels[k] = v
	}

	d.Set("labels", _labels)

	return nil
}

func resourceKubernetesPersistentVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client.Client)
	log.Printf("[DEBUG] preparing to create persistent volume");

	_name := d.Get("name").(string)

	spec := api.PersistentVolumeSpec{}

	if v, ok := d.GetOk("capacity"); ok {
		spec.Capacity = createResourceList(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("persistent_volume_source"); ok {
		source := createPersistentVolumeSource(v.([]interface{}))

		// Persistent volume source is an inline attribute
		spec.HostPath = source.HostPath
		spec.GCEPersistentDisk = source.GCEPersistentDisk
		spec.AWSElasticBlockStore = source.AWSElasticBlockStore
		spec.NFS = source.NFS
		spec.ISCSI = source.ISCSI
		spec.Glusterfs = source.Glusterfs
		spec.Cinder = source.Cinder
		spec.CephFS = source.CephFS
		spec.Flocker = source.Flocker
		spec.FC = source.FC
	}

	if v, ok := d.GetOk("access_modes"); ok {
		spec.AccessModes = createStringList(v.([]interface{}))
	}

	if v, ok := d.GetOk("claim_ref"); ok {
		spec.ClaimRef = createObjectReference(v.([]interface{}))
	}

	if v, ok := d.GetOk("persistent_volume_reclaim_policy"); ok {
		spec.PersitentVolumeReclaimPolicy = v.(string)
	}

	_labels := d.Get("labels").(map[string]interface{})
	labels := make(map[string]string, len(_labels))
	for k, v := range _labels {
		labels[k] = v.(string)
	}

	req := api.PersitentVolume{
		ObjectMeta: api.ObjectMeta{
			Name:   _name,
			Labels: labels,
		},
		Spec: spec,
	}

	_namespace := d.Get("namespace").(string)

	_, err := c.PersistentVolumes().Update(&req)
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to create pod %s: %s", _name, err)
	}

	return resourceKubernetesPersistentVolumeRead(d, meta)
}

func resourceKubernetesPersistentVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client.Client)
	return c.PersistentVolumes().Delete(d.Get("name").(string))
}
