package helpers

import (
	"hash/crc32"
	"log"
	"net"
)

func Crc_hash(data string, table *crc32.Table) uint32 {
	byte_data := []byte(data)
	hashed_url := crc32.Checksum(byte_data, table)
	return hashed_url

}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
