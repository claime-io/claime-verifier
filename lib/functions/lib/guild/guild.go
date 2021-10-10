package guild

import (
	"claime-verifier/lib/functions/lib/common/log"
	"claime-verifier/lib/functions/lib/infrastructure/ssm"
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	validate "github.com/go-playground/validator/v10"

	"github.com/bwmarrin/discordgo"
)

var (
	// SupportedChains supported chains
	SupportedChains = []string{"mainnet", "rinkeby", "polygon"}
)

const (
	discordAuthURL = "https://claime-webfront-git-feature-discord-squard.vercel.app/claim/discord"
)

type (
	KeyResolver interface {
		DiscordPublicKey() (val string, err error)
		DiscordBotToken() (val string, err error)
		ClaimePrivateKey() (val ed25519.PrivateKey, err error)
	}
	GuildInteractor struct {
		discordPublicKey string
		discordBotToken  string
		claimePrivateKey ed25519.PrivateKey
		rep              Repository
		dg               *discordgo.Session
	}
	NFTInfo struct {
		RoleID          string `json:"roleId" validate:"required"`
		ContractAddress string `json:"contract_address" validate:"required"`
		Network         string `json:"network" validate:"required"`
		GuildID         string `json:"guildId" validate:"required"`
	}
	Repository interface {
		RegisterContract(ctx context.Context, in NFTInfo) error
		ListContracts(ctx context.Context, guildID string) ([]NFTInfo, error)
		DeleteContract(ctx context.Context, guildID string, contractAddress common.Address) error
		GetContract(ctx context.Context, guildID string, contractAddress common.Address) (NFTInfo, error)
	}
)

func (in NFTInfo) validate() error {
	err := validate.New().Struct(in)
	if err != nil {
		log.Error("", err)
		return err
	}
	for _, c := range SupportedChains {
		if c == in.Network {
			return nil
		}
	}
	return errors.New(in.Network + " is not currently supported")
}

func New(r KeyResolver, rep Repository) (GuildInteractor, error) {
	pub, err := r.DiscordPublicKey()
	if err != nil {
		return GuildInteractor{}, err
	}
	t, err := r.DiscordBotToken()
	if err != nil {
		return GuildInteractor{}, err
	}
	pri, err := r.ClaimePrivateKey()
	if err != nil {
		return GuildInteractor{}, err
	}
	sess, err := discordgo.New("Bot " + t)
	if err != nil {
		return GuildInteractor{}, err
	}
	return GuildInteractor{
		discordPublicKey: pub,
		discordBotToken:  t,
		claimePrivateKey: pri,
		rep:              rep,
		dg:               sess,
	}, nil
}

func (i GuildInteractor) RegisterContract(ctx context.Context, interaction discordgo.Interaction, in NFTInfo) error {
	if !HasPermissionAdministrator(interaction.Member.Permissions) {
		return i.error(interaction, errors.New("Only administrator can configure contracts"))
	}
	if err := in.validate(); err != nil {
		return i.error(interaction, err)
	}
	if !common.IsHexAddress(in.ContractAddress) {
		return i.error(interaction, errors.New("Contract address should be hex string"))
	}
	if err := i.existsRole(interaction.GuildID, in.RoleID); err != nil {
		return i.error(interaction, err)
	}
	if err := i.rep.RegisterContract(ctx, in); err != nil {
		return i.error(interaction, err)
	}
	return i.notify(interaction, in)
}

func (i GuildInteractor) ListNFTs(ctx context.Context, interaction discordgo.Interaction) error {
	nfts, err := i.rep.ListContracts(ctx, interaction.GuildID)
	if err != nil {
		return i.error(interaction, err)
	}
	return i.notifyNFTs(interaction, nfts)
}

func (i GuildInteractor) DeleteNFT(ctx context.Context, interaction discordgo.Interaction, contractAddress common.Address) error {
	if !HasPermissionAdministrator(interaction.Member.Permissions) {
		return i.error(interaction, errors.New("Only administrator can configure contracts"))
	}
	nft, err := i.rep.GetContract(ctx, interaction.GuildID, contractAddress)
	if err != nil {
		return i.error(interaction, err)
	}
	if nft == (NFTInfo{}) {
		return i.error(interaction, errors.New(fmt.Sprintf("Not registered. contract address: %s", contractAddress.Hex())))
	}
	err = i.rep.DeleteContract(ctx, interaction.GuildID, contractAddress)
	if err != nil {
		return i.error(interaction, err)
	}
	return i.notifyDelete(interaction, contractAddress)
}

func (i GuildInteractor) GrantRole(userID string, in NFTInfo) error {
	err := i.dg.GuildMemberRoleAdd(in.GuildID, userID, in.RoleID)
	if err != nil {
		log.Error("role add failed", err)
	}
	return err
}

func (i GuildInteractor) existsRole(guildID, roleID string) error {
	st, err := i.dg.GuildRoles(guildID)
	if err != nil {
		return err
	}
	for _, s := range st {
		if s.ID == roleID {
			return nil
		}
	}
	return errors.New("Role with ID " + roleID + " does not exist.")
}

func (i GuildInteractor) error(interaction discordgo.Interaction, cause error) error {
	err := i.dg.InteractionRespond(&interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Something went wrong",
					Description: cause.Error(),
					Color:       int(0xFF0000),
				},
			},
		},
	})
	if err != nil {
		log.Error("interaction respond failed", err)
	}
	return err
}

func (i GuildInteractor) respond(interaction discordgo.Interaction, embeds []*discordgo.MessageEmbed) error {
	err := i.dg.InteractionRespond(&interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embeds,
		},
	})
	if err != nil {
		log.Error("notify failed", err)
	}
	return err
}

func (i GuildInteractor) notify(interaction discordgo.Interaction, in NFTInfo) error {
	return i.respond(interaction, []*discordgo.MessageEmbed{
		{
			Title:       "Contract Address Registration Succeeded!",
			Description: "Configure contract address succeeded with following properties:",
			Color:       int(0x0000FF),
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "ContractAddress",
					Value: in.ContractAddress,
				},
				{
					Name:  "Network",
					Value: in.Network,
				},
				{
					Value: in.RoleID,
					Name:  "RoleID",
				},
			},
		},
	})
}

func (i GuildInteractor) notifyNFTs(interaction discordgo.Interaction, nfts []NFTInfo) error {
	if len(nfts) == 0 {
		return i.respond(interaction, []*discordgo.MessageEmbed{
			{
				Title:       "Registered NFTs",
				Description: "No NFTs registered in this guild.",
				Color:       int(0x0000FF),
			},
		})
	}
	var fields []*discordgo.MessageEmbedField
	for _, nft := range nfts {
		fields = append(fields,
			&discordgo.MessageEmbedField{
				Name:  "ContractAddress",
				Value: nft.ContractAddress,
			},
			&discordgo.MessageEmbedField{
				Name:   "Network",
				Value:  nft.Network,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Value:  nft.RoleID,
				Name:   "RoleID",
				Inline: true,
			},
		)
	}
	return i.respond(interaction, []*discordgo.MessageEmbed{
		{
			Title:  "Registered NFTs",
			Color:  int(0x0000FF),
			Fields: fields,
		},
	})
}

func (i GuildInteractor) notifyDelete(interaction discordgo.Interaction, contractAddress common.Address) error {
	return i.respond(interaction, []*discordgo.MessageEmbed{
		{
			Title: "Delete contract address Succeeded!",
			Color: int(0x808080),
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "ContractAddress",
					Value: contractAddress.Hex(),
				},
			},
		},
	})
}

// GuildMemberAdd This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func GuildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	sendSignature(s, m.User.ID, m.GuildID, "Please complete sign to prove you have a NFT:")
}

func (i GuildInteractor) ResendVerifyMessage(userID, guildID string) error {
	return sendSignature(i.dg, userID, guildID, "URL has been expired. Please sign with this url again:")
}

func sendSignature(s *discordgo.Session, userID, guildID, messagePrefix string) error {
	channel, err := s.UserChannelCreate(userID)
	if err != nil {
		log.Error("error creating channel:", err)
		return err
	}
	in := SignatureInput{
		UserID:   userID,
		GuildID:  guildID,
		Validity: time.Now().Add(time.Minute * 10),
	}
	cli := ssm.New(context.Background())
	pk, err := cli.ClaimePrivateKey()
	if err != nil {
		log.Error("get claime public key failed", err)
		return err
	}
	sig := Sign(in, pk)
	_, err = s.ChannelMessageSend(channel.ID, messagePrefix+" "+url(in, sig))
	if err != nil {
		log.Error("error sending DM message:", err)
	}
	return err
}

func url(in SignatureInput, sig string) string {
	return fmt.Sprintf("%s?userId=%s&guildId=%s&validity=%d&signature=%s", discordAuthURL, in.UserID, in.GuildID, in.Validity.UnixNano(), sig)
}
