import api
import pandas as pd
import os
BASE_URL = os.environ.get("COVERITY_BASE_URL")
TOKEN =  os.environ.get("COVERITY_TOKEN")
PROJECT_NAME = os.environ.get("COVERITY_PROJECT_NAME")
USER = os.environ.get("COVERITY_USER")

launcher_search_querry={
        "base_url": BASE_URL,
        "project_name": PROJECT_NAME,
        "password": TOKEN,
        "user": USER,
        "stream": PROJECT_NAME,
        "columns": ["cid", "classification", "severity", "displayType",
                    "displayImpact", "displayFile", "displayFunction", "displayCategory"],
        "snapshot": "last()",
        "extra_filters": [
          {
            "columnKey": "fileLanguage",
            "matchMode": "oneOrMoreMatch",
            "matchers": [
              {
                "key": "Go",
                "type": "keyMatcher"
              }
            ]
          }
        ]
}

def main():
    df = pd.DataFrame()

    raw_response_data = api.get_snapshot_issues(launcher_search_querry)
    converted_response = []
    for row in raw_response_data["rows"]:
      converted_dict ={ item['key']: item['value']  for item in row }
      converted_response.append(converted_dict)
    df_stream = pd.DataFrame(converted_response)
    df = pd.concat([df, df_stream], ignore_index=True)


    impact_order =["High", "Medium", "Low"]
    df["displayImpact"] = pd.Categorical(df["displayImpact"], categories=impact_order)

    df = df.sort_values(by ="displayImpact")
    df.to_csv("report.csv", index=False)


if __name__ == '__main__':
    main()