#!/usr/bin/python
# -*- coding: utf-8 -*-

import sqlite3
import time
import codecs
import re

def cleanWord(word):
    word = word.strip().lower()
    word = re.sub(r'[^A-Za-z0-9áéíóúñ]', '', word)
    return word

def validWord(word):
    global blacklist
    if len(word) < 3 or word in blacklist:
        return False

    if word.find('http') > -1:
        return False

    return True

conn = sqlite3.connect("../db/qtpba.db")

with codecs.open("../static/blacklist.txt", "r", encoding="UTF-8") as f:
    blacklist = f.readlines()

blacklist = map(lambda s: s.strip(), blacklist)
blacklist = set(blacklist)

wordCount = {}

for row in conn.execute("SELECT message FROM tweets ORDER BY rowid DESC LIMIT 1000"):
    message = row[0]
    words = message.split()
    for word in words:
        word = cleanWord(word)

        if not validWord(word):
            continue

        wordCount[word] = wordCount[word]+1 if word in wordCount else 1

import operator
sortedWords = sorted(wordCount.iteritems(), key=operator.itemgetter(1))
N = 20
topN = sortedWords[len(sortedWords)-N:]
print topN
        
