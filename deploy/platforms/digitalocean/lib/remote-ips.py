#!/usr/bin/env python3

import pathlib
import json
import sys
import argparse


class RemoteIPs(object):
    def __init__(self):
        parser = argparse.ArgumentParser(description='Show remote IPs')
        parser.add_argument(
            '-f', '--show-only-first',
            action='store_true',
            dest='show_only_first',
            help='Show only the first IP'
        )
        args = parser.parse_args()

        this_dir = pathlib.Path(__file__).parent.resolve()
        config_file = this_dir.joinpath('..', 'live', 'config.json').resolve()

        with open(config_file, 'r') as stream:
            data = json.load(stream)

        ips = [data['floating_ip']]

        if args.show_only_first:
            print(json.dumps(ips[0]))
            return

        print(json.dumps(ips))


def error(msg):
    print('... ERROR: {}'.format(msg), file=sys.stderr)
    exit(1)


if __name__ == '__main__':
    RemoteIPs()
