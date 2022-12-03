package services

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/arunvelsriram/sodexwoe/internal/config"
	"github.com/arunvelsriram/sodexwoe/internal/constants"
	"github.com/arunvelsriram/sodexwoe/internal/models"
	"github.com/arunvelsriram/sodexwoe/internal/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/gmail/v1"
)

type BillEmailService interface {
	GetLabels(billNames ...string) (models.BillEmailLabels, error)
	GetEmails(billNames []string, year int, month time.Month) (models.BillEmails, error)
}

type billEmailService struct {
	gmailSrv *gmail.Service
	cfg      config.Config
}

func (s billEmailService) GetLabels(billNames ...string) (models.BillEmailLabels, error) {
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
	result := make(models.BillEmailLabels, 0, len(labelsResponse.Labels))
	for i := 0; i < len(labelNames); i++ {
		if label, ok := labelNameToLabel[labelNames[i]]; ok {
			billEmailLabel := models.BillEmailLabel{
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

func (s billEmailService) GetEmails(billNames []string, year int, month time.Month) (models.BillEmails, error) {
	billEmailLabels, err := s.GetLabels(billNames...)
	if err != nil {
		return nil, err
	}

	labelQ := utils.AnyLabelQ(billEmailLabels.LabelNames()...)
	dateRangeQ := utils.WithinMonthQ(year, month)
	q := fmt.Sprintf("%s %s", labelQ, dateRangeQ)
	log.WithField("query", q).Info("listing emails from gmail")
	messagesResponse, err := s.gmailSrv.Users.Messages.List(constants.GMAIL_USER).Q(q).Do()
	if err != nil {
		return nil, err
	}

	messages := messagesResponse.Messages
	log.Debugf("listed emails: %d", len(messages))
	result := make(models.BillEmails, 0, len(messages))
	log.Info("fetching emails from gmail")
	for _, m := range messagesResponse.Messages {
		log.WithField("messageId", m.Id).Debug("fetching email")
		message, err := s.gmailSrv.Users.Messages.Get(constants.GMAIL_USER, m.Id).Do()
		if err != nil {
			return nil, err
		}

		log.WithField("messageId", m.Id).
			WithField("messageLabelIds", message.LabelIds).
			Debug("determining bill label for the message")
		var billEmailLabel *models.BillEmailLabel
		for _, labelId := range message.LabelIds {
			if billEmailLabel = billEmailLabels.FindById(labelId); billEmailLabel != nil {
				break
			}
		}
		if billEmailLabel == nil {
			log.WithField("messageId", message.Id).
				WithField("messageLabelIds", message.LabelIds).
				WithField("billLabelIds", billEmailLabels.LabelIds()).
				WithField("billLabelNames", billEmailLabels.LabelNames()).
				Errorf("unexpected email - email labels not having any of the bill labels")
			return nil, fmt.Errorf("got unexpected email, messageId: %v", message.Id)
		}

		var attachmentId string
		var attachmentFilename string
		for _, p := range message.Payload.Parts {
			if p.Filename != "" && strings.Contains(p.Filename, ".pdf") {
				attachmentId = p.Body.AttachmentId
				attachmentFilename = p.Filename
				log.WithField("messageId", message.Id).
					WithField("attachmentId", attachmentId).
					WithField("filename", attachmentFilename).
					Debug("found pdf attachment in email")
				break
			}
		}
		if attachmentId == "" || attachmentFilename == "" {
			log.WithField("messageId", message.Id).Error("no attachment found in email")
			return nil, fmt.Errorf("no attachment found in email, messageId: %v", message.Id)
		}

		log.WithField("attachmentId", attachmentId).Debug("fetching attachment")
		attachmentRes, err := s.gmailSrv.Users.Messages.Attachments.Get(constants.GMAIL_USER, message.Id, attachmentId).Do()
		if err != nil {
			return nil, err
		}

		log.Debug("decoding attachment content")
		content, err := base64.URLEncoding.DecodeString(attachmentRes.Data)
		if err != nil {
			return nil, err
		}

		billEmail := models.BillEmail{
			BillName: billEmailLabel.BillName,
			Year:     year,
			Month:    month,
			Bill: models.Bill{
				Filename: attachmentFilename,
				Data:     content,
			},
		}
		result = append(result, billEmail)
	}
	log.Debugf("fetched emails: %d", len(result))

	return result, nil
}

func NewBillEmailService(gmailSrv *gmail.Service, cfg config.Config) BillEmailService {
	return billEmailService{gmailSrv, cfg}
}
