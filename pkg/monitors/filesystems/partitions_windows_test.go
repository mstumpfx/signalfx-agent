// +build windows

package filesystems

import (
	"testing"

	gopsutil "github.com/shirou/gopsutil/disk"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllMounts_ShouldInclude_gopsutil_Mounts(t *testing.T) {
	logger := logrus.WithFields(logrus.Fields{"monitorType": monitorType})

	// Drive and folder mounts.
	got := (&Monitor{logger: logger}).getAllMounts()
	require.NotEmpty(t, got, "failed to find any mount points")

	// Mounts from gopsutil are for drives only.
	want, err := getGopsutilMounts()
	require.NoError(t, err)

	require.NotEmpty(t, want, "failed to find any mount points using gopsutil")

	// Asserting `got` getAllMounts() mounts superset of `want` gopsutil mounts.
	assert.Subset(t, got, want)
}

func TestNewPartitionStats_SameAs_gopsutil_PartitionStats(t *testing.T) {
	// Partition stats from gopsutil are for drive mounts only.
	gopsutilStats, err := gopsutil.Partitions(true)
	require.NoError(t, err)

	require.NotEmpty(t, gopsutilStats, "failed to find any partition stats using gopsutil")

	logger := logrus.WithFields(logrus.Fields{"monitorType": monitorType})
	monitor := Monitor{logger: logger}

	// All drive and folder mounts.
	allMounts := monitor.getAllMounts()

	require.NotEmpty(t, allMounts, "failed to find any mount points")

	var got gopsutil.PartitionStat
	for _, want := range gopsutilStats {
		got, err = newPartitionStats(want.Mountpoint)
		require.NoError(t, err)

		// Asserting `got` newPartitionStats() stats equal `want` gopsutil stats.
		assert.Equal(t, got, want)
	}
}

func TestGetPartitions_ShouldInclude_gopsutil_PartitionStats(t *testing.T) {
	// Partition stats from gopsutil are for drive mounts only.
	want, err := gopsutil.Partitions(true)
	require.NoError(t, err)

	require.NotEmpty(t, want, "failed to find any partition stats using gopsutil")

	logger := logrus.WithFields(logrus.Fields{"monitorType": monitorType})
	monitor := Monitor{logger: logger}

	var got []gopsutil.PartitionStat
	// Partition stats for drive and folder mounts.
	got, err = monitor.getPartitions(true)
	require.NoError(t, err)

	require.NotEmpty(t, got, "failed to find any partition stats")

	// Asserting `got` getPartitions stats superset of `want` gopsutil stats.
	assert.Subset(t, got, want)
}

func getGopsutilMounts() ([]string, error) {
	partitionsStats, err := gopsutil.Partitions(true)
	if err != nil {
		return nil, err
	}

	mounts := make([]string, 0)
	for _, stats := range partitionsStats {
		mounts = append(mounts, stats.Mountpoint)
	}

	return mounts, nil
}