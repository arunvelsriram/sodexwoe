# sodexwoe (Sodexo Woe!)

CLI to remove password and service usage (for my privacy) from Mobile and Internet bills before submitting them for reimbursement.

I use email as a single source for downloading all bills using filters and labels to avoid logging in to various service provider websites.

### Install

```
wget https://github.com/arunvelsriram/sodexwoe/releases/download/v1.1.0/sodexwoe.tar.gz
tar -xvzf sodexwoe.tar.gz
chmod +x sodexwoe-<your-cpu-architecture>
sudo mv sodexwoe-<your-cpu-architecture> /usr/local/bin/sodexwoe
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
sodexwoe bill-convert --name jio path/to/bill.pdf
sodexwoe bill-download --name jio
```

## Development

```
go mod tidy -v
go run main.go
```
