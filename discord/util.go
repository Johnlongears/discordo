package discord

import (
	"strings"

	"github.com/ayntgl/astatine"
	"github.com/ayntgl/discordo/ui"
)

func ChannelToString(c *astatine.Channel) string {
	var repr string
	if c.Name != "" {
		repr = "#" + c.Name
	} else if len(c.Recipients) == 1 {
		rp := c.Recipients[0]
		repr = rp.Username + "#" + rp.Discriminator
	} else {
		rps := make([]string, len(c.Recipients))
		for i, r := range c.Recipients {
			rps[i] = r.Username + "#" + r.Discriminator
		}

		repr = strings.Join(rps, ", ")
	}

	return repr
}

func FindMessageByID(ms []*astatine.Message, mID string) (int, *astatine.Message) {
	for i, m := range ms {
		if m.ID == mID {
			return i, m
		}
	}

	return -1, nil
}

func HasPermission(s *astatine.State, cID string, p int64) bool {
	perm, err := s.UserChannelPermissions(s.User.ID, cID)
	if err != nil {
		return false
	}

	return perm&p == p
}

func GetMember(app *App, gID string, uID string) *astatine.Member {
	var member *astatine.Member
	member,_ = app.Session.State.Member(gID,uID)
	if(member == nil){
		member,_ = app.Session.Member(gID,uID)
		if(member == nil){
			return nil
		}
		app.Session.State.MemberAdd(member)
		return member
	}
	return member
}

func GetChannel(app *App, cID string) *astatine.Channel {
	var channel *astatine.Channel
	channel,_ = app.Session.State.Channel(cID)
	if(channel == nil){
		channel,_ = app.Session.Channel(cID)
		if(channel == nil){
			return nil
		}
		app.Session.State.ChannelAdd(channel)
		return channel
	}
	return channel
}

func GetGuild(app *App, gID string) *astatine.Guild {
	var guild *astatine.Guild
	guild,_ = app.Session.State.Guild(gID)
	if(guild == nil){
		guild,_ = app.Session.Guild(gID)
		if(guild == nil){
			return nil
		}
		app.Session.State.GuildAdd(guild)
		return guild
	}
	return guild
}
