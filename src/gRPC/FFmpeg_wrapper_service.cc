#include "CmdPassImpl.h"

CmdPassImpl manager;

void wrapperSignalHandler(int signal) {
  std::cout << "[*] wrapperSignalHandler called SIG : " << signal << std::endl;
  manager.Shutdown();
}

int main(int argc, char *argv[]) {
  if (argc != 3) {
    std::cout
        << "[*] FFmpeg wrapper service takes 2 arguments: interface and port"
        << std::endl;
    return 1;
  }

  std::stringstream ss;
  ss << argv[1] << ":" << argv[2];

  std::signal(SIGTERM, wrapperSignalHandler);
  std::signal(SIGINT, wrapperSignalHandler);

  manager.Run(ss.str());

  std::cout << "[*] Service exited gracefully" << std::endl;

  return 0;
}
