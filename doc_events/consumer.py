import boto3
import json
from datetime import datetime
import time
import requests
my_stream_name = 'TESTE'

ACCESS_KEY='123'
SECRET_KEY='abc'
kc  = boto3.client('kinesis',
                    endpoint_url="http://localhost:4568",
                    use_ssl=False,
                    aws_access_key_id=ACCESS_KEY,
                    region_name='us-east-1',

                    aws_secret_access_key=SECRET_KEY)
 

response = kc.describe_stream(StreamName=my_stream_name)

my_shard_id = response['StreamDescription']['Shards'][0]['ShardId']

shard_iterator = kc.get_shard_iterator(StreamName=my_stream_name,
                                                      ShardId=my_shard_id,
                                                      ShardIteratorType='LATEST')

my_shard_iterator = shard_iterator['ShardIterator']

record_response = kc.get_records(ShardIterator=my_shard_iterator,
                                              Limit=2)

while 'NextShardIterator' in record_response:
    record_response = kc.get_records(ShardIterator=record_response['NextShardIterator'],
                                                  Limit=2)
    if len(record_response["Records"]) >0:
         for i in record_response['Records']:
                data= i['Data']
                d = json.loads(data)
                print d 
                r = requests.post("http://localhost:8010/SF", data=d)
         #record =record_response["Records"]
          
    # wait for 5 seconds
    time.sleep(5)