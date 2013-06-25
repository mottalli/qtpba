#!/usr/bin/python
# -*- coding: utf-8 -*-

import sqlite3

class MapReducer(object):
    def __init__(self):
        self.valuesMap = {}
        self.output = {}

    def map(self, element):
        # yield (key, value)
        raise NotImplemented()

    def reduce(self, key, elements):
        # yield output
        raise NotImplemented()

    def getElements(self):
        raise NotImplemented()

    def processOutput(self):
        raise NotImplemented()

    def run(self):
        for element in self.getElements():
            for (key, value) in self.map(element):
                if not key in self.valuesMap:
                    self.valuesMap[key] = []
                self.valuesMap[key].append(value)

        self.output = {}
        for key in self.valuesMap:
            self.output[key] = []
            for result in self.reduce(key, self.valuesMap[key]):
                self.output[key].append(result)

        self.processOutput()

class DatabaseMapReducer(MapReducer):
    def __init__(self, db, query):
        super(DatabaseMapReducer, self).__init__()
        self.db = db
        self.query = query

    def getElements(self):
        return self.db.execute(self.query)

class UsercountMapReduce(DatabaseMapReducer):
    """def __init__(self, db, query):
        super(UsercountMapReduce, self).__init__(db, query)"""

    def map(self, row):
        yield (int(row[0]), 1)

    def reduce(self, user_id, counts):
        yield len(counts)

    def processOutput(self):
        print self.output


conn = sqlite3.connect("../db/qtpba.db")
mr = UsercountMapReduce(conn, "SELECT user_id FROM tweet ORDER BY id ASC")
mr.run()