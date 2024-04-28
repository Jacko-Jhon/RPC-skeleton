import time
import hashlib
from random import randint

class Register(object):
    _name_to_id = {}
    _services = {}
    _clients = {}
    # generate a new id by MD5
    def id_generate(self):
        key = str(time.now()) + str(randint(10000, 99999))
        id = hashlib.md5(key.encode('utf-8')).hexdigest()
        while id in Register._services: 
            key = str(time.now()) + str(randint(10000, 99999))
            id = hashlib.md5(key.encode('utf-8')).hexdigest()
        return id

    # dump all service info to file
    def dump(self):
        pass

    # load all service info from file
    def load(self):
        pass

    def log(self):
        pass

    def token_authorize(self, token):
        pass

