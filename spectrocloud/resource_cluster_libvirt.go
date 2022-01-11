package spectrocloud

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spectrocloud/gomi/pkg/ptr"
	"github.com/spectrocloud/hapi/models"
	"github.com/spectrocloud/terraform-provider-spectrocloud/pkg/client"
)

func resourceClusterLibvirt() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterVirtCreate,
		ReadContext:   resourceClusterLibvirtRead,
		UpdateContext: resourceClusterVirtUpdate,
		DeleteContext: resourceClusterDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		SchemaVersion: 2,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Set:      schema.HashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cluster_profile_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Switch to cluster_profile",
			},
			"cluster_profile": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"cluster_profile_id", "pack"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"pack": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"tag": {
										Type:     schema.TypeString,
										Required: true,
									},
									"values": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"cloud_account_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cloud_config_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_patch_on_boot": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"os_patch_schedule": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateOsPatchSchedule,
			},
			"os_patch_after": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateOsPatchOnDemandAfter,
			},
			"kubeconfig": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_config": {
				Type:     schema.TypeList,
				ForceNew: true,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ssh_key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"vip": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"pack": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"tag": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"machine_pool": {
				Type:     schema.TypeSet,
				Required: true,
				Set:      resourceMachinePoolLibvirtHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"control_plane": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							//ForceNew: true,
						},
						"control_plane_as_worker": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,

							//ForceNew: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
							//ForceNew: true,
						},
						"count": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"update_strategy": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "RollingUpdateScaleOut",
						},
						"instance_type": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"attached_disks_size_gb": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"cpus_sets": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"disk_size_gb": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"memory_mb": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"cpu": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"placements": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"appliance_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"network_type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"network_names": {
										Type:     schema.TypeString,
										Required: true,
									},
									"image_storage_pool": {
										Type:     schema.TypeString,
										Required: true,
									},
									"target_storage_pool": {
										Type:     schema.TypeString,
										Required: true,
									},
									"data_storage_pool": {
										Type:     schema.TypeString,
										Required: true,
									},
									"network": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},

			"backup_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix": {
							Type:     schema.TypeString,
							Required: true,
						},
						"backup_location_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"schedule": {
							Type:     schema.TypeString,
							Required: true,
						},
						"expiry_in_hour": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"include_disks": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"include_cluster_resources": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"namespaces": {
							Type:     schema.TypeSet,
							Optional: true,
							Set:      schema.HashString,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"scan_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"configuration_scan_schedule": {
							Type:     schema.TypeString,
							Required: true,
						},
						"penetration_scan_schedule": {
							Type:     schema.TypeString,
							Required: true,
						},
						"conformance_scan_schedule": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceClusterVirtCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.V1Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	cluster := toLibvirtCluster(d)

	uid, err := c.CreateClusterLibvirt(cluster)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(uid)

	stateConf := &resource.StateChangeConf{
		Pending:    resourceClusterCreatePendingStates,
		Target:     []string{"Running"},
		Refresh:    resourceClusterStateRefreshFunc(c, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate) - 1*time.Minute,
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceClusterLibvirtRead(ctx, d, m)

	return diags
}

//goland:noinspection GoUnhandledErrorResult
func resourceClusterLibvirtRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.V1Client)

	var diags diag.Diagnostics
	//
	uid := d.Id()
	//
	cluster, err := c.GetCluster(uid)
	if err != nil {
		return diag.FromErr(err)
	} else if cluster == nil {
		// Deleted - Terraform will recreate it
		d.SetId("")
		return diags
	}

	// Update the kubeconfig
	kubeconfig, err := c.GetClusterKubeConfig(uid)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("tags", flattenTags(cluster.Metadata.Labels)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("kubeconfig", kubeconfig); err != nil {
		return diag.FromErr(err)
	}

	if policy, err := c.GetClusterBackupConfig(d.Id()); err != nil {
		return diag.FromErr(err)
	} else if policy != nil && policy.Spec.Config != nil {
		if err := d.Set("backup_policy", flattenBackupPolicy(policy.Spec.Config)); err != nil {
			return diag.FromErr(err)
		}
	}

	if policy, err := c.GetClusterScanConfig(d.Id()); err != nil {
		return diag.FromErr(err)
	} else if policy != nil && policy.Spec.DriverSpec != nil {
		if err := d.Set("scan_policy", flattenScanPolicy(policy.Spec.DriverSpec)); err != nil {
			return diag.FromErr(err)
		}
	}

	return flattenCloudConfigLibvirt(cluster.Spec.CloudConfigRef.UID, d, c)
}

func flattenCloudConfigLibvirt(configUID string, d *schema.ResourceData, c *client.V1Client) diag.Diagnostics {
	d.Set("cloud_config_id", configUID)
	if config, err := c.GetCloudConfigLibvirt(configUID); err != nil {
		return diag.FromErr(err)
	} else {
		mp := flattenMachinePoolConfigsLibvirt(config.Spec.MachinePoolConfig)
		if err := d.Set("machine_pool", mp); err != nil {
			return diag.FromErr(err)
		}
	}

	return diag.Diagnostics{}
}

func flattenMachinePoolConfigsLibvirt(machinePools []*models.V1LibvirtMachinePoolConfig) []interface{} {

	if machinePools == nil {
		return make([]interface{}, 0)
	}

	ois := make([]interface{}, len(machinePools))

	for i, machinePool := range machinePools {
		oi := make(map[string]interface{})

		oi["control_plane"] = machinePool.IsControlPlane
		oi["control_plane_as_worker"] = machinePool.UseControlPlaneAsWorker
		oi["name"] = machinePool.Name
		oi["count"] = machinePool.Size
		oi["update_strategy"] = machinePool.UpdateStrategy.Type

		if machinePool.InstanceType != nil {
			s := make(map[string]interface{})
			additionalDisks := make([]string, 0)

			if machinePool.NonRootDisksInGB != nil && len(machinePool.NonRootDisksInGB) > 0 {
				for _, disk := range machinePool.NonRootDisksInGB {
					additionalDisks = append(additionalDisks, fmt.Sprint(*disk.SizeInGB))
				}
			}
			s["disk_size_gb"] = int(*machinePool.RootDiskInGB)
			s["memory_mb"] = int(*machinePool.InstanceType.MemoryInMB)
			s["cpu"] = int(*machinePool.InstanceType.NumCPUs)

			oi["instance_type"] = []interface{}{s}
			additionalDisksStr := strings.Join(additionalDisks, ",")
			s["attached_disks_size_gb"] = additionalDisksStr
		}

		placements := make([]interface{}, len(machinePool.Placements))
		for j, p := range machinePool.Placements {
			pj := make(map[string]interface{})
			pj["appliance_id"] = p.HostUID
			if p.Networks != nil {
				for _, network := range p.Networks {
					pj["network_type"] = network.NetworkType
					break
				}
			}
			networkNames := make([]string, 0)
			for _, network := range p.Networks {
				networkNames = append(networkNames, *network.NetworkName)
			}
			networkNamesStr := strings.Join(networkNames, ",")

			pj["network_names"] = networkNamesStr
			pj["image_storage_pool"] = p.SourceStoragePool
			pj["target_storage_pool"] = p.TargetStoragePool
			pj["data_storage_pool"] = p.DataStoragePool
			placements[j] = pj
		}
		oi["placements"] = placements

		ois[i] = oi
	}

	return ois
}

func resourceClusterVirtUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.V1Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	cloudConfigId := d.Get("cloud_config_id").(string)

	if d.HasChange("machine_pool") {
		oraw, nraw := d.GetChange("machine_pool")
		if oraw == nil {
			oraw = new(schema.Set)
		}
		if nraw == nil {
			nraw = new(schema.Set)
		}

		os := oraw.(*schema.Set)
		ns := nraw.(*schema.Set)

		osMap := make(map[string]interface{})
		for _, mp := range os.List() {
			machinePool := mp.(map[string]interface{})
			osMap[machinePool["name"].(string)] = machinePool
		}

		for _, mp := range ns.List() {
			machinePoolResource := mp.(map[string]interface{})
			name := machinePoolResource["name"].(string)
			hash := resourceMachinePoolAwsHash(machinePoolResource)

			machinePool := toMachinePoolAws(machinePoolResource)

			var err error
			if oldMachinePool, ok := osMap[name]; !ok {
				log.Printf("Create machine pool %s", name)
				err = c.CreateMachinePoolAws(cloudConfigId, machinePool)
			} else if hash != resourceMachinePoolAwsHash(oldMachinePool) {
				log.Printf("Change in machine pool %s", name)
				err = c.UpdateMachinePoolAws(cloudConfigId, machinePool)
			}

			if err != nil {
				return diag.FromErr(err)
			}

			// Processed (if exists)
			delete(osMap, name)
		}

		// Deleted old machine pools
		for _, mp := range osMap {
			machinePool := mp.(map[string]interface{})
			name := machinePool["name"].(string)
			log.Printf("Deleted machine pool %s", name)
			if err := c.DeleteMachinePoolAws(cloudConfigId, name); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChanges("cluster_profile") {
		if err := updateProfiles(c, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("backup_policy") {
		if err := updateBackupPolicy(c, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("scan_policy") {
		if err := updateScanPolicy(c, d); err != nil {
			return diag.FromErr(err)
		}
	}

	resourceClusterLibvirtRead(ctx, d, m)

	return diags
}

func toLibvirtCluster(d *schema.ResourceData) *models.V1SpectroLibvirtClusterEntity {
	cloudConfig := d.Get("cloud_config").([]interface{})[0].(map[string]interface{})

	cluster := &models.V1SpectroLibvirtClusterEntity{
		Metadata: &models.V1ObjectMeta{
			Name:   d.Get("name").(string),
			UID:    d.Id(),
			Labels: toTags(d),
		},
		Spec: &models.V1SpectroLibvirtClusterEntitySpec{
			Profiles: toProfiles(d),
			Policies: toPolicies(d),
			CloudConfig: &models.V1LibvirtClusterConfig{
				NtpServers: nil,
				SSHKeys:    []string{cloudConfig["ssh_key"].(string)},
				ControlPlaneEndpoint: &models.V1LibvirtControlPlaneEndPoint{
					Host: cloudConfig["vip"].(string),
					Type: "VIP",
				},
			},
		},
	}

	machinePoolConfigs := make([]*models.V1LibvirtMachinePoolConfigEntity, 0)
	for _, machinePool := range d.Get("machine_pool").(*schema.Set).List() {
		mp := toMachinePoolLibvirt(machinePool)
		machinePoolConfigs = append(machinePoolConfigs, mp)
	}

	// sort
	sort.SliceStable(machinePoolConfigs, func(i, j int) bool {
		return machinePoolConfigs[i].PoolConfig.IsControlPlane
	})

	cluster.Spec.Machinepoolconfig = machinePoolConfigs
	cluster.Spec.ClusterConfig = toClusterConfig(d)

	return cluster
}

func toMachinePoolLibvirt(machinePool interface{}) *models.V1LibvirtMachinePoolConfigEntity {
	m := machinePool.(map[string]interface{})

	labels := make([]string, 0)
	controlPlane := m["control_plane"].(bool)
	controlPlaneAsWorker := m["control_plane_as_worker"].(bool)
	if controlPlane {
		labels = append(labels, "master")
	}

	placements := make([]*models.V1LibvirtPlacementEntity, 0)
	for _, pos := range m["placements"].([]interface{}) {
		p := pos.(map[string]interface{})
		networks := getNetworks(p)

		imageStoragePool := p["image_storage_pool"].(string)
		targetStoragePool := p["target_storage_pool"].(string)
		dataStoragePool := p["data_storage_pool"].(string)

		placements = append(placements, &models.V1LibvirtPlacementEntity{
			Networks:          networks,
			SourceStoragePool: imageStoragePool,
			TargetStoragePool: targetStoragePool,
			DataStoragePool:   dataStoragePool,
			HostUID:           ptr.StringPtr(p["appliance_id"].(string)),
		})

	}

	ins := m["instance_type"].([]interface{})[0].(map[string]interface{})
	instanceType := models.V1LibvirtInstanceType{
		Cpuset:     strconv.FormatInt(int64(ins["cpus_sets"].(int)), 10),
		MemoryInMB: ptr.Int32Ptr(int32(ins["memory_mb"].(int))),
		NumCPUs:    ptr.Int32Ptr(int32(ins["cpu"].(int))),
	}

	addDisks := getAdditionalDisks(ins)

	mp := &models.V1LibvirtMachinePoolConfigEntity{
		CloudConfig: &models.V1LibvirtMachinePoolCloudConfigEntity{
			Placements:       placements,
			RootDiskInGB:     ptr.Int32Ptr(int32(ins["disk_size_gb"].(int))),
			NonRootDisksInGB: addDisks,
			InstanceType:     &instanceType,
		},
		PoolConfig: &models.V1MachinePoolConfigEntity{
			IsControlPlane: controlPlane,
			Labels:         labels,
			Name:           ptr.StringPtr(m["name"].(string)),
			Size:           ptr.Int32Ptr(int32(m["count"].(int))),
			UpdateStrategy: &models.V1UpdateStrategy{
				Type: m["update_strategy"].(string),
			},
			UseControlPlaneAsWorker: controlPlaneAsWorker,
		},
	}
	return mp
}

func getAdditionalDisks(ins map[string]interface{}) []*models.V1LibvirtDiskSpec {
	addDisks := make([]*models.V1LibvirtDiskSpec, 0)

	if ins["attached_disks_size_gb"] != nil {
		disks := strings.Split(ins["attached_disks_size_gb"].(string), ",")
		for _, addDisk := range disks {
			x, err := strconv.ParseInt(strings.TrimSpace(addDisk), 10, 32)
			if err != nil {
				return nil
			}
			size := int32(x)
			addDisks = append(addDisks, &models.V1LibvirtDiskSpec{
				SizeInGB: &size,
			})
		}
	}
	return addDisks
}

func getNetworks(p map[string]interface{}) []*models.V1LibvirtNetworkSpec {
	networkType := ""
	networks := make([]*models.V1LibvirtNetworkSpec, 0)

	if p["network_names"] != nil {
		for _, n := range strings.Split(p["network_names"].(string), ",") {
			networkName := strings.TrimSpace(n)
			networkType = p["network_type"].(string)
			network := &models.V1LibvirtNetworkSpec{
				NetworkName: &networkName,
				NetworkType: &networkType,
			}
			networks = append(networks, network)
		}
	}
	return networks
}