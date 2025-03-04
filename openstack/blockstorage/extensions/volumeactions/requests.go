package volumeactions

import (
	"github.com/gophercloud/gophercloud/v2"
)

// AttachOptsBuilder allows extensions to add additional parameters to the
// Attach request.
type AttachOptsBuilder interface {
	ToVolumeAttachMap() (map[string]interface{}, error)
}

// AttachMode describes the attachment mode for volumes.
type AttachMode string

// These constants determine how a volume is attached.
const (
	ReadOnly  AttachMode = "ro"
	ReadWrite AttachMode = "rw"
)

// AttachOpts contains options for attaching a Volume.
type AttachOpts struct {
	// The mountpoint of this volume.
	MountPoint string `json:"mountpoint,omitempty"`

	// The nova instance ID, can't set simultaneously with HostName.
	InstanceUUID string `json:"instance_uuid,omitempty"`

	// The hostname of baremetal host, can't set simultaneously with InstanceUUID.
	HostName string `json:"host_name,omitempty"`

	// Mount mode of this volume.
	Mode AttachMode `json:"mode,omitempty"`
}

// ToVolumeAttachMap assembles a request body based on the contents of a
// AttachOpts.
func (opts AttachOpts) ToVolumeAttachMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "os-attach")
}

// Attach will attach a volume based on the values in AttachOpts.
func Attach(client *gophercloud.ServiceClient, id string, opts AttachOptsBuilder) (r AttachResult) {
	b, err := opts.ToVolumeAttachMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(actionURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// BeginDetaching will mark the volume as detaching.
func BeginDetaching(client *gophercloud.ServiceClient, id string) (r BeginDetachingResult) {
	b := map[string]interface{}{"os-begin_detaching": make(map[string]interface{})}
	resp, err := client.Post(actionURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DetachOptsBuilder allows extensions to add additional parameters to the
// Detach request.
type DetachOptsBuilder interface {
	ToVolumeDetachMap() (map[string]interface{}, error)
}

// DetachOpts contains options for detaching a Volume.
type DetachOpts struct {
	// AttachmentID is the ID of the attachment between a volume and instance.
	AttachmentID string `json:"attachment_id,omitempty"`
}

// ToVolumeDetachMap assembles a request body based on the contents of a
// DetachOpts.
func (opts DetachOpts) ToVolumeDetachMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "os-detach")
}

// Detach will detach a volume based on volume ID.
func Detach(client *gophercloud.ServiceClient, id string, opts DetachOptsBuilder) (r DetachResult) {
	b, err := opts.ToVolumeDetachMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(actionURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Reserve will reserve a volume based on volume ID.
func Reserve(client *gophercloud.ServiceClient, id string) (r ReserveResult) {
	b := map[string]interface{}{"os-reserve": make(map[string]interface{})}
	resp, err := client.Post(actionURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Unreserve will unreserve a volume based on volume ID.
func Unreserve(client *gophercloud.ServiceClient, id string) (r UnreserveResult) {
	b := map[string]interface{}{"os-unreserve": make(map[string]interface{})}
	resp, err := client.Post(actionURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// InitializeConnectionOptsBuilder allows extensions to add additional parameters to the
// InitializeConnection request.
type InitializeConnectionOptsBuilder interface {
	ToVolumeInitializeConnectionMap() (map[string]interface{}, error)
}

// InitializeConnectionOpts hosts options for InitializeConnection.
// The fields are specific to the storage driver in use and the destination
// attachment.
type InitializeConnectionOpts struct {
	IP        string   `json:"ip,omitempty"`
	Host      string   `json:"host,omitempty"`
	Initiator string   `json:"initiator,omitempty"`
	Wwpns     []string `json:"wwpns,omitempty"`
	Wwnns     string   `json:"wwnns,omitempty"`
	Multipath *bool    `json:"multipath,omitempty"`
	Platform  string   `json:"platform,omitempty"`
	OSType    string   `json:"os_type,omitempty"`
}

// ToVolumeInitializeConnectionMap assembles a request body based on the contents of a
// InitializeConnectionOpts.
func (opts InitializeConnectionOpts) ToVolumeInitializeConnectionMap() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "connector")
	return map[string]interface{}{"os-initialize_connection": b}, err
}

// InitializeConnection initializes an iSCSI connection by volume ID.
func InitializeConnection(client *gophercloud.ServiceClient, id string, opts InitializeConnectionOptsBuilder) (r InitializeConnectionResult) {
	b, err := opts.ToVolumeInitializeConnectionMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(actionURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// TerminateConnectionOptsBuilder allows extensions to add additional parameters to the
// TerminateConnection request.
type TerminateConnectionOptsBuilder interface {
	ToVolumeTerminateConnectionMap() (map[string]interface{}, error)
}

// TerminateConnectionOpts hosts options for TerminateConnection.
type TerminateConnectionOpts struct {
	IP        string   `json:"ip,omitempty"`
	Host      string   `json:"host,omitempty"`
	Initiator string   `json:"initiator,omitempty"`
	Wwpns     []string `json:"wwpns,omitempty"`
	Wwnns     string   `json:"wwnns,omitempty"`
	Multipath *bool    `json:"multipath,omitempty"`
	Platform  string   `json:"platform,omitempty"`
	OSType    string   `json:"os_type,omitempty"`
}

// ToVolumeTerminateConnectionMap assembles a request body based on the contents of a
// TerminateConnectionOpts.
func (opts TerminateConnectionOpts) ToVolumeTerminateConnectionMap() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "connector")
	return map[string]interface{}{"os-terminate_connection": b}, err
}

// TerminateConnection terminates an iSCSI connection by volume ID.
func TerminateConnection(client *gophercloud.ServiceClient, id string, opts TerminateConnectionOptsBuilder) (r TerminateConnectionResult) {
	b, err := opts.ToVolumeTerminateConnectionMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(actionURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ExtendSizeOptsBuilder allows extensions to add additional parameters to the
// ExtendSize request.
type ExtendSizeOptsBuilder interface {
	ToVolumeExtendSizeMap() (map[string]interface{}, error)
}

// ExtendSizeOpts contains options for extending the size of an existing Volume.
// This object is passed to the volumes.ExtendSize function.
type ExtendSizeOpts struct {
	// NewSize is the new size of the volume, in GB.
	NewSize int `json:"new_size" required:"true"`
}

// ToVolumeExtendSizeMap assembles a request body based on the contents of an
// ExtendSizeOpts.
func (opts ExtendSizeOpts) ToVolumeExtendSizeMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "os-extend")
}

// ExtendSize will extend the size of the volume based on the provided information.
// This operation does not return a response body.
func ExtendSize(client *gophercloud.ServiceClient, id string, opts ExtendSizeOptsBuilder) (r ExtendSizeResult) {
	b, err := opts.ToVolumeExtendSizeMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(actionURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UploadImageOptsBuilder allows extensions to add additional parameters to the
// UploadImage request.
type UploadImageOptsBuilder interface {
	ToVolumeUploadImageMap() (map[string]interface{}, error)
}

// UploadImageOpts contains options for uploading a Volume to image storage.
type UploadImageOpts struct {
	// Container format, may be bare, ofv, ova, etc.
	ContainerFormat string `json:"container_format,omitempty"`

	// Disk format, may be raw, qcow2, vhd, vdi, vmdk, etc.
	DiskFormat string `json:"disk_format,omitempty"`

	// The name of image that will be stored in glance.
	ImageName string `json:"image_name,omitempty"`

	// Force image creation, usable if volume attached to instance.
	Force bool `json:"force,omitempty"`

	// Visibility defines who can see/use the image.
	// supported since 3.1 microversion
	Visibility string `json:"visibility,omitempty"`

	// whether the image is not deletable.
	// supported since 3.1 microversion
	Protected bool `json:"protected,omitempty"`
}

// ToVolumeUploadImageMap assembles a request body based on the contents of a
// UploadImageOpts.
func (opts UploadImageOpts) ToVolumeUploadImageMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "os-volume_upload_image")
}

// UploadImage will upload an image based on the values in UploadImageOptsBuilder.
func UploadImage(client *gophercloud.ServiceClient, id string, opts UploadImageOptsBuilder) (r UploadImageResult) {
	b, err := opts.ToVolumeUploadImageMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(actionURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ForceDelete will delete the volume regardless of state.
func ForceDelete(client *gophercloud.ServiceClient, id string) (r ForceDeleteResult) {
	resp, err := client.Post(actionURL(client, id), map[string]interface{}{"os-force_delete": ""}, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ImageMetadataOptsBuilder allows extensions to add additional parameters to the
// ImageMetadataRequest request.
type ImageMetadataOptsBuilder interface {
	ToImageMetadataMap() (map[string]interface{}, error)
}

// ImageMetadataOpts contains options for setting image metadata to a volume.
type ImageMetadataOpts struct {
	// The image metadata to add to the volume as a set of metadata key and value pairs.
	Metadata map[string]string `json:"metadata"`
}

// ToImageMetadataMap assembles a request body based on the contents of a
// ImageMetadataOpts.
func (opts ImageMetadataOpts) ToImageMetadataMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "os-set_image_metadata")
}

// SetImageMetadata will set image metadata on a volume based on the values in ImageMetadataOptsBuilder.
func SetImageMetadata(client *gophercloud.ServiceClient, id string, opts ImageMetadataOptsBuilder) (r SetImageMetadataResult) {
	b, err := opts.ToImageMetadataMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(actionURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// BootableOpts contains options for setting bootable status to a volume.
type BootableOpts struct {
	// Enables or disables the bootable attribute. You can boot an instance from a bootable volume.
	Bootable bool `json:"bootable"`
}

// ToBootableMap assembles a request body based on the contents of a
// BootableOpts.
func (opts BootableOpts) ToBootableMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "os-set_bootable")
}

// SetBootable will set bootable status on a volume based on the values in BootableOpts
func SetBootable(client *gophercloud.ServiceClient, id string, opts BootableOpts) (r SetBootableResult) {
	b, err := opts.ToBootableMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(actionURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// MigrationPolicy type represents a migration_policy when changing types.
type MigrationPolicy string

// Supported attributes for MigrationPolicy attribute for changeType operations.
const (
	MigrationPolicyNever    MigrationPolicy = "never"
	MigrationPolicyOnDemand MigrationPolicy = "on-demand"
)

// ChangeTypeOptsBuilder allows extensions to add additional parameters to the
// ChangeType request.
type ChangeTypeOptsBuilder interface {
	ToVolumeChangeTypeMap() (map[string]interface{}, error)
}

// ChangeTypeOpts contains options for changing the type of an existing Volume.
// This object is passed to the volumes.ChangeType function.
type ChangeTypeOpts struct {
	// NewType is the name of the new volume type of the volume.
	NewType string `json:"new_type" required:"true"`

	// MigrationPolicy specifies if the volume should be migrated when it is
	// re-typed. Possible values are "on-demand" or "never". If not specified,
	// the default is "never".
	MigrationPolicy MigrationPolicy `json:"migration_policy,omitempty"`
}

// ToVolumeChangeTypeMap assembles a request body based on the contents of an
// ChangeTypeOpts.
func (opts ChangeTypeOpts) ToVolumeChangeTypeMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "os-retype")
}

// ChangeType will change the volume type of the volume based on the provided information.
// This operation does not return a response body.
func ChangeType(client *gophercloud.ServiceClient, id string, opts ChangeTypeOptsBuilder) (r ChangeTypeResult) {
	b, err := opts.ToVolumeChangeTypeMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(actionURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ReImageOpts contains options for Re-image a volume.
type ReImageOpts struct {
	// New image id
	ImageID string `json:"image_id"`
	// set true to re-image volumes in reserved state
	ReImageReserved bool `json:"reimage_reserved"`
}

// ToReImageMap assembles a request body based on the contents of a ReImageOpts.
func (opts ReImageOpts) ToReImageMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "os-reimage")
}

// ReImage will re-image a volume based on the values in ReImageOpts
func ReImage(client *gophercloud.ServiceClient, id string, opts ReImageOpts) (r ReImageResult) {
	b, err := opts.ToReImageMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(actionURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ResetStatusOptsBuilder allows extensions to add additional parameters to the
// ResetStatus request.
type ResetStatusOptsBuilder interface {
	ToResetStatusMap() (map[string]interface{}, error)
}

// ResetStatusOpts contains options for resetting a Volume status.
// For more information about these parameters, please, refer to the Block Storage API V3,
// Volume Actions, ResetStatus volume documentation.
type ResetStatusOpts struct {
	// Status is a volume status to reset to.
	Status string `json:"status"`
	// MigrationStatus is a volume migration status to reset to.
	MigrationStatus string `json:"migration_status,omitempty"`
	// AttachStatus is a volume attach status to reset to.
	AttachStatus string `json:"attach_status,omitempty"`
}

// ToResetStatusMap assembles a request body based on the contents of a
// ResetStatusOpts.
func (opts ResetStatusOpts) ToResetStatusMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "os-reset_status")
}

// ResetStatus will reset the existing volume status. ResetStatusResult contains only the error.
// To extract it, call the ExtractErr method on the ResetStatusResult.
func ResetStatus(client *gophercloud.ServiceClient, id string, opts ResetStatusOptsBuilder) (r ResetStatusResult) {
	b, err := opts.ToResetStatusMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(actionURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
