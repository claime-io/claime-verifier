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
	SupportedChains = []string{"mainnet", "rinkeby"}
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

func (i GuildInteractor) RegisterContract(ctx context.Context, channelID, guildID string, permission int64, in NFTInfo) error {
	if err := in.validate(); err != nil {
		return i.error(channelID, err)
	}
	if !common.IsHexAddress(in.ContractAddress) {
		return i.error(channelID, errors.New("Contract address should be hex string"))
	}
	if !HasPermissionAdministrator(permission) {
		return i.error(channelID, errors.New("Only administrator can configure contracts"))
	}

	if err := i.rep.RegisterContract(ctx, in); err != nil {
		return i.error(channelID, err)
	}
	return i.notify(channelID, in)
}

func (i GuildInteractor) ListNFTs(ctx context.Context, channelID, guildID string) error {
	nfts, err := i.rep.ListContracts(ctx, guildID)
	if err != nil {
		return i.error(channelID, err)
	}
	return i.notifyNFTs(channelID, nfts)
}

func (i GuildInteractor) DeleteNFT(ctx context.Context, channelID, guildID string, contractAddress common.Address) error {
	nft, err := i.rep.GetContract(ctx, guildID, contractAddress)
	if err != nil {
		return i.error(channelID, err)
	}
	if nft == (NFTInfo{}) {
		return i.error(channelID, errors.New(fmt.Sprintf("Not registered. contract address: %s", contractAddress.Hex())))
	}
	err = i.rep.DeleteContract(ctx, guildID, contractAddress)
	if err != nil {
		return i.error(channelID, err)
	}
	return i.notifyDelete(channelID, contractAddress)
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

func (i GuildInteractor) error(channelID string, cause error) error {
	_, err := i.dg.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title:       "Something went wrong",
			Description: cause.Error(),
			Color:       int(0xFF0000),
		},
	})
	return err
}

func (i GuildInteractor) notify(channelID string, in NFTInfo) error {
	_, err := i.dg.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title:       "Set contract address Succeeded!",
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
	return err
}

func (i GuildInteractor) notifyNFTs(channelID string, nfts []NFTInfo) error {
	var err error
	if len(nfts) == 0 {
		_, err := i.dg.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title:       "Registered NFTs",
				Description: "No NFTs registered in this guild.",
				Color:       int(0x0000FF),
			},
		})
		return err
	}
	for _, nft := range nfts {
		_, err = i.dg.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title: "Registered NFTs",
				Color: int(0x0000FF),
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "ContractAddress",
						Value: nft.ContractAddress,
					},
					{
						Name:  "Network",
						Value: nft.Network,
					},
					{
						Value: nft.RoleID,
						Name:  "RoleID",
					},
				},
			},
		})
	}
	return err
}

func (i GuildInteractor) notifyDelete(channelID string, contractAddress common.Address) error {
	_, err := i.dg.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
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
	return err
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
