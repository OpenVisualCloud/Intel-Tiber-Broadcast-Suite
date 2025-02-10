import os
class Config():
   def __init__(self):
      self.base_url = os.environ.get("COVERITY_URL")
      self.token =  os.environ.get("COVERITY_TOKEN")
      self.project_name = os.environ.get("COVERITY_PROJECT_NAME")