#! /usr/bin/env bash

set -eo pipefail

case $(uname) in
	Linux)
		os_identified="Linux"
    ;;
  Darwin)
		os_identified="Mac"
		;;
	Mac)
		os_identified="Mac"
		;;
	*)
		echo "Could not determine the Operating System"
		exit 1
		;;
esac

arch_identified=$(uname -m)

SODEXWOE_OS=${SODEXWOE_OS:-${os_identified}}
SODEXWOE_ARCH=${SODEXWOE_ARCH:-${arch_identified}}

echo ">>> OS: ${SODEXWOE_OS}"
echo ">>> ARCH: ${SODEXWOE_ARCH}"

tag=$(curl -s https://api.github.com/repos/arunvelsriram/sodexwoe/releases/latest | jq -r '.tag_name')
echo ">>> Latest release: ${tag}"
filename=sodexwoe_${tag}_${SODEXWOE_OS}_${SODEXWOE_ARCH}.tar.gz
echo ">> Downloading: ${filename}"
curl -SL https://github.com/arunvelsriram/sodexwoe/releases/download/${tag}/${filename} -o /tmp/${filename}
tar -C /tmp -xzf /tmp/${filename}
sudo mv /tmp/sodexwoe /usr/local/bin/sodexwoe
echo ">>> Installation completed"
