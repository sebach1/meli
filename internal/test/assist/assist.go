/*
Package assist contains testing helpers but for init()
-not for testing init func, but for using init on testing-

Because of that, that type of funcs couldnt contain testing.T arg on any function,
and then the errs are notified by log.Fatal (due every bad init would causes bad processing).
*/
package assist
