from Register import Register

class ServiceProvider(Register):

    
    # find service by name
    def find_service(self, name):
        services_id = Register._name_to_id[name]
        services_list = []
        for id in services_id:
            if Register._services[id].status >= 1:
                services_list.append(id)
        return services_list