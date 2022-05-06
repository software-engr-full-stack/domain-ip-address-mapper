#!/usr/bin/env python3

import pathlib
import json
import sys
import subprocess


class SSH(object):
    def __init__(self):
        remote_user = 'remote'
        this_dir = pathlib.Path(__file__).parent.resolve()
        config_file = this_dir.joinpath('..', '..', 'config.json').resolve()

        with open(config_file, 'r') as stream:
            data = json.load(stream)

        subprocess.run(['echo', '-ne', '\033[22;0t'])

        subprocess.run([
            'ssh', '-o', 'StrictHostKeyChecking=no', '{}@{}'.format(
                remote_user, data['floating_ip']
            )
        ])

        subprocess.run(['echo', '-ne', '\033[23;0t'])


def error(msg):
    print('... ERROR: {}'.format(msg), file=sys.stderr)
    exit(1)


if __name__ == '__main__':
    SSH()
