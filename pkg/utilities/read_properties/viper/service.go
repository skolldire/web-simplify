package viper

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/skolldire/web-simplify/pkg/utilities/app_profile"
	rp "github.com/skolldire/web-simplify/pkg/utilities/read_properties"
	"github.com/spf13/viper"
)

const errorLoadingConfiguration = "error loading configuration - %w"

type service struct {
	propertyFiles []string
	path          string
}

func NewService() *service {
	return &service{
		propertyFiles: []string{
			"application.yaml",
			"application-local.yaml",
			"application-pt.yaml",
			"application-qa.yaml",
			"application-prod.yaml",
		},
		path: getConfigPath(),
	}
}

// Apply loads the configuration and maps it to a struct
func (s *service) Apply() (rp.Config, error) {
	if err := s.validateRequiredFiles(); err != nil {
		return rp.Config{}, err
	}

	mergedConfig, err := s.loadAndMergeConfigs()
	if err != nil {
		return rp.Config{}, fmt.Errorf(errorLoadingConfiguration, err)
	}

	return s.mapConfigToStruct(mergedConfig)
}

// validateRequiredFiles ensures all necessary property files are present
func (s *service) validateRequiredFiles() error {
	missingFiles := getMissingFiles(s.propertyFiles, listFiles(s.path))
	if len(missingFiles) > 0 {
		return fmt.Errorf("missing required files: %v", missingFiles)
	}
	return nil
}

// loadAndMergeConfigs loads and merges the base and environment-specific configurations
func (s *service) loadAndMergeConfigs() (*viper.Viper, error) {
	baseConfig, err := loadConfig(s.path, "application")
	if err != nil {
		return nil, err
	}

	envConfig, err := loadConfig(s.path, s.getPropertyFileName())
	if err != nil {
		return nil, err
	}

	if err := baseConfig.MergeConfigMap(envConfig.AllSettings()); err != nil {
		return nil, fmt.Errorf("failed to merge configurations: %w", err)
	}

	return baseConfig, nil
}

// mapConfigToStruct decodes and validates the configuration into a struct
func (s *service) mapConfigToStruct(v *viper.Viper) (rp.Config, error) {
	configMap, err := unmarshalConfig(v)
	if err != nil {
		return rp.Config{}, err
	}

	processedConfig, err := processConfigValues(configMap)
	if err != nil {
		return rp.Config{}, err
	}

	return decodeToStruct(processedConfig)
}

// getPropertyFileName determines the name of the environment-specific configuration file
func (s *service) getPropertyFileName() string {
	scopeFile := fmt.Sprintf("application-%s", app_profile.GetScopeValue())
	profileFile := fmt.Sprintf("application-%s", app_profile.GetProfileByScope())
	files := listFiles(s.path)

	if contains(files, scopeFile) {
		return scopeFile
	}
	return profileFile
}

// Helper Functions
func getConfigPath() string {
	if path := os.Getenv("CONF_DIR"); path != "" {
		return path
	}
	return "kit/config"
}

func loadConfig(path, filename string) (*viper.Viper, error) {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName(filename)
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	return v, nil
}

func unmarshalConfig(v *viper.Viper) (map[string]interface{}, error) {
	var configMap map[string]interface{}
	if err := v.Unmarshal(&configMap); err != nil {
		return nil, err
	}
	return configMap, nil
}

func processConfigValues(config map[string]interface{}) (map[string]interface{}, error) {
	for key, value := range config {
		switch v := value.(type) {
		case map[string]interface{}:
			processed, err := processConfigValues(v)
			if err != nil {
				return nil, err
			}
			config[key] = processed
		case []interface{}:
			processed, err := processSliceValues(v)
			if err != nil {
				return nil, err
			}
			config[key] = processed
		case string:
			config[key] = resolveEnvValue(v)
		}
	}
	return config, nil
}

func processSliceValues(slice []interface{}) ([]interface{}, error) {
	for i, elem := range slice {
		switch v := elem.(type) {
		case map[string]interface{}:
			processed, err := processConfigValues(v)
			if err != nil {
				return nil, err
			}
			slice[i] = processed
		case []interface{}:
			processed, err := processSliceValues(v)
			if err != nil {
				return nil, err
			}
			slice[i] = processed
		case string:
			slice[i] = resolveEnvValue(v)
		}
	}
	return slice, nil
}

func resolveEnvValue(value string) string {
	if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
		return os.Getenv(strings.Trim(value, "${}"))
	}
	return value
}

func decodeToStruct(config map[string]interface{}) (rp.Config, error) {
	var result rp.Config
	if err := mapstructure.Decode(config, &result); err != nil {
		return rp.Config{}, err
	}
	return result, nil
}

func listFiles(path string) []string {
	files, err := os.ReadDir(path)
	if err != nil {
		panic(fmt.Errorf("error reading configuration directory: %w", err))
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	return fileNames
}

func getMissingFiles(required, available []string) []string {
	var missing []string
	for _, file := range required {
		if !contains(available, file) {
			missing = append(missing, file)
		}
	}
	return missing
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
