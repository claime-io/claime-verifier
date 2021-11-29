package evmnetwork

type (
	network string
)

const (
	Rinkeby = network("rinkeby")
	Mumbai  = network("mumbai")
	Polygon = network("polygon")
	Mainnet = network("mainnet")
)

func (n network) Equals(str string) bool {
	return string(n) == str
}
