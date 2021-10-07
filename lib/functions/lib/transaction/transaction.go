package transaction

import (
	"claime-verifier/lib/functions/lib/contracts"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
)

const (
	messagePrefix = "\x19Ethereum Signed Message:\n"
)

func RecoverAddress(rawTx string, signature string) (string, error) {
	txBytes, err := hexutil.Decode(rawTx)
	if err != nil {
		return "", err
	}
	return recover(crypto.Keccak256(txBytes), signature)
}

func RecoverClaimFromTx(rawTx string) (contracts.IClaimRegistryClaim, error) {
	var result []interface{}
	arr, err := hexutil.Decode(rawTx)
	if err != nil {
		return contracts.IClaimRegistryClaim{}, err
	}
	err = rlp.DecodeBytes(arr[1:], &result)
	if err != nil {
		return contracts.IClaimRegistryClaim{}, err
	}
	data, ok := result[7].([]byte)
	if !ok {
		return contracts.IClaimRegistryClaim{}, errors.New(fmt.Sprintf("expected []byte but got %T", data))
	}
	val, err := abi.Arguments{
		abi.Argument{Type: abi.Type{T: abi.StringTy}},
		abi.Argument{Type: abi.Type{T: abi.StringTy}},
		abi.Argument{Type: abi.Type{T: abi.StringTy}},
		abi.Argument{Type: abi.Type{T: abi.StringTy}},
	}.UnpackValues(data[4:])
	if err != nil {
		return contracts.IClaimRegistryClaim{}, err
	}
	return contracts.IClaimRegistryClaim{
		PropertyType: val[0].(string),
		PropertyId:   val[1].(string),
		Evidence:     val[2].(string),
		Method:       val[3].(string),
	}, nil
}

func RecoverAddressFromMessage(message string, signature string) (string, error) {
	return recover(signHash([]byte(message)), signature)
}

func RecoverClaimFromMessage(message string) (contracts.IClaimRegistryClaim, error) {
	var val map[string]string
	err := json.Unmarshal([]byte(message), &val)
	return contracts.IClaimRegistryClaim{
		PropertyType: val["propertyType"],
		PropertyId:   val["propertyId"],
		Evidence:     val["evidence"],
		Method:       val["method"],
	}, err
}

func recover(hash []byte, signature string) (string, error) {
	sigArr, err := hexutil.Decode(signature)
	if err != nil {
		return "", err
	}
	sigArr[64] -= 27
	rpk, err := crypto.Ecrecover(hash, sigArr)
	if err != nil {
		return "", err
	}
	pubKey, err := crypto.UnmarshalPubkey(rpk)
	if err != nil {
		return "", err
	}
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	return recoveredAddr.Hex(), err
}

func signHash(data []byte) []byte {
	msg := fmt.Sprintf(messagePrefix+"%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}
