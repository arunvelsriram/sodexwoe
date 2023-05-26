package services

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/arunvelsriram/sodexwoe/internal/config"
	"github.com/arunvelsriram/sodexwoe/internal/utils"
	pdfcpuapi "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	log "github.com/sirupsen/logrus"
)

type BillConverterService interface {
	ConvertFile(billName, input, output string) error
	Convert(billName string, input io.ReadSeeker, output io.Writer) error
}

type billConverterService struct {
	cfg       config.Config
	pdfCpuCfg *pdfcpu.Configuration
}

func (s billConverterService) ConvertFile(billName, input, output string) error {
	billConfig, ok := s.cfg.BillConfigs[billName]
	if !ok {
		return fmt.Errorf("billName: %s not found in config", billName)
	}

	s.pdfCpuCfg.UserPW = billConfig.Password

	log.WithField("input", input).Info("opening input bill")
	inputFile, err := os.Open(input)
	if err != nil {
		return err
	}
	defer func() {
		if err = inputFile.Close(); err != nil {
			log.Error(err)
		}
	}()

	output = filepath.Join(s.cfg.DownloadDir, billName, output)
	log.WithField("output", output).Info("creating output file")
	outputFile, err := utils.CreateFile(output)
	if err != nil {
		return err
	}

	return s.Convert(billName, inputFile, outputFile)
}

func (s billConverterService) Convert(billName string, input io.ReadSeeker, output io.Writer) error {
	billConfig, ok := s.cfg.BillConfigs[billName]
	if !ok {
		return fmt.Errorf("billName: %s not found in config", billName)
	}

	s.pdfCpuCfg.UserPW = billConfig.Password
	buffer := bytes.NewBuffer([]byte{})

	log.Info("removing password from bill")
	err := pdfcpuapi.Decrypt(input, buffer, s.pdfCpuCfg)
	if err != nil {
		return err
	}
	decryptedBill := bytes.NewReader(buffer.Bytes())
	buffer.Reset()

	log.Info("writing addditional text in the bill")
	textSpec := "sc:0.5 abs, points:14, pos:br, rot:0, offset:-5 5, color:Black"
	wm, err := pdfcpuapi.TextWatermark(fmt.Sprintf("%1s", billConfig.AdditionalText), textSpec, true, true, pdfcpu.POINTS)
	if err != nil {
		return err
	}
	err = pdfcpuapi.AddWatermarks(decryptedBill, buffer, []string{"1"}, wm, s.pdfCpuCfg)
	if err != nil {
		return err
	}
	gstBill := bytes.NewReader(buffer.Bytes())
	buffer.Reset()

	log.Info("removing unnecesssary pages from bill")
	pagesToRemove := []string{fmt.Sprintf("%d-", billConfig.KeepPages+1)}
	err = pdfcpuapi.RemovePages(gstBill, buffer, pagesToRemove, s.pdfCpuCfg)
	if err != nil {
		return err
	}

	log.Info("writing bill output")
	_, err = output.Write(buffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func NewBillConverterService(cfg config.Config) BillConverterService {
	return billConverterService{cfg, pdfcpu.NewDefaultConfiguration()}
}
