package services

import (
	"bytes"
	"fmt"
	"os"

	"github.com/arunvelsriram/sodexwoe/internal/config"
	pdfcpuapi "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	log "github.com/sirupsen/logrus"
)

type BillConverterService interface {
	ConvertFile(billName, input, output string) error
}

type billConverterService struct {
	cfg       config.Config
	pdfCpuCfg *pdfcpu.Configuration
}

func (s billConverterService) ConvertFile(billName, input, output string) error {
	billConfig, ok := s.cfg[billName]
	if !ok {
		return fmt.Errorf("billName: %s not found in config", billName)
	}

	s.pdfCpuCfg.UserPW = billConfig.Password
	buffer := bytes.NewBuffer([]byte{})

	log.WithField("input", input).Info("opening input bill")
	inputFile, err := os.Open(input)
	if err != nil {
		return err
	}

	log.Info("removing password from bill")
	err = pdfcpuapi.Decrypt(inputFile, buffer, s.pdfCpuCfg)
	if err != nil {
		return err
	}

	log.Info("removing unnecesssary pages from bill")
	decryptedBill := bytes.NewReader(buffer.Bytes())
	buffer.Reset()
	pagesToRemove := []string{fmt.Sprintf("%d-", billConfig.KeepPages+1)}
	err = pdfcpuapi.RemovePages(decryptedBill, buffer, pagesToRemove, s.pdfCpuCfg)
	if err != nil {
		return err
	}

	log.WithField("output", output).Info("writing converted bill")
	err = os.WriteFile(output, buffer.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func NewBillConverterService(cfg config.Config) BillConverterService {
	return billConverterService{cfg, pdfcpu.NewDefaultConfiguration()}
}
