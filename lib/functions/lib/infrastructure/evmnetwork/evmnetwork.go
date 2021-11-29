package evmnetwork

type (
	Network string
)

const (
	Rinkeby = Network("rinkeby")
	Mumbai  = Network("mumbai")
	Polygon = Network("polygon")
	Mainnet = Network("mainnet")
)

func (n Network) ToString() string {
	return string(n)
}

func (n Network) Equals(str string) bool {
	return n.ToString() == str
}

func ValueOf(str string) Network {
	if Rinkeby.Equals(str) {
		return Rinkeby
	}
	if Mumbai.Equals(str) {
		return Mumbai
	}
	if Polygon.Equals(str) {
		return Polygon
	}
	if Mainnet.Equals(str) {
		return Mainnet
	}
	return ""
}
