# sodexwoe (Sodexo Woe!)

CLI to remove password and service usage (for my privacy) from Mobile and Internet bills before submitting them for reimbursement.

I use email as a single source for downloading all bills using filters and labels to avoid logging in to various service provider websites.

## Setup

### Install

```
go install github.com/arunvelsriram/sodexwoe@latest
```

### Running

Place your configuration in `~/.config/sodexwoe/config.yaml`
Sample configuration for reference: [config.sample.yaml](config.sample.yaml)

```
sodexwoe --help
sodexwoe bill-convert --name work_mobile path/to/bill.pdf
```

## Development

```
go mod tidy -v
go run main.go
```
