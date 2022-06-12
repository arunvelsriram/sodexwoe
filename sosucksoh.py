import argparse
from pathlib import Path
import pikepdf
import os

_config = {
  'airtel_mobile': {
    'password': os.environ['AIRTEL_MOBILE_PASSWORD'],
    'no_pages_to_keep': 4,
  },
  'jio_mobile': {
    'password': os.environ['JIO_MOBILE_PASSWORD'],
    'no_pages_to_keep': 5,
  },
  'jio_fiber': {
    'password': os.environ['JIO_FIBER_PASSWORD'],
    'no_pages_to_keep': 5,
  }
}

def main(args: argparse.Namespace):
  print(f'Bill type is: {args.bill_type}')
  config = _config[args.bill_type]
  in_bill_path = Path(args.in_bill_path)
  out_bill_path = in_bill_path.parent / f'final_{args.bill_type}_{in_bill_path.name}'
  print(f'Reading input {in_bill_path}')
  with pikepdf.open(str(in_bill_path), password=config['password']) as pdf:
    no_pages_to_keep = config['no_pages_to_keep']
    page_nos = list(range(0, len(pdf.pages)+1))
    print(f'Deleting pages: {page_nos[no_pages_to_keep:]}')
    del pdf.pages[no_pages_to_keep:]
    print(f'Saving output {out_bill_path}')
    pdf.save(out_bill_path)


if __name__ == "__main__":
  parser = argparse.ArgumentParser(description='sosucksoh')
  parser.add_argument('in_bill_path', help='path to PDF file')
  parser.add_argument('--bill-type', dest='bill_type', help='type of bill (airtel_mobile, jio_mobile, jio_fiber)')
  
  args = parser.parse_args()
  print('Parsed CLI Arguments:', args)
  
  main(args)
