package driver_location

import (
	"flag"
	"fmt"
	"github.com/farzadrastegar/simple-cab/driver_location/config"
	"github.com/spf13/viper"
)

// appName contains the name of this app.
const appName = "driver_location"

// viper config parameters.
var viperAlreadyInitialized	bool

// CreateCabService creates a service through input client.
func CreateCabService(c Client) CabService {
	return c.Connect()
}

//// driverLocationConfigFilename contains configuration parameters in yaml format.
//var driverLocationConfigFilename string
//
//// SetYamlFilename sets the configurations file name.
//func SetConfigFilename(yamlFilename string) {
//	driverLocationConfigFilename = yamlFilename
//}
//
//// GetYamlFilename returns the configurations file name.
//func GetConfigFilename() string {
//	return driverLocationConfigFilename
//}

// LoadConfigurationFromBranch loads config into viper.
func LoadConfigurationFromBranch() {
	// Read config once.
	if viperAlreadyInitialized {
		return
	} else {
		viperAlreadyInitialized = true
	}

	// Initialize configurations.
	initConfiguration()

	config.LoadConfigurationFromBranch(
		viper.GetString("configServerUrl"),
		appName,
		viper.GetString("profile"),
		viper.GetString("configBranch"))
}

// initConfiguration initializes viper with the profile, configServerUrl, and configBranch flags.
func initConfiguration() {
	viper.Reset()

	profile := flag.String("profile", "test", "Environment profile, something similar to spring profiles")
	configServerUrl := flag.String("configServerUrl", "http://localhost:8888", "Address to config server")
	configBranch := flag.String("configBranch", "master", "git branch to fetch configuration from")

	flag.Parse()

	fmt.Println("Specified configBranch is " + *configBranch)

	viper.Set("profile", *profile)
	viper.Set("configServerUrl", *configServerUrl)
	viper.Set("configBranch", *configBranch)
}
