from datetime import datetime
import time

class ServiceInfo(object):
    __slots__ = ('id','name' ,'pbType', 'description', 'url', 'status', 'last_update', 'HeartbeatTime')
    def __init__(self, id, name, ip, port, pbType = None, description = None, last_update = None):
        self.id = id
        self.name = name
        self.url = (ip,port)
        self.pbType = pbType
        self.description = description
        self.HeartbeatTime = time.time()
        if last_update != None: self.last_update = datetime.fromtimestamp(last_update)
        else: self.last_update = datetime.now()
        self.status = 0
    
    def to_dict(self):
        return {
            'id': self.id,
            'name': self.name,
            'ip': self.url[0],
            'port': self.url[1],
            'pbType': self.pbType,
            'description': self.description,
            'last_update': time.mktime(self.last_update.timetuple)
        }
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

    def update_url(self, ip, port):
        if (ip,port) == self.url:
            return False
        else:
            self.url = (ip,port)
            self.last_update = datetime.now()
            return True

    def update_all(self, ip, port, pbType = None, description = None):
        self.url = (ip,port)
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
        self.status += 1
