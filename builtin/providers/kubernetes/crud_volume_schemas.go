package kubernetes

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/resource"
	"k8s.io/kubernetes/pkg/util"
)

func createPersistentVolumeSource(_volume_sources []interface{}) *api.PersistentVolumeSource {
	if len(_volume_sources) == 0 {
		return nil
	} else {
		_volume_source := _volume_sources[0].(map[string]interface{})
		volumeSource := &api.VolumeSource{}

		if val, ok := _volume_source["host_path"]; ok {
			volumeSource.HostPath = createHostPathVolumeSource(val.([]interface{}))
		}

		if val, ok := _volume_source["gce_persistent_disk"]; ok {
			volumeSource.GCEPersistentDisk = createGcePersistentDiskVolumeSource(val.([]interface{}))
		}

		if val, ok := _volume_source["aws_elastic_block_store"]; ok {
			volumeSource.AWSElasticBlockStore = createAwsElasticBlockStoreVolumeSource(val.([]interface{}))
		}

		if val, ok := _volume_source["nfs"]; ok {
			volumeSource.NFS = createNfsVolumeSource(val.([]interface{}))
		}

		if val, ok := _volume_source["iscsi"]; ok {
			volumeSource.ISCSI = createIscsiVolumeSource(val.([]interface{}))
		}

		if val, ok := _volume_source["gluster_fs"]; ok {
			volumeSource.Glusterfs = createGlusterfsVolumeSource(val.([]interface{}))
		}

		if val, ok := _volume_source["cinder"]; ok {
			volumeSource.Cinder = createCinderVolumeSource(val.([]interface{}))
		}

		if val, ok := _volume_source["ceph_fs"]; ok {
			volumeSource.CephFS = createCephFsVolumeSource(val.([]interface{}))
		}

		if val, ok := _volume_source["flocker"]; ok {
			volumeSource.Flocker = createFlockerVolumeSource(val.([]interface{}))
		}

		if val, ok := _volume_source["fc"]; ok {
			volumeSource.FC = createFcVolumeSource(val.([]interface{}))
		}

		return volumeSource
	}
}

func createObjectReference(_object_references []interface{}) *api.ObjectReference {
	if len(_object_references) == 0 {
		return nil
	} else {
		_object_reference := _object_references[0].(map[string]interface{})
		objectReference := &api.ObjectReference{}

		if val, ok := _object_reference["kind"]; ok {
			objectReference.Kind = val.(string)
		}

		if val, ok := _object_reference["namespace"]; ok {
			objectReference.Namespace = val.(string)
		}

		if val, ok := _object_reference["name"]; ok {
			objectReference.Name = val.(string)
		}

		if val, ok := _object_reference["uid"]; ok {
			objectReference.UID = val.(string)
		}

		if val, ok := _object_reference["api_version"]; ok {
			objectReference.APIVersion = val.(string)
		}

		if val, ok := _object_reference["resource_version"]; ok {
			objectReference.ResourceVersion = val.(string)
		}

		if val, ok := _object_reference["field_path"]; ok {
			objectReference.FieldPath = val.(string)
		}

		return objectReference
	}
}

func readPersistentVolumeSource(volume []api.Volume) []interface{} {
	_volumes := make([]interface{}, 1)

	_volume := make(map[string]interface{})

	if volume.HostPath != nil {
		_volume_source["host_path"] = readHostPathVolumeSource(volume.HostPath)
	}

	if volume.GCEPersistentDisk != nil {
		_volume_source["gce_persistent_disk"] = readGcePersistentDiskVolumeSource(volume.GCEPersistentDisk)
	}

	if volume.AWSElasticBlockStore != nil {
		_volume_source["aws_elastic_block_store"] = readAwsElasticBlockStoreVolumeSource(volume.AWSElasticBlockStore)
	}

	if volume.NFS != nil {
		_volume_source["nfs"] = readNfsVolumeSource(volume.NFS)
	}

	if volume.ISCSI != nil {
		_volume_source["iscsi"] = readIscsiVolumeSource(volume.ISCSI)
	}

	if volume.Glusterfs != nil {
		_volume_source["gluster_fs"] = readGlusterfsVolumeSource(volume.Glusterfs)
	}

	if volume.Cinder != nil {
		_volume_source["cinder"] = readCinderVolumeSource(volume.Cinder)
	}

	if volume.CephFS != nil {
		_volume_source["ceph_fs"] = readCephFsVolumeSource(volume.CephFS)
	}

	if volume.Flocker != nil {
		_volume_source["flocker"] = readFlockerVolumeSource(volume.Flocker)
	}

	if volume.FC != nil {
		_volume_source["fc"] = readFcVolumeSource(volume.FC)
	}

	_volumes[0] = _volume

	return _volumes
}

func readStringList(values []api.PersistentVolumeAccessMode) []interface{} {
	_values := make([]interface{}, len(values))
	for i, v := range values {
		_values[i] = v
	}
	return _values
}

func readObjectReference(objectRefs *api.ObjectRef) []interface{} {
	if objectRefs == nil {
		return nil
	} else {
		_object_ref := make(map[string]interface{})
		_object_ref["kind"] = objectRefs.Kind
		_object_ref["name"] = objectRefs.Name
		_object_ref["namespace"] = objectRefs.Namespace
		_object_ref["uid"] = objectRefs.UID
		_object_ref["api_version"] = objectRefs.APIVersion
		_object_ref["resource_version"] = objectRefs.ResourceVersion
		_object_ref["field_path"] = objectRefs.FieldPath

		_object_refs := make([]interface{}, 1)
		_object_refs[0] = _object_ref

		return _object_refs
	}
}

