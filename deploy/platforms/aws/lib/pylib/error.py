import sys


def error(msg):
    print('... ERROR: {msg}'.format(msg=msg), file=sys.stderr)
    exit(1)
