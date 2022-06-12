# sosucksoh (So-Sucks-Oh!)

CLI to remove password and service usage (for my privacy) from Mobile and Internet bills before submitting it for reimbursement.

I use email as a single source for downloading all bills using filters and labels to avoid logging in to various service provider websites.

## Setup

### Dependencies

  1. Python 3
  2. miniconda3 `brew install miniconda`

### Running

```
# create conda environment
conda env create -f environment.yml

# configs
cp .envrc.template .envrc
source .envrc

# run
python sosucksoh.py --bill-type airtel_mobile /Users/arunvelsriram/Downloads/my-airtel-bill.pdf
```
