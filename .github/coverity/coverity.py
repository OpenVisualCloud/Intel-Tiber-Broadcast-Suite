import subprocess
import concurrent.futures

exec = lambda cmd: subprocess.run(cmd, shell=True, check=True)

def run_coverity():
  with concurrent.futures.ThreadPoolExecutor() as executor:
    builds = ("nmos", "nmos-cpp", "grpc", "launcher")
    futures = [executor.submit(exec, f"./cov-build.sh {name} | tee {name}.log") for name in builds ]
    for future in concurrent.futures.as_completed(futures):
      future.result()


if __name__ == "__main__": 
  run_coverity() 