package gattService

import (
	"github.com/paypal/gatt"
	"log"
	"time"
)

var UART_SERVICE_UUID = "6E400001-B5A3-F393-E0A9-E50E24DCCA9E"
var UART_RX_CHAR_UUID = "6E400002-B5A3-F393-E0A9-E50E24DCCA9E" // RX Characteristic (Property = Notify)
var UART_TX_CHAR_UUID = "6E400003-B5A3-F393-E0A9-E50E24DCCA9E" // TX Characteristic (Property = Write without response)

var echoData = make(chan string)

func NewUartService() *gatt.Service {

	s := gatt.NewService(gatt.MustParseUUID(UART_SERVICE_UUID))

	// Mobile -> Raspberry
	s.AddCharacteristic(gatt.MustParseUUID(UART_RX_CHAR_UUID)).HandleWriteFunc(
		func(r gatt.Request, data []byte) (status byte) {
			log.Println("[RX] Wrote:", string(data))

			echoData <- string(data)

			return gatt.StatusSuccess
		})

	// Raspberry -> Mobile
	s.AddCharacteristic(gatt.MustParseUUID(UART_TX_CHAR_UUID)).HandleNotifyFunc(
		func(r gatt.Request, n gatt.Notifier) {
			log.Println("notification Start!")
			for !n.Done() {
				select {
				case d := <-echoData:
					sentLen, err := n.Write([]byte(d))
					if err != nil {
						log.Println(err)
					}
					log.Printf("[TX] send Data... : %s, length : %d", d, sentLen)
				default:
					time.Sleep(time.Second)
				}
			}
			log.Println("notification Done!")
		})

	return s
}
