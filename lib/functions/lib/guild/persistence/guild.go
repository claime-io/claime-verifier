package guildrep

import (
	"claime-verifier/lib/functions/lib/common/log"
	"claime-verifier/lib/functions/lib/guild"
	"context"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ethereum/go-ethereum/common"
	"github.com/guregu/dynamo"
)

const (
	pkPrefix = "Discord:Guild:"
	skPrefix = "Discord:Contract:"
)

type (
	Repository struct {
		ddb *dynamo.DB
		ctx context.Context
	}

	Contract struct {
		PK      string `dynamo:"PK"`
		SK      string `dynamo:"SK"`
		Network string `dynamo:"Network"`
		RoleID  string `json:"RoleId"`
		//ContractAddress string `json:"contract_address"`
		//GuildID         string `dynamo:"PK"`
	}
)

func New(ctx context.Context) Repository {
	return Repository{
		ddb: dynamo.New(session.New()),
		ctx: ctx,
	}
}

func (r Repository) RegisterContract(in guild.NFTInfo) error {
	item := toContract(in)
	err := r.ddb.Table(table()).Put(&item).RunWithContext(r.ctx)
	if err != nil {
		log.Error("put item failed", err)
	}
	return err
}

func (r Repository) ListContracts(guildID string) ([]guild.NFTInfo, error) {
	res := []Contract{}
	err := r.ddb.Table(table()).Get("PK", toPK(guildID)).AllWithContext(r.ctx, &res)
	if err != nil {
		log.Error("query failed", err)
	}
	return fromDDB(res), err
}

func (r Repository) DeleteContract(guildID string, contractAddress common.Address) error {
	err := r.ddb.Table(table()).Delete("PK", toPK(guildID)).Range("SK", toSK(contractAddress.Hex())).RunWithContext(r.ctx)
	if err != nil {
		log.Error("delete failed", err)
	}
	return err
}

func (r Repository) GetContract(guildID string, contractAddress common.Address) (guild.NFTInfo, error) {
	res := []Contract{}
	err := r.ddb.Table(table()).Get("PK", toPK(guildID)).Range("SK", dynamo.Equal, toSK(contractAddress.Hex())).AllWithContext(r.ctx, &res)
	if err != nil {
		log.Error("query failed", err)
	}
	if len(res) == 0 {
		return guild.NFTInfo{}, err
	}
	return fromDDB(res)[0], nil
}

func fromDDB(vals []Contract) []guild.NFTInfo {
	res := []guild.NFTInfo{}
	for _, v := range vals {
		res = append(res, toNFTInfo(v))
	}
	return res
}

func toNFTInfo(in Contract) guild.NFTInfo {
	return guild.NFTInfo{
		GuildID:         fromPK(in.PK),
		ContractAddress: fromSK(in.SK),
		Network:         in.Network,
		RoleID:          in.RoleID,
	}
}

func table() string {
	return "claime-verifier-main-" + os.Getenv("EnvironmentId")
}

func toContract(in guild.NFTInfo) Contract {
	return Contract{
		PK:      toPK(in.GuildID),
		SK:      toSK(in.ContractAddress),
		Network: in.Network,
		RoleID:  in.RoleID,
	}
}

func toPK(guildID string) string {
	return pkPrefix + guildID
}

func toSK(contractaddress string) string {
	return skPrefix + contractaddress
}

func fromPK(pk string) (guildID string) {
	return from(pk, pkPrefix)
}

func fromSK(sk string) (contractaddress string) {
	return from(sk, skPrefix)
}

func from(key, prefix string) (contractval string) {
	return strings.Replace(key, prefix, "", 1)
}
