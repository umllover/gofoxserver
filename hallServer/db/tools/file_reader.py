#!/usr/bin/env python
#-* coding:UTF-8 -*

import codecs


class FileReader(object):

    def __init__(self, path):
        self.path = path
        self.data = []
        self.load()

    def load(self):
        with codecs.open(self.path, 'r', 'utf-8') as f:
            lines = f.readlines()
        headers = [line.strip('\n') for line in lines[0].split(',')]
        print 'headers', headers
        lines = [line for line in lines[1:] if not line.startswith('#')]
        self.data = [dict(zip(headers, line.split(','))) for line in lines]
        print 'load %s data  %s' % (self.path, self.data)
