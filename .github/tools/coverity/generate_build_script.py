from dockerfile_parse import DockerfileParser
import json
import os
import re
OUTPUT_FILE="coverity_build.sh"
CMD_OFFSET=10
builds= []
invalid_dirs=[
    "/tmp/svt-av1/Build",
    "/tmp/vulkan-headers",
    "/tmp/dpdk",
    "/tmp"
]

dp = DockerfileParser(path="Dockerfile")

WORKDIR_CMD_LIST = list(filter(lambda s: s['instruction'] == "WORKDIR", dp.structure))
RUN_CMD_LIST= list(filter(lambda s: s['instruction'] == "RUN", dp.structure))
for workdir in WORKDIR_CMD_LIST:
  for cmd in RUN_CMD_LIST:
    # find nearest next RUN command
    if cmd['startline'] > workdir["endline"] and cmd['startline'] < (workdir["endline"] + CMD_OFFSET):
      if  "BUILD" in cmd["content"] and workdir["value"] not in invalid_dirs:
        builds.append({
          "dir": workdir["value"],
          "cmd": cmd["value"]
        })

with open(OUTPUT_FILE, "w") as script_file:

    script_file.write("#!/bin/bash\n\n")

    for build in builds:
      if "dpdk" not in build['dir']:
        script_file.write(f"cd {build['dir']}\n")
        # remove build dirs to force rebuild
        # Regular expression to match content between && and && that starts with git
        filter_pattern = r'&&\s*(git|curl).*?&&'

        # Replace the matched content with '&&'
        filtered_cmd = modified_string = re.sub(filter_pattern, '&&', build['cmd'])

        build_subdirs="{build,Build,BUILD}"
        script_file.write(f"rm -rf {build['dir']}/{build_subdirs}\n")
        script_file.write(f"{filtered_cmd}\n\n")
print("Bash script generated successfully.")