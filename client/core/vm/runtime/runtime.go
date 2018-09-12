package runtime

import (
	"math"
	"math/big"
	"time"

	"github.com/kowala-tech/kcoin/client/common"
	"github.com/kowala-tech/kcoin/client/core/state"
	"github.com/kowala-tech/kcoin/client/core/vm"
	"github.com/kowala-tech/kcoin/client/crypto"
	"github.com/kowala-tech/kcoin/client/kcoindb"
	"github.com/kowala-tech/kcoin/client/params"
)

// Config is a basic type specifying certain configuration flags for running
// the VM.
type Config struct {
	ChainConfig      *params.ChainConfig
	Origin           common.Address
	Coinbase         common.Address
	BlockNumber      *big.Int
	Time             *big.Int
	ComputeLimit     uint64
	ComputeUnitPrice *big.Int
	Value            *big.Int
	Debug            bool
	VMConfig         vm.Config

	State     *state.StateDB
	GetHashFn func(n uint64) common.Hash
}

// sets defaults on the config
func setDefaults(cfg *Config) {
	if cfg.ChainConfig == nil {
		cfg.ChainConfig = &params.ChainConfig{
			ChainID: big.NewInt(1),
		}
	}

	if cfg.Time == nil {
		cfg.Time = big.NewInt(time.Now().Unix())
	}
	if cfg.ComputeLimit == 0 {
		cfg.ComputeLimit = math.MaxUint64
	}
	if cfg.ComputeUnitPrice == nil {
		cfg.ComputeUnitPrice = new(big.Int)
	}
	if cfg.Value == nil {
		cfg.Value = new(big.Int)
	}
	if cfg.BlockNumber == nil {
		cfg.BlockNumber = new(big.Int)
	}
	if cfg.GetHashFn == nil {
		cfg.GetHashFn = func(n uint64) common.Hash {
			return common.BytesToHash(crypto.Keccak256([]byte(new(big.Int).SetUint64(n).String())))
		}
	}
}

// Execute executes the code using the input as call data during the execution.
// It returns the EVM's return value, the new state and an error if it failed.
//
// Executes sets up a in memory, temporarily, environment for the execution of
// the given code. It makes sure that it's restored to it's original state afterwards.
func Execute(code, input []byte, cfg *Config) ([]byte, *state.StateDB, error) {
	if cfg == nil {
		cfg = new(Config)
	}
	setDefaults(cfg)

	if cfg.State == nil {
		cfg.State, _ = state.New(common.Hash{}, state.NewDatabase(kcoindb.NewMemDatabase()))
	}
	var (
		address = common.BytesToAddress([]byte("contract"))
		vmenv   = NewEnv(cfg)
		sender  = vm.AccountRef(cfg.Origin)
	)
	cfg.State.CreateAccount(address)
	// set the receiver's (the executing contract) code for execution.
	cfg.State.SetCode(address, code)
	// Call the code with the given configuration.
	ret, _, err := vmenv.Call(
		sender,
		common.BytesToAddress([]byte("contract")),
		input,
		cfg.ComputeLimit,
		cfg.Value,
	)

	return ret, cfg.State, err
}

// Create executes the code using the EVM create method
func Create(input []byte, cfg *Config) ([]byte, common.Address, uint64, error) {
	if cfg == nil {
		cfg = new(Config)
	}
	setDefaults(cfg)

	if cfg.State == nil {
		cfg.State, _ = state.New(common.Hash{}, state.NewDatabase(kcoindb.NewMemDatabase()))
	}
	var (
		vmenv  = NewEnv(cfg)
		sender = vm.AccountRef(cfg.Origin)
	)

	// Call the code with the given configuration.
	code, address, leftOverResource, err := vmenv.Create(
		sender,
		input,
		cfg.ComputeLimit,
		cfg.Value,
	)
	return code, address, leftOverResource, err
}

// Call executes the code given by the contract's address. It will return the
// EVM's return value or an error if it failed.
//
// Call, unlike Execute, requires a config and also requires the State field to
// be set.
func Call(address common.Address, input []byte, cfg *Config) ([]byte, uint64, error) {
	setDefaults(cfg)

	vmenv := NewEnv(cfg)

	sender := cfg.State.GetOrNewStateObject(cfg.Origin)
	// Call the code with the given configuration.
	ret, leftOverResource, err := vmenv.Call(
		sender,
		address,
		input,
		cfg.ComputeLimit,
		cfg.Value,
	)

	return ret, leftOverResource, err
}
