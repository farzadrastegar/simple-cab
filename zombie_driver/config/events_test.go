package config

import (
        "fmt"
        "testing"
)
import (
        . "github.com/smartystreets/goconvey/convey"
        "github.com/spf13/viper"
        "gopkg.in/h2non/gock.v1"
)

// Test parameters
const (
        serviceName = "zombie_driver"
        configServerUrl = "localhost:8888"
        httpConfigServerUrl = "http://" + configServerUrl
        profile = "test"
        configBranch = "master"
)

func TestHandleRefreshEvent(t *testing.T) {
        // Configure initial viper values
        viper.Set("configServerUrl", httpConfigServerUrl)
        viper.Set("profile", profile)
        viper.Set("configBranch", configBranch)

        // Mock the expected outgoing request for new config
        defer gock.Off()
        gock.New(httpConfigServerUrl).
                Get(fmt.Sprintf("/%s/%s/%s", serviceName, profile, configBranch)).
                Reply(200).
                BodyString (fmt.Sprintf(`{"name":"%s-%s","profiles":["%s"],"label":null,"version":null,"propertySources":[{"name":"file:/config-repo/%s-%s.yml","source":{"urls.zombieStatus.definition.distance":1000.0}}]}`, serviceName, profile, profile, serviceName, profile))
                //BodyString(`{"name":"driver_location-test","profiles":["test"],"label":null,"version":null,"propertySources":[{"name":"file:/config-repo/driver_location-test.yml","source":{"urls.driverLocations.nsq.topic":"changed_locations"}}]}`)

        Convey("Given a refresh event received, targeting our application", t, func() {
                var body = fmt.Sprintf(`{"type":"RefreshRemoteApplicationEvent","timestamp":1494514362123,"originService":"config-server:%s","destinationService":"%s:**","id":"53e61c71-cbae-4b6d-84bb-d0dcc0aeb4dc"}`, configServerUrl, serviceName)
                //var body = `{"type":"RefreshRemoteApplicationEvent","timestamp":1494514362123,"originService":"config-server:localhost:8888","destinationService":"driver_location:**","id":"53e61c71-cbae-4b6d-84bb-d0dcc0aeb4dc"}`

                Convey("When handled", func() {
                        handleRefreshEvent([]byte(body), serviceName)

                        Convey("Then Viper should have been re-populated with values from Source", func() {
                                So(viper.GetFloat64("urls.zombieStatus.definition.distance"), ShouldEqual, 1000.0)
                        })
                })
        })
}

func TestHandleRefreshEventForOtherApplication(t *testing.T) {

        gock.Intercept()
        defer gock.Off()
        Convey("Given a refresh event received, targeting another application", t, func() {
                var body = fmt.Sprintf(`{"type":"RefreshRemoteApplicationEvent","timestamp":1494514362123,"originService":"config-server:%s","destinationService":"vipservice:**","id":"53e61c71-cbae-4b6d-84bb-d0dcc0aeb4dc"}`, configServerUrl)
                //var body = `{"type":"RefreshRemoteApplicationEvent","timestamp":1494514362123,"originService":"config-server:docker:8888","destinationService":"vipservice:**","id":"53e61c71-cbae-4b6d-84bb-d0dcc0aeb4dc"}`
                Convey("When parsed", func() {
                        handleRefreshEvent([]byte(body), serviceName)

                        Convey("Then no outgoing HTTP requests should have been intercepted", func() {
                                So(gock.HasUnmatchedRequest(), ShouldBeFalse)
                        })
                })
        })
}
