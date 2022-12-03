package service

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/arunvelsriram/sodexwoe/internal/config"
	"github.com/arunvelsriram/sodexwoe/internal/constants"
	"github.com/arunvelsriram/sodexwoe/internal/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/gmail/v1"
)

type BillEmails []BillEmail

type BillEmail struct {
	BillName  string
	LabelName string
	Message   *gmail.Message
	Year      int
	Month     time.Month
}

type BillEmailAttachments []BillEmailAttachment

type BillEmailAttachment struct {
	BillName  string
	LabelName string
	Year      int
	Month     time.Month
	Filename  string
	Data      []byte
}

type BillEmailLabels []BillEmailLabel

type BillEmailLabel struct {
	*gmail.Label
	BillName string
}

func (billEmailLabels BillEmailLabels) LabelNames() []string {
	labelNames := make([]string, 0, len(billEmailLabels))
	for _, it := range billEmailLabels {
		labelNames = append(labelNames, it.Name)
	}

	return labelNames
}

func (billEmailLabels BillEmailLabels) FindById(id string) *BillEmailLabel {
	for _, it := range billEmailLabels {
		if it.Id == id {
			return &it
		}
	}

	return nil
}

type BillEmailService struct {
	gmailSrv *gmail.Service
	cfg      config.Config
}

func (s BillEmailService) getLabels(billNames ...string) (BillEmailLabels, error) {
	labelNames, err := s.cfg.Labels(billNames)
	if err != nil {
		return nil, err
	}

	log.Info("listing labels from gmail")
	labelsResponse, err := s.gmailSrv.Users.Labels.List(constants.GMAIL_USER).Do()
	if err != nil {
		return nil, err
	}
	log.Debugf("listed labels: %v", len(labelsResponse.Labels))
	labelNameToLabel := make(map[string]*gmail.Label, len(labelsResponse.Labels))
	for _, label := range labelsResponse.Labels {
		labelNameToLabel[label.Name] = label
	}

	log.Debug("filtering listed labels")
	result := make(BillEmailLabels, 0, len(labelsResponse.Labels))
	for i := 0; i < len(labelNames); i++ {
		if label, ok := labelNameToLabel[labelNames[i]]; ok {
			billEmailLabel := BillEmailLabel{
				BillName: billNames[i],
				Label:    label,
			}
			result = append(result, billEmailLabel)
		} else {
			return nil, fmt.Errorf("label not found: %v", labelNames[i])
		}
	}
	log.Debugf("filtered labels: %v", len(result))

	return result, nil
}

func (s BillEmailService) GetEmails(billNames []string, year int, month time.Month) (BillEmails, error) {
	billEmailLabels, err := s.getLabels(billNames...)
	if err != nil {
		return nil, err
	}

	labelQs := make([]string, 0, len(billEmailLabels.LabelNames()))
	for _, labelName := range billEmailLabels.LabelNames() {
		labelQs = append(labelQs, fmt.Sprintf("label:\"%s\"", labelName))
	}
	labelQ := strings.Join(labelQs, " OR ")

	layout := "2006/01/02"
	afterDate := utils.StartOfMonth(year, month)
	beforeDate := utils.EndOfMonth(year, month)
	dateRangeQ := fmt.Sprintf("after:%s before:%s", afterDate.Format(layout), beforeDate.Format(layout))

	q := fmt.Sprintf("%s %s", labelQ, dateRangeQ)
	log.WithField("query", q).Info("listing emails from gmail")
	messagesResponse, err := s.gmailSrv.Users.Messages.List(constants.GMAIL_USER).Q(q).Do()
	if err != nil {
		return nil, err
	}

	messages := messagesResponse.Messages
	log.Debugf("listed emails: %d", len(messages))
	result := make(BillEmails, 0, len(messages))
	log.Info("fetching emails from gmail")
	for _, m := range messagesResponse.Messages {
		log.WithField("messageId", m.Id).Debug("fetching email")
		message, err := s.gmailSrv.Users.Messages.Get(constants.GMAIL_USER, m.Id).Do()
		if err != nil {
			return nil, err
		}

		var billEmailLabel *BillEmailLabel
		for _, labelId := range message.LabelIds {
			if billEmailLabel = billEmailLabels.FindById(labelId); billEmailLabel != nil {
				break
			}
		}
		if billEmailLabel == nil {
			return nil, fmt.Errorf("got unexpected email: %v", err)
		}

		billEmail := BillEmail{
			BillName:  billEmailLabel.BillName,
			LabelName: billEmailLabel.Name,
			Message:   message,
			Year:      year,
			Month:     month,
		}
		result = append(result, billEmail)
	}
	log.Debugf("fetched emails: %d", len(result))

	return result, nil
}

func (s BillEmailService) GetAttachments(billEmails BillEmails) (BillEmailAttachments, error) {
	billEmailAttachments := make(BillEmailAttachments, 0, len(billEmails))
	log.Info("searching pdf attachment in messages and fetching the attachment")
	for _, b := range billEmails {
		var attachmentId string
		var filename string
		message := b.Message
		log.WithField("messageId", message.Id).Debug("searching pdf attachment in the message")
		for _, p := range message.Payload.Parts {
			if p.Filename != "" && strings.Contains(p.Filename, ".pdf") {
				attachmentId = p.Body.AttachmentId
				filename = p.Filename
				log.WithField("messageId", message.Id).
					WithField("attachmentId", attachmentId).
					WithField("filename", filename).
					Debug("found pdf attachment in email")
				break
			}
		}

		log.WithField("attachmentId", attachmentId).Debug("fatching attachment")
		attachmentRes, err := s.gmailSrv.Users.Messages.Attachments.Get(constants.GMAIL_USER, message.Id, attachmentId).Do()
		if err != nil {
			return nil, err
		}

		log.Debug("decoding attachment content")
		content, err := base64.URLEncoding.DecodeString(attachmentRes.Data)
		if err != nil {
			return nil, err
		}

		billEmailAttachment := BillEmailAttachment{
			BillName:  b.BillName,
			LabelName: b.LabelName,
			Year:      b.Year,
			Month:     b.Month,
			Filename:  filename,
			Data:      content,
		}
		billEmailAttachments = append(billEmailAttachments, billEmailAttachment)
	}
	log.Debug("total attachments: %v", len(billEmailAttachments))

	return billEmailAttachments, nil
}

func NewBillEmailService(gmailSrv *gmail.Service, cfg config.Config) BillEmailService {
	return BillEmailService{gmailSrv, cfg}
}
