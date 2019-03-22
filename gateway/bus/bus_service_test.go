package bus_test

import (
	"fmt"
	"github.com/farzadrastegar/simple-cab/gateway"
	"os/exec"
	"strings"
	"testing"
)

// Ensure producing a message through a real bus (if exists) does not generate any error.
func TestBusService_ProduceNoError(t *testing.T) {
	b := NewBus()
	var strBuilder strings.Builder
	errMsg := "NSQ is not up and running"

	// Check NSQLookupdAddress port is listening
	checkPort1 := fmt.Sprintf("lsof -i -n -P | grep %s | grep LISTEN | tail -n1", b.Params.NSQLookupdAddress)
	out, _ := exec.Command("/bin/bash", "-c", checkPort1).Output()
	strBuilder.Write(out)
	if strings.Index(strBuilder.String(), b.Params.NSQLookupdAddress) == -1 {
		t.Skip(errMsg)
	}

	// Check NSQDAddress is listening
	checkPort2 := fmt.Sprintf("lsof -i -n -P | grep %s | grep LISTEN | tail -n1", b.Params.NSQDAddress)
	out, _ = exec.Command("/bin/bash", "-c", checkPort2).Output()
	strBuilder.Write(out)
	if strings.Index(strBuilder.String(), b.Params.NSQDAddress) == -1 {
		t.Skip(errMsg)
	}

	if strings.Count(strBuilder.String(), "\n") != 2 {
		t.Skip(errMsg)
	}

	s, err := b.Initialize()
	if err != nil {
		t.Fatal(err)
	}

	id := gateway.DriverID(100)
	message := &gateway.Data{Latitude: 35.7766061, Longitude: 51.3963186}

	err = s.Produce(id, message)
	if err != nil {
		t.Fatal(err)
	}
}
