// from https://www.nayuki.io/page/fast-sha1-hash-implementation-in-x86-assembly
#include <stdint.h>

#define BLOCK_LEN 64  // In bytes
#define STATE_LEN 5  // In words

void sha1_sum(const uint8_t message[], size_t len, uint32_t hash[static STATE_LEN]);
void sha1_compress(const uint8_t block[static BLOCK_LEN], uint32_t state[static STATE_LEN]);
