package exbench

/*
 * @brief Test util package
 *
 * @file testutil.go
 *
 */

/*
#cgo pkg-config: atmisrvinteg

#include <ndebug.h>
#include <string.h>
#include <stdlib.h>
#include <nstdutil.h>


*/
import "C"
/*
import "unsafe"
import "fmt"
import "runtime"
*/

/**
 * SUCCEED/FAIL flags
 */
const (
        /** succeed status */
	SUCCEED = 0
        /** fail status */
	FAIL    = -1
)

/**
 * Write benchmark results
 * @param[in] message size at benchmark (bytes)
 * @param[in] callspersec number of calls per second
 * @return 0 - succeed; -1 - fail
 */
func ndrx_bench_write_stats(msgsize float64, callspersec float64) int {
        return int(C.ndrx_bench_write_stats(C.double(msgsize), C.double(callspersec)))
}
