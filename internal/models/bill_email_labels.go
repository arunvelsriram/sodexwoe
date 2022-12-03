package models

import "google.golang.org/api/gmail/v1"

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

func (billEmailLabels BillEmailLabels) LabelIds() []string {
	labelIds := make([]string, 0, len(billEmailLabels))
	for _, it := range billEmailLabels {
		labelIds = append(labelIds, it.Id)
	}

	return labelIds
}

func (billEmailLabels BillEmailLabels) FindById(id string) *BillEmailLabel {
	for _, it := range billEmailLabels {
		if it.Id == id {
			return &it
		}
	}

	return nil
}
