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
    
    # raw_response_data = api.get_snapshot_issues({
    #     "base_url": BASE_URL,
    #     "project_name": PROJECT_NAME,
    #     "password": TOKEN,
    #     "user": USER,
    #     "stream": PROJECT_NAME,
    #     "columns": ["cid", "classification", "severity", "displayType",
    #                 "displayImpact", "displayFile", "displayFunction", "displayCategory"],
    #     "snapshot": snapshot}
    #         })
    # df = api.issues_to_pandas(raw_response_data)
    # df.to_csv(f"grpc-report.csv", index=False)
if __name__ == '__main__':
    main()