#!/usr/bin/env python3

import pathlib
import argparse
import glob
import os
import subprocess

from pylib.tls import TLS


class ExtractFetchedTLSArchive(object):
    def __init__(self):
        super(ExtractFetchedTLSArchive, self).__init__()

        parser = argparse.ArgumentParser(description='Extract fetched TLS archive to main TLS dir')
        parser.add_argument(
            '-f', '--fetched-tls-files-dir',
            dest='fetched_tls_files_dir',
            required=True
        )

        parser.add_argument(
            '-d', '--dest-dir',
            dest='dest_dir',
            required=True
        )

        parser.add_argument(
            '-c', '--domain-config-file',
            dest='domain_config_file',
            required=True
        )

        args = parser.parse_args()

        tls = TLS(domains_dir=args.dest_dir, domain_config_file=args.domain_config_file)

        match = self.__find_matching_archive_file(tls.domain, args.fetched_tls_files_dir)

        if match and not self.__was_extracted(tls.domain, args.dest_dir):
            self.__extract(match, to=args.dest_dir, domain=tls.domain)

    def __find_matching_archive_file(self, domain, fetched_tls_files_dir):
        arch_name = '.'.join([domain, 'tar', 'xz'])

        for file in glob.glob('/'.join([fetched_tls_files_dir, '*.tar.xz'])):
            bname = pathlib.Path(file).name

            if bname == arch_name:
                return file

    def __was_extracted(self, domain, dest_dir):
        for file in os.listdir(dest_dir):
            if file == domain:
                return True

    def __extract(self, archive, to, domain):
        os.chdir(to)
        subprocess.run(['tar', 'xJvf', archive])
        subprocess.run(['mv', '--verbose', 'etc', domain])


ExtractFetchedTLSArchive()
