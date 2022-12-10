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
	"github.com/urfave/cli/v2"
)

var GoogleAPICredentials string

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
				Name:  "config",
				Usage: "Sodexwoe config",
				Subcommands: []*cli.Command{
					{
						Name:  "view",
						Usage: "View configuration",
						Action: func(ctx *cli.Context) error {
							return config.DumpConfig()
						},
					},
				},
			},
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
					input := ctx.Args().Get(0)
					output := fmt.Sprintf("%s--%s", billName, filepath.Base(input))
					billConverterSrv := services.NewBillConverterService(cfg)
					err := billConverterSrv.ConvertFile(billName, input, output)
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
						Name:        "names",
						Aliases:     []string{"n"},
						Usage:       fmt.Sprintf("Comma separated bill names from: %v", strings.Join(billNames, ", ")),
						Value:       cli.NewStringSlice(billNames...),
						DefaultText: strings.Join(billNames, ","),
						Required:    false,
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
					billConverterSrv := services.NewBillConverterService(cfg)

					emails, err := billEmailSrv.GetEmails(billNames, year, month)
					if err != nil {
						return err
					}

					for _, email := range emails {
						log.WithField("billName", email.BillName).WithField("filename", email.Bill.Filename).Info("converting file")
						outputFilename := fmt.Sprintf("%s_%s_%d--%s", email.BillName, email.Month.String(), email.Year, filepath.Base(email.Bill.Filename))
						output := filepath.Join(cfg.DownloadDir, email.BillName, outputFilename)
						log.WithField("output", output).Info("creating output file")
						outputFile, err := utils.CreateFile(output)
						if err != nil {
							return err
						}

						err = billConverterSrv.Convert(email.BillName, bytes.NewReader(email.Bill.Data), outputFile)
						if err != nil {
							return err
						}
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
