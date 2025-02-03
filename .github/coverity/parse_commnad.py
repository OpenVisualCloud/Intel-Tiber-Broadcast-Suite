from dockerfile_parse import DockerfileParser
import os
import re

GEN_BUILD_DIRS = ".github/coverity"
CMD_OFFSET = 10
builds = []
invalid_dirs = ["/tmp/svt-av1/Build", "/tmp/vulkan-headers", "/tmp/dpdk", "/tmp"]

dp = DockerfileParser(path="Dockerfile")

WORKDIR_CMD_LIST = list(filter(lambda s: s["instruction"] == "WORKDIR", dp.structure))
RUN_CMD_LIST = list(filter(lambda s: s["instruction"] == "RUN", dp.structure))
for workdir in WORKDIR_CMD_LIST:
    for cmd in RUN_CMD_LIST:
        # find first RUN command after WORKDIR
        if cmd["startline"] > workdir["endline"] and cmd["startline"] < (
            workdir["endline"] + CMD_OFFSET
        ):
            if "BUILD" in cmd["content"] and workdir["value"] not in invalid_dirs:
                builds.append({"dir": workdir["value"], "cmd": cmd["value"]})


for build in builds:
    if "dpdk" not in build["dir"]:
        # Remove build non relevant code from command
        filtered_cmd = modified_string = re.sub(
            r"&&\s*(git|curl).*?&&", "&&", build["cmd"]
        )

        if "make" in build["cmd"]:
            make_pattern = r"\bmake\b"
            # Replace 'make' with 'make -B'
            filtered_cmd = re.sub(make_pattern, "make -B", filtered_cmd)

        repo = build["dir"].split("/")[2]  # ["", "tmp", "repo",...]
        filename = f"{GEN_BUILD_DIRS}/{repo}_build_cmd.sh"
        with open(filename, "w") as script_file:
            build_subdirs = "{build,Build,BUILD,sdk/out}"
            script_file.write("#!/bin/bash\n")
            script_file.write(f"cd {build['dir']}\n")
            script_file.write(f"rm -rf {build['dir']}/{build_subdirs}\n")
            script_file.write(f"{filtered_cmd}\n\n")
        os.chmod(filename, 0o755)


grpc_script = f"{GEN_BUILD_DIRS}/gRPC_build_cmd.sh"
pod_launcher_script = f"{GEN_BUILD_DIRS}/pod_launcher_build_cmd.sh"

with open(grpc_script, "w") as script_file:
    script_file.write("#!/bin/bash\n")
    script_file.write('echo "**** BUILD gRPC ****"\n')
    script_file.write("cd gRPC\n")
    script_file.write("sed -i '$s/make/make -B/' compile.sh\n")
    script_file.write("./compile.sh\n")
os.chmod(grpc_script, 0o755)

with open(pod_launcher_script, "w") as script_file:
    script_file.write("#!/bin/bash\n")
    script_file.write('echo "**** BUILD pod Launcher ****"\n')
    script_file.write("cd launcher/cmd/\n")
    script_file.write("go build main.go\n")
os.chmod(pod_launcher_script, 0o755)

print("coverity build scripts generated successfully.")