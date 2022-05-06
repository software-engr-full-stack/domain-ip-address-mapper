#!/usr/bin/env python3

import argparse

from pylib.tls import TLS


class TLSExec(object):
    def __init__(self):
        parser = argparse.ArgumentParser(
            description='Find path of where TLS files are stored for first domain listed in given domain config file'
        )
        parser.add_argument(
            '-c', '--domain-config-file',
            dest='domain_config_file',
            help='path to the domain config file'
        )

        parser.add_argument(
            '-d', '--dir-to-list-of-domains',
            dest='domains_dir',
            help='path to the dir containing list of directories containing TLS files',
            required=True
        )
        args = parser.parse_args()

        tls = TLS(domains_dir=args.domains_dir, domain_config_file=args.domain_config_file)

        print(tls.data)


TLSExec()
