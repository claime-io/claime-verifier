package registry

import (
	"claime-verifier/lib/functions/lib/infrastructure/evmnetwork"
	"os"

	"github.com/pkg/errors"
)

const (
	claimeRegistryAddressRinkeby = "0xD721AF405fb939fFeBF7B44b294D9D02A232b359"
	claimeRegistryAddressMumbai  = "0x9b67374857503dA14209844598B0e65fA022Ac1B"
	claimeRegistryAddressMainnet = "0xb52E96533528eD66AbFC3a9680A998a4eBe0E35a"
)

func registryAddress(network string) (string, error) {
	env := os.Getenv("EnvironmentId")
	if evmnetwork.Mainnet.Equals(network) {
		return claimeRegistryAddressMainnet, nil
	}
	if env != "prod" {
		if evmnetwork.Rinkeby.Equals(network) {
			return claimeRegistryAddressRinkeby, nil
		}
		if evmnetwork.Mumbai.Equals(network) {
			return claimeRegistryAddressMumbai, nil
		}
	}
	return "", errors.Errorf("Unsupported network: %s", network)
}
