package core_test

import (
	"strings"
	"testing"

	"github.com/andreyvit/diff"
	mole "github.com/thalesfsp/mole/core"
	"github.com/thalesfsp/mole/fsutils"
)

const expectedInstance string = `id = "id1"
tunnel-type = ""
verbose = false
insecure = false
detach = false
key = ""
key-value = ""
keep-alive-interval = 0
connection-retries = 0
wait-and-retry = 0
ssh-agent = ""
timeout = 0
ssh-config = ""
rpc = false
rpc-address = ""

[server]
  user = ""
  host = ""
  port = ""`

const expectedMultipleInstances string = `[instances]
  [instances.id1]
    id = "id1"
    tunnel-type = ""
    verbose = false
    insecure = false
    detach = false
    key = ""
    key-value = ""
    keep-alive-interval = 0
    connection-retries = 0
    wait-and-retry = 0
    ssh-agent = ""
    timeout = 0
    ssh-config = ""
    rpc = false
    rpc-address = ""
    [instances.id1.server]
      user = ""
      host = ""
      port = ""
  [instances.id2]
    id = "id2"
    tunnel-type = ""
    verbose = false
    insecure = false
    detach = false
    key = ""
    key-value = ""
    keep-alive-interval = 0
    connection-retries = 0
    wait-and-retry = 0
    ssh-agent = ""
    timeout = 0
    ssh-config = ""
    rpc = false
    rpc-address = ""
    [instances.id2.server]
      user = ""
      host = ""
      port = ""`

func TestFormatRuntimeToML(t *testing.T) {
	instances := []mole.Runtime{
		mole.Runtime{Id: "id1"},
		mole.Runtime{Id: "id2"},
	}

	runtimes := mole.InstancesRuntime(instances)

	tests := []struct {
		formatter mole.Formatter
		expected  string
	}{
		{formatter: mole.Runtime{Id: "id1"}, expected: expectedInstance},
		{formatter: runtimes, expected: expectedMultipleInstances},
	}

	for _, test := range tests {
		out, err := test.formatter.Format("toml")

		if err != nil {
			t.Errorf(err.Error())
		}

		if a, e := strings.TrimSpace(out), strings.TrimSpace(test.expected); a != e {
			t.Errorf("Result not as expected:\n%v", diff.LineDiff(e, a))
		}
	}
}

func TestClientRunning(t *testing.T) {
	id := "test-client-running"

	// Mock the pid file using the process id of the program running the test
	_, err := fsutils.CreateInstanceDir(id)
	if err != nil {
		t.Errorf(err.Error())
	}

	conf := &mole.Configuration{Id: id}
	client := mole.Client{Conf: conf}

	running, err := client.Running()
	if err != nil {
		t.Errorf(err.Error())
	}

	if !running {
		t.Errorf("client was supposed to be running")
	}
}
