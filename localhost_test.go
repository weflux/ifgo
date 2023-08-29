package localhost

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetLocalIP(t *testing.T) {
	ip, err := IP()
	require.NoError(t, err)
	require.EqualValues(t, "192.168.12.73", ip)
}

func TestMacAddr(t *testing.T) {
	mac, err := MacAddr()
	require.NoError(t, err)
	require.EqualValues(t, "a8:5e:45:e0:f0:62", mac)
}
