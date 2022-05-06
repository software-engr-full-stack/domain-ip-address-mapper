#!/usr/bin/env python3

import argparse

from pylib.aws.rds_by_tag_name import RDSByTagName
from pylib.db_config_file_builder import DbConfigFileBuilder

from secrets_data import SecretsData


class Exec(object):
    def __init__(self):
        parser = argparse.ArgumentParser(description='Build secret database config YAML file')
        parser.add_argument(
            '-t', '--tag-name-value',
            dest='tag_name_value',
            required=True
        )

        parser.add_argument(
            '-o', '--output-file',
            dest='output_file',
            required=True
        )

        args = parser.parse_args()

        rds_by_tag_name = RDSByTagName(args.tag_name_value)
        dbpw = SecretsData().dbpw()

        endpoint = rds_by_tag_name.data['Endpoint']
        DbConfigFileBuilder(
            dbpw,
            endpoint['Address'],
            endpoint['Port'],
            args.output_file
        )


Exec()
