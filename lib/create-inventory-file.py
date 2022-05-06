#!/usr/bin/env python3

import argparse
import sys
import json
import pathlib

from pylib.error import error


class CreateInventoryFile(object):
    def __init__(self):
        parser = argparse.ArgumentParser(description='Create inventory file')
        parser.add_argument(
            '-i', '--remote-ips',
            dest='remote_ips',
            required=True
        )

        parser.add_argument(
            '-o', '--inventory-file',
            dest='inventory_file',
            required=True,
            help='the output file'
        )

        args = parser.parse_args()

        ips = json.loads(args.remote_ips)

        this_dir = pathlib.Path(__file__).parent.resolve()
        template_file = this_dir.joinpath('cm', 'inventory.template')

        inventory_file = pathlib.Path(args.inventory_file).resolve()
        inventory_file.parent.mkdir(parents=True, exist_ok=True)

        with open(template_file, 'r') as fh:
            template_str = fh.read()

        with open(args.inventory_file, 'w') as fh:
            fh.write(template_str.format(remote_ips="\n".join(ips)))


if __name__ == '__main__':
    if len(sys.argv) < 2:
        error('must pass remote IPs')

    CreateInventoryFile()
