package main_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/url"
	"time"

	"github.com/farzadrastegar/simple-cab/driver_location"
	DL_bus "github.com/farzadrastegar/simple-cab/driver_location/bus"
	DL_http "github.com/farzadrastegar/simple-cab/driver_location/http"
	DL_redis "github.com/farzadrastegar/simple-cab/driver_location/redis"
	"github.com/farzadrastegar/simple-cab/gateway"
	GW_bus "github.com/farzadrastegar/simple-cab/gateway/bus"
	GW_client "github.com/farzadrastegar/simple-cab/gateway/client"
	GW_http "github.com/farzadrastegar/simple-cab/gateway/http"
	"github.com/farzadrastegar/simple-cab/zombie_driver"
	ZD_client "github.com/farzadrastegar/simple-cab/zombie_driver/client"
	ZD_http "github.com/farzadrastegar/simple-cab/zombie_driver/http"

	"testing"
)

// Constants for tests.
var (
	sendingLocationsStartedAtThisTime = time.Date(2019, 01, 01, 00, 00, 00, 00, time.UTC)
	sendingLocationsIntervals         = 5 //seconds
	numberOfLocationsFetched          = 4
)

// Servers variables.
var (
	gwSrv    *GW_http.Server
	zdSrv    *ZD_http.Server
	dlSrv    *DL_http.Server
	dlClient *DL_redis.Client

	gwHttpClient *GW_http.Client
	//zdHttpClient *ZD_http.Client
	//dlHttpClient *DL_http.Client
)

var _ = Describe("Store locations and report zombie status", func() {

	var doItOnce = 0

	BeforeSuite(func() {
		// Set configurations parameters.
		gateway.SetConfigFilename("config_gateway.yaml")
		zombie_driver.SetConfigFilename("config_zombie_driver.yaml")
		driver_location.SetConfigFilename("config_driver_location.yaml")
		//})
		//
		//BeforeEach(func() {
		// Open servers.
		_, gwSrv = OpenGatewayServer()
		_, zdSrv = OpenZombieDriverServer()
		dlClient, dlSrv = OpenDriverLocationServer()

		// Open http clients.
		gwHttpClient = OpenGatewayHttpClient()
		//zdHttpClient = OpenZombieDriverHttpClient()
		//dlHttpClient = OpenDriverLocationHttpClient()
	})

	//AfterEach(func() {
	AfterSuite(func() {
		// Close servers.
		if gwSrv != nil {
			Expect(gwSrv.Close()).NotTo(HaveOccurred())
		}
		if zdSrv != nil {
			Expect(zdSrv.Close()).NotTo(HaveOccurred())
		}
		if dlSrv != nil {
			Expect(dlSrv.Close()).NotTo(HaveOccurred())
			Expect(dlClient.Close()).NotTo(HaveOccurred())
		}
	})

	JustBeforeEach(func() {
		if doItOnce != 1 {
			flushDB()
			doItOnce = 1
		}
	})

	Describe("Driver is zombie", func() {
		Context("Driver has moved to some locations but the driven distance does not satisfy the false zombie definition", func() {

			It("should return driver id and zombie=true", func() {
				// Prepare expected output.
				id := 1234
				zombie := true
				driverId := gateway.DriverID(id)
				expectedOut := gateway.Status{ID: driverId, Zombie: zombie}

				idStr := fmt.Sprintf("%d", id)

				By("Sending a bunch of PATCH requests through Gateway to store locations")
				timeNow, err := sendDriverLocations(idStr, zombie)
				Expect(err).NotTo(HaveOccurred())

				By("Making sure the sent locations can be fetched through driver_location")
				locations, err := getDriverLocations(idStr, timeNow)
				Expect(err).NotTo(HaveOccurred())
				expectLocationsToEqual(locations, zombie)

				By("Getting zombie status through Gatway's GET requests and comparing with expected output, i.e. zombie=true")
				out, err := getDriverStatus(idStr)
				Expect(err).NotTo(HaveOccurred())

				Expect(*out).To(Equal(expectedOut))
			})
		})
	})

	Describe("Driver is not zombie", func() {
		Context("Driver has moved to some locations and the driven distance satisfies the false zombie definition", func() {

			It("should return driver id and zombie=false", func() {
				// Prepare expected output.
				id := 56789
				zombie := false
				driverId := gateway.DriverID(id)
				expectedOut := gateway.Status{ID: driverId, Zombie: zombie}

				idStr := fmt.Sprintf("%d", id)

				By("Sending a bunch of PATCH requests through Gateway to store locations")
				timeNow, err := sendDriverLocations(idStr, zombie)
				Expect(err).NotTo(HaveOccurred())

				By("Making sure the sent locations can be fetched through driver_location")
				locations, err := getDriverLocations(idStr, timeNow)
				Expect(err).NotTo(HaveOccurred())
				expectLocationsToEqual(locations, zombie)

				By("Getting zombie status through zombie_driver and comparing with expected output, i.e. zombie=false")
				out, err := getDriverStatus(idStr)
				Expect(err).NotTo(HaveOccurred())

				Expect(*out).To(Equal(expectedOut))
			})
		})
	})

})

// getNonZombieLocations returns locations of a non-zombie driver.
func getNonZombieLocations() []gateway.Data {
	var locations []gateway.Data

	locations = append(locations, gateway.Data{Latitude: 35.721758, Longitude: 51.399568})
	locations = append(locations, gateway.Data{Latitude: 35.728064, Longitude: 51.397552}) //720m
	locations = append(locations, gateway.Data{Latitude: 35.733815, Longitude: 51.398257}) //640m
	locations = append(locations, gateway.Data{Latitude: 35.745386, Longitude: 51.401593}) //1320m
	locations = append(locations, gateway.Data{Latitude: 35.757446, Longitude: 51.403617}) //1350m
	locations = append(locations, gateway.Data{Latitude: 35.757731, Longitude: 51.409763}) //560m

	return locations
}

// getZombieLocations returns locations of a zombie driver.
func getZombieLocations() []gateway.Data {
	var locations []gateway.Data

	locations = append(locations, gateway.Data{Latitude: 35.721905, Longitude: 51.400402})
	locations = append(locations, gateway.Data{Latitude: 35.721766, Longitude: 51.401303}) //80m
	locations = append(locations, gateway.Data{Latitude: 35.723926, Longitude: 51.408166}) //660m
	locations = append(locations, gateway.Data{Latitude: 35.727176, Longitude: 51.413998}) //640m
	locations = append(locations, gateway.Data{Latitude: 35.728569, Longitude: 51.418502}) //440m

	return locations
}

// sendPatchRequest sends an http PATCH request to store a location through the Gateway server.
func sendPatchRequest(id string, loc *gateway.Data) error {
	var u url.URL
	u.Scheme = "HTTP"
	u.Host = fmt.Sprintf("localhost:%d", gwSrv.Port())

	gwHttpClient.URL = u
	return gwHttpClient.Connect().StoreLocation(id, loc)
}

type patchResponse struct {
	Err string `json:"err,omitempty"`
}

// sendDriverLocations sends driver's locations one by one through PATCH requests.
func sendDriverLocations(id string, zombie bool) (int64, error) {
	var locations []gateway.Data

	// Get sample locations.
	if zombie {
		locations = getZombieLocations()
		Expect(len(locations) > numberOfLocationsFetched).To(Equal(true))
	} else {
		locations = getNonZombieLocations()
		Expect(len(locations) > numberOfLocationsFetched).To(Equal(true))
	}

	// Send an http PATCH request for each location element.
	dlClient.Interval = 0
	dlClient.Now = sendingLocationsStartedAtThisTime.Unix()
	timeNow := sendingLocationsStartedAtThisTime.Unix()
	for _, elem := range locations {
		// Record time now.
		timeNow += int64(sendingLocationsIntervals)

		// Send request.
		err := sendPatchRequest(id, &elem)
		if err != nil {
			return 0, err
		}

		fmt.Println(">>>>> taking 1-second break")
		time.Sleep(1 * time.Second)

		dlClient.Interval = int64(sendingLocationsIntervals)
	}

	return timeNow, nil
}

// getDriverLocations returns last numberOfLocationsFetched locations.
func getDriverLocations(id string, timeNow int64) ([]gateway.Data, error) {
	// Set time.
	dlClient.Interval = 0
	secondsBackFromNow := numberOfLocationsFetched * sendingLocationsIntervals
	dlClient.Now = timeNow //sendingLocationsStartedAtThisTime.Unix() + int64(secondsBackFromNow)

	// Get locations based on numberOfLocationsFetched.
	minutes := float64(secondsBackFromNow) / 60.0
	locations, err := dlClient.Connect().GetDriverLocations(id, minutes)
	if err != nil {
		return nil, err
	}
	Expect(len(locations.Locations)).To(Equal(numberOfLocationsFetched))

	// Convert locations' type to []gateway.Data.
	var returnedLoc []gateway.Data
	for _, elem := range locations.Locations {
		returnedLoc = append(returnedLoc, gateway.Data{Latitude: elem.Latitude, Longitude: elem.Longitude})
	}

	return returnedLoc, nil
}

// expectLocationsToEqual makes sure the last numberOfLocationsFetched original locations match input loc.
func expectLocationsToEqual(loc []gateway.Data, zombie bool) {
	var locations []gateway.Data

	// Get sample locations.
	if zombie {
		locations = getZombieLocations()
	} else {
		locations = getNonZombieLocations()
	}

	// get last numberOfLocationsFetched locations.
	locations = locations[len(locations)-numberOfLocationsFetched : len(locations)]
	Expect(len(locations)).To(Equal(numberOfLocationsFetched))

	Expect(locations).To(Equal(loc))
}

// sendGetRequest sends an http GET request through Gateway.
func sendGetRequest(id string) (*gateway.Status, error) {
	var u url.URL
	u.Scheme = "HTTP"
	u.Host = fmt.Sprintf("localhost:%d", gwSrv.Port())

	gwHttpClient.URL = u
	return gwHttpClient.Connect().CheckZombieStatus(id)
}

// getDriverStatus returns the output of sendGetRequest.
func getDriverStatus(id string) (*gateway.Status, error) {
	return sendGetRequest(id)
}

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "use gateway, zombie_driver, and driver_location microservices")
}

// flushDB executes the flushdb command in redis.
func flushDB() {
	//if dbPortIsListening() {
	Expect(dlClient.GetDB().FlushDB().Err()).NotTo(HaveOccurred())
	//}
}

// OpenDriverLocationServer returns a working driver_location server and client.
func OpenDriverLocationServer() (*DL_redis.Client, *DL_http.Server) {
	// Create a cab service through redis.
	c := DL_redis.NewClient()
	if err := c.Open(); err != nil {
		panic(err)
	}
	cabService := driver_location.CreateCabService(c)

	// Create a bus consumer.
	b := DL_bus.NewBus()
	busService, err := b.Initialize(cabService)
	if err != nil {
		panic(err)
	}
	err = busService.Consume()
	if err != nil {
		panic(err)
	}

	// Attach cabService to HTTP handler.
	h := DL_http.NewDataHandler()
	h.CabService = cabService

	// Start an HTTP server.
	srv := DL_http.NewServer()
	srv.Handler = &DL_http.Handler{DataHandler: h}
	if err := srv.Open(); err != nil {
		panic(err)
	}

	return c, srv
}

// OpenZombieDriverServer returns a working zombie_driver server.
func OpenZombieDriverServer() (*ZD_client.Client, *ZD_http.Server) {
	// Create a client for managing services.
	c := ZD_client.NewClient()

	// Create cab services.
	s := zombie_driver.CreateCabService(c)

	// Attach services to HTTP handler.
	h := ZD_http.NewDataHandler()
	h.CabService = s

	// Start an HTTP server.
	srv := ZD_http.NewServer()
	srv.Handler = &ZD_http.Handler{DataHandler: h}
	if err := srv.Open(); err != nil {
		panic(err)
	}

	return c, srv
}

// OpenGatewayServer returns a working gateway server.
func OpenGatewayServer() (*GW_client.Client, *GW_http.Server) {
	// Create a bus service.
	b := GW_bus.NewBus()
	busService, err := b.Initialize()
	if err != nil {
		panic(err)
	}

	// Create a client and attach bus service to it.
	c := GW_client.NewClient()
	c.Handler = GW_client.NewHandler()
	c.Handler.BusService = busService

	// Create cab services.
	s := gateway.CreateCabService(c)

	// Attach services to HTTP handler.
	h := GW_http.NewDataHandler()
	h.CabService = s

	// Start an HTTP server.
	srv := GW_http.NewServer()
	srv.Handler = &GW_http.Handler{DataHandler: h}
	if err := srv.Open(); err != nil {
		panic(err)
	}

	return c, srv
}

// OpenGatewayHttpClient returns a gatway http client.
func OpenGatewayHttpClient() *GW_http.Client {
	return GW_http.NewClient()
}

// OpenZombieDriverHttpClient returns a zombie_driver http client.
func OpenZombieDriverHttpClient() *ZD_http.Client {
	return ZD_http.NewClient()
}

// OpenDriverLocationHttpClient returns a driver_location http client.
func OpenDriverLocationHttpClient() *DL_http.Client {
	return DL_http.NewClient()
}
