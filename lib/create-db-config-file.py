#!/usr/bin/env python3

import argparse
import yaml
import json

from pylib.error import error


class Exec(object):
    def __init__(self):
        parser = argparse.ArgumentParser(description='Build secret database config YAML file')
        parser.add_argument(
            '-i', '--input-config-file',
            dest='input_config_file',
            required=True
        )

        parser.add_argument(
            '-m', '--modifications-to-default',
            dest='modifications',
            help='JSON string of hash of modifications to default block of original input config file',
            required=True
        )

        parser.add_argument(
            '-o', '--output-config-file',
            dest='output_config_file',
            required=True
        )

        args = parser.parse_args()

        with open(args.input_config_file, 'r') as stream:
            try:
                data = yaml.safe_load(stream)
            except yaml.YAMLError as exc:
                error('YAML error: {}'.format(exc))

        if args.modifications.strip() == '':
            error('modifications JSON arg should not be empty')

        modifications = json.loads(args.modifications)
        for key in modifications:
            data['default'][key] = modifications[key]
            # Copy to all environments because it doesn't work if we don't
            data['development'][key] = modifications[key]
            data['production'][key] = modifications[key]
            data['test'][key] = modifications[key]

        with open(args.output_config_file, 'w') as stream:
            stream.write(yaml.dump(data))


Exec()
