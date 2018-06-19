package atmi
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

}

*/
import "C"
import "unsafe"
import "fmt"
import "runtime"

/*
 * SUCCEED/FAIL flags
 */
const (
	SUCCEED = 0
	FAIL    = -1
)

