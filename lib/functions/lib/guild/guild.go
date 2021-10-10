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
	ContractInfo struct {
		RoleID          string `json:"roleId" validate:"required"`
		ContractAddress string `json:"contract_address" validate:"required"`
		Network         string `json:"network" validate:"required"`
		GuildID         string `json:"guildId" validate:"required"`
	}
	Repository interface {
		RegisterContract(in ContractInfo) error
		ListContracts(guildID string) ([]ContractInfo, error)
	}
)

func (in ContractInfo) validate() error {
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

func (i GuildInteractor) RegisterContract(channelID, guildID string, permission int64, in NFTInfo) error {
	if err := in.validate(); err != nil {
		return i.error(channelID, err)
	}
	if !common.IsHexAddress(in.ContractAddress) {
		return i.error(channelID, errors.New("Contract address should be hex string"))
	}
	if !HasPermissionAdministrator(permission) {
		return i.error(channelID, errors.New("Only administrator can configure contracts"))
	}

	if err := i.rep.RegisterContract(in); err != nil {
		return i.error(channelID, err)
	}
	return i.notify(channelID, in)
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

func (i GuildInteractor) notify(channelID string, in ContractInfo) error {
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

// GuildMemberAdd This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func GuildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	fmt.Println(m.GuildID)
	fmt.Println(m.User.ID)
	channel, err := s.UserChannelCreate(m.User.ID)
	if err != nil {
		// If an error occurred, we failed to create the channel.
		//
		// Some common causes are:
		// 1. We don't share a server with the user (not possible here).
		// 2. We opened enough DM channels quickly enough for Discord to
		//    label us as abusing the endpoint, blocking us from opening
		//    new ones.
		log.Error("error creating channel:", err)
		return
	}
	in := SignatureInput{
		UserID:   m.User.ID,
		GuildID:  m.GuildID,
		Validity: time.Now().Add(time.Minute * 10),
	}
	cli := ssm.New(context.Background())
	pk, err := cli.ClaimePrivateKey()
	sig := Sign(in, pk)

	_, err = s.ChannelMessageSend(channel.ID, "Please complete sign to prove you have a NFT: "+url(in, sig))
	if err != nil {
		// If an error occurred, we failed to send the message.
		//
		// It may occur either when we do not share a server with the
		// user (highly unlikely as we just received a message) or
		// the user disabled DM in their settings (more likely).
		log.Error("error sending DM message:", err)
	}

	fmt.Println(m)
	s.UserChannels()
}

func url(in SignatureInput, sig string) string {
	return fmt.Sprintf("%s?userId=%s&guildId=%s&validity=%d&signature=%s", discordAuthURL, in.UserID, in.GuildID, in.Validity.UnixNano(), sig)
}
