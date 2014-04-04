#! /usr/bin/python

import sys
import imp
import json


if len(sys.argv) != 2:
    print "convert python data file to json file"
    print 'Usage:', sys.argv[0], "prob_emit|prob_start|prob_trans|prov_prev"
    sys.exit(1)

#print imp.find_module(sys.argv[1])
m = imp.load_module(sys.argv[1], *imp.find_module(sys.argv[1]))
print json.dumps(m.P)
