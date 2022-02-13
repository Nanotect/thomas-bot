package members

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"text/template"
	"time"

	"github.com/itfactory-tm/thomas-bot/pkg/util/slash"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
	"github.com/itfactory-tm/thomas-bot/pkg/db"
)

var welcomeEmoji = []string{"💗", "✨", "🤩", "👾", "💫", "🎀", "🌟"}

// MemberCommands contains the tm!role command and welcome messages
type MemberCommands struct {
	db db.Database
}

// NewMemberCommand gives a new MemberCommands
func NewMemberCommand(conn db.Database) *MemberCommands {
	return &MemberCommands{
		db: conn,
	}
}

// Register registers the handlers
func (m *MemberCommands) Register(registry command.Registry, server command.Server) {
	registry.RegisterInteractionCreate("role", m.roleSlashCommand)
	registry.RegisterInteractionCreate("rolereq", m.handleRoleRequest)
	registry.RegisterInteractionCreate("roleresponse", m.handleRolePermissionResponse)

	registry.RegisterGuildMemberAddHandler(m.onGuildMemberAdd)
}

// InstallSlashCommands registers the slash commands
func (m *MemberCommands) InstallSlashCommands(session *discordgo.Session) error {
	return slash.InstallSlashCommand(session, "", discordgo.ApplicationCommand{
		Name:        "role",
		Description: "Request a new role om this server",
		Options:     []*discordgo.ApplicationCommandOption{},
	})
}

func (m *MemberCommands) onGuildMemberAdd(s *discordgo.Session, g *discordgo.GuildMemberAdd) {
	conf, err := m.db.ConfigForGuild(g.GuildID)
	if err != nil {
		log.Println(err)
		return
	}

	if conf.WelcomeChannelID == "" {
		// no welcome channel set!
		return
	}

	t, err := template.New("welcome").Parse(conf.WelcomeText)
	if err != nil {
		log.Println(err)
		return
	}
	var welcomeText bytes.Buffer
	err = t.Execute(&welcomeText, g)
	if err != nil {
		log.Println(err)
		return
	}

	welcome, _ := s.ChannelMessageSend(conf.WelcomeChannelID, welcomeText.String())
	go func() {
		// waving back is not essential and should not delay other actions
		// plus the students want to race against the bot in waving at new users so let's give a head start
		time.Sleep(5 * time.Minute)
		err = s.MessageReactionAdd(conf.WelcomeChannelID, welcome.ID, "👋")
		if err != nil {
			log.Println(err)
		}
		// seed the random number generator
		rand.Seed(time.Now().UnixNano())

		// pick random emoji from welcomeEmoji
		err = s.MessageReactionAdd(conf.WelcomeChannelID, welcome.ID, welcomeEmoji[rand.Intn(len(welcomeEmoji))])
		if err != nil {
			log.Println(err)
		}
	}()

	if conf.RoleManagement.DefaultRole != "" {
		err := s.GuildMemberRoleAdd(g.GuildID, g.Member.User.ID, conf.RoleManagement.DefaultRole)
		if err != nil {
			log.Printf("Cannot set role for user %s: %q\n", g.Member.User.ID, err)
		}
	}

	if len(conf.WelcomeDM) > 0 {
		c, err := s.UserChannelCreate(g.Member.User.ID)
		if err != nil {
			log.Printf("Cannot DM user %s\n", g.Member.User.ID)
			return
		}

		s.ChannelMessageSend(c.ID, fmt.Sprintf("Hi %s", g.User.Username))
		time.Sleep(time.Second)

		for _, msg := range conf.WelcomeDM {
			s.ChannelMessageSend(c.ID, msg)
			time.Sleep(time.Second)
		}

		if len(conf.RoleManagement.RoleSets) > 0 {
			m.SendRoleDM(s, g.GuildID, g.Member.User.ID)
		}
	}

}

// Info return the commands in this package
func (m *MemberCommands) Info() []command.Command {
	return []command.Command{
		command.Command{
			Name:        "role",
			Category:    command.CategoryAlgemeen,
			Description: "Modify your ITFactory Discord role",
			Hidden:      false,
		},
	}
}
