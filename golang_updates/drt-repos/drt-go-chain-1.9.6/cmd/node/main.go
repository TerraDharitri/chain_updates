package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	logger "github.com/TerraDharitri/drt-go-chain-logger"
	"github.com/TerraDharitri/drt-go-chain-logger/file"
	"github.com/TerraDharitri/drt-go-chain/cmd/node/factory"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/config/overridableConfig"
	"github.com/TerraDharitri/drt-go-chain/node"
	"github.com/klauspost/cpuid/v2"
	"github.com/urfave/cli"
	// test point 1 for custom profiler
)

const (
	defaultLogsPath = "logs"
	logFilePrefix   = "drt-go-chain"
)

var (
	memoryBallastObject []byte
	nodeHelpTemplate    = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}
VERSION:
   {{.Version}}
   {{end}}
`
)

// appVersion should be populated at build time using ldflags
// Usage examples:
// linux/mac:
//
//	go build -v -ldflags="-X main.appVersion=$(git describe --tags --long --dirty)"
//
// windows:
//
//	for /f %i in ('git describe --tags --long --dirty') do set VERS=%i
//	go build -v -ldflags="-X main.appVersion=%VERS%"
var appVersion = common.UnVersionedAppString

func main() {
	_ = logger.SetDisplayByteSlice(logger.ToHexShort)
	log := logger.GetOrCreate("main")

	// test point 2 for custom profiler

	app := cli.NewApp()
	cli.AppHelpTemplate = nodeHelpTemplate
	app.Name = "Dharitri Node CLI App"
	machineID := core.GetAnonymizedMachineID(app.Name)

	baseVersion := fmt.Sprintf("%s/%s/%s-%s", appVersion, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	app.Version = fmt.Sprintf("%s/%s", baseVersion, machineID)
	app.Usage = "This is the entry point for starting a new Dharitri node - the app will start after the genesis timestamp"
	app.Flags = getFlags()
	app.Authors = []cli.Author{
		{
			Name:  "The Dharitri Team",
			Email: "contact@dharitri.org",
		},
	}

	app.Action = func(c *cli.Context) error {
		return startNodeRunner(c, log, baseVersion, app.Version)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func startNodeRunner(c *cli.Context, log logger.Logger, baseVersion string, version string) error {
	flagsConfig := getFlagsConfig(c, log)

	fileLogging, errLogger := attachFileLogger(log, flagsConfig)
	if errLogger != nil {
		return errLogger
	}

	cfgs, errCfg := readConfigs(c, log)
	if errCfg != nil {
		return errCfg
	}

	errCfgOverride := overridableConfig.OverrideConfigValues(cfgs.PreferencesConfig.Preferences.OverridableConfigTomlValues, cfgs)
	if errCfgOverride != nil {
		return errCfgOverride
	}

	if !check.IfNil(fileLogging) {
		timeLogLifeSpan := time.Second * time.Duration(cfgs.GeneralConfig.Logs.LogFileLifeSpanInSec)
		sizeLogLifeSpanInMB := uint64(cfgs.GeneralConfig.Logs.LogFileLifeSpanInMB)

		err := fileLogging.ChangeFileLifeSpan(timeLogLifeSpan, sizeLogLifeSpanInMB)
		if err != nil {
			return err
		}
	}

	err := applyFlags(c, cfgs, flagsConfig, log)
	if err != nil {
		return err
	}

	memBallastValue := c.GlobalUint64(memBallast.Name)
	if memBallastValue > 0 {
		// memory ballast is an optimization for golang's garbage collector. If set to a high value, it can decrease
		// the number of times when GC performs STW processes, that results is a better performance over high load
		memoryBallastObject = make([]byte, memBallastValue*core.MegabyteSize)
		log.Debug("initialized memory ballast object", "size", core.ConvertBytes(uint64(len(memoryBallastObject))))
	}

	cfgs.FlagsConfig.BaseVersion = baseVersion
	cfgs.FlagsConfig.Version = version

	err = checkHardwareRequirements(cfgs.GeneralConfig.HardwareRequirements)
	if err != nil {
		return fmt.Errorf("Hardware Requirements checks failed: %s", err.Error())
	}

	nodeRunner, errRunner := node.NewNodeRunner(cfgs)
	if errRunner != nil {
		return errRunner
	}

	err = nodeRunner.Start()
	if err != nil {
		log.Error(err.Error())
	}

	if !check.IfNil(fileLogging) {
		err = fileLogging.Close()
		log.LogIfError(err)
	}

	return err
}

func readConfigs(ctx *cli.Context, log logger.Logger) (*config.Configs, error) {
	log.Trace("reading Configs")

	configurationPaths := &config.ConfigurationPathsHolder{}

	configurationPaths.MainConfig = ctx.GlobalString(configurationFile.Name)
	generalConfig, err := common.LoadMainConfig(configurationPaths.MainConfig)
	if err != nil {
		return nil, err
	}
	log.Debug("config", "file", configurationPaths.MainConfig)

	configurationPaths.ApiRoutes = ctx.GlobalString(configurationApiFile.Name)
	apiRoutesConfig, err := common.LoadApiConfig(configurationPaths.ApiRoutes)
	if err != nil {
		return nil, err
	}
	log.Debug("config", "file", configurationPaths.ApiRoutes)

	configurationPaths.Economics = ctx.GlobalString(configurationEconomicsFile.Name)
	economicsConfig, err := common.LoadEconomicsConfig(configurationPaths.Economics)
	if err != nil {
		return nil, err
	}
	log.Debug("config", "file", configurationPaths.Economics)

	configurationPaths.SystemSC = ctx.GlobalString(configurationSystemSCFile.Name)
	systemSCConfig, err := common.LoadSystemSmartContractsConfig(configurationPaths.SystemSC)
	if err != nil {
		return nil, err
	}
	log.Debug("config", "file", configurationPaths.SystemSC)

	configurationPaths.Ratings = ctx.GlobalString(configurationRatingsFile.Name)
	ratingsConfig, err := common.LoadRatingsConfig(configurationPaths.Ratings)
	if err != nil {
		return nil, err
	}
	log.Debug("config", "file", configurationPaths.Ratings)

	configurationPaths.Preferences = ctx.GlobalString(configurationPreferencesFile.Name)
	preferencesConfig, err := common.LoadPreferencesConfig(configurationPaths.Preferences)
	if err != nil {
		return nil, err
	}
	log.Debug("config", "file", configurationPaths.Preferences)

	configurationPaths.External = ctx.GlobalString(externalConfigFile.Name)
	externalConfig, err := common.LoadExternalConfig(configurationPaths.External)
	if err != nil {
		return nil, err
	}
	log.Debug("config", "file", configurationPaths.External)

	configurationPaths.MainP2p = ctx.GlobalString(p2pConfigurationFile.Name)
	mainP2PConfig, err := common.LoadP2PConfig(configurationPaths.MainP2p)
	if err != nil {
		return nil, err
	}
	log.Debug("config", "file", configurationPaths.MainP2p)

	configurationPaths.FullArchiveP2p = ctx.GlobalString(fullArchiveP2PConfigurationFile.Name)
	fullArchiveP2PConfig, err := common.LoadP2PConfig(configurationPaths.FullArchiveP2p)
	if err != nil {
		return nil, err
	}
	log.Debug("config", "file", configurationPaths.FullArchiveP2p)

	configurationPaths.Epoch = ctx.GlobalString(epochConfigurationFile.Name)
	epochConfig, err := common.LoadEpochConfig(configurationPaths.Epoch)
	if err != nil {
		return nil, err
	}
	log.Debug("config", "file", configurationPaths.Epoch)

	configurationPaths.RoundActivation = ctx.GlobalString(roundConfigurationFile.Name)
	roundConfig, err := common.LoadRoundConfig(configurationPaths.RoundActivation)
	if err != nil {
		return nil, err
	}
	log.Debug("config", "file", configurationPaths.RoundActivation)

	var nodesSetup config.NodesConfig
	configurationPaths.Nodes = ctx.GlobalString(nodesFile.Name)
	err = core.LoadJsonFile(&nodesSetup, configurationPaths.Nodes)
	if err != nil {
		return nil, err
	}
	log.Debug("config", "file", configurationPaths.Nodes)

	if ctx.IsSet(port.Name) {
		mainP2PConfig.Node.Port = ctx.GlobalString(port.Name)
	}
	if ctx.IsSet(fullArchivePort.Name) {
		fullArchiveP2PConfig.Node.Port = ctx.GlobalString(fullArchivePort.Name)
	}
	if ctx.IsSet(destinationShardAsObserver.Name) {
		preferencesConfig.Preferences.DestinationShardAsObserver = ctx.GlobalString(destinationShardAsObserver.Name)
	}
	if ctx.IsSet(nodeDisplayName.Name) {
		preferencesConfig.Preferences.NodeDisplayName = ctx.GlobalString(nodeDisplayName.Name)
	}
	if ctx.IsSet(identityFlagName.Name) {
		preferencesConfig.Preferences.Identity = ctx.GlobalString(identityFlagName.Name)
	}

	return &config.Configs{
		GeneralConfig:            generalConfig,
		ApiRoutesConfig:          apiRoutesConfig,
		EconomicsConfig:          economicsConfig,
		SystemSCConfig:           systemSCConfig,
		RatingsConfig:            ratingsConfig,
		PreferencesConfig:        preferencesConfig,
		ExternalConfig:           externalConfig,
		MainP2pConfig:            mainP2PConfig,
		FullArchiveP2pConfig:     fullArchiveP2PConfig,
		ConfigurationPathsHolder: configurationPaths,
		EpochConfig:              epochConfig,
		RoundConfig:              roundConfig,
		NodesConfig:              &nodesSetup,
	}, nil
}

func attachFileLogger(log logger.Logger, flagsConfig *config.ContextFlagsConfig) (factory.FileLoggingHandler, error) {
	var fileLogging factory.FileLoggingHandler
	var err error
	if flagsConfig.SaveLogFile {
		args := file.ArgsFileLogging{
			WorkingDir:      flagsConfig.LogsDir,
			DefaultLogsPath: defaultLogsPath,
			LogFilePrefix:   logFilePrefix,
		}
		fileLogging, err = file.NewFileLogging(args)
		if err != nil {
			return nil, fmt.Errorf("%w creating a log file", err)
		}
	}

	err = logger.SetDisplayByteSlice(logger.ToHex)
	log.LogIfError(err)
	logger.ToggleCorrelation(flagsConfig.EnableLogCorrelation)
	logger.ToggleLoggerName(flagsConfig.EnableLogName)
	logLevelFlagValue := flagsConfig.LogLevel

	err = logger.SetLogLevel(logLevelFlagValue)
	if err != nil {
		return nil, err
	}

	if flagsConfig.DisableAnsiColor {
		err = logger.RemoveLogObserver(os.Stdout)
		if err != nil {
			return nil, err
		}

		err = logger.AddLogObserver(os.Stdout, &logger.PlainFormatter{})
		if err != nil {
			return nil, err
		}
	}
	log.Trace("logger updated", "level", logLevelFlagValue, "disable ANSI color", flagsConfig.DisableAnsiColor)

	return fileLogging, nil
}

func checkHardwareRequirements(cfg config.HardwareRequirementsConfig) error {
	cpuFlags, err := parseFeatures(cfg.CPUFlags)
	if err != nil {
		return err
	}

	if !cpuid.CPU.Supports(cpuFlags...) {
		return fmt.Errorf("CPU Flags: Streaming SIMD Extensions 4 required")
	}

	return nil
}

func parseFeatures(features []string) ([]cpuid.FeatureID, error) {
	flags := make([]cpuid.FeatureID, 0)

	for _, cpuFlag := range features {
		featureID := cpuid.ParseFeature(cpuFlag)
		if featureID == cpuid.UNKNOWN {
			return nil, fmt.Errorf("CPU Flags: cpu flag %s not found", cpuFlag)
		}
	}

	return flags, nil
}
