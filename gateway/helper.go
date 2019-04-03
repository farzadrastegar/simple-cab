package gateway

import (
	"flag"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/farzadrastegar/simple-cab/gateway/config"
)

// appName contains the name of this app.
const appName = "gateway"

// viper config parameters.
var viperAlreadyInitialized	bool

// init runs before any other method.
func init() {
	logger.SetFormatter(&logger.JSONFormatter{})
}

// CreateCabService creates a service through input client.
func CreateCabService(c Client) CabService {
	return c.Connect()
}

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

	// Load configurations into viper.
	config.LoadConfigurationFromBranch(
		viper.GetString("configServerUrl"),
		appName,
		viper.GetString("profile"),
		viper.GetString("configBranch"))

	// Make auto-update available for parameters.
	//go config.StartListener(appName, viper.GetString("amqpServerUrl"), viper.GetString("configEventBus"))
}

// initConfiguration initializes viper with the profile, configServerUrl, and configBranch flags.
func initConfiguration() {
	viper.Reset()

	profile := flag.String("profile", "test", "Environment profile, something similar to spring profiles")
	configServerUrl := flag.String("configServerUrl", "http://localhost:8888", "Address to config server")
	configBranch := flag.String("configBranch", "master", "git branch to fetch configuration from")

	flag.Parse()

	logger.Println("Specified configBranch is " + *configBranch)

	viper.Set("profile", *profile)
	viper.Set("configServerUrl", *configServerUrl)
	viper.Set("configBranch", *configBranch)
}

// GetAppName returns app's name.
func GetAppName() string {
	return appName
}
