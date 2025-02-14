package rds

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func DataSourceSnapshot() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSnapshotRead,

		Schema: map[string]*schema.Schema{
			//selection criteria
			"db_instance_identifier": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"db_snapshot_identifier": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"snapshot_type": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"include_shared": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"include_public": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"most_recent": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			//Computed values returned
			"allocated_storage": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_snapshot_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"encrypted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"iops": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"license_model": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"option_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"source_db_snapshot_identifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RDSConn

	instanceIdentifier, instanceIdentifierOk := d.GetOk("db_instance_identifier")
	snapshotIdentifier, snapshotIdentifierOk := d.GetOk("db_snapshot_identifier")
	snapshotType, snapshotTypeOk := d.GetOk("snapshot_type")

	if !instanceIdentifierOk && !snapshotIdentifierOk {
		return fmt.Errorf("One of db_snapshot_identifier or db_instance_identifier must be assigned")
	}

	params := &rds.DescribeDBSnapshotsInput{
		IncludePublic: aws.Bool(d.Get("include_public").(bool)),
		IncludeShared: aws.Bool(d.Get("include_shared").(bool)),
	}

	if snapshotTypeOk {
		params.SnapshotType = aws.String(snapshotType.(string))
	}

	// Don't set DBInstanceIdentifier for public or shared snapshots as it will
	// never match (only instance identifiers from the current account will match).
	// The filtering will need to be done client-side.
	if instanceIdentifierOk {
		if snapshotType == "public" || snapshotType == "shared" || d.Get("include_public").(bool) || d.Get("include_shared").(bool) {
			log.Printf("[DEBUG] Not combining DBInstanceIdentifier with SnapshotType %s in query; filtering client-side instead", snapshotType)
		} else {
			params.DBInstanceIdentifier = aws.String(instanceIdentifier.(string))
		}
	}

	if snapshotIdentifierOk {
		params.DBSnapshotIdentifier = aws.String(snapshotIdentifier.(string))
	}

	log.Printf("[DEBUG] Reading DB Snapshot: %s", params)
	var snapshots []*rds.DBSnapshot
	err := conn.DescribeDBSnapshotsPages(params, func(page *rds.DescribeDBSnapshotsOutput, lastPage bool) bool {
		snapshots = append(snapshots, page.DBSnapshots...)
		return true
	})
	if err != nil {
		return err
	}

	if (snapshotType == "public" || snapshotType == "shared" || d.Get("include_public").(bool) || d.Get("include_shared").(bool)) && instanceIdentifierOk {
		snapshots = filterSnapshotsByInstanceId(snapshots, instanceIdentifier.(string))
	}

	if len(snapshots) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	var snapshot *rds.DBSnapshot
	if len(snapshots) > 1 {
		recent := d.Get("most_recent").(bool)
		log.Printf("[DEBUG] aws_db_snapshot - multiple results found and `most_recent` is set to: %t", recent)
		if recent {
			snapshot = mostRecentDbSnapshot(snapshots)
		} else {
			return fmt.Errorf("Your query returned more than one result. Please try a more specific search criteria.")
		}
	} else {
		snapshot = snapshots[0]
	}

	return dbSnapshotDescriptionAttributes(d, snapshot)
}

type rdsSnapshotSort []*rds.DBSnapshot

func (a rdsSnapshotSort) Len() int      { return len(a) }
func (a rdsSnapshotSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a rdsSnapshotSort) Less(i, j int) bool {
	// Snapshot creation can be in progress
	if a[i].SnapshotCreateTime == nil {
		return true
	}
	if a[j].SnapshotCreateTime == nil {
		return false
	}

	return (*a[i].SnapshotCreateTime).Before(*a[j].SnapshotCreateTime)
}

func mostRecentDbSnapshot(snapshots []*rds.DBSnapshot) *rds.DBSnapshot {
	sortedSnapshots := snapshots
	sort.Sort(rdsSnapshotSort(sortedSnapshots))
	return sortedSnapshots[len(sortedSnapshots)-1]
}

// Client-side filtering of snapshots by DB Instance Identifier.
// This is needed for shared or public snapshots i.e. situations where the
// DB Instance Identifier isn't valid for the current account.
func filterSnapshotsByInstanceId(snapshots []*rds.DBSnapshot, instanceId string) []*rds.DBSnapshot {
	results := make([]*rds.DBSnapshot, 0, len(snapshots))

	for _, s := range snapshots {
		if *s.DBInstanceIdentifier == instanceId {
			results = append(results, s)
		}
	}

	log.Printf("[DEBUG] Of %d snapshots, %d had DBInstanceIdentifer == %s", len(snapshots), len(results), instanceId)

	return results
}

func dbSnapshotDescriptionAttributes(d *schema.ResourceData, snapshot *rds.DBSnapshot) error {
	d.SetId(aws.StringValue(snapshot.DBSnapshotIdentifier))
	d.Set("db_instance_identifier", snapshot.DBInstanceIdentifier)
	d.Set("db_snapshot_identifier", snapshot.DBSnapshotIdentifier)
	d.Set("snapshot_type", snapshot.SnapshotType)
	d.Set("storage_type", snapshot.StorageType)
	d.Set("allocated_storage", snapshot.AllocatedStorage)
	d.Set("availability_zone", snapshot.AvailabilityZone)
	d.Set("db_snapshot_arn", snapshot.DBSnapshotArn)
	d.Set("encrypted", snapshot.Encrypted)
	d.Set("engine", snapshot.Engine)
	d.Set("engine_version", snapshot.EngineVersion)
	d.Set("iops", snapshot.Iops)
	d.Set("kms_key_id", snapshot.KmsKeyId)
	d.Set("license_model", snapshot.LicenseModel)
	d.Set("option_group_name", snapshot.OptionGroupName)
	d.Set("port", snapshot.Port)
	d.Set("source_db_snapshot_identifier", snapshot.SourceDBSnapshotIdentifier)
	d.Set("source_region", snapshot.SourceRegion)
	d.Set("status", snapshot.Status)
	d.Set("vpc_id", snapshot.VpcId)
	if snapshot.SnapshotCreateTime != nil {
		d.Set("snapshot_create_time", snapshot.SnapshotCreateTime.Format(time.RFC3339))
	}

	return nil
}
