from Register import Register

class ServiceProvider(Register):

    
    # find service by name
    def find_service(self, name):
        services_list = []
        if name not in Register._name_to_id: return services_list
        services_id = Register._name_to_id[name]
        for id in services_id:
            if Register._services[id].status >= 1:
                services_list.append(id)
        services_list = [Register._services[id] for id in services_id if Register._services[id].status >= 1]
        return services_list