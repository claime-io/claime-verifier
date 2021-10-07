package transaction

import (
	"claime-verifier/lib/functions/lib/common/log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func RecoverAddress(rawTx string, signature string) string {
	arr, err := hexutil.Decode(rawTx)
	if err != nil {
		log.Error("failed to decode rawTx", err)
		return ""
	}
	hash := crypto.Keccak256Hash(arr)
	sigArr, err := hexutil.Decode(signature)
	if err != nil {
		log.Error("failed to decode signature", err)
		return ""
	}
	sigArr[64] -= 27
	rpk, err := crypto.Ecrecover(hash.Bytes(), sigArr)
	if err != nil {
		log.Error("failed to recover address", err)
		return ""
	}
	pubKey, err := crypto.UnmarshalPubkey(rpk)
	if err != nil {
		log.Error("failed to unmarshal pubkey", err)
		return ""
	}
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	return recoveredAddr.Hex()
}
