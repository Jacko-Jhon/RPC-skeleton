import hashlib
import pandas as pd
import numpy as np
from ServiceInfo import ServiceInfo

class Register(object):
    _name_to_id = {}
    _services = {}
    _clients = {}
    # generate a new id by MD5
    def id_generate(self):
        key = len(Register._services)+100000000 
        id = hashlib.md5(str(key).encode('utf-8')).hexdigest()
        return id

    # dump all service info to file
    def dump(self):
        data = [it.to_dict() for it in Register._services.values()]
        df = pd.DataFrame(data)
        df.to_csv('../services.csv', index=False, header=True)
    # load all service info from file
    def load(self):
        df = pd.read_csv('../services.csv')
        data = np.array(df)
        for it in data:
            if it[1] not in Register._name_to_id: Register._name_to_id[it[1]] = [it[0]]
            else: Register._name_to_id[it[1]].append(it[0])
            Register._services[it[0]] = ServiceInfo(it[0], it[1], it[2], int(it[3]), it[4], it[5], float(it[6]))

    def log(self):
        pass

    def token_authorize(self, token):
        pass

