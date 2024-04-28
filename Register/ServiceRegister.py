import time
import hashlib
from random import randint

import Message
from ServiceInfo import ServiceInfo

MessageToService = Message.MessageToService

class ServiceRegister:
    __name_to_id = {}
    __services = {}
    __healthTime = 3

    def setHealthTime(self, seconds):
        ServiceRegister.__healthTime = seconds
    
    def getHealthTime(self):
        return ServiceRegister.__healthTime
    # generate a new id by MD5
    def id_generate(self):
        key = str(time.now()) + str(randint(10000, 99999))
        id = hashlib.md5(key.encode('utf-8')).hexdigest()
        while id in ServiceRegister.__services: 
            key = str(time.now()) + str(randint(10000, 99999))
            id = hashlib.md5(key.encode('utf-8')).hexdigest()
        return id
    
    # regist a new service
    def regist(self, name, ip, port, pbType = None, description = None):
        if name in ServiceRegister.__name_to_id:
            info = 'Service name already exists'
            msg = MessageToService(Message.failure, info)
            return msg
        else:
            # this may be a problem that the id is not unique because of non-lock
            id = self.id_generate()
            ServiceRegister.__name_to_id[name] = [id]
            ServiceRegister.__services[id] = ServiceInfo(id, name, ip, port, pbType, description)
            info = 'Service regist success'
            msg = MessageToService(Message.success, info, id)
            return msg
    
    # update all info of service by id
    # Not recommended 
    def update_all(self, id, ip, port, pbType = None, description = None):
        if id in ServiceRegister.__services:
            ServiceRegister.__services[id].update_all(ip, port, pbType, description)
            info = 'Service update success'
            msg = MessageToService(Message.success, info, id)
            return msg
        else:
            info = 'Service id not found'
            msg =MessageToService(Message.failure, info)
        
    # update one info of service by id
    # @address is a tuple like ('127.0.0.1', 5000)
    def update_one(self, id, address = None, pbType = None, description = None):
        if id in ServiceRegister.__services:
            if address != None:
                ServiceRegister.__services[id].update_address(address)
                info = 'Service update success: your service address is '+address[0] + ':' + str(address[1])
                msg = MessageToService(Message.success, info)
            elif pbType != None:
                ServiceRegister.__services[id].update_pbType(pbType)
                info = 'Service update success: your service pbType is '+str(pbType)
                msg = MessageToService(Message.success, info)
            elif description !=None:
                ServiceRegister.__services[id].update_description(description)
                info = 'Service update success: your service description is '+description
                msg = MessageToService(Message.success, info)
            else :
                info = 'Service update failed: no update'
                msg = MessageToService(Message.failure, info)
                return msg
        else:
            info = 'Service id not found'
            msg = MessageToService(Message.failure, info)
            return msg
    
    # delete service by id
    def unregist(self, name, id):
        if name in ServiceRegister.__name_to_id.keys() and id in ServiceRegister.__name_to_id[name]:
            del ServiceRegister.__services[id]
            if len(ServiceRegister.__name_to_id[name]) == 1:
                del ServiceRegister.__name_to_id[name]
            else:
                ServiceRegister.__name_to_id[name].remove(id)
            info = 'Service unregist success'
            msg = MessageToService(Message.success, info)
            return msg
        else:
            info = 'Service not found'
            msg = MessageToService(Message.failure, info)
            return msg

    # delete all service by name
    def unregist_all(self, name):
        for id in ServiceRegister.__name_to_id[name]:
            del ServiceRegister.__services[id]
        del ServiceRegister.__name_to_id[name]
        info = 'Service unregist success'
        msg = MessageToService(Message.success, info)
        return msg

    # check the health of service, if not healthy, change the status to 0
    def check_health(self):
        T = time.time()
        for value in ServiceRegister.__services.values():
            if value.HeartbeatTime + ServiceRegister.__healthTime < T:
                value.status = 0
    # service heartbeat by id
    def service_heartbeat(self, id):
        if id in ServiceRegister.__services:
            ServiceRegister.__services[id].heartbeat()
            info = 'Beat'
            msg = MessageToService(Message.success, info, id)
            return msg
        else:
            info = 'Service id not found'
            msg = MessageToService(Message.failure, info)
            return msg

    # add a new service by name
    def add_server_by_name(self, name, authId, ip, port, pbType = None, description = None):
        if name in ServiceRegister.__name_to_id.keys() and authId in ServiceRegister.__name_to_id[name]:
            id = self.id_generate()
            ServiceRegister.__name_to_id[name].append(id)
            ServiceRegister.__services[id] = ServiceInfo(id, ip, port, pbType, description)
            info = 'Service regist success'
            msg = MessageToService(Message.success, info, id)
            return msg
        else:
            info = 'Service not found or auth id not found'
            msg = MessageToService(Message.failure, info)
            return msg
    
    # find service by name
    def find_service(self, name):
        services_id = ServiceRegister.__name_to_id[name]
        services_list = []
        for id in services_id:
            if ServiceRegister.__services[id].status == 1:
                services_list.append(id)
        return services_list

    # dump all service info to file
    def dump(self):
        pass

    # load all service info from file
    def load(self):
        pass

    # def server_communicate(self, msg):
    #     pass

    # def client_communicate(self, msg):
    #     pass

    # def server_msg_handler(self, msg):
    #     pass

    # def client_msg_handler(self, msg):
    #     pass
