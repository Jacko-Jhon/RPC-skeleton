success = "SUCCESS"
failure = "FAILURE"

class MessageToService(object):
    __slots__ = ('status','id','info')
    def __init__(self, status, info, id = ''):
        self.status = status # (SCCESS, FAILURE)
        self.id = id
        self.info = info # a specific message like: (ip, port, prbuff, description)