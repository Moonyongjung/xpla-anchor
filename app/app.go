package app

import (
	xtypes "github.com/Moonyongjung/xpla.go/types"
)

// Structure of the app.yaml file.
type AppType struct {
	Home     string     `yaml:"Home"`
	Config   ConfigType `yaml:"Config"`
	Contract Contract   `yaml:"Contract"`
}

type ConfigType struct {
	Anchor       Anchor       `yaml:"Anchor"`
	PublicChain  PublicChain  `yaml:"PublicChain"`
	PrivateChain PrivateChain `yaml:"PrivateChain"`
}

type Anchor struct {
	CollectBlockCount int `yaml:"CollectBlockCount"`
	RequestPeriod     int `yaml:"RequestPeriod"`
	DB                DB  `yaml:"DB"`
}

type DB struct {
	DBUserName string `yaml:"DBUserName"`
	DBPassword string `yaml:"DBPassword"`
	DBHost     string `yaml:"DBHost"`
	DBPort     string `yaml:"DBPort"`
	DBName     string `yaml:"DBName"`
}

type Contract struct {
	ContractFilePath string `yaml:"ContractPath"`
	CodeID           string `yaml:"CodeID"`
	Address          string `yaml:"Address"`
}

type PublicChain struct {
	ChainID       string `yaml:"ChainID"`
	LCD           string `yaml:"LCD"`
	GasAdj        string `yaml:"GasAdj"`
	GasLimit      string `yaml:"GasLimit"`
	BroadcastMode string `yaml:"BroadcastMode"`
}

type PrivateChain struct {
	ChainID string `yaml:"ChainID"`
	LCD     string `yaml:"LCD"`
}

// Generate default app.yaml.
func GenDefaultApp(home string, appFilePath string, configType ConfigType) error {
	appType := &AppType{
		Home:   home,
		Config: configType,
	}

	err := saveAppFile(appFilePath, appType)
	if err != nil {
		return err
	}

	return nil
}

// Modify the config parameters in the existing app.yaml file.
func (a *AppType) SetConfig(config ConfigType, appFilePath string) error {
	a.Config = config

	err := saveAppFile(appFilePath, a)
	if err != nil {
		return err
	}

	return nil
}

// Record the code ID of the anchor contract after executing store command.
func SetCodeId(appFilePath string, contractPath string, res *xtypes.TxRes) error {
	appType, err := readAppFile(appFilePath)
	if err != nil {
		return err
	}

	cwRes := res.Response
	codeId := cwRes.Logs[0].Events[1].Attributes[0].Value
	appType.Contract.CodeID = codeId
	appType.Contract.ContractFilePath = contractPath

	err = saveAppFile(appFilePath, appType)
	if err != nil {
		return err
	}

	return nil
}

// Record the contract address of the anchor contract after execute instantiate command.
func SetContractAddress(appFilePath string, res *xtypes.TxRes) error {
	appType, err := readAppFile(appFilePath)
	if err != nil {
		return err
	}

	cwRes := res.Response
	contractAddress := cwRes.Logs[0].Events[0].Attributes[0].Value
	appType.Contract.Address = contractAddress

	err = saveAppFile(appFilePath, appType)
	if err != nil {
		return err
	}

	return nil
}
