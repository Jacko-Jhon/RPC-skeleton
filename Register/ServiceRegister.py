import time

import Message
from ServiceInfo import ServiceInfo
from Register import Register

MessageToService = Message.MessageToService

class ServiceRegister(Register):
    __healthTime = 30

    def setHealthTime(self, seconds):
        ServiceRegister.__healthTime = seconds
    
    def getHealthTime(self):
        return ServiceRegister.__healthTime
    
    # regist a new service
    def regist(self, name, ip, port, pbType = None, description = None):
        if name in Register._name_to_id:
            info = 'Service name already exists'
            msg = MessageToService(Message.failure, info)
            return msg
        else:
            # this may be a problem that the id is not unique because of non-lock
            id = self.id_generate()
            Register._name_to_id[name] = [id]
            Register._services[id] = ServiceInfo(id, name, ip, port, pbType, description)
            info = 'Service regist success'
            msg = MessageToService(Message.success, info, id)
            return msg
    
    # update all info of service by id
    # Not recommended 
    def update_all(self, id, ip, port, pbType = None, description = None):
        if id in Register._services:
            Register._services[id].update_all(ip, port, pbType, description)
            info = 'Service update success'
            msg = MessageToService(Message.success, info, id)
            return msg
        else:
            info = 'Service id not found'
            msg =MessageToService(Message.failure, info)
        
    # update one info of service by id
    # @url is a tuple like ('127.0.0.1', 5000)
    def update_one(self, id, url: tuple = None, pbType = None, description = None):
        if id in Register._services:
            if url != None:
                Register._services[id].update_address(url)
                info = 'Service update success: your service address is '+url[0] + ':' + str(url[1])
                msg = MessageToService(Message.success, info)
            elif pbType != None:
                Register._services[id].update_pbType(pbType)
                info = 'Service update success: your service pbType is '+str(pbType)
                msg = MessageToService(Message.success, info)
            elif description !=None:
                Register._services[id].update_description(description)
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
        if name in Register._name_to_id.keys() and id in Register._name_to_id[name]:
            del Register._services[id]
            if len(Register._name_to_id[name]) == 1:
                del Register._name_to_id[name]
            else:
                Register._name_to_id[name].remove(id)
            info = 'Service unregist success'
            msg = MessageToService(Message.success, info)
            return msg
        else:
            info = 'Service not found'
            msg = MessageToService(Message.failure, info)
            return msg

    # delete all service by name
    def unregist_all(self, name):
        for id in Register._name_to_id[name]:
            del Register._services[id]
        del Register._name_to_id[name]
        info = 'Service unregist success'
        msg = MessageToService(Message.success, info)
        return msg

    # check the health of service, if not healthy, change the status to 0
    def check_health(self):
        T = time.time()
        for value in Register._services.values():
            if value > 0 and value.HeartbeatTime + ServiceRegister.__healthTime < T:
                value.status = 0
            if value.status > 1:
                value.status -= 1
    # service heartbeat by id
    def service_heartbeat(self, id):
        if id in Register._services:
            Register._services[id].heartbeat()
            info = 'Beat'
            msg = MessageToService(Message.success, info, id)
            return msg
        else:
            info = 'Service id not found'
            msg = MessageToService(Message.failure, info)
            return msg

    # add a new service by name
    def add_server_by_name(self, name, authId, ip, port, pbType = None, description = None):
        if name in Register._name_to_id.keys() and authId in Register._name_to_id[name]:
            id = self.id_generate()
            Register._name_to_id[name].append(id)
            Register._services[id] = ServiceInfo(id, ip, port, pbType, description)
            info = 'Service regist success'
            msg = MessageToService(Message.success, info, id)
            return msg
        else:
            info = 'Service not found or auth id not found'
            msg = MessageToService(Message.failure, info)
            return msg