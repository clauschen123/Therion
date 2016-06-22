package main

//	"io"
//	"log"
//	"net"
//	"os"

/*
func Connect(addr string) error {

	go Run()

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
		return err
	}

	MakeSession(conn)

	//	done := make(chan struct{})
	//	go func() {
	//		//recv
	//		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
	//		log.Println("done")
	//		done <- struct{}{} // signal the main goroutine
	//	}()

	//	mustCopy(conn, os.Stdin)
	//	conn.Close()
	//	<-done // wait for background goroutine to finish
	return err
}
*/

//send
//func mustCopy(dst io.Writer, src io.Reader) {
//	if _, err := io.Copy(dst, src); err != nil {
//		log.Fatal(err)
//	}
//}
