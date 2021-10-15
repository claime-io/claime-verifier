package registry

import "os"

const (
	claimeRegistryAddressRinkeby = "0x886A9591D624b5360d22B32C5866387340803593"
	claimeRegistryAddressMainnet = "0xb52E96533528eD66AbFC3a9680A998a4eBe0E35a"
)

func registryAddress() string {
	env := os.Getenv("EnvironmentId")
	if env == "prod" {
		return claimeRegistryAddressMainnet
	}
	return claimeRegistryAddressRinkeby
}
