package discord

import (
	"strings"
	
	"github.com/ayntgl/astatine"
)

var missingMembers = make(map[string]bool)

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

func GetMember(session *astatine.Session, gID string, uID string) *astatine.Member {
	member,_ := session.State.Member(gID,uID)
	if(member == nil){
		if(missingMembers[gID+":"+uID] == true){
			return nil
		}
		member,_ = session.GuildMember(gID,uID)
		if(member == nil){
			missingMembers[gID+":"+uID] = true
			return nil
		}
		session.State.MemberAdd(member)
		return member
	}
	if(missingMembers[gID+":"+uID] == true){
		delete(missingMembers,gID+":"+uID)
	}
	return member
}

func GetChannel(session *astatine.Session, cID string) *astatine.Channel {
	channel,_ := session.State.Channel(cID)
	if(channel == nil){
		channel,_ = session.Channel(cID)
		if(channel == nil){
			return nil
		}
		session.State.ChannelAdd(channel)
		return channel
	}
	return channel
}

func GetGuild(session *astatine.Session, gID string) *astatine.Guild {
	guild,_ := session.State.Guild(gID)
	if(guild == nil){
		guild,_ = session.Guild(gID)
		if(guild == nil){
			return nil
		}
		session.State.GuildAdd(guild)
		return guild
	}
	return guild
}
