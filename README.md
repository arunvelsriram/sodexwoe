# sodexwoe (Sodexo Woe!)

CLI to download Mobile and Internet bills from Gmail, remove password and service usage (for my privacy) before submitting for reimbursement.

I use email as a single source for downloading all bills using filters and labels to avoid signing in to various service provider websites.

### Install

#### Auto

```
curl https://raw.githubusercontent.com/arunvelsriram/sodexwoe/main/install.sh | bash
```

#### Manual
Get the download URL based on your OS and Platform from [releases page](https://github.com/arunvelsriram/sodexwoe/releases/latest).

```
wget <download-url>
tar -xvzf sodexwoe_*.tar.gz
sudo mv sodexwoe /usr/local/bin/sodexwoe
```

### Usage

### Configure

```
mkdir -p ~/.config/sodexwoe/
cp ./config.sample.yaml ~/.config/sodexwoe/config.yaml

# update the config
vim ~/.config/sodexwoe/config.yaml
```

#### Config location: `~/.config/sodexwoe/config.yaml`
#### Sample configuration for reference: [config.sample.yaml](config.sample.yaml)

### Run

```
sodexwoe --help
sodexwoe config view
sodexwoe bill-convert --name jio path/to/bill.pdf
sodexwoe bill-download --name jio
```

## Development

```
go mod tidy -v
go run main.go
```
