package registry

import (
	"claime-verifier/lib/functions/lib/infrastructure/evmnetwork"
	"os"

	"github.com/pkg/errors"
)

const (
	claimeRegistryAddressRinkeby = "0xD721AF405fb939fFeBF7B44b294D9D02A232b359"
	claimeRegistryAddressMumbai  = "0x9b67374857503dA14209844598B0e65fA022Ac1B"
	claimeRegistryAddressPolygon = "0x7Cac4b4a233849b301b4b651666C3f8cCcb834e2"
	claimeRegistryAddressMainnet = "0x7cac4b4a233849b301b4b651666c3f8cccb834e2"
)

func registryAddress(network string) (string, error) {
	env := os.Getenv("EnvironmentId")
	if evmnetwork.Mainnet.Equals(network) {
		return claimeRegistryAddressMainnet, nil
	}
	if evmnetwork.Polygon.Equals(network) {
		return claimeRegistryAddressPolygon, nil
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
