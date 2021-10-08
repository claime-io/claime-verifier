package guild

import (
	"claime-verifier/lib/functions/lib/common/log"
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"time"

	validate "github.com/go-playground/validator/v10"

	"github.com/bwmarrin/discordgo"
)

var (
	privateKey ed25519.PrivateKey
	// SupportedChains supported chains
	SupportedChains = []string{"mainnet", "rinkeby"}
)

const (
	discordAuthURL = "https://claime-webfront-k6p1srx99-squard.vercel.app/claim/discord"
)

type (
	KeyResolver interface {
		DiscordPublicKey(ctx context.Context) (val string, err error)
		DiscordBotToken(ctx context.Context) (val string, err error)
		ClaimePrivateKey(ctx context.Context) (val ed25519.PrivateKey, err error)
	}
	GuildInteractor struct {
		discordPublicKey string
		discordBotToken  string
		claimePrivateKey ed25519.PrivateKey
		rep              Repository
	}
	ContractInfo struct {
		RoleID          string `json:"roleId" validate:"required"`
		ContractAddress string `json:"contract_address" validate:"required"`
		Network         string `json:"network" validate:"required"`
		GuildID         string `json:"guildId" validate:"required"`
	}
	Repository interface {
		RegisterContract(ctx context.Context, in ContractInfo) error
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

func New(ctx context.Context, r KeyResolver, rep Repository) (GuildInteractor, error) {
	pub, err := r.DiscordPublicKey(ctx)
	if err != nil {
		return GuildInteractor{}, err
	}
	t, err := r.DiscordBotToken(ctx)
	if err != nil {
		return GuildInteractor{}, err
	}
	pri, err := r.ClaimePrivateKey(ctx)
	if err != nil {
		return GuildInteractor{}, err
	}
	return GuildInteractor{
		discordPublicKey: pub,
		discordBotToken:  t,
		claimePrivateKey: pri,
		rep:              rep,
	}, nil
}

func (i GuildInteractor) RegisterContract(ctx context.Context, in ContractInfo) error {
	if err := in.validate(); err != nil {
		return err
	}
	return i.rep.RegisterContract(ctx, in)
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
		UserID:    m.User.ID,
		GuildID:   m.GuildID,
		Validity:  time.Now().Add(time.Minute * 10),
		Timestamp: time.Now(),
	}
	sig := Sign(in, privateKey)

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
	return discordAuthURL + "?userId=" + in.UserID + "&guildId=" + in.GuildID + "&validity=" + fmt.Sprint(in.Validity.UnixNano()) + "&timestamp=" + fmt.Sprint(in.Timestamp.UnixNano()) + "&signature=" + sig
}
