package main

import (
	"fmt"
	"gatt-echo-server/gattService"
	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
	"github.com/paypal/gatt/examples/service"
)

func main() {

	gattServer, err := gatt.NewDevice(option.DefaultServerOptions...)
	if err != nil {
		panic(err)
	}

	// central conn/disconn handler
	gattServer.Handle(
		gatt.CentralConnected(func(c gatt.Central) { fmt.Println("Connect: ", c.ID()) }),
		gatt.CentralDisconnected(func(c gatt.Central) {
			fmt.Println("Disconnect: ", c.ID())
		}),
	)

	gattServer.Init(onStateChanged)

	select {}
}

func onStateChanged(gattServer gatt.Device, s gatt.State) {

	fmt.Printf("State: %s\n", s)
	switch s {
	case gatt.StatePoweredOn:
		// Setup GAP and GATT services
		gattServer.AddService(service.NewGapService("HYERIN001")) // 'Generic Access' Service
		gattServer.AddService(service.NewGattService())           // 'Generic Attribute' Service

		// add uart service
		uartService := gattService.NewUartService()
		gattServer.AddService(uartService)

		// Advertise device name and service's UUIDs.
		gattServer.AdvertiseNameAndServices("HYERIN001", []gatt.UUID{uartService.UUID()})

	default:
		fmt.Print("default????")
	}

}
