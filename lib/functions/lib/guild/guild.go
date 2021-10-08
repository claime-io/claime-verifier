package guild

import (
	"claime-verifier/lib/functions/lib/common/log"
	"context"
	"crypto/ed25519"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	privateKey ed25519.PrivateKey
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
	}
	RegisterContractInput struct {
		RoleID          string `json:"roleId"`
		ContractAddress string `json:"contract_address"`
		Network         string `json:"network"`
	}
)

func New(ctx context.Context, r KeyResolver) (GuildInteractor, error) {
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
	}, nil
}

func (i GuildInteractor) RegisterContract(ctx context.Context, in RegisterContractInput) error {
	return nil
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
