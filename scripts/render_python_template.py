#!/usr/bin/python3

import sys
from jinja2 import Environment, BaseLoader

#print("Command: ")
#print(sys.argv)
branch = sys.argv[2]
version = sys.argv[3]
template = sys.argv[1]
#Example template:
#'{%- if branch in [ "p10", "c10f1", "c10f2"] -%} php8.2 {%- else -%} php8.3 {%- endif -%}'
 
print(Environment(loader=BaseLoader()).from_string(template).render(branch=branch,version=version))