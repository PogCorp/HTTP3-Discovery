#include <lsquic.h>
#include <stdio.h>
#include <stdlib.h>

void logger(const char *tag, const char *message) {
  time_t now;
  time(&now);
  printf("%s [%s]: %s\n", ctime(&now), tag, message);
}

int main(int argc, char *argv[]) {
  if (lsquic_global_init(LSQUIC_GLOBAL_CLIENT) != 0) {
    logger("ERROR", "unable to initialize lsquic");
    exit(EXIT_FAILURE);
  }
  // lsquic_packets_out_f pogcorp_out = ;
  //  lsquic_engine_api engine_api = {.ea_packets_out}

  lsquic_global_cleanup();
  return EXIT_SUCCESS;
}
