package redis_test

import (
	"fmt"
	"github.com/farzadrastegar/simple-cab/driver_location"
	"github.com/farzadrastegar/simple-cab/driver_location/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

var port string
var paramsReady bool
var randSeed = time.Date(2019, 01, 01, 02, 00, 00, 00, time.UTC)

type Location struct {
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Updated_at string  `json:"updated_at,omitempty"`
}

func genRand() int {
	return rand.Intn(1000000)
}

func OnAvailableRedisIt(description string, f interface{}) {
	if dbPortIsListening() {
		It(description, f)
	} else {
		//	PIt(description, f)
	}
}

var _ = Describe("basic functionality", func() {
	var client *Client
	entryLen := 12
	pathID := "1000"
	latitude := 1.1
	longitude := 2.2

	BeforeSuite(func() {
		rand.Seed(randSeed.Unix())
		flushDB()
	})

	AfterSuite(func() {
	})

	BeforeEach(func() {
		driver_location.SetConfigFilename("../cmd/config.yaml")
	})

	AfterEach(func() {
		if client != nil {
			Expect(client.Close()).NotTo(HaveOccurred())
		}
	})

	OnAvailableRedisIt("should store", func() {
		client = MustOpenClient()
		client.Now = randSeed.Unix()
		s := client.Connect()

		for i := 0; i < entryLen; i++ {
			loc := driver_location.Location{Latitude: latitude + float64(i), Longitude: longitude + float64(i)}
			err := s.StoreLocation(pathID, &loc)
			Expect(err).NotTo(HaveOccurred())
			client.Now += 5
		}

		vals, err := client.GetDB().XRange("drivers:"+pathID, "-", "+").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(len(vals)).To(Equal(entryLen))
	})

	OnAvailableRedisIt("should query", func() {
		client = MustOpenClient()
		client.Now = randSeed.Unix()
		s := client.Connect()
		minutes := 1.0
		client.Now += 60

		locations, err := s.GetDriverLocations(pathID, minutes)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(locations.Locations)).To(Equal(entryLen))

		for i := 0; i < len(locations.Locations); i++ {
			Expect(locations.Locations[i].Latitude).To(Equal(latitude + float64(i)))
			Expect(locations.Locations[i].Longitude).To(Equal(longitude + float64(i)))
		}
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "database")
}

func readDBPort() {
	driver_location.SetConfigFilename("../cmd/config.yaml")

	// Read server address and port from config.yaml.
	logger := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
	configHandler := config.NewConfig(logger)
	configHandler.ReadYaml(driver_location.GetConfigFilename())
	port = configHandler.GetYamlValueStr("database", "port")
}

func dbPortIsListening() bool {
	// Read port (if needed)
	if !paramsReady {
		readDBPort()
		paramsReady = true
	}

	// Port must be available here.
	if port == "" {
		return false
	}

	// Check DB port is listening.
	checkPort := fmt.Sprintf("lsof -i -n -P | grep %s | grep LISTEN | tail -n1", port)
	cmdOut, _ := exec.Command("/bin/bash", "-c", checkPort).Output()
	portIsListening := true
	var strBuilder strings.Builder
	strBuilder.Write(cmdOut)
	if strings.Index(strBuilder.String(), port) == -1 {
		portIsListening = false
	}

	return portIsListening
}

func flushDB() {
	if dbPortIsListening() {
		//driver_location.SetConfigFilename("../cmd/config.yaml")
		client := MustOpenClient()
		Expect(client.GetDB().FlushDB().Err()).NotTo(HaveOccurred())
		Expect(client.Close()).NotTo(HaveOccurred())
	}
}
