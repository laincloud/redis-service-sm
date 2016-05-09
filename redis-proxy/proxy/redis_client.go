package proxy

// import (
// 	"fmt"
// 	"net"
// 	"os"
// 	"time"
// )

// func tcpConnect(host string, port string) (*net.TCPConn, error) {
// 	server := host + ":" + port
// 	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
// 		return nil, err
// 	}
// 	conn, err := net.DialTCP("tcp", nil, tcpAddr)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
// 		return nil, err
// 	}

// 	return conn, nil
// }

// func tcpSend(conn *net.TCPConn, words string, time_out time.Duration) error {
// 	if conn == nil {
// 		return nil
// 	}
// 	err := conn.SetWriteDeadline(time.Now().Add(time_out))
// 	if err != nil {
// 		return err
// 	}
// 	_, err = conn.Write([]byte(words))
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func tcpReceiver(conn *net.TCPConn, time_out time.Duration) (string, error) {
// 	if conn == nil {
// 		return "", nil
// 	}
// 	err := conn.SetReadDeadline(time.Now().Add(time_out))
// 	if err != nil {
// 		return "", err
// 	}
// 	for {
// 		buffer := make([]byte, BUFFER_SIZE)
// 		res := ""
// 		for {
// 			n, err := conn.Read(buffer)
// 			if err != nil {
// 				fmt.Println(conn.RemoteAddr().String(), " connection error: ", err.Error())
// 				return res, err
// 			}
// 			res += string(buffer[:n])
// 			if n < BUFFER_SIZE {
// 				break
// 			}
// 		}
// 		return res, nil
// 	}
// }

// func tcpClose(conn *net.TCPConn) {
// 	if conn == nil {
// 		return
// 	}
// 	connFromAddr := conn.RemoteAddr().String()
// 	conn.Close()
// 	fmt.Println("Closed connection:", connFromAddr)
// }
