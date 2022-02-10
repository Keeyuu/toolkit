import pymongo

class Mongo:
    def __init__(self, url: str, database: str):
        self.url = url
        self.database = database
        self.client = pymongo.MongoClient(url)

    def get_collection(self, collection: str):
        return self.client[self.database][collection]
