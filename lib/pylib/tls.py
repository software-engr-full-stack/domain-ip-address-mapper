import pathlib
import yaml
import os
import json

from pylib.error import error


class TLS(object):
    def __init__(self, domains_dir, domain_config_file=None):
        super(TLS, self).__init__()

        no_domain = self.__build_json({'domain': ''})

        if domain_config_file is None:
            self.data = no_domain
            return

        domains = self.__get_domains(domain_config_file)

        if domains is None:
            self.data = no_domain
            return

        # Check the deploy config file. The root domain, if any, should be first.
        # This is because it's the name Let's Encrypt will use to build the TLS files.
        main_domain = domains[0]

        match = self.__find_match_for(main_domain, inside=domains_dir)

        self.domain = main_domain
        self.cert_files_dir = match
        self.data = self.__build_json({'domain': main_domain, 'cert_files_dir': match})

    def __get_domains(self, domain_config_file):
        with open(domain_config_file, 'r') as stream:
            try:
                data = yaml.safe_load(stream)
            except yaml.YAMLError as exc:
                error('YAML error: {}'.format(exc))

        domains_key = 'domains'
        if domains_key in data:
            return [item.strip() for item in data[domains_key]]

    def __find_match_for(self, domain, inside):
        ftable = {
            filename: pathlib.Path(inside).joinpath(filename).resolve()
            for filename in os.listdir(inside)
        }

        if domain in ftable:
            return ftable[domain].as_posix()

        return ''

    def __build_json(self, result):
        return json.dumps({'result': result})
