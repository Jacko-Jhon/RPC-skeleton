from datetime import datetime
import time

class ServiceInfo(object):
    __slots__ = ('id', 'pbType', 'description', 'address', 'status', 'last_update', 'HeartbeatTime')
    def __init__(self, id, ip, port, pbType = None, description = None):
        self.id = id
        self.address = (ip,port)
        self.pbType = pbType
        self.description = description
        self.HeartbeatTime = time.time()
        self.last_update = datetime.now()
        self.status = 1

    # def delete_server(self, ip, port):
    #     if (ip,port) in self.address:
    #         self.address.remove((ip,port))
    #         self.last_update = datetime.now()
    #         return True
    #     return False

    # def add_new_server(self, ip, port):
    #     if (ip,port) in self.address:
    #         return False
    #     self.address.add((ip,port))
    #     self.last_update = datetime.now()
    #     return True

    def update_address(self, ip, port):
        if (ip,port) == self.address:
            return False
        else:
            self.address = (ip,port)
            self.last_update = datetime.now()
            return True

    def update_all(self, ip, port, pbType = None, description = None):
        self.address = (ip,port)
        self.pbType = pbType
        self.description = description
        self.last_update = datetime.now()

    def update_pbType(self, pbType):
        self.pbType = pbType
        self.last_update = datetime.now()

    def update_description(self, description):
        self.description = description
        self.last_update = datetime.now()

    def heartbeat(self):
        self.HeartbeatTime = time.time()
        self.status = 1
