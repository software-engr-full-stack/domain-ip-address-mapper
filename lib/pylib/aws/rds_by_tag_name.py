import subprocess
import json

from pylib.error import error


class RDSByTagName(object):
    def __init__(self, tag_name_value):
        arn = self.__get_arn(tag_name_value)
        rds_match = self.__find_rds_match(arn)
        self.data = rds_match
        self.json = json.dumps(rds_match)

    def __get_arn(self, tag_name_value):
        result = subprocess.run([
            'aws',
            'resourcegroupstaggingapi',
            'get-resources',
            '--resource-type-filters',
            'rds',
            '--tag-filters',
            'Key=Name,Values={tag_name_value}'.format(tag_name_value=tag_name_value)
        ], stdout=subprocess.PIPE)

        rtml = json.loads(result.stdout)['ResourceTagMappingList']
        if len(rtml) != 1:
            error(
                "resource tag mapping of tag name '{tag_name_value}' list count '{count}' should be == 1".
                format(tag_name_value=tag_name_value, count=len(rtml))
            )

        return rtml[0]['ResourceARN']

    def __find_rds_match(self, arn):
        result = subprocess.run(['aws', 'rds', 'describe-db-instances'], stdout=subprocess.PIPE)
        db_instances = json.loads(result.stdout)['DBInstances']

        db_instances_by_arn = {item['DBInstanceArn']: item for item in db_instances}

        if arn not in db_instances_by_arn:
            error("ARN '{arn}' not found in list of database instances".format(arn=arn))

        return db_instances_by_arn[arn]
