package transaction

import (
	"claime-verifier/lib/functions/lib/contracts"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
)

const (
	messagePrefix = "\x19Ethereum Signed Message:\n"
)

type (
	// EOAInput eoa input
	EOAInput struct {
		Signature string `json:"signature"`
		Message   string `json:"message"`
		RawTx     string `json:"rawTx"`
	}
)

// Recover recover EOA And Claim information from input
func Recover(in EOAInput) (common.Address, contracts.IClaimRegistryClaim, error) {
	if in.RawTx != "" {
		return recoverFromRawTx(in.RawTx, in.Signature)
	}
	return recoverFromMessage(in.Message, in.Signature)
}

func recoverFromRawTx(raw, sig string) (common.Address, contracts.IClaimRegistryClaim, error) {
	return recoverFrom(raw, sig, recoverAddressFromTx, recoverClaimFromTx)
}

func recoverFrom(from, sig string, recoverAddress func(from, sig string) (common.Address, error), recoverClaim func(from string) (contracts.IClaimRegistryClaim, error)) (common.Address, contracts.IClaimRegistryClaim, error) {
	recovered, err := recoverAddress(from, sig)
	if err != nil {
		return common.Address{}, contracts.IClaimRegistryClaim{}, err
	}
	claim, err := recoverClaim(from)
	if err != nil {
		return common.Address{}, contracts.IClaimRegistryClaim{}, err
	}
	return recovered, claim, err
}

func recoverFromMessage(message, sig string) (common.Address, contracts.IClaimRegistryClaim, error) {
	return recoverFrom(message, sig, recoverAddressFromMessage, recoverClaimFromMessage)
}

func recoverAddressFromTx(rawTx string, signature string) (common.Address, error) {
	txBytes, err := hexutil.Decode(rawTx)
	if err != nil {
		return common.Address{}, err
	}
	return recover(crypto.Keccak256(txBytes), signature)
}

func recoverClaimFromTx(rawTx string) (contracts.IClaimRegistryClaim, error) {
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

func recoverAddressFromMessage(message string, signature string) (common.Address, error) {
	return recover(signHash([]byte(message)), signature)
}

func recoverClaimFromMessage(message string) (contracts.IClaimRegistryClaim, error) {
	var val map[string]string
	err := json.Unmarshal([]byte(message), &val)
	return contracts.IClaimRegistryClaim{
		PropertyType: val["propertyType"],
		PropertyId:   val["propertyId"],
		Evidence:     val["evidence"],
		Method:       val["method"],
	}, err
}

func recover(hash []byte, signature string) (common.Address, error) {
	sigArr, err := hexutil.Decode(signature)
	if err != nil {
		return common.Address{}, err
	}
	sigArr[64] -= 27
	rpk, err := crypto.Ecrecover(hash, sigArr)
	if err != nil {
		return common.Address{}, err
	}
	pubKey, err := crypto.UnmarshalPubkey(rpk)
	if err != nil {
		return common.Address{}, err
	}
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	return recoveredAddr, err
}

func signHash(data []byte) []byte {
	msg := fmt.Sprintf(messagePrefix+"%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}
