package imap

import (
	"log"
	"mime/quotedprintable"
	"strings"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/lnk00/prosp/models"
	"github.com/spf13/viper"
)

type Imap struct {
	Client *imapclient.Client
}

func New() Imap {

	client, err := imapclient.DialInsecure(viper.GetString("host")+":"+viper.GetString("port"), nil)
	if err != nil {
		log.Fatalf("failed to dial IMAP server: %v", err)
	}

	return Imap{
		Client: client,
	}
}

func (I Imap) Login() {
	if err := I.Client.Login(viper.GetString("user"), viper.GetString("password")).Wait(); err != nil {
		log.Fatalf("failed to login: %v", err)
	}
}

func (I Imap) GetMailboxes() []*imap.ListData {

	mailboxes, err := I.Client.List("", "%", nil).Collect()
	if err != nil {
		log.Fatalf("failed to list mailboxes: %v", err)
	}

	return mailboxes
}

func (I Imap) Logout() {
	if err := I.Client.Logout().Wait(); err != nil {
		log.Fatalf("failed to logout: %v", err)
	}

	I.Client.Close()
}

func (I Imap) FetchFrom(email string) []models.Message {
	messageList := []models.Message{}

	I.Client.Select("INBOX", nil).Wait()

	data, err := I.Client.UIDSearch(&imap.SearchCriteria{
		Header: []imap.SearchCriteriaHeaderField{{Key: "FROM", Value: email}},
	}, nil).Wait()
	if err != nil {
		log.Fatalf("FETCH command failed: %v", err)
	}

	if len(data.AllUIDs()) == 0 {
		return messageList
	}

	fetchOptions := &imap.FetchOptions{
		Flags:    true,
		Envelope: true,
		BodySection: []*imap.FetchItemBodySection{
			{Specifier: imap.PartSpecifierText},
		},
	}
	messages, err := I.Client.Fetch(data.All, fetchOptions).Collect()
	if err != nil {
		log.Fatalf("FETCH command failed: %v", err)
	}

	for _, msg := range messages {
		var body []byte
		for _, buf := range msg.BodySection {
			body = buf
			break
		}

		messageList = append(messageList, models.Message{
			Date:   msg.Envelope.Date,
			Reader: quotedprintable.NewReader(strings.NewReader(string(body))),
		})
	}

	return messageList

}
