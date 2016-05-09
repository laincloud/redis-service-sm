package monitor

import (
	"fmt"
	"github.com/mijia/sweb/log"
	"net"
	"time"
)

const (
	PLAINTEXT_FORMAT = "%s %s %s\n"
	TIMEOUT_SEC      = 5 * time.Second
)

var (
	need_report = true
)

func ReportData(key, value, timestamp string) {
	if !need_report {
		return
	}

	conn, err := buildConn()
	if err != nil {
		log.Warnf("Connect graphite fatal error: %s\n", err.Error())
		return
	}
	defer closeConn(conn)

	data := fmt.Sprintf(PLAINTEXT_FORMAT, key, value, timestamp)
	report(conn, data)
}

func ReportDatas(datas map[string]string, timestamp string) {
	if !need_report {
		return
	}
	conn, err := buildConn()
	if err != nil {
		log.Warnf("Connect graphite fatal error: %s\n", err.Error())
		return
	}
	defer closeConn(conn)
	for key, value := range datas {
		data := fmt.Sprintf(PLAINTEXT_FORMAT, key, value, timestamp)
		report(conn, data)
	}
}

func buildConn() (net.Conn, error) {
	server := GRAPHITE_DOMAIN + ":" + GRAPHITE_PORT
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func closeConn(conn net.Conn) error {
	if conn == nil {
		return nil
	}
	return conn.Close()
}

func report(conn net.Conn, data string) error {
	if conn == nil {
		return nil
	}

	err := conn.SetWriteDeadline(time.Now().Add(TIMEOUT_SEC))
	if err != nil {
		return err
	}
	_, err = conn.Write([]byte(data))
	if err != nil {
		return err
	}
	return nil
}
