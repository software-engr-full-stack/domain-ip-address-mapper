#!/usr/bin/env python3

import sys
import subprocess

from pylib.error import error


class RemoteIP(object):
    def __init__(self):
        if len(sys.argv) < 2:
            error('must provide tag name value')

        tag_name_value = sys.argv[1]

        result = subprocess.run([
            'aws', 'ec2', 'describe-addresses',
            '--filters', 'Name=tag:Name,Values={}'.format(tag_name_value),
            '--query', 'Addresses[0].PublicIp',
            '--output', 'text'
        ], stdout=subprocess.PIPE)

        ip = result.stdout.decode('utf-8').strip()
        if ip == 'None':
            error("no IP found for query using tag name value '{}'".format(tag_name_value))

        print(ip)


RemoteIP()
