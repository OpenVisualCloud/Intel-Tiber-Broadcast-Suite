import api
import os
import sys


#TODO: error handling, branch/version required

branch  = sys.argv[1]
version = sys.argv[2]

BASE_URL = os.environ.get("COVERITY_BASE_URL")
TOKEN =  os.environ.get("COVERITY_TOKEN")
PROJECT_NAME = os.environ.get("COVERITY_PROJECT_NAME")
USER = os.environ.get("COVERITY_USER")

def main():
    base_config = {"base_url": BASE_URL,
              "project_name": PROJECT_NAME,
              "password": TOKEN,
              "user": USER,
              "stream": PROJECT_NAME,
    }
    snapshot = api.find_snapshot(base_config ,branch , version )
    print(snapshot)
    
if __name__ == '__main__':
    main()