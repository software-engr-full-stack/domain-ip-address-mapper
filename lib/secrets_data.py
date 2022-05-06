#!/usr/bin/env python3

import pathlib
import sys
import yaml

from pylib.error import error


class SecretsData(object):
    def __init__(self):
        this_dir = pathlib.Path(__file__).parent.resolve()

        self.__secrets_file = this_dir.joinpath('..', 'secrets', 'secrets.yml').resolve()
        self.__secrets_dir = self.__secrets_file.parent

        self.__data = self.__get_data()

    def dbpw(self):
        database_template = self.get_value('database_template')

        return database_template['default']['password']

    def get_value(self, key):
        if key not in self.__data:
            error("key '{}' not found in data".format(key))

        item = self.__data[key]
        if self.__is_relative_path(item):
            value = self.__secrets_dir.joinpath(item['value'])
            if self.__do_resolve(item):
                return value.resolve()

            return value

        return item

    def __get_data(self):
        with open(self.__secrets_file, 'r') as stream:
            try:
                data = yaml.safe_load(stream)
            except yaml.YAMLError as exc:
                error('YAML error: {}'.format(exc))

        return data

    def __is_relative_path(self, item):
        irp_key = 'is_relative_path'
        if irp_key in item:
            return item[irp_key]

    def __do_resolve(self, item):
        dr_key = 'do_resolve'
        if dr_key in item:
            return item[dr_key]


if __name__ == '__main__':
    if len(sys.argv) < 2:
        error('must pass key')
    print(SecretsData().get_value(sys.argv[1]))
