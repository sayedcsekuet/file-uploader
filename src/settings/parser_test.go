package settings

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestInit(t *testing.T) {
	file, _ := ioutil.TempFile("/tmp", "test")
	defer file.Close()
	file.WriteString(`notifications_channel_url: test`)
	tr := Init(file.Name())
	assert.Nil(t, tr)
	assert.Equal(t, "test", NotificationsChannelUrl())
}
