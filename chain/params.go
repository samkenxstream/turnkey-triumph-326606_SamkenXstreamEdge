package chain

import (
	"errors"
	"sort"

	"github.com/0xPolygon/polygon-edge/forkmanager"
	"github.com/0xPolygon/polygon-edge/types"
)

var (
	// ErrBurnContractAddressMissing is the error when a contract address is not provided
	ErrBurnContractAddressMissing = errors.New("burn contract address missing")
)

// Params are all the set of params for the chain
type Params struct {
	Forks          *Forks                 `json:"forks"`
	ChainID        int64                  `json:"chainID"`
	Engine         map[string]interface{} `json:"engine"`
	BlockGasTarget uint64                 `json:"blockGasTarget"`

	// BaseFeeChangeDenom is the value to bound the amount the base fee can change between blocks
	BaseFeeChangeDenom uint64 `json:"baseFeeChangeDenom,omitempty"`

	// Access control configuration
	AccessListsOwner          *types.Address     `json:"accessListsOwner,omitempty"`
	ContractDeployerAllowList *AddressListConfig `json:"contractDeployerAllowList,omitempty"`
	ContractDeployerBlockList *AddressListConfig `json:"contractDeployerBlockList,omitempty"`
	TransactionsAllowList     *AddressListConfig `json:"transactionsAllowList,omitempty"`
	TransactionsBlockList     *AddressListConfig `json:"transactionsBlockList,omitempty"`
	BridgeAllowList           *AddressListConfig `json:"bridgeAllowList,omitempty"`
	BridgeBlockList           *AddressListConfig `json:"bridgeBlockList,omitempty"`

	// Governance contract where the token will be sent to and burn in london fork
	BurnContract map[uint64]types.Address `json:"burnContract"`
	// Destination address to initialize default burn contract with
	BurnContractDestinationAddress types.Address `json:"burnContractDestinationAddress,omitempty"`
}

type AddressListConfig struct {
	// AdminAddresses is the list of the initial admin addresses
	AdminAddresses []types.Address `json:"adminAddresses,omitempty"`

	// EnabledAddresses is the list of the initial enabled addresses
	EnabledAddresses []types.Address `json:"enabledAddresses,omitempty"`
}

// CalculateBurnContract calculates burn contract address for the given block number
func (p *Params) CalculateBurnContract(block uint64) (types.Address, error) {
	blocks := make([]uint64, 0, len(p.BurnContract))

	for startBlock := range p.BurnContract {
		blocks = append(blocks, startBlock)
	}

	if len(blocks) == 0 {
		return types.ZeroAddress, ErrBurnContractAddressMissing
	}

	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i] < blocks[j]
	})

	for i := 0; i < len(blocks)-1; i++ {
		if block >= blocks[i] && block < blocks[i+1] {
			return p.BurnContract[blocks[i]], nil
		}
	}

	return p.BurnContract[blocks[len(blocks)-1]], nil
}

func (p *Params) GetEngine() string {
	// We know there is already one
	for k := range p.Engine {
		return k
	}

	return ""
}

// predefined forks
const (
	Homestead           = "homestead"
	Byzantium           = "byzantium"
	Constantinople      = "constantinople"
	Petersburg          = "petersburg"
	Istanbul            = "istanbul"
	London              = "london"
	EIP150              = "EIP150"
	EIP158              = "EIP158"
	EIP155              = "EIP155"
	QuorumCalcAlignment = "quorumcalcalignment"
	LondonFix           = "londonfix"
	Governance          = "governance"
	DoubleSignSlashing  = "doubleSignSlashing"
)

// Forks is map which contains all forks and their starting blocks from genesis
type Forks map[string]Fork

// IsActive returns true if fork defined by name exists and defined for the block
func (f *Forks) IsActive(name string, block uint64) bool {
	ff, exists := (*f)[name]

	return exists && ff.Active(block)
}

// SetFork adds/updates fork defined by name
func (f *Forks) SetFork(name string, value Fork) {
	(*f)[name] = value
}

func (f *Forks) RemoveFork(name string) {
	delete(*f, name)
}

// At returns ForksInTime instance that shows which supported forks are enabled for the block
func (f *Forks) At(block uint64) ForksInTime {
	return ForksInTime{
		Homestead:           f.IsActive(Homestead, block),
		Byzantium:           f.IsActive(Byzantium, block),
		Constantinople:      f.IsActive(Constantinople, block),
		Petersburg:          f.IsActive(Petersburg, block),
		Istanbul:            f.IsActive(Istanbul, block),
		London:              f.IsActive(London, block),
		EIP150:              f.IsActive(EIP150, block),
		EIP158:              f.IsActive(EIP158, block),
		EIP155:              f.IsActive(EIP155, block),
		QuorumCalcAlignment: f.IsActive(QuorumCalcAlignment, block),
		LondonFix:           f.IsActive(LondonFix, block),
		Governance:          f.IsActive(Governance, block),
		DoubleSignSlashing:  f.IsActive(DoubleSignSlashing, block),
	}
}

type Fork struct {
	Block  uint64                  `json:"block"`
	Params *forkmanager.ForkParams `json:"params,omitempty"`
}

func NewFork(n uint64) Fork {
	return Fork{Block: n}
}

func (f Fork) Active(block uint64) bool {
	return block >= f.Block
}

// ForksInTime should contain all supported forks by current edge version
type ForksInTime struct {
	Homestead,
	Byzantium,
	Constantinople,
	Petersburg,
	Istanbul,
	London,
	EIP150,
	EIP158,
	EIP155,
	QuorumCalcAlignment,
	LondonFix,
	Governance,
	DoubleSignSlashing bool
}

// AllForksEnabled should contain all supported forks by current edge version
var AllForksEnabled = &Forks{
	Homestead:           NewFork(0),
	EIP150:              NewFork(0),
	EIP155:              NewFork(0),
	EIP158:              NewFork(0),
	Byzantium:           NewFork(0),
	Constantinople:      NewFork(0),
	Petersburg:          NewFork(0),
	Istanbul:            NewFork(0),
	London:              NewFork(0),
	QuorumCalcAlignment: NewFork(0),
	LondonFix:           NewFork(0),
	Governance:          NewFork(0),
	DoubleSignSlashing:  NewFork(0),
}
