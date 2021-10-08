package guildrep

import (
	"claime-verifier/lib/functions/lib/common/log"
	"claime-verifier/lib/functions/lib/guild"
	"context"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

const (
	pkPrefix = "Discord:Guild:"
	skPrefix = "Discord:Contract:"
)

type (
	Repository struct {
		ddb *dynamo.DB
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

func New() Repository {
	return Repository{
		ddb: dynamo.New(session.New()),
	}
}

func (r Repository) RegisterContract(ctx context.Context, in guild.RegisterContractInput) error {
	item := toContract(in)
	err := r.ddb.Table(table()).Put(&item).RunWithContext(ctx)
	if err != nil {
		log.Error("put item failed", err)
	}
	return err
}

func table() string {
	return "claime-verifier-main-" + os.Getenv("EnvironmentId")
}

func toContract(in guild.RegisterContractInput) Contract {
	return Contract{
		PK:      pkPrefix + in.GuildID,
		SK:      skPrefix + in.ContractAddress,
		Network: in.Network,
		RoleID:  in.RoleID,
	}
}
