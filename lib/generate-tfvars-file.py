#!/usr/bin/env python3

import json
import sys

from secrets_data import SecretsData


class Exec(object):
    def __init__(self, output_file):
        sec = SecretsData()
        tfvars = sec.get_value('tfvars')
        tfv_secrets_key = 'secrets'
        tfvars[tfv_secrets_key]['database']['password'] = sec.dbpw()

        with open(output_file, 'w') as stream:
            stream.write(json.dumps(tfvars))


Exec(output_file=sys.argv[1])
