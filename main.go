package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/arunvelsriram/sodexwoe/internal/config"
	"github.com/arunvelsriram/sodexwoe/internal/services"
	"github.com/arunvelsriram/sodexwoe/internal/utils"
	pdfcpuapi "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/urfave/cli/v2"
)

var GoogleAPICredentials string

func billConvert(cfg config.Config, billName, inputBill, outputBill string) error {
	billConfig, ok := cfg[billName]
	if !ok {
		log.WithField("billName", billName).Debug("bill name not found in config")
		return fmt.Errorf("billName: %s not found in config", billName)
	}

	pdfConfig := pdfcpu.NewDefaultConfiguration()
	pdfConfig.UserPW = billConfig.Password
	buffer := bytes.NewBuffer([]byte{})

	inputFile, err := os.Open(inputBill)
	if err != nil {
		log.WithField("inputBill", inputBill).Debug("failed to open input bill")
		return err
	}

	log.WithField("inputBill", inputBill).Info("removing password from bill")
	err = pdfcpuapi.Decrypt(inputFile, buffer, pdfConfig)
	if err != nil {
		log.WithField("inputBill", inputBill).Debug("failed to remove password from bill")
		return err
	}

	log.WithField("inputBill", inputBill).Info("removing unnecesssary pages from bill")
	decryptedBill := bytes.NewReader(buffer.Bytes())
	buffer.Reset()
	pagesToRemove := []string{fmt.Sprintf("%d-", billConfig.KeepPages+1)}
	err = pdfcpuapi.RemovePages(decryptedBill, buffer, pagesToRemove, pdfConfig)
	if err != nil {
		log.WithField("pagesToRemove", pagesToRemove).Debug("failed to remove pages")
		return err
	}

	log.WithField("outputBiii", outputBill).Info("writing converted bill")
	err = os.WriteFile(outputBill, buffer.Bytes(), 0644)
	if err != nil {
		log.WithField("outputBill", outputBill).Debug("failed to write converted output")
		return err
	}

	return nil
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("unable to load configuration: %v", err)
	}

	billNames := cfg.BillNames()

	app := &cli.App{
		Name:  "sodexwoe",
		Usage: "Sodexo Woe!",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log-level",
				Aliases: []string{"l"},
				Usage: fmt.Sprintf("Log level. Could be one of: %v", strings.Join([]string{log.PanicLevel.String(), log.FatalLevel.String(), log.ErrorLevel.String(),
					log.WarnLevel.String(), log.InfoLevel.String(), log.DebugLevel.String(), log.TraceLevel.String()}, ", ")),
				Required: false,
				Value:    log.InfoLevel.String(),
			},
		},
		Before: func(ctx *cli.Context) error {
			logLevel := ctx.String("log-level")
			level, err := log.ParseLevel(logLevel)
			if err != nil {
				return err
			}
			log.SetLevel(level)

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "bill-convert",
				Aliases: []string{"bc"},
				Usage:   "Convert bill for uploading to Sodexo",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "name",
						Aliases:  []string{"n"},
						Usage:    fmt.Sprintf("Bill name. Could be one of: %v", strings.Join(billNames, ", ")),
						Required: true,
					},
				},
				Action: func(ctx *cli.Context) error {
					if ctx.NArg() == 0 {
						return errors.New("path to bill file is required")
					}

					billName := ctx.String("name")
					inputBill := ctx.Args().Get(0)
					outputBill := filepath.Join(filepath.Dir(inputBill), fmt.Sprintf("converted_%s_%s", billName, filepath.Base(inputBill)))
					err = billConvert(cfg, billName, inputBill, outputBill)
					if err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name:    "bill-download",
				Aliases: []string{"bd"},
				Usage:   "Download bills from Gmail and convert them for uploading to Sodexo",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "year",
						Aliases:  []string{"y"},
						Usage:    "Year",
						Value:    time.Now().Local().Year(),
						Required: false,
					},
					&cli.StringFlag{
						Name:     "month",
						Aliases:  []string{"m"},
						Usage:    "Case-insensitive short or long month name",
						Value:    time.Now().Local().Month().String(),
						Required: false,
					},
					&cli.StringSliceFlag{
						Name:     "names",
						Aliases:  []string{"n"},
						Usage:    fmt.Sprintf("Comma separated bill names from: %v", strings.Join(billNames, ", ")),
						Value:    cli.NewStringSlice(billNames...),
						Required: false,
					},
				},
				Action: func(ctx *cli.Context) error {
					billNames := ctx.StringSlice("names")
					year := ctx.Int("year")
					monthFlagValue := ctx.String("month")
					month, err := utils.GetMonthByName(monthFlagValue)
					if err != nil {
						return err
					}

					gmailSrv, err := services.NewGmailService(GoogleAPICredentials)
					if err != nil {
						return err
					}
					billEmailSrv := services.NewBillEmailService(gmailSrv, cfg)
					emails, err := billEmailSrv.GetEmails(billNames, year, month)
					if err != nil {
						return err
					}
					fmt.Printf("emails: %v", emails)

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
