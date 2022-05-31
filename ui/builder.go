package ui

import (
	"fmt"
	"strings"
	"time"
	"regexp"

	"github.com/ayntgl/astatine"
	"github.com/craftxbox/discordo/discord"
	"github.com/rivo/tview"
)

func buildEdit(app *App, e *astatine.MessageUpdate) []byte {
	var b strings.Builder
	// Define a new region and assign message ID as the region ID.
	// Learn more:
	// https://pkg.go.dev/github.com/rivo/tview#hdr-Regions_and_Highlights
	m := e
	b.WriteString("[\"")
	b.WriteString(m.ID)
	b.WriteString("\"]")
	// Build the message associated with crosspost, channel follow add, pin, or a reply.
	buildReferencedMessage(&b, e.BeforeUpdate, app.Session.State.User.ID, app)

	if app.Config.General.Timestamps {
		b.WriteString("[::d]")
		b.WriteString(m.Timestamp.Format(time.Stamp))
		b.WriteString("[::-]")
		b.WriteByte(' ')
	}
	b.WriteString(tview.Escape("[EDIT] "))
	if e.BeforeUpdate == nil {
		// Build the author of this message.
		if m.Member == nil {
			if len(m.GuildID) > 0 {
				member := discord.GetMember(app.Session, m.GuildID, m.Author.ID)
				buildAuthor(&b, m.Author, app.Session.State.User.ID, member, app)
			} else {
				c := app.SelectedChannel
				if c != nil && len(c.GuildID) > 0 {
					member := discord.GetMember(app.Session, c.GuildID, m.Author.ID)
					buildAuthor(&b, m.Author, app.Session.State.User.ID, member, app)
				} else {
					buildAuthor(&b, m.Author, app.Session.State.User.ID, nil, app) // dm channel, probably.
				}
			}
		} else {
			buildAuthor(&b, m.Author, app.Session.State.User.ID, m.Member, app)
		}
	}

	// Build the contents of the message.
	buildContent(&b, m.Message, app.Session.State.User.ID)

	// Build the embeds associated with the message.
	buildEmbeds(&b, m.Embeds)

	// Build the message attachments (attached files to the message).
	buildAttachments(&b, m.Attachments)

	// Tags with no region ID ([""]) do not start new regions. They can
	// therefore be used to mark the end of a region.
	b.WriteString("[\"\"]")

	b.WriteByte('\n')
	if str := b.String(); str != "" {
		b := make([]byte, len(str)+1)
		copy(b, str)

		return b
	}

	return nil
}

func buildDelete(app *App, d *astatine.MessageDelete) []byte {
	var b strings.Builder
	if d.BeforeDelete != nil {
		m := d.BeforeDelete
		b.WriteString("[\"")
		b.WriteString(m.ID)
		b.WriteString("\"]")
		buildReferencedMessage(&b, m, app.Session.State.User.ID, app)
		if app.Config.General.Timestamps {
			b.WriteString("[::d]")
			b.WriteString(m.Timestamp.Format(time.Stamp))
			b.WriteString("[::-]")
			b.WriteByte(' ')
		}
		b.WriteString("Message was Deleted")
		b.WriteString("[\"\"]")

		b.WriteByte('\n')
	}
	if str := b.String(); str != "" {
		b := make([]byte, len(str)+1)
		copy(b, str)

		return b
	}

	return nil
}

func buildMessage(app *App, m *astatine.Message) []byte {
	var b strings.Builder

	switch m.Type {
	case astatine.MessageTypeDefault, astatine.MessageTypeReply:
		// Define a new region and assign message ID as the region ID.
		// Learn more:
		// https://pkg.go.dev/github.com/rivo/tview#hdr-Regions_and_Highlights
		b.WriteString("[\"")
		b.WriteString(m.ID)
		b.WriteString("\"]")
		// Build the message associated with crosspost, channel follow add, pin, or a reply.
		if m.Type == astatine.MessageTypeReply {
			buildReferencedMessage(&b, m.ReferencedMessage, app.Session.State.User.ID, app)
		}

		if app.Config.General.Timestamps {
			b.WriteString("[::d]")
			b.WriteString(m.Timestamp.Format(time.Stamp))
			b.WriteString("[::-]")
			b.WriteByte(' ')
		}

		// Build the author of this message.
		if m.Member == nil {
			if len(m.GuildID) > 0 {
				member := discord.GetMember(app.Session, m.GuildID, m.Author.ID)
				buildAuthor(&b, m.Author, app.Session.State.User.ID, member, app)
			} else {
				c := app.SelectedChannel
				if c != nil && len(c.GuildID) > 0 {
					member := discord.GetMember(app.Session, c.GuildID, m.Author.ID)
					buildAuthor(&b, m.Author, app.Session.State.User.ID, member, app)
				} else {
					buildAuthor(&b, m.Author, app.Session.State.User.ID, nil, app) // dm channel, probably.
				}
			}
		} else {
			buildAuthor(&b, m.Author, app.Session.State.User.ID, m.Member, app)
		}

		// Build the contents of the message.
		buildContent(&b, m, app.Session.State.User.ID)

		if m.EditedTimestamp != nil {
			b.WriteString(" [::d](edited)[::-]")
		}

		// Build the embeds associated with the message.
		buildEmbeds(&b, m.Embeds)

		// Build the message attachments (attached files to the message).
		buildAttachments(&b, m.Attachments)

		// Tags with no region ID ([""]) do not start new regions. They can
		// therefore be used to mark the end of a region.
		b.WriteString("[\"\"]")

		b.WriteByte('\n')
	case astatine.MessageTypeGuildMemberJoin:
		b.WriteString("[#5865F2]")
		b.WriteString(m.Author.Username)
		b.WriteString("[-] joined the server.")

		b.WriteByte('\n')
	case astatine.MessageTypeCall:
		b.WriteString("[#5865F2]")
		b.WriteString(m.Author.Username)
		b.WriteString("[-] started a call.")

		b.WriteByte('\n')
	case astatine.MessageTypeChannelPinnedMessage:
		b.WriteString("[#5865F2]")
		b.WriteString(m.Author.Username)
		b.WriteString("[-] pinned a message.")

		b.WriteByte('\n')
	}

	if str := b.String(); str != "" {
		b := make([]byte, len(str)+1)
		copy(b, str)

		return b
	}

	return nil
}

func buildReferencedMessage(b *strings.Builder, rm *astatine.Message, clientID string, app *App) {
	if rm != nil {
		b.WriteString(" ╭ ")
		b.WriteString("[::d]")
		if rm.Member == nil {
			if len(rm.GuildID) > 0 {
				member := discord.GetMember(app.Session, rm.GuildID, rm.Author.ID)
				buildAuthor(b, rm.Author, clientID, member, app)
			} else {
				c := app.SelectedChannel
				if c != nil && len(c.GuildID) > 0 {
					member := discord.GetMember(app.Session, c.GuildID, rm.Author.ID)
					buildAuthor(b, rm.Author, clientID, member, app)
				} else {
					buildAuthor(b, rm.Author, clientID, nil, app) // dm channel, probably.
				}
			}
		} else {
			buildAuthor(b, rm.Author, clientID, rm.Member, app)
		}

		if rm.Content != "" {
			rm.Content = buildMentions(rm.Content, rm.Mentions, clientID)
			var content string
			if len(rm.Content) > 100 {
				content = rm.Content[:100] + "..."
			} else {
				content = rm.Content
			}
			content = strings.ReplaceAll(content, "\n", " ")
			b.WriteString(discord.ParseMarkdown(content))
		}

		b.WriteString("[::-]")
		b.WriteByte('\n')
	} else {
		b.WriteString(" ╭ [::d]Original message deleted[::-]")
		b.WriteByte('\n')
	}
}

func buildContent(b *strings.Builder, m *astatine.Message, clientID string) {
	if m.Content != "" {
		content := buildMentions(m.Content, m.Mentions, clientID)
		re := regexp.MustCompile(`\n{3,}`)
		content = re.ReplaceAllString(content, "\n\n")
		b.WriteString(discord.ParseMarkdown(content))
	}
}

func buildEmbeds(b *strings.Builder, es []*astatine.MessageEmbed) {
	for _, e := range es {
		if e.Type != astatine.EmbedTypeRich {
			continue
		}

		var (
			embedBuilder strings.Builder
			hasHeading   bool
		)
		prefix := fmt.Sprintf("[#%06X]▐[-] ", e.Color)

		b.WriteByte('\n')
		embedBuilder.WriteString(prefix)

		if e.Author != nil {
			hasHeading = true
			embedBuilder.WriteString("[::u]")
			embedBuilder.WriteString(e.Author.Name)
			embedBuilder.WriteString("[::-]")
		}

		if e.Title != "" {
			if hasHeading {
				embedBuilder.WriteByte('\n')
				embedBuilder.WriteByte('\n')
			}

			embedBuilder.WriteString("[::b]")
			embedBuilder.WriteString(e.Title)
			embedBuilder.WriteString("[::-]")
		}

		if e.Description != "" {
			if hasHeading {
				embedBuilder.WriteByte('\n')
				embedBuilder.WriteByte('\n')
			}

			embedBuilder.WriteString(discord.ParseMarkdown(e.Description))
		}

		if len(e.Fields) != 0 {
			if hasHeading || e.Description != "" {
				embedBuilder.WriteByte('\n')
				embedBuilder.WriteByte('\n')
			}

			for i, ef := range e.Fields {
				embedBuilder.WriteString("[::b]")
				embedBuilder.WriteString(ef.Name)
				embedBuilder.WriteString("[::-]")
				embedBuilder.WriteByte('\n')
				embedBuilder.WriteString(discord.ParseMarkdown(ef.Value))

				if i != len(e.Fields)-1 {
					embedBuilder.WriteString("\n\n")
				}
			}
		}

		if e.Footer != nil {
			if hasHeading {
				embedBuilder.WriteString("\n\n")
			}

			embedBuilder.WriteString(e.Footer.Text)
		}

		b.WriteString(strings.ReplaceAll(embedBuilder.String(), "\n", "\n"+prefix))
	}
}

func buildAttachments(b *strings.Builder, as []*astatine.MessageAttachment) {
	for _, a := range as {
		b.WriteByte('\n')
		b.WriteString("FILE: ")
		b.WriteString(a.URL)
	}
}

func buildMentions(content string, mentions []*astatine.User, clientID string) string {
	for _, mUser := range mentions {
		var color string
		if mUser.ID == clientID {
			color = "[:#FFA500]"
		} else {
			color = "[#EB459E]"
		}

		content = strings.NewReplacer(
			// <@USER_ID>
			"<@"+mUser.ID+">",
			color+"@"+mUser.Username+"[-:-]",
			// <@!USER_ID>
			"<@!"+mUser.ID+">",
			color+"@"+mUser.Username+"[-:-]",
		).Replace(content)
	}

	return content
}

func buildAuthor(b *strings.Builder, u *astatine.User, clientID string, m *astatine.Member, app *App) {
	if m != nil && len(m.Nick) > 0 {
		b.WriteString("! ")
	}
	var gotRoleColor bool = false
	/*if app != nil && m != nil && len(m.Roles) >= 1 {
		//r, err := app.Session.State.Role(m.GuildID)
		//if r != nil {
		//TODO
		//}
	}*/
	if !gotRoleColor {
		if u.ID == clientID {
			b.WriteString("[#57F287]")
		} else {
			b.WriteString("[#ED4245]")
		}
	}

	if m != nil && len(m.Nick) > 0 {
		b.WriteString(m.Nick)
	} else {
		b.WriteString(u.Username)
	}
	b.WriteString("[-] ")
	if m != nil && m.CommunicationDisabledUntil != nil && m.CommunicationDisabledUntil.After(time.Now()) {
		b.WriteString("[#FFFF00]MUTED[-] ")
	}
	// If the message author is a bot account, render the message with bot label
	// for distinction.
	if u.Bot {
		b.WriteString("[#EB459E]BOT[-] ")
	}
}
