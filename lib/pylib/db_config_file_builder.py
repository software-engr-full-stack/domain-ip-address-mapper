import inspect


class DbConfigFileBuilder(object):
    def __init__(self, dbpw, host_address, port, output_file):
        data = inspect.cleandoc(f'''
            default: &default
              adapter: 'postgres'
              encoding: 'utf-8'
              host: '{host_address}'
              port: {port}
              user: 'postgres'
              password: '{dbpw}'
              pool: 50
              # max_open_conns: 0
              # search_path: '???'

            development:
              <<: *default
              name: 'demo_dev'

            test:
              <<: *default
              name: 'demo_test'

            production:
              <<: *default
              name: 'demo'
        ''')

        with open(output_file, 'w') as stream:
            stream.write(data)
