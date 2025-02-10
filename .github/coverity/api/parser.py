import pandas as pd


def issues_to_pandas(raw_response_data)-> pd.DataFrame:

  df = pd.DataFrame()
  converted_response = []

  for row in raw_response_data["rows"]:
    converted_dict ={ item['key']: item['value']  for item in row }
    converted_response.append(converted_dict)
  df_stream = pd.DataFrame(converted_response)
  df = pd.concat([df, df_stream], ignore_index=True)

  impact_order =["High", "Medium", "Low"]
  df["displayImpact"] = pd.Categorical(df["displayImpact"], categories=impact_order)
  df = df.sort_values(by ="displayImpact")

  return df